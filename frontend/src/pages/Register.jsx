import React, { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import axios from "axios";

import "../component/Login.css";
import bg from "../assets/bg.jpg";
import logo from "../assets/logo.png";

import { MdOutlineAlternateEmail, MdLockOutline } from "react-icons/md";
import { VscEye, VscEyeClosed } from "react-icons/vsc";
import { BiUser } from "react-icons/bi";

const Register = () => {
  const [formData, setForm] = useState({
    userEmail: "",
    username: "",
    password: "",
    confirmpassword: "",
  });

  const navigate = useNavigate();
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setshowConfirmPassword] = useState(false);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const validateEmail = (email) => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);

  const handleChange = (e) => {
    setForm({
      ...formData,
      [e.target.name]: e.target.value,});
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");

    const email = formData.userEmail.trim().toLowerCase();
    const username = formData.username.trim();

    if (!email) {
      setError("กรุณากรอกอีเมล");
      return;
    }
    if (!validateEmail(email)) {
      setError("รูปแบบอีเมลไม่ถูกต้อง");
      return;
    }
    if (!username) {
      setError("กรุณากรอกชื่อผู้ใช้");
      return;
    }
    if (!formData.password) {
      setError("กรุณากรอกรหัสผ่าน");
      return;
    }
    if (formData.password !== formData.confirmpassword) {
      setError("รหัสผ่านไม่ตรงกัน");
      return;
    }

    try {
      setLoading(true);
      const res = await axios.post("/auth/register", {
        email,
        username,
        password: formData.password,
      });

      alert("สมัครสมาชิกสำเร็จ!");
      navigate("/home", { replace: true });
      
    } catch (err) {
      const msg =
        err.response?.data?.error ||
        err.response?.data?.detail ||
        "เกิดข้อผิดพลาดในการสมัครสมาชิก";
      setError(msg);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="register-page">
      <div
        className="register-container"
        style={{
          backgroundImage: `url(${bg})`,
          backgroundSize: "cover",
          backgroundPosition: "center",
          backgroundRepeat: "no-repeat",
        }}
      >
        <div className="login-box">
          <img src={logo} alt="Logo" />
          <h2>สมัครสมาชิก</h2>

          <form onSubmit={handleSubmit}>
            <div className="input-group">
              <span className="icon">
                <MdOutlineAlternateEmail />
              </span>
              <input
                type="email"
                name="userEmail"
                value={formData.userEmail}
                onChange={handleChange}
                placeholder="Email"
                required
              />
            </div>

            <div className="input-group">
              <span className="icon">
                <BiUser />
              </span>
              <input
                type="text"
                name="username"
                value={formData.username}
                onChange={handleChange}
                placeholder="Username"
                required
              />
            </div>

            <div className="input-group" style={{ position: "relative" }}>
              <span className="icon">
                <MdLockOutline />
              </span>
              <input
                type={showPassword ? "text" : "password"}
                name="password"
                value={formData.password}
                onChange={handleChange}
                placeholder="Password"
                required
              />
              <span
                className="icon-right"
                onClick={() => setShowPassword(!showPassword)}
                style={{ cursor: "pointer" }}
              >
                {showPassword ? <VscEyeClosed /> : <VscEye />}
              </span>
            </div>

            <div className="input-group" style={{ position: "relative" }}>
              <span className="icon">
                <MdLockOutline />
              </span>
              <input
                type={showConfirmPassword ? "text" : "password"}
                name="confirmpassword"
                value={formData.confirmpassword}
                onChange={handleChange}
                placeholder="Confirm password"
                required
              />
              <span
                className="icon-right"
                onClick={() => setshowConfirmPassword(!showConfirmPassword)}
                style={{ cursor: "pointer" }}
              >
                {showConfirmPassword ? <VscEyeClosed /> : <VscEye />}
              </span>
            </div>

            {error && (
              <p style={{ color: "red", marginBottom: "10px" }}>{error}</p>
            )}

            <button
              type="submit"
              className="mb-3 p-2 border border-gray-300 rounded"
            >
              สมัครสมาชิก
            </button>
          </form>

          <div className="ClickToRegis">
            <p>คุณมีบัญชีแล้ว?</p>
            <Link to="/">เข้าสู่ระบบ</Link>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Register;
