import React, { useEffect, useMemo, useState } from "react";
import { useParams, Link, useNavigate } from "react-router-dom";
import axios from "axios";
import Sidebar from "./Sidebar";
import PostCard from "../component/Postcard";
import "../component/Profile.css";

const API_HOST = "http://localhost:8080";
 const toAbsUrl = (p) => {
   if (!p) return "";
   if (p.startsWith("http")) return p;
   const clean = p.replace(/^\./, "");
   return `${API_HOST}${clean.startsWith("/") ? clean : `/${clean}`}`;
 };

const Profile = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [myId, setMyId] = useState(null);
  const isOwn = useMemo(
    () => (myId == null ? null : !id || String(id) === String(myId)),
    [id, myId]
  );
  const targetId = id || null;
  const ownerId = useMemo(() => {
    if (isOwn == null) return null;
    return isOwn ? myId : targetId;
  }, [isOwn, myId, targetId]);

  const [activeTab, setActiveTab] = useState("posts");
  const [profile, setProfile] = useState({
    name: "",
    bio: "",
    avatar: "/img/author2.jpg",
    posts: 0,
    followers: 0,
    following: 0,
  });

  // เพิ่มฟังก์ชันไว้ใต้ useState ตรงนี้เลย
  const handleEditProfile = () => {
    const newBio = prompt("กรอกข้อความแนะนำตัวใหม่:", profile.bio);
    if (newBio !== null && newBio.trim() !== "") {
      setProfile({ ...profile, bio: newBio });
    }
  };

  const [posts, setPosts] = useState([]);
  const [savedPosts, setSavedPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [err, setErr] = useState("");
  const [followStatus, setFollowStatus] = useState("idle");
  const [friendStatus, setFriendStatus] = useState("idle");

  const goToPostDetail = (post) => {
    if (post?.id) navigate(`/posts/${post.id}`);
  };

  // โหลด myId ของผู้ใช้ที่ล็อกอินก่อน
  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const me = await axios.get("/profile");
        if (!cancelled) setMyId(me?.data?.user_id ?? me?.data?.id ?? null);
      } catch (e) {
        if (e?.response?.status === 401) navigate("/login", { replace: true });
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [navigate]);

  useEffect(() => {
    if (isOwn == null || !ownerId) return;
    setLoading(true);
    setErr("");
    const fetchData = async () => {
      try {
        const prof = isOwn
          ? await axios.get("/profile", {
            params: { with: "stats,followers,following,rel" },
          })
          : await axios.get(`/profile/${targetId}`, {
            params: { with: "stats,followers,following,rel" },
          });

        // 2) โพสต์ของเจ้าของหน้า
        const postsRes = isOwn
          ? await axios.get("/posts", { params: { mine: 1 } })
          : await axios.get("/posts", { params: { user_id: ownerId } });

        // 3) รายการที่บันทึกไว้ (เฉพาะของตัวเอง)
        let savedRes = { data: [] };
        if (isOwn) {
          try {
            savedRes = await axios.get("/posts/saved");
          } catch { }
        }

        const format = (list) =>
          Array.isArray(list)
            ? list.map((p) => ({
              id: p.post_id,
              img: (() => {
                const raw = p.file_url || "";
                const isPdf = /\.pdf$/i.test(raw);
                if (!raw || isPdf) return "/img/pdf-placeholder.jpg";
                return toAbsUrl(raw);
              })(),
              isPdf: /\.pdf$/i.test(p.file_url || ""),
              likes: p.like_count ?? 0,
              title: p.post_title,
              tags: Array.isArray(p.tags)
                ? p.tags
                  .map((t) => (t.startsWith("#") ? t : `#${t}`))
                  .join(" ")
                : "",
              authorId: p.author_id ?? p.post_user_id ?? p.user_id,
              authorName: prof?.data?.username || (isOwn ? "ฉัน" : "ผู้ใช้"),
              authorImg: prof?.data?.avatar_url || "/img/author2.jpg",
            }))
            : [];

        setProfile((prev) => ({
          ...prev,
          name: prof?.data?.username ?? prev.name,
          bio: prof?.data?.bio ?? prev.bio,
          avatar: prof?.data?.avatar_url || prev.avatar,
          posts: prof?.data?.posts_count ?? 0,
          followers: prof?.data?.followers_count ?? 0,
          following: prof?.data?.following_count ?? 0,
        }));

        // รองรับทั้งรูปแบบ {data: [...]} และ [...]
        const postRows = Array.isArray(postsRes?.data?.data)
          ? postsRes.data.data
          : Array.isArray(postsRes?.data)
            ? postsRes.data
            : [];
        const savedRows = Array.isArray(savedRes?.data?.data)
          ? savedRes.data.data
          : Array.isArray(savedRes?.data)
            ? savedRes.data
            : [];

        const ownerRows = postRows.filter(
          (p) =>
            String(p.author_id ?? p.post_user_id ?? p.user_id) ===
            String(ownerId)
        );

        console.log("profile ownerId=", ownerId, "posts:", postRows);
        setPosts(format(ownerRows));
        setSavedPosts(isOwn ? format(savedRows) : []);
        if (!isOwn) {
          if (typeof prof?.data?.is_following === "boolean") {
            setFollowStatus(prof.data.is_following ? "following" : "idle");
          }
          if (typeof prof?.data?.is_friend === "boolean") {
            setFriendStatus(prof.data.is_friend ? "friends" : "idle");
          } else if (prof?.data?.friend_request_outgoing) {
            setFriendStatus("requested");
          }
        } else {
          setFollowStatus("idle");
          setFriendStatus("idle");
        }
      } catch (e) {
        setErr(e?.response?.data?.error || e.message);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [isOwn, ownerId, targetId]);

  // ถ้าเป็นโปรไฟล์คนอื่น บังคับแท็บให้เป็น "posts" เสมอ
  useEffect(() => {
    if (isOwn === false && activeTab !== "posts") setActiveTab("posts");
  }, [isOwn, activeTab]);

  // ปุ่มติดตาม / เลิกติดตาม (มุมมองคนอื่น)
  const onToggleFollow = async () => {
    if (isOwn || !targetId) return;
    try {
      if (followStatus === "following") {
        await axios.delete(`/follows/${targetId}`);
        setFollowStatus("idle");
        setProfile((p) => ({ ...p, followers: Math.max(0, p.followers - 1) }));
      } else {
        await axios.post(`/follows/${targetId}`);
        setFollowStatus("following");
        setProfile((p) => ({ ...p, followers: p.followers + 1 }));
      }
    } catch (e) {
      console.error(e);
    }
  };

  // ปุ่มเพิ่มเพื่อน (ตัวอย่าง: ส่งคำขออย่างเดียว)
  const onAddFriend = async () => {
    if (isOwn || !targetId || friendStatus !== "idle") return;
    try {
      await axios.post(`/friend-requests/${targetId}`);
      setFriendStatus("requested");
    } catch (e) {
      console.error(e);
    }
  };

  const showing = useMemo(
    () => (activeTab === "posts" ? posts : savedPosts),
    [activeTab, posts, savedPosts]
  );

  // ระหว่างยังไม่รู้ว่าเป็นของตัวเองไหม แสดงโหลดไว้ก่อน
  if (isOwn == null) {
    return (
      <div className="profile-page">
        <div className="profile-container">
          <Sidebar />
          <main className="profile-content">
            <div className="profile-shell">
              <p className="profile-msg">กำลังโหลด...</p>
            </div>
          </main>
        </div>
      </div>
    );
  }

  return (
    <div className="profile-page">
      <div className="profile-container">
        <Sidebar />

        <main className="profile-content">
          <div className="profile-shell">
            {/* Header */}
            <section className="profile-header">
              <img
                src={profile.avatar}
                alt={profile.name}
                className="profile-avatar"
              />
              <div className="profile-info">
                <div className="profile-toprow">
                  {isOwn ? (
                    <h2 className="profile-name">{profile.name || "—"}</h2>
                  ) : (
                    <h2 className="profile-name">
                      <Link to={`/profile/${targetId}`}>
                        {profile.name || "—"}
                      </Link>
                    </h2>
                  )}
                  {isOwn ? (
                    <button
                      className="profile-btn-edit"
                      onClick={handleEditProfile}
                    >
                      แก้ไขโปรไฟล์
                    </button>
                  ) : (
                    <div className="profile-actions">
                      <button className="btn-follow" onClick={onToggleFollow}>
                        {followStatus === "following"
                          ? "กำลังติดตาม"
                          : "ติดตาม"}
                      </button>
                      <button
                        className="btn-friend"
                        onClick={onAddFriend}
                        disabled={friendStatus !== "idle"}
                      >
                        {friendStatus === "friends"
                          ? "เป็นเพื่อนกันแล้ว"
                          : friendStatus === "requested"
                            ? "ส่งคำขอแล้ว"
                            : "+ เพิ่มเพื่อน"}
                      </button>
                    </div>
                  )}
                </div>
                <p className="profile-bio">{profile.bio || ""}</p>
                <div className="profile-stats">
                  <span>{profile.posts} โพสต์</span>
                  <span>{profile.followers} ผู้ติดตาม</span>
                  <span>{profile.following} กำลังติดตาม</span>
                </div>
              </div>
            </section>

            {/* Tabs */}
            <div className="profile-tabs">
              <div className="profile-tabs-track">
                <button
                  className={`profile-tab ${activeTab === "posts" ? "active" : ""
                    }`}
                  onClick={() => setActiveTab("posts")}
                >
                  โพสต์
                </button>
                {isOwn && (
                  <button
                    className={`profile-tab ${activeTab === "saved" ? "active" : ""
                      }`}
                    onClick={() => setActiveTab("saved")}
                  >
                    รายการที่บันทึกไว้
                  </button>
                )}
                <div
                  className={`profile-tab-underline ${activeTab === "saved" ? "saved" : ""
                    }`}
                />
              </div>
            </div>

            {/* Posts */}
            <section className="profile-body">
              {loading ? (
                <p className="profile-msg">กำลังโหลด...</p>
              ) : err ? (
                <p className="profile-msg error">เกิดข้อผิดพลาด: {err}</p>
              ) : showing.length === 0 ? (
                <p className="profile-msg muted">ไม่มีโพสต์ที่จะแสดง</p>
              ) : (
                <div className="card-list">
                  {showing.map((post) => (
                    <div
                      key={post.id}
                      onClick={() => goToPostDetail(post)}
                      style={{ cursor: "pointer" }}
                    >
                      <PostCard post={post} />
                    </div>
                  ))}
                </div>
              )}
            </section>
          </div>
        </main>
      </div>
    </div>
  );
};

export default Profile;
