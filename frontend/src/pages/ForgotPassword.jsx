import React, { useState } from "react";
// import { useNavigate, Link } from "react-router-dom";
import { Link } from "react-router-dom";
import axios from "axios";

import { MdOutlineAlternateEmail } from "react-icons/md";
import { VscArrowLeft } from "react-icons/vsc";

import "../component/Login.css"; // ใช้ CSS เดียวกับหน้า Login ให้หน้าตาเหมือนกัน
import bg from "../assets/bg.jpg";
import logo from "../assets/logo.png";

const ForgotPassword = () => {
  const [email, setEmail] = useState("");
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [loading, setLoading] = useState(false);
//   const navigate = useNavigate();

  const validateEmail = (val) => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(val);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setSuccess("");

    const normalized = email.trim().toLowerCase();
    if (!normalized) return setError("กรุณากรอกอีเมล");
    if (!validateEmail(normalized)) return setError("รูปแบบอีเมลไม่ถูกต้อง");

    try {
      setLoading(true);
      const res = await axios.post("/auth/forgot-password", { email: normalized });
      if (res.status === 200) {
        setSuccess("ส่งรหัส OTP ไปยังอีเมลของคุณแล้ว โปรดตรวจสอบกล่องจดหมาย");
        // navigate("/verify_otp", { state: { email: normalized } });
      }
    } catch (err) {
      const msg = err.response?.data?.error || err.response?.data?.detail || "ไม่สามารถส่งรหัส OTP ได้";
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
          {/* ปุ่มย้อนกลับชิดซ้ายบนการ์ด */}
          <button
            className="back-btn"
            onClick={() => window.history.back()}
            aria-label="ย้อนกลับ"
          >
            <VscArrowLeft aria-hidden="true" />
          </button>

          <img src={logo} alt="Logo" />
          <h2>ลืมรหัสผ่าน</h2>

          <form onSubmit={handleSubmit} noValidate>
            <div className="input-group">
              <span className="icon">
                <MdOutlineAlternateEmail />
              </span>
              <input
                type="email"
                name="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="อีเมล"
                autoComplete="email"
                required
              />
            </div>

            <p style={{ color: "#6b7280", fontSize: 13, marginTop: -4 }}>
              ในขั้นตอนถัดไป รหัส OTP จะถูกส่งไปทางอีเมลของท่าน โปรดตรวจสอบอีเมล
            </p>

            <button type="submit" disabled={loading}>
              {loading ? "กำลังส่งรหัส OTP..." : "ส่ง"}
            </button>

            {error && (
              <p style={{ color: "red", fontSize: 15, marginTop: "0.75rem" }}>
                {error}
              </p>
            )}
            {success && (
              <p
                style={{
                  color: "#0f5132",
                  background: "#d1e7dd",
                  border: "1px solid #badbcc",
                  padding: "8px 10px",
                  borderRadius: 8,
                  fontSize: 14,
                  marginTop: "0.75rem",
                }}
              >
                {success}
              </p>
            )}
          </form>

          <div className="ClickToRegis" style={{ marginTop: 12 }}>
            <p>จำรหัสผ่านได้แล้ว?</p>
            <Link to="/">กลับไปเข้าสู่ระบบ</Link>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ForgotPassword;