// หน้า Home.jsx (ทำ prefix แล้ว)

import React, { useMemo, useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import Sidebar from "./Sidebar";
import { IoSearch } from "react-icons/io5";

import PostCard from "../component/Postcard";
import RankingCard from "../component/RankingCard";
import author2 from "../assets/author2.jpg";
import one from "../assets/one.jpg";
import two from "../assets/two.jpg";
import three from "../assets/three.jpg";

import "../component/Home.css";

const API_HOST = "http://localhost:8080";
const toAbsUrl = (p) => {
  if (!p) return "";
  if (p.startsWith("http")) return p;
  const clean = p.replace(/^\./, "");
  return `${API_HOST}${clean.startsWith("/") ? clean : `/${clean}`}`;
};

const Home = () => {
  // ข้อมูลโพสต์ยอดนิยม
  const [popularPosts, setPopularPosts] = useState([
    {
      img: one,
      likes: 123,
      title: "UML",
      tags: "#SE #softwareengineer #UML",
      authorName: "Anchalee",
      authorImg: author2,
    },
    {
      img: two,
      likes: 350,
      title: "PM - Project Management",
      tags: "#IT #PM #ProjectManagement",
      authorName: "Benjaporn",
      authorImg: author2,
    },
    {
      img: three,
      likes: 2890,
      title: "Software Testing",
      tags: "#SWtest #Req #functionalTesting",
      authorName: "Chaiwat",
      authorImg: author2,
    },
  ]);

  // เรียงโพสต์ยอดนิยมจากไลก์มาก→น้อย แล้วแปะ rank 1..N
  const rankedPopular = useMemo(() => {
    return popularPosts
      .slice() // กัน side-effect ไม่แก้ array เดิม
      .sort((a, b) => b.likes - a.likes);
  }, [popularPosts]);

  // ข้อมูลแนะนำสรุปน่าอ่าน
  const [recommendedPosts, setRecommendedPosts] = useState([]);
  const [loadingRec, setLoadingRec] = useState(true);
  const [recErr, setRecErr] = useState("");
  const navigate = useNavigate();

  useEffect(() => {
    let cancelled = false;

    (async () => {
      try {
        setLoadingRec(true);
        setRecErr("");

        // ต้องล็อกอินถึงจะได้ 200
        const res = await axios.get("/posts");
        const rows = Array.isArray(res?.data?.data)
          ? res.data.data
          : Array.isArray(res?.data)
          ? res.data
          : [];

        const mapped = rows.map((p) => {
          const rawUrl = p.file_url || null;
          const isPdf = /\.pdf$/i.test(rawUrl || "");
          const isImage = /\.(jpg|jpeg|png|gif|webp)$/i.test(rawUrl || "");
          return {
            id: p.post_id,
            img:
              rawUrl && !isPdf && isImage
                ? toAbsUrl(rawUrl) // ถ้าเป็นรูปภาพ ก็ใช้
                : "/img/pdf-placeholder.jpg", // ถ้าไม่ใช่ (เช่น /@preview) ก็ใช้ placeholder
            isPdf,
            document_url: isPdf ? toAbsUrl(rawUrl) : null,
            documentId: p.post_document_id,
            likes: p.like_count ?? 0,
            title: p.post_title,
            tags: Array.isArray(p.tags)
              ? p.tags.map((t) => (t.startsWith("#") ? t : `#${t}`)).join(" ")
              : "",
            authorName: p.author_name || "ไม่ระบุ",
            authorImg: author2,
          };
        });

        if (!cancelled) setRecommendedPosts(mapped);
      } catch (e) {
        if (!cancelled) {
          if (e?.response?.status === 401) {
            navigate("/login", { replace: true });
            return;
          }
          setRecErr(
            e?.response?.data?.error || e.message || "โหลดข้อมูลล้มเหลว"
          );
        }
      } finally {
        if (!cancelled) setLoadingRec(false);
      }
    })();

    return () => {
      cancelled = true;
    };
  }, [navigate]);

  const goToPostDetail = (post) => {
    if (post?.id) navigate(`/posts/${post.id}`);
  };

  return (
    <div className="home-page">
      <div className="home-container">
        {/* Sidebar */}
        <Sidebar />

        {/* เนื้อหาหลัก */}
        <div className="home">
          {/* Search bar */}
          <div className="search-bar">
            <input type="text" placeholder="ค้นหาความสนใจของคุณ" />
            <IoSearch />
          </div>

          {/* โพสต์ยอดนิยม */}
          <h3>โพสต์สรุปยอดเยี่ยมประจำเดือน</h3>
          <div className="card-list">
            {rankedPopular.map((post, index) => (
              <div
                key={index}
                onClick={() => goToPostDetail(post)}
                style={{ cursor: "default" }}
              >
                <RankingCard post={post} rank={index + 1} />
              </div>
            ))}
          </div>

          {/* แนะนำสรุปน่าอ่าน */}
          <h3>แนะนำสรุปน่าอ่าน</h3>
          {loadingRec && <div>กำลังโหลด...</div>}
          {recErr && <div style={{ color: "#b00020" }}>{recErr}</div>}
          <div className="card-list">
            {!loadingRec &&
              !recErr &&
              recommendedPosts.map((post, index) => (
                <div
                  key={post.id || index}
                  onClick={() => goToPostDetail(post)}
                  style={{ cursor: "pointer" }}
                >
                  <PostCard post={post} />
                </div>
              ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default Home;
