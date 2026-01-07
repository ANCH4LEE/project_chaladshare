import React, { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import axios from "axios";

import { MdOutlineAlternateEmail, MdLockOutline } from "react-icons/md";
import { VscEye, VscEyeClosed } from "react-icons/vsc";

import "../component/Login.css";
import bg from "../assets/bg.jpg";
import logo from "../assets/logo.png";

const Login = () => {
  const [formData, setFormData] = useState({email: "",password: "",});
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  // ฟังก์ชันตรวจสอบรูปแบบอีเมล
  const validateEmail = (email) => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = async(e) => {
    e.preventDefault();
    setError("");

  const email = formData.email.trim().toLowerCase();
  const password = formData.password;

    if (!email) {
      setError("กรุณากรอกอีเมล");
      return; 
    }

    if (!validateEmail(email)) {
      setError("รูปแบบอีเมลไม่ถูกต้อง");
      return;
    }

    if (!password) {
      setError("กรุณากรอกรหัสผ่าน");
      return;
    }
  
    try {
    setLoading(true);
    const res = await axios.post("/auth/login", { email, password });
     if (res.status === 200) navigate("/home");
    } catch (err) {
      const msg =
        err.response?.data?.error ||
        err.response?.data?.detail ||
        "เข้าสู่ระบบไม่สำเร็จ";
      setError(msg);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-page">
      <div
        className="login-container"
        style={{
          backgroundImage: `url(${bg})`,
          backgroundSize: "cover",
          backgroundPosition: "center",
          backgroundRepeat: "no-repeat",
        }}
      >
        <div className="login-box">
          <img src={logo} alt="Logo" />
          <h2>เข้าสู่ระบบ</h2>

          <form onSubmit={handleSubmit}>
            {/* email */}
            <div className="input-group">
              <span className="icon">
                <MdOutlineAlternateEmail />
              </span>
              <input
                type="email"
                name="email"
                value={formData.email}
                onChange={handleChange}
                placeholder="Email"
                required
                className="mb-3 p-2 border border-gray-300 rounded"
              />
            </div>

            {/* password */}
            <div className="input-group">
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
              >
                {showPassword ? <VscEyeClosed /> : <VscEye />}
              </span>
            </div>

            <div className="forgot-password">
              <Link to="/forgot_password">ลืมรหัสผ่าน?</Link>
            </div>

            <button
              type="submit"
              className="mb-3 p-2 border border-gray-300 rounded"
              disabled={loading}
            >
              {loading ? "กำลังเข้าสู่ระบบ..." : "เข้าสู่ระบบ"}
            </button>
            {/* error message */}
            {error && (
              <p
                style={{ color: "red", fontSize: "15px", marginBottom: "1rem" }}
              >
                {error}
              </p>
            )}
          </form>

          <div className="ClickToRegis">
            <p>มีบัญชีแล้วหรือยัง?</p>
            <Link to="/register">สมัครสมาชิก</Link>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Login;
