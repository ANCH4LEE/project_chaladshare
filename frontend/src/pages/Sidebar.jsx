import React from "react";
import { useNavigate } from "react-router-dom";
import { RiUser6Line } from "react-icons/ri";
import { RiUserAddLine } from "react-icons/ri";
import { HiOutlineSparkles } from "react-icons/hi2";
import { RiHome2Line } from "react-icons/ri";
import { IoMdAddCircleOutline } from "react-icons/io";
import { RiLogoutCircleRLine } from "react-icons/ri";

import '../component/Sidebar.css';
import logo from "../assets/logo.png";

const Sidebar = () => {
  const navigate = useNavigate();

  return (
    <div className="sidebar">
      <div className="logo" onClick={() => navigate("/home")} style={{ cursor: "pointer" }}>
        <img src={logo} alt="Chalad Share logo" />
      </div>
      <ul className="menu">
        <li onClick={() => navigate("/home")} style={{ cursor: "pointer" }}>
          <RiHome2Line /> หน้าหลัก
        </li>
        <li onClick={() => navigate("/newpost")} style={{ cursor: "pointer" }}>
          <IoMdAddCircleOutline /> สร้าง</li>
        <li onClick={() => navigate("/home")} style={{ cursor: "pointer" }}>
          <HiOutlineSparkles /> AI ช่วยสรุป</li>
        <li onClick={() => navigate("/friends")} style={{ cursor: "pointer" }}>
          <RiUserAddLine /> เพื่อน</li>
        <li onClick={() => navigate("/home")} style={{ cursor: "pointer" }}>
          <RiUser6Line /> โปรไฟล์</li>
        <li onClick={() => navigate("/")} style={{ cursor: "pointer" }}>
          <RiLogoutCircleRLine /> ออกจากระบบ
        </li>
      </ul>
    </div>
  );
};

export default Sidebar;
