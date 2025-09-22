import React, { useState } from "react";
import { useNavigate,Link } from "react-router-dom";
import "../component/login.css";
import { MdOutlineAlternateEmail, MdLockOutline } from "react-icons/md";
import { BsEye, BsEyeSlash } from "react-icons/bs";

const Login = () => {
  const [formData, setForm] = useState({
    userEmail: "",
    // username: "",
    password: "",
  });

  const navigate = useNavigate();
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState(""); // เก็บข้อความ error

  // ฟังก์ชันตรวจสอบรูปแบบอีเมล
  const validateEmail = (email) => {
    const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return re.test(email);
  };

  const handleChange = (e) => {
    setForm({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!formData.userEmail) {
      setError("กรุณากรอกอีเมล");
      return;
    }
    if (!validateEmail(formData.userEmail)) {
      setError("รูปแบบอีเมลไม่ถูกต้อง");
      return;
    }
    if (!formData.password) {
      setError("กรุณากรอกรหัสผ่าน");
      return;
    }

    try {
      const response = await fetch("http://localhost:8080/api/v1/auth/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          userEmail: formData.userEmail,
          password: formData.password,
        }),
      });

      const data = await response.json();
      if (!response.ok) {
        setError(data.error || "เกิดข้อผิดพลาด");
      } else {
        setError("");
        navigate("/home");
      }
    } catch (error) {
      setError("เชื่อมต่อไม่สำเร็จ");
      console.error("Login error:", error);
    }
  };

  return (
    <div
      className="container"
      style={{
        backgroundImage: 'url("/img/bg.jpg")',
        backgroundSize: "cover",
        backgroundPosition: "center",
        backgroundRepeat: "no-repeat",
      }}
    >
      <div className="login-box">
        <img src="/img/chalad share.png" alt="Logo" />
        <h2>เข้าสู่ระบบ</h2>

        <form onSubmit={handleSubmit}>
          {/* email */}
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
              className="mb-3 p-2 border border-gray-300 rounded"
            />
          </div>

          {/* password */}
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
              {showPassword ? <BsEyeSlash /> : <BsEye />}
            </span>
          </div>


          <div className="forgot-password">
            <a href="#">ลืมรหัสผ่าน?</a>
          </div>

          <button
            type="submit"
            className="mb-3 p-2 border border-gray-300 rounded"
          >
            เข้าสู่ระบบ
          </button>
            {/* error message */}
          {error && (
            <p style={{ color: "red", fontsize: "15px", marginBottom: "1rem" }}>
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
  );
};

export default Login;