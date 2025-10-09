// หน้า Home.jsx (ทำ prefix แล้ว)

import React, { useMemo, useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import Sidebar from "./Sidebar";
import { IoSearch } from "react-icons/io5";

import PostCard from "../component/Postcard";
import RankingCard from "../component/RankingCard";
import author2 from "../assets/author2.jpg"
import one from "../assets/one.jpg"
import two from "../assets/two.jpg"
import three from "../assets/three.jpg"

import "../component/Home.css";

const API_BASE = "http://localhost:8080";

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
  const [recommendedPosts, setRecommendedPosts] = useState([])
  const [loadingRec, setLoadingRec] = useState(true);
  const [recErr, setRecErr] = useState("");

  useEffect(() => {
    (async () => {
      try {
        const res = await axios.get(`${API_BASE}/api/v1/posts`);
        // map ข้อมูลจาก BE -> รูปแบบที่ PostCard ใช้
        const list = (res.data || []).map(p => ({
          id: p.post_id,
          img: p.file_url || "/img/pdf-placeholder.jpg", // ถ้าเป็น .pdf เดี๋ยวค่อยทำ thumbnail ภายหลัง
          likes: p.like_count ?? 0,
          title: p.post_title,
          // แปะ # ทุกแท็ก (ถ้า BE คืนมาเป็น array)
          tags: Array.isArray(p.tags) ? p.tags.map(t => (t.startsWith("#") ? t : `#${t}`)).join(" ") : "",
          authorName: p.author_name || "ไม่ระบุ",
          authorImg: "img/author2.jpg",
        }));
        setRecommendedPosts(list);
      } catch (e) {
        setRecErr(e?.response?.data?.error || e.message);
      } finally {
        setLoadingRec(false);
      }
    })();
  }, []);

  const navigate = useNavigate();
  const goToPostDetail = (post) => {
    navigate(`/post/${post.id ?? post.title}`);
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
              style={{ cursor: "pointer" }}
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
          {!loadingRec && !recErr && recommendedPosts.map((post, index) => (
            <div key={post.id || index} onClick={() => goToPostDetail(post)} style={{ cursor: "pointer" }}>
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
