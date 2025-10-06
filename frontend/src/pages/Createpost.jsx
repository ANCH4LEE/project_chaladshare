import React, { useState } from "react";
import axios from "axios";
import "../component/Createpost.css";
import Sidebar from "./Sidebar";

const CreatePost = () => {
  const [title, setTitle] = useState("");
  const [visibility, setVisibility] = useState("public");
  const [description, setDescription] = useState("");
  const [tags, setTags] = useState("");
  const [file, setFile] = useState(null);
  const [isLoading, setIsLoading] = useState(false);


  // choose file
  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
  };

  // กดโพสต์
  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!file) {
      alert("กรุณาอัปโหลดไฟล์");
      return;
    }

    try {
        setIsLoading(true);

        //อัปไฟล์ไป be /api/v1/files/upload
        const formData = new FormData();
        formData.append("file", file);

        const uploadRes = await axios.post(
            "http://localhost:8080/api/v1/files/upload",
            formData,
            {headers: {"Content-Type": "multipart/form-data"}}
        );

        const fileUrl = uploadRes.data.file_url; //be sent URL back

        //ส่งข้อมูลโพสต์ไป /api/v1/posts/
        const postData = {
            title: title,
            visibility: visibility, // public หรือ friends
            description: description,
            tags: tags,
            file_url: fileUrl,
            };

            await axios.post("http://localhost:8080/api/v1/posts", postData);

            alert("โพสต์สำเร็จ");
            setTitle("");
            setDescription("");
            setTags("");
            setFile(null);
        
        }catch (err) {
            console.error(err);
            alert("เกิดข้อผิดพลาดในการโพสต์");
        } finally {
            setIsLoading(false);
        }
  };

  return (
    <div className="create-post-container">
        <Sidebar />
      <h2 className="create-title">สร้างโพสต์ใหม่</h2>
      <form className="create-form" onSubmit={handleSubmit}>
        {/* หัวข้อ */}
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
            />
            <select
              value={visibility}
              onChange={(e) => setVisibility(e.target.value)}
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
              accept=".pdf"
            />
            <label htmlFor="file-upload" className="upload-label">
              {file ? (
                <span>{file.name}</span>
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
            placeholder="เพิ่มรายละเอียดเกี่ยวกับโพสต์ของคุณ..."
            value={description}
            onChange={(e) => setDescription(e.target.value)}
          ></textarea>
        </div>

        {/* แท็ก */}
        <div className="form-group">
          <label>
            แท็ก<span className="required">*</span>
          </label>
          <input
            type="text"
            placeholder="เพิ่มแท็ก..."
            value={tags}
            onChange={(e) => setTags(e.target.value)}
          />
        </div>

        {/* ปุ่ม */}
        <div className="button-group">
          <button type="button" className="btn-cancel">
            ยกเลิก
          </button>
          <button type="submit" className="btn-submit">
            โพสต์
          </button>
        </div>
      </form>
    </div>
  );
};

export default CreatePost;
