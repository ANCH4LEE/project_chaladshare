import React, { useState } from "react";
import { AiFillHeart, AiOutlineHeart } from "react-icons/ai";

const PostCard = ({ post }) => {
  const [liked, setLiked] = useState(false);
  const [likes, setLikes] = useState(post.likes);

  const toggleLike = (e) => {
    e.stopPropagation(); // กันไม่ให้ไป trigger click ที่ parent
    if (liked) {
      setLikes(likes - 1);
    } else {
      setLikes(likes + 1);
    }
    setLiked(!liked);
  };

  return (
    <div className="card">
      <div className="card-header">
        <img src={post.authorImg} alt="author" className="author-img" />
        <span>{post.authorName}</span>
      </div>
      <img src={post.img} alt="summary" />
      <div className="card-body">
        <span className="likes" onClick={toggleLike} style={{ cursor: "pointer" }}>
          {liked ? (
            <AiFillHeart style={{ color: "red", fontSize: "20px" }} />
          ) : (
            <AiOutlineHeart style={{ color: "black", fontSize: "20px" }} />
          )}
          {likes}
        </span>
        <h4>{post.title}</h4>
        <p>{post.tags}</p>
      </div>
    </div>
  );
};

export default PostCard;
