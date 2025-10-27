// หน้า Createpost.jsx (ทำ prefix แล้ว)

import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";

import "../component/Createpost.css";
import Sidebar from "./Sidebar";

const MAX_FILE_MB = 20; // จำกัดขนาดไฟล์ 20MB
const ACCEPTED_MIME = ["application/pdf"];

function parseTags(input) {
  return input
    .split(/[,\s]+/g)
    .map((t) => t.trim().replace(/^#/, "").toLowerCase())
    .filter(Boolean);
}

const CreatePost = () => {
  const [formData, setForm] = useState({
    title: "",
    description: "",
    tags: "",
    visibility: "public",
    file: null,
  });

  const [isLoading, setIsLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");
  const navigate = useNavigate();

  const handleChange = (e) => {
    setForm({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  // เลือกไฟล์ + ตรวจสอบชนิด+ขนาด
  const handleFileChange = (e) => {
    setErrorMsg("");
    const f = e.target.files?.[0];
    if (!f) return;

    const sizeMB = f.size / (1024 * 1024);
    if (!ACCEPTED_MIME.includes(f.type)) {
      setErrorMsg("รองรับเฉพาะไฟล์ .pdf เท่านั้น");
      e.target.value = "";
      return;
    }
    if (sizeMB > MAX_FILE_MB) {
      setErrorMsg(`ไฟล์ใหญ่เกินไป (สูงสุด ${MAX_FILE_MB} MB)`);
      e.target.value = "";
      return;
    }
    setForm({ ...formData, file: f });
  };

  // โพสต์
  const handleSubmit = async (e) => {
    e.preventDefault();
    setErrorMsg("");

    if (!formData.title.trim()) {
      setErrorMsg("กรุณากรอกหัวข้อ");
      return;
    }
    if (!formData.file) {
      setErrorMsg("กรุณาอัปโหลดไฟล์ .pdf");
      return;
    }

    try {
      setIsLoading(true);

      // อัปโหลดไฟล์ PDF ก่อน
      const fileData = new FormData();
      fileData.append("file", formData.file);
      const uploadRes = await axios.post("/files/upload", fileData, {
        withCredentials: true,
        headers: { "Content-Type": "multipart/form-data" },
      });

      const documentId = uploadRes.data?.document_id;
      if (!documentId) throw new Error("ไม่พบ document_id จากการอัปโหลด");

      // สร้างโพสต์
      const postData = {
        post_title: formData.title.trim(),
        post_description: formData.description.trim(),
        post_visibility: formData.visibility,
        document_id: documentId,
        // post_summary_id: null,
        tags: parseTags(formData.tags),
      };

      await axios.post("/posts", postData,
        { withCredentials: true }
      );
      alert("โพสต์สำเร็จ!");
      handleCancel();
    } catch (err) {
      if (err?.response?.status === 401) {
        // คุกกี้หมดอายุ/ไม่ได้ล็อกอิน
        return navigate("/login", { replace: true });
      }
      console.error("Create post error:", err);
      setErrorMsg(
        err?.response?.data?.error || err?.message || "เกิดข้อผิดพลาดในการโพสต์"
      );
    } finally {
      setIsLoading(false);
    }
  };

  const handleCancel = () => {
    setForm({
      title: "",
      description: "",
      tags: "",
      visibility: "public",
      file: null,
    });
    setErrorMsg("");
    navigate("/home");
  };

  return (
    <div className="create-page">
      <Sidebar />
      <div className="create-post-container">
        <h2 className="create-title">สร้างโพสต์ใหม่</h2>

        <form className="create-form" onSubmit={handleSubmit}>
          {/* หัวข้อ + visibility */}
          <div className="form-group">
            <label>
              หัวข้อ<span className="required">*</span>
            </label>
            <div className="title-row">
              <input
                type="text"
                name="title"
                placeholder="พิมพ์หัวข้อของคุณ..."
                value={formData.title}
                onChange={handleChange}
                disabled={isLoading}
              />
              <select
                name="visibility"
                value={formData.visibility}
                onChange={handleChange}
                disabled={isLoading}
              >
                <option value="public">สาธารณะ</option>
                <option value="friends">เฉพาะเพื่อน</option>
              </select>
            </div>
          </div>

          {/* อัปโหลดไฟล์ */}
          <div className="form-group">
            <label>
              อัปโหลดไฟล์<span className="required">*</span>
            </label>
            <div className="upload-box">
              <input
                type="file"
                id="file-upload"
                onChange={handleFileChange}
                accept=".pdf,application/pdf"
                disabled={isLoading}
              />
              <label htmlFor="file-upload" className="upload-label">
                {formData.file ? (
                  <span>{formData.file.name}</span>
                ) : (
                  <>
                    <img
                      src="https://cdn-icons-png.flaticon.com/512/864/864685.png"
                      alt="upload"
                      className="upload-icon"
                    />
                    <p>เพิ่มไฟล์</p>
                  </>
                )}
              </label>
            </div>
          </div>

          {/* คำอธิบาย */}
          <div className="form-group">
            <label>คำอธิบาย</label>
            <textarea
              name="description"
              placeholder="เพิ่มรายละเอียดเกี่ยวกับโพสต์ของคุณ..."
              value={formData.description}
              onChange={handleChange}
              disabled={isLoading}
            />
          </div>

          {/* แท็ก */}
          <div className="form-group">
            <label>
              แท็ก<span className="required">*</span>
            </label>
            <input
              type="text"
              name="tags"
              placeholder="เช่น #uml #se หรือ uml,se"
              value={formData.tags}
              onChange={handleChange}
              disabled={isLoading}
            />
          </div>

          {/* ปุ่ม */}
          <div className="button-group">
            <button
              type="button"
              className="btn-cancel"
              onClick={handleCancel}
              disabled={isLoading}
            >
              ยกเลิก
            </button>
            <button
              type="submit"
              className="btn-submit"
              disabled={isLoading || !formData.title.trim() || !formData.file}
            >
              {isLoading ? "กำลังโพสต์..." : "โพสต์"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default CreatePost;
