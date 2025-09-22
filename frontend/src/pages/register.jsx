import React, { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import "../component/login.css";
import { MdOutlineAlternateEmail, MdLockOutline } from "react-icons/md";
import { BsEye, BsEyeSlash } from "react-icons/bs";
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
    if (formData.password !== formData.confirmpassword) {
      setError("รหัสผ่านไม่ตรงกัน");
      return;
    }

    try {
      const response = await fetch("http://localhost:8080/api/v1/auth/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email: formData.userEmail,
          username: formData.username,
          password: formData.password,
        }),
      });

      const data = await response.json();
      if (!response.ok) {
        setError(data.error || "เกิดข้อผิดพลาดในการสมัครสมาชิก");
      } else {
        alert("สมัครสมาชิกสำเร็จ");
        navigate("/home");
      }
    } catch (err) {
      setError("ไม่สามารถเชื่อมต่อกับเซิร์ฟเวอร์ได้");
      console.error("Register error:", err);
    }
  };

  return (
    <div className="container" style={{
      backgroundImage: 'url("/img/bg.jpg")',
      backgroundSize: "cover",
      backgroundPosition: "center",
      backgroundRepeat: "no-repeat",
    }}>
      <div className="login-box">
        <img src="/img/chalad share.png" alt="Logo" />
        <h2>สมัครสมาชิก</h2>

        <form onSubmit={handleSubmit}>
          <div className="input-group">
            <span className="icon"><MdOutlineAlternateEmail /></span>
            <input type="email" name="userEmail" value={formData.userEmail} onChange={handleChange} placeholder="Email" required />
          </div>

          <div className="input-group">
            <span className="icon"><BiUser /></span>
            <input type="text" name="username" value={formData.username} onChange={handleChange} placeholder="Username" required />
          </div>

          <div className="input-group" style={{ position: "relative" }}>
            <span className="icon"><MdLockOutline /></span>
            <input type={showPassword ? "text" : "password"} name="password" value={formData.password} onChange={handleChange} placeholder="Password" required />
            <span className="icon-right" onClick={() => setShowPassword(!showPassword)} style={{ cursor: "pointer" }}>
              {showPassword ? <BsEyeSlash /> : <BsEye />}
            </span>
          </div>

          <div className="input-group" style={{ position: "relative" }}>
            <span className="icon"><MdLockOutline /></span>
            <input type={showConfirmPassword ? "text" : "password"} name="confirmpassword" value={formData.confirmpassword} onChange={handleChange} placeholder="Confirm password" required />
            <span className="icon-right" onClick={() => setshowConfirmPassword(!showConfirmPassword)} style={{ cursor: "pointer" }}>
              {showConfirmPassword ? <BsEyeSlash /> : <BsEye />}
            </span>
          </div>

          {error && <p style={{ color: "red", marginBottom: "10px" }}>{error}</p>}

          <button type="submit" className="mb-3 p-2 border border-gray-300 rounded">
            สมัครสมาชิก
          </button>
        </form>

        <div className="ClickToRegis">
          <p>คุณมีบัญชีแล้ว?</p>
          <Link to="/">เข้าสู่ระบบ</Link>
        </div>
      </div>
    </div>
  );
};

export default Register;
