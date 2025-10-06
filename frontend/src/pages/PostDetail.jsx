import { useParams } from "react-router-dom";
import { useNavigate } from "react-router-dom";
import React, { useState } from "react";
import "../component/PostDetail.css";
import { FaArrowLeft, FaChevronLeft, FaChevronRight, } from "react-icons/fa";
import { AiFillHeart, AiOutlineHeart } from "react-icons/ai";

import Sidebar from "./Sidebar";

const PostDetail = () => {
  const { id } = useParams();
  const navigate = useNavigate();

  // mock data: array ของรูป
  const images = ["/img/1.jpg", "/img/12.jpg", "/img/13.jpg"];

  const [currentIndex, setCurrentIndex] = useState(0);
  const [liked, setLiked] = useState(false);
  const [likes, setLikes] = useState(1006);

  const prevImage = () => {
    setCurrentIndex((prev) => (prev === 0 ? images.length - 1 : prev - 1));
  };

  const nextImage = () => {
    setCurrentIndex((prev) => (prev === images.length - 1 ? 0 : prev + 1));
  };

  const toggleLike = () => {
    if (liked) setLikes(likes - 1);
    else setLikes(likes + 1);
    setLiked(!liked);
  };

  return (
    <div className="post-detail">
      <Sidebar />

      <div
        className="back-btn"
        onClick={() => navigate("/home")}
        style={{ cursor: "pointer" }}
      >
        <FaArrowLeft />
      </div>

      {/* โปรไฟล์ */}
      <div className="user-info">
        <img src="/img/author2.jpg" alt="profile" className="profile-img" />
        <div className="user-details">
          <h4>Apinya Saeaeung</h4>
          <p className="status">สาธารณะ 🌐</p>
        </div>
      </div>

      {/* รูปโพสต์ */}
      <div className="post-image">
        <img src={images[currentIndex]} alt="summary" />

        <button className="nav-btn left" onClick={prevImage}>
          <FaChevronLeft />
        </button>
        <button className="nav-btn right" onClick={nextImage}>
          <FaChevronRight />
        </button>

        {/* จุดบอกตำแหน่ง */}
        <div className="dots">
          {images.map((_, i) => (
            <span
              key={i}
              className={`dot ${i === currentIndex ? "active" : ""}`}
              onClick={() => setCurrentIndex(i)}
            ></span>
          ))}
        </div>
      </div>

      {/* like */}
      <div className="detail-likes" onClick={toggleLike}>
        {liked ? (
          <AiFillHeart style={{ color: "red", fontSize: "20px" }} />
        ) : (
          <AiOutlineHeart style={{ color: "black", fontSize: "20px" }} />
        )}
        <span>{likes}</span>
      </div>

      {/* title */}
      <h3 className="post-title">SE - UML</h3>
      <p className="description">
        วิชา SE (Software engineer) สรุปเกี่ยวกับ UML ที่มี Class Diagram,
        Use Case Diagram, Sequence Diagram
      </p>
    </div>
  );
};

export default PostDetail;
