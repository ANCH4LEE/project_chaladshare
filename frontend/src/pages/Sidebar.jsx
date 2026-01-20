import React, { useEffect, useState } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import { RiUser6Line, RiUserAddLine, RiLogoutCircleRLine, RiHome2Line } from "react-icons/ri";
import { HiOutlineSparkles } from "react-icons/hi2";
import { IoMdAddCircleOutline } from "react-icons/io";

import "../component/Sidebar.css";
import logo from "../assets/logo.png";

const Sidebar = () => {
  const navigate = useNavigate();
  const location = useLocation();

  //จำเมนูที่เลือกไว้ (ค้างได้ทีละอัน)
  const [activeKey, setActiveKey] = useState("home");

  //ซิงค์ตาม path สำหรับหน้าที่มี path ชัดเจน (ไม่ชนกัน)
  useEffect(() => {
    const p = location.pathname;

    if (p === "/newpost") return setActiveKey("newpost");
    if (p === "/friends") return setActiveKey("friends");
    if (p === "/profile") return setActiveKey("profile");
    if (p === "/ai") return setActiveKey("ai");
    if (p === "/home") return setActiveKey("home");
  }, [location.pathname]);


  const go = (key, path) => {
    setActiveKey(key);
    navigate(path);
  };

  return (
    <div className="sidebar">
      <div
        className="logo"
        onClick={() => go("home", "/home")}
        style={{ cursor: "pointer" }}
      >
        <img src={logo} alt="Chalad Share logo" />
      </div>

      <ul className="menu">
        <li
          className={activeKey === "home" ? "active" : ""}
          onClick={() => go("home", "/home")}
          style={{ cursor: "pointer" }}
        >
          <RiHome2Line /> หน้าหลัก
        </li>

        <li
          className={activeKey === "newpost" ? "active" : ""}
          onClick={() => go("newpost", "/newpost")}
          style={{ cursor: "pointer" }}
        >
          <IoMdAddCircleOutline /> สร้าง
        </li>

        <li
          className={activeKey === "ai" ? "active" : ""}
          onClick={() => go("ai", "/ai")}
        >
          <HiOutlineSparkles /> AI ช่วยสรุป
        </li>


        <li
          className={activeKey === "friends" ? "active" : ""}
          onClick={() => go("friends", "/friends")}
          style={{ cursor: "pointer" }}
        >
          <RiUserAddLine /> เพื่อน
        </li>

        <li
          className={activeKey === "profile" ? "active" : ""}
          onClick={() => go("profile", "/profile")}
          style={{ cursor: "pointer" }}
        >
          <RiUser6Line /> โปรไฟล์
        </li>

        <li
          className={activeKey === "profile" ? "active" : ""}
          onClick={() => go("profile", "/profile")}
          style={{ cursor: "pointer" }}
        >
          <RiUser6Line /> ช่วยเหลือ
        </li>

        <li onClick={() => navigate("/")} style={{ cursor: "pointer" }}>
          <RiLogoutCircleRLine /> ออกจากระบบ
        </li>
      </ul>
    </div>
  );
};

export default Sidebar;
