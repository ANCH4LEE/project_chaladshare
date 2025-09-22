import React from "react";
import { useNavigate } from "react-router-dom";
import "../component/Home.css"; // import CSS ของ Home
import Sidebar from "./Sidebar";

const Home = () => {
  // ข้อมูลโพสต์ยอดนิยม
  const popularPosts = [
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
  ];

  // ข้อมูลแนะนำสรุปน่าอ่าน
  const recommendedPosts = [
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
  ];

  const navigate = useNavigate(); // เพิ่มตรงนี้ข้างบน goToPostDetail
  const goToPostDetail = (post) => {
    if (post.title === "UML") {
      navigate(`/post/${post.title}`); // หรือใช้ postId ก็ได้
    }
  };

  return (
    <div className="home-container">
      {/* Sidebar */}
      <Sidebar />

      {/* เนื้อหาหลัก */}
      <div className="home">
        {/* Search bar */}
        <div className="search-bar">
          <input type="text" placeholder="ค้นหาความสนใจของคุณ" />
          <i className="fas fa-search"></i>
        </div>

        {/* โพสต์ยอดนิยม */}
        <h3>โพสต์สรุปยอดเยี่ยมประจำเดือน</h3>
        <div className="card-list">
          {popularPosts.map((post, index) => (
            <div
              className="card"
              key={index}
              onClick={() => goToPostDetail(post)}
              style={{ cursor: post.title === "UML" ? "pointer" : "default" }}
            >
              {/* โปรไฟล์อยู่ด้านบน */}
              <div className="card-header">
                <img src={post.authorImg} alt="author" className="author-img" />
                <span>{post.authorName}</span>
              </div>
              <img src={post.img} alt="summary" />
              <div className="card-body">
                <span className="likes">❤️ {post.likes}</span>
                <h4>{post.title}</h4>
                <p>{post.tags}</p>
              </div>
            </div>
          ))}
        </div>

        {/* แนะนำสรุปน่าอ่าน */}
        <h3>แนะนำสรุปน่าอ่าน</h3>
        <div className="card-list">
          {recommendedPosts.map((post, index) => (
            <div className="card" key={index}>
              <img src={post.img} alt="summary" />
              <div className="card-body">
                <span className="likes">❤️ {post.likes}</span>
                <h4>{post.title}</h4>
                <p>{post.tags}</p>
                {/* โปรไฟล์อยู่ด้านล่าง */}
                <div className="card-footer">
                  <img
                    src={post.authorImg}
                    alt="author"
                    className="author-img"
                  />
                  <span>{post.authorName}</span>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default Home;
