import React, { useEffect, useMemo, useState } from "react";
import { useParams, Link, useNavigate } from "react-router-dom";
import { VscEye, VscEyeClosed } from "react-icons/vsc";
import { FiEdit2, FiTrash2 } from "react-icons/fi";
import axios from "axios";
import Sidebar from "./Sidebar";
import PostCard from "../component/Postcard";
import Avatar from "../assets/default.png";
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
    email: "",
    avatar: Avatar,
    posts: 0,
    followers: 0,
    following: 0,
  });

  // โหมดแก้ไขโปรไฟล์
  const [isEditing, setIsEditing] = useState(false);
  const [avatarPreview, setAvatarPreview] = useState(null);
  const [avatarFile, setAvatarFile] = useState(null);
  const [editForm, setEditForm] = useState({
    name: "",
    email: "",
    bio: "",
  });
  const [pwdForm, setPwdForm] = useState({
    current: "",
    newPwd: "",
    confirm: "",
  });
  const [saving, setSaving] = useState(false);
  const [saveErr, setSaveErr] = useState("");

  // state เปิด/ปิดการมองรหัสผ่าน
  const [showPwd, setShowPwd] = useState({
    current: false,
    newPwd: false,
    confirm: false,
  });
  const toggleShow = (key) => setShowPwd((s) => ({ ...s, [key]: !s[key] }));

  const openEdit = () => {
    setEditForm({
      name: profile.name || "",
      email: profile.email || "",
      bio: profile.bio || "",
    });
    setPwdForm({ current: "", newPwd: "", confirm: "" });
    setAvatarPreview(null);
    setAvatarFile(null);
    setSaveErr("");
    setIsEditing(true);
  };

  const onPickAvatar = (e) => {
    const f = e.target.files?.[0];
    if (!f) return;
    setAvatarFile(f);
    const url = URL.createObjectURL(f);
    setAvatarPreview(url);
  };

  const submitEdit = async () => {
    try {
      setSaving(true);
      setSaveErr("");

      // อัปโหลดรูป
      let avatarUrl = null;
      let avatarStorage = null;

      if (avatarFile) {
        const fd = new FormData();
        fd.append("file", avatarFile);
        const res = await axios.post("/files/avatar", fd, {
          headers: { "Content-Type": "multipart/form-data" },
        });

        avatarUrl = res?.data?.avatar_url || null;
        avatarStorage = res?.data?.avatar_storage || "local";
      }

      // edit profile
      await axios.put("/profile", {
        username: editForm.name,
        email: editForm.email,
        bio: editForm.bio,
        ...(avatarUrl && {
          avatar_url: avatarUrl,
          avatar_storage: avatarStorage,
        }),
      });

      // เปลี่ยนรหัสผ่าน (ถ้ากรอกครบ)
      if (pwdForm.current && pwdForm.newPwd && pwdForm.confirm) {
        await axios.post("/profile/password", {
          current_password: pwdForm.current,
          new_password: pwdForm.newPwd,
          confirm_password: pwdForm.confirm,
        });
      }

      // อัปเดตสถานะบนหน้า
      const fullAvatarUrl = avatarPreview
        ? avatarPreview
        : avatarUrl
        ? toAbsUrl(avatarUrl)
        : profile.avatar;

      setProfile((p) => ({
        ...p,
        name: editForm.name,
        email: editForm.email,
        bio: editForm.bio,
        avatar: avatarPreview || fullAvatarUrl || p.avatar,
      }));
      setIsEditing(false);
    } catch (e) {
      setSaveErr(e?.response?.data?.error || e.message || "บันทึกล้มเหลว");
    } finally {
      setSaving(false);
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

  // ---------- ฟังก์ชัน “แก้ไข / ลบโพสต์ของฉัน” ----------

  const handleEditPost = (event, post) => {
    event.stopPropagation();
    if (!post?.id) return;
    // ไปหน้าฟอร์มแก้ไขโพสต์ (สมมติ route แบบนี้; ถ้าโปรเจ็กต์ใช้ path อื่น ค่อยเปลี่ยนเองได้)
    navigate(`/posts/${post.id}/edit`);
  };

  const handleDeletePost = async (event, post) => {
    event.stopPropagation();
    if (!post?.id) return;

    const ok = window.confirm("คุณต้องการลบโพสต์นี้หรือไม่?");
    if (!ok) return;

    try {
      await axios.delete(`/posts/${post.id}`);
      setPosts((list) => list.filter((p) => p.id !== post.id));
      setProfile((p) => ({
        ...p,
        posts: Math.max(0, (p.posts ?? 0) - 1),
      }));
    } catch (e) {
      alert(e?.response?.data?.error || "ลบโพสต์ไม่สำเร็จ");
    }
  };

  // โหลด myId
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

        const statsUserId = isOwn
          ? prof?.data?.user_id ?? prof?.data?.id ?? myId
          : ownerId;
        const statsRes = await axios.get(`/social/stats/${statsUserId}`);
        const stats = statsRes?.data ?? {};

        const postsRes = isOwn
          ? await axios.get("/posts", { params: { mine: 1 } })
          : await axios.get("/posts", { params: { user_id: ownerId } });

        let savedRes = { data: [] };
        if (isOwn) {
          try {
            savedRes = await axios.get("/posts/save");
          } catch {}
        }

        const rawAvatar = prof?.data?.avatar_url || "";
        // const avatarForCards = rawAvatar ? toAbsUrl(rawAvatar) : Avatar;

        const format = (list) =>
          Array.isArray(list)
            ? list.map((p) => {
                const fileRaw = p.file_url || "";
                const coverRaw = p.cover_url || "";
                const isPdf = /\.pdf$/i.test(fileRaw);

                const imgSrc = coverRaw
                  ? toAbsUrl(coverRaw)
                  : !fileRaw || isPdf
                  ? "/img/pdf-placeholder.jpg"
                  : toAbsUrl(fileRaw);

                const rawAuthorAvatar = p.avatar_url || "";
                const authorImg = rawAuthorAvatar
                  ? toAbsUrl(rawAuthorAvatar)
                  : Avatar;

                const authorName =
                  p.author_name ||
                  p.username ||
                  (isOwn && profile.name) ||
                  "ผู้ใช้";

                return {
                  id: p.post_id,
                  post: p.post_id,
                  img: imgSrc,
                  isPdf,

                  likes: p.like_count ?? 0,
                  like_count: p.like_count ?? 0,
                  is_liked: !!p.is_liked,
                  is_saved: !!p.is_saved,

                  title: p.post_title,
                  tags: Array.isArray(p.tags)
                    ? p.tags
                        .map((t) => (t.startsWith("#") ? t : `#${t}`))
                        .join(" ")
                    : "",
                  authorId: p.author_id ?? p.post_user_id ?? p.user_id,
                  authorName,
                  authorImg,
                };
              })
            : [];

        setProfile((prev) => {
          const avatarFull = rawAvatar
            ? toAbsUrl(rawAvatar)
            : prev.avatar || Avatar;

          return {
            ...prev,
            name: prof?.data?.username ?? prev.name,
            email: prof?.data?.email ?? prev.email,
            bio: prof?.data?.bio ?? prev.bio,
            avatar: avatarFull,
            posts: prof?.data?.posts_count ?? prev.posts ?? 0,
            followers: stats.followers ?? prev.followers ?? 0,
            following: stats.following ?? prev.following ?? 0,
          };
        });

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

        setPosts(format(ownerRows));
        setSavedPosts(isOwn ? format(savedRows) : []);
        if (!isOwn) {
          const rel = prof?.data ?? {};

          // follow
          if (typeof rel.is_following === "boolean") {
            setFollowStatus(rel.is_following ? "following" : "idle");
          } else if (rel.is_following) {
            setFollowStatus("following");
          } else {
            setFollowStatus("idle");
          }

          // friend
          if (
            rel.is_friend === true ||
            rel.is_friend === 1 ||
            rel.is_friend === "1"
          ) {
            setFriendStatus("friends");
          } else if (rel.friend_request_outgoing) {
            setFriendStatus("requested");
          } else {
            setFriendStatus("idle");
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
  }, [isOwn, ownerId, targetId, myId, profile.name]);

  useEffect(() => {
    // ถ้าเป็นโปรไฟล์ตัวเอง ไม่ต้องเช็ค
    if (isOwn || !myId || !targetId) return;

    // ถ้ารู้อยู่แล้วว่าเป็นเพื่อน/ส่งคำขอแล้ว ก็ไม่ต้องยิงซ้ำ
    if (friendStatus !== "idle") return;

    const checkFriendRelation = async () => {
      try {
        // 1) เช็คว่าเป็นเพื่อนกันแล้วหรือยัง
        const friendsRes = await axios.get(`/social/friends/${myId}`, {
          params: { page: 1, size: 200, search: "" },
        });
        const friendItems = friendsRes.data.items || [];
        const isFriend = friendItems.some(
          (f) => String(f.user_id) === String(targetId)
        );

        if (isFriend) {
          setFriendStatus("friends");
          return;
        }

        // 2) ถ้ายังไม่ใช่เพื่อน → เช็ค outgoing requests
        const outgoingRes = await axios.get("/social/requests/outgoing", {
          params: { page: 1, size: 50 },
        });
        const outgoingItems = outgoingRes.data.items || [];
        const hasOutgoing = outgoingItems.some(
          (r) => String(r.addressee_user_id) === String(targetId)
        );

        if (hasOutgoing) {
          setFriendStatus("requested");
        } else {
          setFriendStatus("idle");
        }
      } catch (e) {
        console.error("checkFriendRelation failed:", e);
      }
    };

    checkFriendRelation();
  }, [isOwn, myId, targetId, friendStatus]);

  useEffect(() => {
    if (isOwn === false && activeTab !== "posts") setActiveTab("posts");
  }, [isOwn, activeTab]);

  const onToggleFollow = async () => {
    if (isOwn || !targetId) return;
    try {
      if (followStatus === "following") {
        await axios.delete(`/social/follow/${targetId}`);
        setFollowStatus("idle");
        setProfile((p) => ({
          ...p,
          followers: Math.max(0, (p.followers ?? 0) - 1),
        }));
      } else {
        await axios.post(`/social/follow`, {
          followed_user_id: Number(targetId),
        });
        setFollowStatus("following");
        setProfile((p) => ({ ...p, followers: (p.followers ?? 0) + 1 }));
      }
    } catch (e) {
      console.error(e);
    }
  };

  const onAddFriend = async () => {
    if (isOwn || !targetId || friendStatus !== "idle") return;
    try {
      await axios.post("/social/requests", { to_user_id: Number(targetId) });
      setFriendStatus("requested");
    } catch (e) {
      console.error(e);
    }
  };

  const showing = useMemo(
    () => (activeTab === "posts" ? posts : savedPosts),
    [activeTab, posts, savedPosts]
  );

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
                    <button className="profile-btn-edit" onClick={openEdit}>
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

            {/* ===== Edit Card (เฉพาะเจ้าของ/โหมดแก้ไข) ===== */}
            {isOwn && isEditing && (
              <section className="edit-card">
                <div className="edit-grid">
                  <div className="edit-avatar-col">
                    <img
                      src={avatarPreview || profile.avatar}
                      alt="avatar-preview"
                      className="edit-avatar"
                    />
                    <label className="edit-upload-btn">
                      เปลี่ยนรูปโปรไฟล์
                      <input
                        type="file"
                        accept="image/*"
                        onChange={onPickAvatar}
                        hidden
                      />
                    </label>
                  </div>

                  <div className="edit-form-col">
                    <div className="edit-field">
                      <label>ชื่อผู้ใช้</label>
                      <input
                        type="text"
                        value={editForm.name}
                        onChange={(e) =>
                          setEditForm((f) => ({ ...f, name: e.target.value }))
                        }
                      />
                    </div>

                    <div className="edit-field">
                      <label>อีเมล</label>
                      <input
                        type="email"
                        value={editForm.email}
                        onChange={(e) =>
                          setEditForm((f) => ({
                            ...f,
                            email: e.target.value,
                          }))
                        }
                      />
                    </div>

                    <div className="edit-field">
                      <label>คำอธิบาย</label>
                      <textarea
                        rows={3}
                        value={editForm.bio}
                        onChange={(e) =>
                          setEditForm((f) => ({ ...f, bio: e.target.value }))
                        }
                      />
                    </div>

                    {/* แถวเดียว: รหัสผ่านปัจจุบัน + ลิงก์ลืมรหัสผ่าน */}
                    <div className="edit-row">
                      <div className="edit-field flex-1">
                        <label>รหัสผ่านปัจจุบัน</label>
                        <div className="pw-field">
                          <input
                            className="pw-input"
                            type={showPwd.current ? "text" : "password"}
                            value={pwdForm.current}
                            onChange={(e) =>
                              setPwdForm((p) => ({
                                ...p,
                                current: e.target.value,
                              }))
                            }
                            placeholder="••••••••"
                          />
                          <button
                            type="button"
                            className="pw-toggle"
                            onClick={() => toggleShow("current")}
                            aria-label={
                              showPwd.current ? "ซ่อนรหัสผ่าน" : "แสดงรหัสผ่าน"
                            }
                            title={
                              showPwd.current ? "ซ่อนรหัสผ่าน" : "แสดงรหัสผ่าน"
                            }
                          >
                            {showPwd.current ? <VscEyeClosed /> : <VscEye />}
                          </button>
                        </div>
                      </div>

                      <button
                        type="button"
                        className="forgot-link"
                        onClick={() => navigate("/forgot_password")}
                        aria-label="ไปหน้าลืมรหัสผ่าน"
                        title="ลืมรหัสผ่าน?"
                      >
                        ลืมรหัสผ่าน?
                      </button>
                    </div>

                    {/* แถว: รหัสผ่านใหม่ / ยืนยันรหัสผ่านใหม่ */}
                    <div className="edit-2col">
                      <div className="edit-field">
                        <label>รหัสผ่านใหม่</label>
                        <div className="pw-field">
                          <input
                            className="pw-input"
                            type={showPwd.newPwd ? "text" : "password"}
                            value={pwdForm.newPwd}
                            onChange={(e) =>
                              setPwdForm((p) => ({
                                ...p,
                                newPwd: e.target.value,
                              }))
                            }
                          />
                          <button
                            type="button"
                            className="pw-toggle"
                            onClick={() => toggleShow("newPwd")}
                            aria-label={
                              showPwd.newPwd ? "ซ่อนรหัสผ่าน" : "แสดงรหัสผ่าน"
                            }
                          >
                            {showPwd.newPwd ? <VscEyeClosed /> : <VscEye />}
                          </button>
                        </div>
                      </div>

                      <div className="edit-field">
                        <label>ยืนยันรหัสผ่านใหม่</label>
                        <div className="pw-field">
                          <input
                            className="pw-input"
                            type={showPwd.confirm ? "text" : "password"}
                            value={pwdForm.confirm}
                            onChange={(e) =>
                              setPwdForm((p) => ({
                                ...p,
                                confirm: e.target.value,
                              }))
                            }
                          />
                          <button
                            type="button"
                            className="pw-toggle"
                            onClick={() => toggleShow("confirm")}
                            aria-label={
                              showPwd.confirm ? "ซ่อนรหัสผ่าน" : "แสดงรหัสผ่าน"
                            }
                          >
                            {showPwd.confirm ? <VscEyeClosed /> : <VscEye />}
                          </button>
                        </div>
                      </div>
                    </div>

                    {saveErr && <p className="edit-error">{saveErr}</p>}

                    <div className="edit-actions">
                      <button
                        className="btn-cancel"
                        onClick={() => setIsEditing(false)}
                        disabled={saving}
                      >
                        ยกเลิก
                      </button>
                      <button
                        className="btn-save"
                        onClick={submitEdit}
                        disabled={saving}
                      >
                        {saving ? "กำลังบันทึก…" : "บันทึกการแก้ไข"}
                      </button>
                    </div>
                  </div>
                </div>
              </section>
            )}
            {/* ===== /Edit Card ===== */}

            {/* Tabs */}
            <div className="profile-tabs">
              <div className="profile-tabs-track">
                <button
                  className={`profile-tab ${
                    activeTab === "posts" ? "active" : ""
                  }`}
                  onClick={() => setActiveTab("posts")}
                >
                  โพสต์
                </button>
                {isOwn && (
                  <button
                    className={`profile-tab ${
                      activeTab === "saved" ? "active" : ""
                    }`}
                    onClick={() => setActiveTab("saved")}
                  >
                    รายการที่บันทึกไว้
                  </button>
                )}
                <div
                  className={`profile-tab-underline ${
                    activeTab === "saved" ? "saved" : ""
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
                      className="profile-card-wrapper"
                      onClick={() => goToPostDetail(post)}
                      style={{ cursor: "pointer" }}
                    >
                      <PostCard post={post} />

                      {/* ปุ่มแก้ไข / ลบ เฉพาะโพสต์ของตัวเองในแท็บ "โพสต์" */}
                      {isOwn && activeTab === "posts" && (
                        <div className="profile-manage-row">
                          <button
                            type="button"
                            className="profile-manage-btn profile-manage-btn-edit"
                            onClick={(e) => handleEditPost(e, post)}
                          >
                            <FiEdit2 />
                            แก้ไข
                          </button>
                          <button
                            type="button"
                            className="profile-manage-btn profile-manage-btn-delete"
                            onClick={(e) => handleDeletePost(e, post)}
                          >
                            <FiTrash2 />
                            ลบ
                          </button>
                        </div>
                      )}
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
