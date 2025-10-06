import React, { useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import { IoSearch } from "react-icons/io5";
import "../component/Home.css";
import Sidebar from "./Sidebar";
import PostCard from "../component/Postcard";
import RankingCard from "../component/RankingCard";

const Home = () => {
  // ข้อมูลโพสต์ยอดนิยม
  const [popularPosts, setPopularPosts] = useState([
    {
      img: "img/1.jpg",
      likes: 123,
      title: "UML",
      tags: "#SE #softwareengineer #UML",
      authorName: "Anchalee",
      authorImg: "img/author2.jpg",
    },
    {
      img: "img/2.jpg",
      likes: 350,
      title: "PM - Project Management",
      tags: "#IT #PM #ProjectManagement",
      authorName: "Benjaporn",
      authorImg: "img/author2.jpg",
    },
    {
      img: "img/3.jpg",
      likes: 2890,
      title: "Software Testing",
      tags: "#SWtest #Req #functionalTesting",
      authorName: "Chaiwat",
      authorImg: "img/author2.jpg",
    },
  ]);

  // ข้อมูลแนะนำสรุปน่าอ่าน
  const [recommendedPosts, setRecommendedPosts] = useState([
    {
      img: "img/4.jpg",
      likes: 1006,
      title: "Security - planning",
      tags: "#ISS #plannimg #Security",
      authorName: "Benjaporn",
      authorImg: "img/author2.jpg",
    },
    {
      img: "img/5.jpg",
      likes: 875,
      title: "basic storytelling",
      tags: "#storytelling #intro #JavaScript",
      authorName: "Chaiwat",
      authorImg: "img/author2.jpg",
    },
    {
      img: "img/6.jpg",
      likes: 875,
      title: "basic JavaScript",
      tags: "#js #FE #frontend",
      authorName: "Chaiwat",
      authorImg: "img/author2.jpg",
    },
  ]);

  const navigate = useNavigate();
  const goToPostDetail = (post) => {
    navigate(`/post/${post.title}`); // ตอนนี้ใช้ title เป็น id mock
  };

  // ✅ เรียงโพสต์ยอดนิยมจากไลก์มาก→น้อย แล้วแปะ rank 1..N
  const rankedPopular = useMemo(() => {
    return popularPosts
      .slice() // กัน side-effect ไม่แก้ array เดิม
      .sort((a, b) => b.likes - a.likes);
  }, [popularPosts]);

  return (
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
        <div className="card-list">
          {recommendedPosts.map((post, index) => (
            <div
              key={index}
              onClick={() => goToPostDetail(post)}
              style={{ cursor: "pointer" }}
            >
              <PostCard post={post} />
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default Home;
