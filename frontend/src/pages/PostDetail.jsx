// หน้า PostDetail.jsx (ทำ prefix แล้ว)

import React, { useState, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
import axios from "axios";
import { FaArrowLeft } from "react-icons/fa";
import { AiFillHeart, AiOutlineHeart } from "react-icons/ai";

import Sidebar from "./Sidebar";
import "../component/PostDetail.css";

const API_HOST = "http://localhost:8080";

const toAbsUrl = (p) => {
  if (!p) return "";
  if (p.startsWith("http")) return p;
  const clean = p.replace(/^\./, "");
  return `${API_HOST}${clean.startsWith("/") ? clean : `/${clean}`}`;
};

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
        setLoading(true);
        setErr("");

        const res = await axios.get(`/posts/${id}`);
        const payload = res?.data?.data ?? res?.data ?? {};
        const data = payload?.post ?? payload ?? {};
        if (!data || (!data.post_id && !data.id)) {
          setErr("ไม่พบโพสต์");
          return;
        }

        const mapped = {
          id: data.post_id,
          title: data.post_title,
          description: data.post_description,
          visibility: data.post_visibility,
          file_url: data.file_url ? toAbsUrl(data.file_url) : null,
          author_name: data.author_name,
          like_count: data.like_count,
          is_liked: data.is_liked,
          is_saved: data.is_saved,
          tags: data.tags || [],
          post_document_id: data.post_document_id,
        };

        setPost(mapped);
        setLikes(mapped.like_count || 0);
        setLiked(!!mapped.is_liked);
      } catch (e) {
        const st = e?.response?.status;
        if (st === 403) setErr("คุณไม่มีสิทธิ์ดูโพสต์นี้");
        else if (st === 404) setErr("ไม่พบโพสต์");
        else
          setErr(e?.response?.data?.error || e.message || "โหลดโพสต์ล้มเหลว");
      } finally {
        setLoading(false);
      }
    })();
  }, [id, navigate]);

  // toggle like
  const toggleLike = async () => {
    try {
      if (liked) {
        await axios.delete(`/posts/${id}/like`);
        setLikes((prev) => prev - 1);
      } else {
        await axios.post(`/posts/${id}/like`);
        setLikes((prev) => prev + 1);
      }
      setLiked(!liked);
    } catch (e) {
      const st = e?.response?.status;
      if (st === 403) setErr("ไม่มีสิทธิ์กดไลก์โพสต์นี้");
      else if (st === 404) setErr("ไม่พบโพสต์");
      else console.error("Like toggle failed:", e);
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

  const isPdf =
    Boolean(post.post_document_id) || /\.pdf$/i.test(post.file_url || "");
  const visibilityText =
    post.visibility === "friends" ? "เฉพาะเพื่อน" : "สาธารณะ";

  return (
    <div className="post-detail-page">
      <div className="post-detail">
        <Sidebar />

        {/* button back */}
        <div
          className="back-btn"
          onClick={() => navigate(-1)}
          style={{ cursor: "pointer" }}
        >
          <FaArrowLeft />
        </div>

        {/* profile */}
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

        {/* post */}
        <div className="post-image">
          {post.file_url ? (
            <object
              data={post.file_url}
              type="application/pdf"
              width="100%"
              height="800"
            >
              <iframe
                src={`${post.file_url}#view=FitH`}
                width="100%"
                height="800"
                title="pdf"
              />
            </object>
          ) : (
            <img src="/img/no-image.png" alt={post.title} />
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
