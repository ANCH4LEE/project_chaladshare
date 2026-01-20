// หน้า Home.jsx (ทำ prefix แล้ว)

import React, { useMemo, useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import Sidebar from "./Sidebar";
import { IoSearch } from "react-icons/io5";

import PostCard from "../component/Postcard";
import RankingCard from "../component/RankingCard";
import Avatar from "../assets/default.png";
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

const mapToCardPost = (p) => {
  const coverRaw = p.cover_url || p.post_cover_url || "";
  const coverImg = coverRaw ? toAbsUrl(coverRaw) : "/img/pdf-placeholder.jpg";

  const rawUrl = p.file_url || p.document_url || "";
  const isPdf = /\.pdf$/i.test(rawUrl || "");

  const avatarRaw = p.avatar_url || p.author_img || "";
  const authorImg = avatarRaw ? toAbsUrl(avatarRaw) : Avatar;

  const tagsRaw = p.tags ?? "";
  const tags =
    Array.isArray(tagsRaw)
      ? tagsRaw.map((t) => (t.startsWith("#") ? t : `#${t}`)).join(" ")
      : typeof tagsRaw === "string"
      ? tagsRaw
          .split(",")
          .map((t) => t.trim())
          .filter(Boolean)
          .map((t) => (t.startsWith("#") ? t : `#${t}`))
          .join(" ")
      : "";

  return {
    id: p.post_id ?? p.PostID ?? p.id,
    post_id: p.post_id ?? p.PostID ?? p.id,

    img: coverImg,

    isPdf,
    document_url: isPdf && rawUrl ? toAbsUrl(rawUrl) : null,
    documentId: p.post_document_id ?? p.document_id,

    likes: p.like_count ?? p.likeCount ?? 0,
    like_count: p.like_count ?? p.likeCount ?? 0,
    is_liked: !!(p.is_liked ?? p.isLiked),
    is_saved: !!(p.is_saved ?? p.isSaved),

    title: p.post_title ?? p.title ?? p.Title ?? "",
    tags,

    authorName: p.author_name ?? p.authorName ?? "ไม่ระบุ",
    authorImg,
  };
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

  // แนะนำ
  const [recommendedPosts, setRecommendedPosts] = useState([]);
  const [loadingRec, setLoadingRec] = useState(true);
  const [recErr, setRecErr] = useState("");

  // โพสต์ทั้งหมด
  const [allPosts, setAllPosts] = useState([]);
  const [loadingAll, setLoadingAll] = useState(true);
  const [allErr, setAllErr] = useState("");
  const navigate = useNavigate();

  // โหลดแนะนำ
  useEffect(() => {
    let cancelled = false;

    (async () => {
      try {
        setLoadingRec(true);
        setRecErr("");

        const res = await axios.get("/recommend", {
          params: { limit: 3 },
          withCredentials: true,
        });

        const rows = Array.isArray(res?.data?.data) ? res.data.data : [];
        const mapped = rows.map(mapToCardPost);

        if (!cancelled) setRecommendedPosts(mapped);
      } catch (e) {
        if (!cancelled) {
          if (e?.response?.status === 401) {
            navigate("/login", { replace: true });
            return;
          }
          setRecErr(
            e?.response?.data?.error || e.message || "โหลดแนะนำล้มเหลว"
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

  // โหลดโพสต์ทั้งหมด
  useEffect(() => {
    let cancelled = false;

    (async () => {
      try {
        setLoadingAll(true);
        setAllErr("");

        const res = await axios.get("/posts", { withCredentials: true });

        const rows = Array.isArray(res?.data?.data)
          ? res.data.data
          : Array.isArray(res?.data)
          ? res.data
          : [];

        const mapped = rows.map(mapToCardPost);

        if (!cancelled) setAllPosts(mapped);
      } catch (e) {
        if (!cancelled) {
          if (e?.response?.status === 401) {
            navigate("/login", { replace: true });
            return;
          }
          setAllErr(
            e?.response?.data?.error || e.message || "โหลดโพสต์ทั้งหมดล้มเหลว"
          );
        }
      } finally {
        if (!cancelled) setLoadingAll(false);
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
          <h3>โพสต์สรุปยอดเยี่ยม</h3>
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
          {loadingRec && <div>กำลังโหลดแนะนำ...</div>}
          {recErr && <div style={{ color: "#b00020" }}>{recErr}</div>}
          <div className="card-list">
            {!loadingRec &&
              !recErr &&
              recommendedPosts.map((post) => (
                <div
                  key={post.id}
                  onClick={() => goToPostDetail(post)}
                  style={{ cursor: "pointer" }}
                >
                  <PostCard post={post} />
                </div>
              ))}
          </div>

          {/* โพสต์ทั้งหมด*/}
          <h3>โพสต์ทั้งหมด</h3>
          {loadingAll && <div>กำลังโหลดโพสต์ทั้งหมด...</div>}
          {allErr && <div style={{ color: "#b00020" }}>{allErr}</div>}
          <div className="card-list">
            {!loadingAll &&
              !allErr &&
              allPosts.map((post) => (
                <div
                  key={post.id}
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
