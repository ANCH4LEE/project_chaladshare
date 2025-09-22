import { useParams } from "react-router-dom";
import { useNavigate } from "react-router-dom";
import React, { useState } from "react";
import "../component/PostDetail.css";
import {FaHeart,FaArrowLeft,FaShareAlt,FaBookmark,FaChevronLeft,FaChevronRight,} from "react-icons/fa";
import Sidebar from "./Sidebar";

const PostDetail = () => {
  const { id } = useParams();
  const navigate = useNavigate();

  // mock data: array ของรูป
  const images = ["/img/1.jpg", "/img/12.jpg", "/img/13.jpg"];

  const [currentIndex, setCurrentIndex] = useState(0);

  const prevImage = () => {
    setCurrentIndex((prev) => (prev === 0 ? images.length - 1 : prev - 1));
  };

  const nextImage = () => {
    setCurrentIndex((prev) => (prev === images.length - 1 ? 0 : prev + 1));
  };

  return (
    <div className="post-detail">
      <Sidebar />

      <div className="back-btn" onClick={() => navigate("/home")} style={{ cursor: "pointer" }}>
      <FaArrowLeft />
    </div>

      {/* โปรไฟล์ + ปุ่ม share/bookmark */}
      <div className="user-info">
        <img src="/img/author2.jpg" alt="profile" className="profile-img" />
        <div className="user-details">
          <h4>Apinya Saeaeung</h4>
          <p className="status">สาธารณะ 🌐</p>
        </div>
      </div>

      {/* รูปโพสต์ (carousel) */}
      <div className="post-image">
        <img src={images[currentIndex]} alt="summary" />

        {/* ปุ่มซ้าย/ขวา */}
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
      <div className="likes">
        <FaHeart color="red" /> <span>1006</span>
      </div>
      <div className="post-actions">
        <FaShareAlt size={20} className="action-icon" />
        <FaBookmark size={20} className="action-icon" />
      </div>

      {/* title */}
      <h3 className="post-title">SE - UML</h3>
      <p className="description">วิชา SE (Software engineer) สรุปเกี่ยวกับ UML ที่มี Class Diagram, Use Case Diagram, Sequence Diagram</p>
    </div>
  );
};

export default PostDetail;
