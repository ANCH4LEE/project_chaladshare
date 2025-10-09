// หน้า PostDetail.jsx (ทำ prefix แล้ว)

import React, { useState, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
import axios from "axios";
import { FaArrowLeft } from "react-icons/fa";
import { AiFillHeart, AiOutlineHeart } from "react-icons/ai";

import Sidebar from "./Sidebar";
import "../component/PostDetail.css";

const API_BASE = "http://localhost:8080/api/v1";

const PostDetail = () => {
  const { id } = useParams();
  const navigate = useNavigate();

  const [post, setPost] = useState(null);
  const [liked, setLiked] = useState(false);
  const [likes, setLikes] = useState(0);
  const [loading, setLoading] = useState(true);
  const [err, setErr] = useState("");

  // ข้อมูลโพสต์
  useEffect(() => {
    (async () => {
      try {
        const res = await axios.get(`${API_BASE}/posts/${id}`);
        const data = res.data;

        const mapped = {
        id: data.post_id,
        title: data.post_title,
        description: data.post_description,
        visibility: data.post_visibility,
        file_url: data.file_url,
        author_name: data.author_name,
        like_count: data.like_count,
        is_liked: data.is_liked,
        is_saved: data.is_saved,
        tags: data.tags || [],
      };

        setPost(mapped);
        setLikes(mapped.like_count || 0);
        setLiked(mapped.is_liked || false); // ถ้ามีส่งจาก backend
      } catch (e) {
        setErr(e?.response?.data?.error || e.message);
      } finally {
        setLoading(false);
      }
    })();
  }, [id]);

  // toggle like
  const toggleLike = async () => {
    try {
      if (liked) {
        await axios.delete(`${API_BASE}/posts/${id}/like`);
        setLikes((prev) => prev - 1);
      } else {
        await axios.post(`${API_BASE}/posts/${id}/like`);
        setLikes((prev) => prev + 1);
      }
      setLiked(!liked);
    } catch (e) {
      console.error("Like toggle failed:", e);
    }
  };

  if (loading)
    return (
      <div className="post-detail">
        <Sidebar />
        <div style={{ padding: 24 }}>กำลังโหลด…</div>
      </div>
    );

  if (err)
    return (
      <div className="post-detail">
        <Sidebar />
        <div style={{ padding: 24, color: "#b00020" }}>{err}</div>
      </div>
    );

  if (!post)
    return (
      <div className="post-detail">
        <Sidebar />
        <div style={{ padding: 24 }}>ไม่พบโพสต์</div>
      </div>
    );

  const isPdf = post.file_url?.toLowerCase().endsWith(".pdf");
  const visibilityText =
    post.post_visibility === "friends" ? "เฉพาะเพื่อน" : "สาธารณะ";

  return (
    <div className="post-detail-page">
    <div className="post-detail">
      <Sidebar />

      {/* ปุ่มย้อนกลับ */}
      <div
        className="back-btn"
        onClick={() => navigate("/home")}
        style={{ cursor: "pointer" }}
      >
        <FaArrowLeft />
      </div>

      {/* โปรไฟล์ */}
      <div className="user-info">
        <img
          src={post.author_profile || "/img/default-profile.png"}
          alt="profile"
          className="profile-img"
        />
        <div className="user-details">
          <h4>{post.author_name || "ไม่ระบุ"}</h4>
          <p className="status">{visibilityText}</p>
        </div>
      </div>

      {/* รูปหรือไฟล์โพสต์ */}
      <div className="post-image">
        {isPdf ? (
          <iframe
            src={post.file_url}
            title="PDF Viewer"
            width="100%"
            height="500px"
          ></iframe>
        ) : (
          <img
            src={post.file_url || "/img/no-image.png"}
            alt="summary"
            className="post-img"
          />
        )}
      </div>

      {/* ปุ่มไลก์ */}
      <div className="detail-likes" onClick={toggleLike}>
        {liked ? (
          <AiFillHeart style={{ color: "red", fontSize: "20px" }} />
        ) : (
          <AiOutlineHeart style={{ color: "black", fontSize: "20px" }} />
        )}
        <span>{likes}</span>
      </div>

      {/* title + description */}
      <h3 className="post-title">{post.title}</h3>
      <p className="description">{post.description}</p>
    </div>
    </div>
  );
};

export default PostDetail;

