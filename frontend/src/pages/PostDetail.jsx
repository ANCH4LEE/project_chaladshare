import React, { useState, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { FaArrowLeft } from "react-icons/fa";
import { AiFillHeart, AiOutlineHeart } from "react-icons/ai";
import { FiShare2 } from "react-icons/fi";
import { BsBookmark, BsBookmarkFill } from "react-icons/bs";
import axios from "axios";
import Sidebar from "./Sidebar";
import Avatar from "../assets/default.png";
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
  const [saved, setSaved] = useState(false);
  const [loading, setLoading] = useState(true);
  const [err, setErr] = useState("");

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

        const avatarRaw = data.avatar_url || "";
        const hasRealAvatar = avatarRaw.startsWith("/uploads/");
        const authorImg = hasRealAvatar ? toAbsUrl(avatarRaw) : Avatar;

        const mapped = {
          id: data.post_id,
          title: data.post_title,
          description: data.post_description,
          visibility: data.post_visibility,
          file_url: data.file_url ? toAbsUrl(data.file_url) : null,
          author_name: data.author_name,
          author_id: data.author_id,
          authorImg,
          like_count: data.like_count,
          is_liked: data.is_liked,
          is_saved: data.is_saved,
          tags: data.tags || [],
          post_document_id: data.post_document_id,
        };

        setPost(mapped);
        setLikes(mapped.like_count || 0);
        setLiked(!!mapped.is_liked);
        setSaved(!!mapped.is_saved);
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

  // ✅ Optimistic like
  const toggleLike = async (e) => {
    e?.preventDefault?.();
    e?.stopPropagation?.();
    const next = !liked;
    setLiked(next);
    setLikes((p) => p + (next ? 1 : -1));
    try {
      if (next) await axios.post(`/posts/${id}/like`);
      else await axios.delete(`/posts/${id}/like`);
    } catch (error) {
      // revert เมื่อพลาด
      setLiked(!next);
      setLikes((p) => p + (next ? -1 : 1));
      const st = error?.response?.status;
      if (st === 403) setErr("ไม่มีสิทธิ์กดไลก์โพสต์นี้");
      else if (st === 404) setErr("ไม่พบโพสต์");
      console.error("Like toggle failed:", error);
    }
  };

  // ✅ Optimistic save
  const toggleSave = async (e) => {
    e?.preventDefault?.();
    e?.stopPropagation?.();
    const next = !saved;
    setSaved(next);
    try {
      if (next) await axios.post(`/posts/${id}/save`);
      else await axios.delete(`/posts/${id}/save`);
    } catch (error) {
      // revert เมื่อพลาด
      setSaved(!next);
      console.error("Save toggle failed:", error);
    }
  };

  const sharePost = async (e) => {
    e?.preventDefault?.();
    e?.stopPropagation?.();
    const url = window.location.href;
    try {
      if (navigator.share) {
        await navigator.share({
          title: post?.title || "ChaladShare",
          text: post?.description || "",
          url,
        });
      } else {
        await navigator.clipboard.writeText(url);
        alert("คัดลอกลิงก์แล้ว");
      }
    } catch {
      /* ผู้ใช้ยกเลิกการแชร์ */
    }
  };

  if (loading)
    return (
      <div className="post-detail-page">
        <Sidebar />
        <main className="post-detail">
          <div style={{ padding: 24 }}>กำลังโหลด…</div>
        </main>
      </div>
    );

  if (err)
    return (
      <div className="post-detail-page">
        <Sidebar />
        <main className="post-detail">
          <div style={{ padding: 24, color: "#b00020" }}>{err}</div>
        </main>
      </div>
    );

  if (!post)
    return (
      <div className="post-detail-page">
        <Sidebar />
        <main className="post-detail">
          <div style={{ padding: 24 }}>ไม่พบโพสต์</div>
        </main>
      </div>
    );

  const isPdf =
    Boolean(post.post_document_id) || /\.pdf$/i.test(post.file_url || "");
  const visibilityText =
    post.visibility === "friends" ? "เฉพาะเพื่อน" : "สาธารณะ";

  return (
    <div className="post-detail-page">
      <Sidebar />

      <main className="post-detail">
        {/* Header */}
        <header className="post-header">
          <button
            type="button"
            className="back-btn"
            onClick={() => navigate(-1)}
            aria-label="ย้อนกลับ"
          >
            <FaArrowLeft />
          </button>

          <div
            className="user-info"
            style={{ cursor: post.author_id ? "pointer" : "default" }}
            onClick={() =>
              post.author_id && navigate(`/profile/${post.author_id}`)
            }
            title={post.author_id ? "ไปที่โปรไฟล์ผู้เขียน" : undefined}
          >
            <img
              src={post.authorImg}
              alt="profile"
              className="profile-img"
            />
            <div className="user-details">
              <h4>{post.author_name || "ไม่ระบุ"}</h4>
              <p className="status">{visibilityText}</p>
            </div>
          </div>
        </header>

        {/* Viewer */}
        <section className="post-image">
          <div className="pdf-slide-wrapper">
            {post.file_url ? (
              isPdf ? (
                <iframe
                  className="pdf-frame"
                  src={`${post.file_url}#zoom=page-width`}
                  title="pdf"
                />
              ) : (
                <img
                  className="pdf-page-img active"
                  src={post.file_url}
                  alt={post.title}
                />
              )
            ) : (
              <img
                className="pdf-page-img active"
                src="/img/no-image.png"
                alt={post.title}
              />
            )}
          </div>
        </section>

        {/* Actions + Meta */}
        <section className="post-footer">
          <div
            className="actions-row"
            onClickCapture={(e) => e.stopPropagation()} // กัน bubble ออกไปนอกแถว
          >
            {/* ไลก์ (ซ้าย) แบบเดียวกับ Home */}
            <button
              type="button"
              className={`likes-btn ${liked ? "active" : ""}`}
              onClick={toggleLike}
              aria-label="กดไลก์"
            >
              {liked ? <AiFillHeart /> : <AiOutlineHeart />}
              <span>{likes}</span>
            </button>

            {/* ปุ่มขวา: แชร์ + บันทึก */}
            <div className="action-right">
              <button
                type="button"
                className="icon-btn"
                onClick={sharePost}
                title="แชร์"
                aria-label="แชร์"
              >
                <FiShare2 />
              </button>

              <button
                type="button"
                className={`icon-btn ${saved ? "active" : ""}`}
                onClick={toggleSave}
                title="บันทึก"
                aria-label="บันทึก"
              >
                {saved ? <BsBookmarkFill /> : <BsBookmark />}
              </button>
            </div>
          </div>

          <h3 className="post-title">{post.title}</h3>
          <p className="description">{post.description}</p>

          {/* แท็ก */}
          {post.tags && post.tags.length > 0 && (
            <div className="post-tags">
              {post.tags.map((t, i) => (
                <span className="tag" key={`${t}-${i}`}>
                  #{t}
                </span>
              ))}
            </div>
          )}
        </section>
      </main>
    </div>
  );
};

export default PostDetail;
