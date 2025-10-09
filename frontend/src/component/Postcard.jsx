import React, { useState } from "react";
import { AiFillHeart, AiOutlineHeart } from "react-icons/ai";
import { BsBookmark, BsBookmarkFill } from "react-icons/bs";
import { FiShare2 } from "react-icons/fi";

const PostCard = ({ post }) => {
  const [liked, setLiked] = useState(false);
  const [likes, setLikes] = useState(post.likes);
  const [saved, setSaved] = useState(false);
  const [toast, setToast] = useState("");

  const toggleLike = (e) => {
    e.stopPropagation();
    setLikes((n) => (liked ? n - 1 : n + 1));
    setLiked((v) => !v);
  };

  const handleSave = (e) => {
    e.stopPropagation();
    const next = !saved;
    setSaved(next);
    setToast(next ? "✔️ รายการที่บันทึก" : "❌ ยกเลิกการบันทึก");
    setTimeout(() => setToast(""), 2500);
  };

  const sharePost = async (e) => {
    e.stopPropagation();
    const url = window.location.origin + "/post/" + encodeURIComponent(post.title);
    try {
      if (navigator.share) {
        await navigator.share({ title: post.title, text: "ดูสรุปนี้บน ChaladShare", url });
      } else {
        await navigator.clipboard.writeText(url);
        setToast("📋 คัดลอกลิงก์แล้ว");
        setTimeout(() => setToast(""), 3000);
      }
    } catch {}
  };

  return (
    <div className="card">
      <div className="card-header">
        <img src={post.authorImg} alt="author" className="author-img" />
        <span>{post.authorName}</span>
      </div>

      <img src={post.img} alt="summary" className="card-image" />


      <div className="card-body">
        <div className="actions-row" onClick={(e) => e.stopPropagation()}>
          <span className="likes" onClick={toggleLike} style={{ cursor: "pointer" }}>
            {liked ? (
              <AiFillHeart style={{ color: "red", fontSize: "20px" }} />
            ) : (
              <AiOutlineHeart style={{ color: "black", fontSize: "20px" }} />
            )}
            {likes}
          </span>

          <div className="action-right">
            <button
              className={`icon-btn ${saved ? "active" : ""}`}
              onClick={handleSave}
              aria-label="save"
              title={saved ? "ยกเลิกบันทึก" : "บันทึก"}
            >
              {saved ? <BsBookmarkFill /> : <BsBookmark />}
            </button>

            <button className="icon-btn" onClick={sharePost} aria-label="share" title="แชร์">
              <FiShare2 />
            </button>
          </div>
        </div>

        <h4>{post.title}</h4>
        <p>{post.tags}</p>
      </div>

      {/* ✅ Toast แยกออกจาก flow การ hover */}
      <div className="toast-container">
        {toast && <div className="mini-toast">{toast}</div>}
      </div>
    </div>
  );
};

export default PostCard;

