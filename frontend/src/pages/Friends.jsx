// หน้า .jsx (ทำ prefix แล้ว)

// src/pages/Friends.jsx
// src/pages/Friends.jsx
import React, { useMemo, useState } from "react";
import { IoSearch } from "react-icons/io5";
import Sidebar from "./Sidebar";
import "../component/Friends.css";

const Friends = () => {
  const [friends, setFriends] = useState([
    { id: 1, name: "Friend1", username: "xxxxxx", avatar: "img/av1.jpg" },
    { id: 2, name: "Friend2", username: "xxxxxx", avatar: "img/av2.jpg" },
    { id: 3, name: "Friend3", username: "xxxxxx", avatar: "img/av3.jpg" },
    { id: 4, name: "Friend4", username: "xxxxxx", avatar: "img/av4.jpg" },
    { id: 5, name: "Friend5", username: "xxxxxx", avatar: "img/av5.jpg" },
    { id: 6, name: "Friend6", username: "xxxxxx", avatar: "img/av6.jpg" },
    { id: 7, name: "Friend7", username: "xxxxxx", avatar: "img/av7.jpg" },
    { id: 8, name: "Friend8", username: "xxxxxx", avatar: "img/av8.jpg" },
    { id: 9, name: "Friend9", username: "xxxxxx", avatar: "img/av9.jpg" },
  ]);

  // แท็บสำหรับเปลี่ยนหน้าฟีเจอร์ (ไม่แสดงเม็ดยา "เพื่อนของฉัน")
  const [activeTab, setActiveTab] = useState("my"); // 'my'|'add'|'requests'
  const [query, setQuery] = useState("");

  const filteredFriends = useMemo(() => {
    const q = query.trim().toLowerCase();
    if (!q) return friends;
    return friends.filter(
      (f) =>
        f.name.toLowerCase().includes(q) ||
        (f.username || "").toLowerCase().includes(q)
    );
  }, [query, friends]);

  const removeFriend = (id) => {
    setFriends((prev) => prev.filter((f) => f.id !== id));
  };

  return (
    <div className="friends-page">
      <div className="friends-container">
        <Sidebar />

        <main className="friends-main">
          {/* ===== Top bar: หัวข้อ + ปุ่ม + ค้นหา ===== */}
          <div className="friends-topbar">
            <div className="friends-top-left">
              <h2 className="friends-title">เพื่อนของฉัน</h2>

              <div className="friends-actions">
                <button
                  type="button"
                  className={`friends-pill friends-pill--green ${
                    activeTab === "add" ? "is-active" : ""
                  }`}
                  onClick={() => setActiveTab("add")}
                >
                  เพิ่มเพื่อน
                </button>

                <button
                  type="button"
                  className={`friends-pill friends-pill--outline ${
                    activeTab === "requests" ? "is-active" : ""
                  }`}
                  onClick={() => setActiveTab("requests")}
                >
                  คำขอ (3)
                </button>
              </div>
            </div>

            <div className="friends-search">
              <input
                type="text"
                placeholder="ค้นหาเพื่อน"
                value={query}
                onChange={(e) => setQuery(e.target.value)}
              />
              <IoSearch className="friends-search-icon" />
            </div>
          </div>

          {/* ===== รายการเพื่อน (แท็บ my) ===== */}
          {activeTab === "my" && (
            <ul className="friends-list">
              {filteredFriends.map((f) => (
                <li key={f.id} className="friends-item">
                  <div className="friends-left">
                    <img
                      className="friends-avatar"
                      src={f.avatar}
                      alt={`${f.name} avatar`}
                      onError={(e) => (e.currentTarget.src = "img/author2.jpg")}
                    />
                    <div className="friends-name">
                      <span className="friends-name-main">{f.name}</span>
                      &nbsp;
                      <span className="friends-name-sub">{f.username}</span>
                    </div>
                  </div>

                  <button
                    className="friends-remove"
                    onClick={() => removeFriend(f.id)}
                  >
                    ลบเพื่อน
                  </button>
                </li>
              ))}
            </ul>
          )}

          {activeTab === "add" && (
            <div className="friends-placeholder">หน้าค้นหา/เพิ่มเพื่อน</div>
          )}
          {activeTab === "requests" && (
            <div className="friends-placeholder">หน้าคำขอเป็นเพื่อน</div>
          )}
        </main>
      </div>
    </div>
  );
};

export default Friends;