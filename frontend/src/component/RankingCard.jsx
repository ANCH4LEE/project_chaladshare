import React, { useState } from "react";
import { AiFillHeart, AiOutlineHeart } from "react-icons/ai";
import { BsBookmark, BsBookmarkFill } from "react-icons/bs";
import { FiShare2 } from "react-icons/fi";

import award  from "../icon/award.png";
import medal  from "../icon/medal.png";
import trophy from "../icon/trophy.png";

const RankingCard = ({ post, rank }) => {
  const [liked, setLiked]   = useState(false);
  const [likes, setLikes]   = useState(post.likes);
  const [saved, setSaved]   = useState(false);
  const [toast, setToast]   = useState("");

  const badge = rank === 1 ? award : rank === 2 ? trophy : medal;

  const toggleLike = (e) => {
    e.stopPropagation();
    setLikes(n => (liked ? n - 1 : n + 1));
    setLiked(v => !v);
  };

  const handleSave = (e) => {
    e.stopPropagation();
    const next = !saved;
    setSaved(next);
    setToast(next ? "✔️  บันทึกแล้วในรายการดูย้อนหลัง" : "❌  ยกเลิกการบันทึก");
    setTimeout(() => setToast(""), 3000);
    // TODO: call API ที่นี่ถ้ามี
  };

  const sharePost = async (e) => {
    e.stopPropagation();
    const url = window.location.origin + "/post/" + encodeURIComponent(post.title);
    try {
      if (navigator.share) {
        await navigator.share({ title: post.title, text: "ดูสรุปนี้บน ChaladShare", url });
      } else {
        await navigator.clipboard.writeText(url);
        setToast("📋  คัดลอกลิงก์แล้ว");
        setTimeout(() => setToast(""), 3000);
      }
    } catch {}
  };

  return (
    <div className="card ranking-card">
      {/* เหรียญรางวัลมุมขวาบน */}
      {rank <= 3 && (
        <img src={badge} alt={`อันดับ ${rank}`} className="rank-badge" />
      )}

      <div className="card-header">
        <img src={post.authorImg} alt="author" className="author-img" />
        <span>{post.authorName}</span>
      </div>

      {/* รูปสรุป: ใส่ className="card-image" เพื่อไม่ให้ชนกับเหรียญ */}
      <img src={post.img} alt="summary" className="card-image" />

      <div className="card-body">
        {/* ✅ แถว Reaction: ไลก์ซ้าย / บันทึก+แชร์ ขวา */}
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

      {/* Toast */}
      {toast && <div className="mini-toast">{toast}</div>}
    </div>
  );
};

export default RankingCard;
