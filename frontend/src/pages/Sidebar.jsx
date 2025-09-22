import React from "react";
import { useNavigate } from "react-router-dom";
import { FaRegUser } from "react-icons/fa6";
import { BsPersonAdd } from "react-icons/bs";
import { PiSparkle } from "react-icons/pi";
import { GoHome } from "react-icons/go";
import { IoMdAddCircleOutline } from "react-icons/io";
import { RiLogoutCircleRLine } from "react-icons/ri";
import '../component/Sidebar.css';

const Sidebar = () => {
  const navigate = useNavigate();

  return (
    <div className="sidebar">
      <div className="logo" onClick={() => navigate("/home")} style={{ cursor: "pointer" }}>
        <img src="/img/chalad share.png" alt="Chalad Share logo" />
      </div>
      <ul className="menu">
        <li onClick={() => navigate("/home")} style={{ cursor: "pointer" }}>
          <GoHome /> หน้าหลัก
        </li>
        <li><IoMdAddCircleOutline /> สร้าง</li>
        <li><PiSparkle /> AI ช่วยสรุป</li>
        <li><BsPersonAdd /> เพื่อน</li>
        <li><FaRegUser /> โปรไฟล์</li>
        <li onClick={() => navigate("/")} style={{ cursor: "pointer" }}>
          <RiLogoutCircleRLine /> ออกจากระบบ
        </li>
      </ul>
    </div>
  );
};

export default Sidebar;
