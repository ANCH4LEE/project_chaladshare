// หน้า Createpost.jsx (ทำ prefix แล้ว)

import React, { useState } from "react";
import { useNavigate } from "react-router-dom"; 
import axios from "axios";

import "../component/Createpost.css";
import Sidebar from "./Sidebar";


const UPLOAD_URL = "http://localhost:8080/api/v1/files/upload";
const POSTS_URL  = "http://localhost:8080/api/v1/posts";


const MAX_FILE_MB = 20; // จำกัดขนาดไฟล์ 20MB
const ACCEPTED_MIME = ["application/pdf"];

function parseTags(input) {
  // แปลง "uml, se #doc" -> ["#uml","#se","#doc"]
  return input
    .split(/[,\s]+/g)
    .map((t) => t.trim())
    .filter(Boolean)
    .map((t) => (t.startsWith("#") ? t : `#${t}`));
}

const CreatePost = () => {
  const [title, setTitle] = useState("");
  const [visibility, setVisibility] = useState("public"); // "public" | "friends"
  const [description, setDescription] = useState("");
  const [tags, setTags] = useState("");
  const [file, setFile] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");

  const navigate = useNavigate(); 

  // เลือกไฟล์ + ตรวจสอบชนิด+ขนาด
  const handleFileChange = (e) => {
    setErrorMsg("");
    const f = e.target.files?.[0];
    if (!f) {
      setFile(null);
      return;
    }

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
    setFile(f);
  };

 
  const handleCancel = () => {
    setTitle("");
    setDescription("");
    setTags("");
    setFile(null);
    setVisibility("public");
    setErrorMsg("");
    navigate("/home"); 
  };

  // โพสต์
  const handleSubmit = async (e) => {
    e.preventDefault();
    setErrorMsg("");

    if (!title.trim()) {
      setErrorMsg("กรุณากรอกหัวข้อ");
      return;
    }
    if (!file) {
      setErrorMsg("กรุณาอัปโหลดไฟล์ .pdf");
      return;
    }

    try {
      setIsLoading(true);

      // อัปโหลดไฟล์
      const formData = new FormData();
      formData.append("file", file); 

      console.log("UPLOAD_URL =", UPLOAD_URL);
      const uploadRes = await axios.post(UPLOAD_URL, formData);
      console.log("uploadRes =", uploadRes.data);

      const documentId = uploadRes.data?.document_id;
      if (!documentId) throw new Error("ไม่พบ document_id จากการอัปโหลดไฟล์");

      // สร้างโพสต์
      const postData = {
        author_user_id: 1, // TODO: เปลี่ยนมาใช้จาก JWT ภายหลัง
        post_title: title.trim(),
        post_description: description.trim(),
        post_visibility: visibility,
        post_document_id: documentId,
        tags: parseTags(tags),
      };

      console.log("POSTS_URL =", POSTS_URL);
      await axios.post(POSTS_URL, postData);

      alert("โพสต์สำเร็จ!");
      handleCancel();
    } catch (err) {
      console.error(err);
      setErrorMsg(
        err?.response?.data?.error ||
          err?.message ||
          "เกิดข้อผิดพลาดในการโพสต์"
      );
    } finally {
      setIsLoading(false);
    }
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
              placeholder="พิมพ์หัวข้อของคุณ..."
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              disabled={isLoading}
            />
            <select
              value={visibility}
              onChange={(e) => setVisibility(e.target.value)}
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
              {file ? <span>{file.name}</span> : (
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
            placeholder="เพิ่มรายละเอียดเกี่ยวกับโพสต์ของคุณ..."
            value={description}
            onChange={(e) => setDescription(e.target.value)}
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
            placeholder="เช่น #uml #se หรือ uml,se"
            value={tags}
            onChange={(e) => setTags(e.target.value)}
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
          >ยกเลิก</button>
          <button
            type="submit"
            className="btn-submit"
            disabled={isLoading || !title.trim() || !file}
          >
            {isLoading ? "กำลังโพสต์" : "โพสต์"}
          </button>
        </div>
      </form>
    </div>
    </div>
  );
};

export default CreatePost;
