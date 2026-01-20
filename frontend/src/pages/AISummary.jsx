import "../component/AISummary.css";
import React, { useRef, useState } from "react";
import Sidebar from "./Sidebar";

const UploadIcon = () => (
  <svg width="48" height="48" viewBox="0 0 64 64" aria-hidden="true">
    <rect x="10" y="14" width="28" height="22" rx="3" fill="none" stroke="#0b5394" strokeWidth="2" />
    <rect x="26" y="26" width="28" height="22" rx="3" fill="none" stroke="#0b5394" strokeWidth="2" />
    <circle cx="20" cy="22" r="3" fill="#0b5394" opacity="0.9" />
    <path d="M14 34l7-7 6 6 5-4 6 5" fill="none" stroke="#0b5394" strokeWidth="2" strokeLinejoin="round" />
    <path d="M30 44l7-7 6 6 5-4 6 5" fill="none" stroke="#0b5394" strokeWidth="2" strokeLinejoin="round" />
  </svg>
);

const SparkleIcon = () => (
  <svg width="28" height="28" viewBox="0 0 64 64" aria-hidden="true">
    <path d="M32 6l4.5 16.5L53 27l-16.5 4.5L32 48l-4.5-16.5L11 27l16.5-4.5L32 6z" fill="#6ec1ff" />
    <path d="M50 38l2.6 9.2L62 50l-9.4 2.8L50 62l-2.6-9.2L38 50l9.4-2.8L50 38z" fill="#ff7aa2" />
  </svg>
);

const AISummary = () => {
  const inputRef = useRef(null);
  const [file, setFile] = useState(null);

  const onPickFile = () => inputRef.current?.click();

  const onFileChange = (e) => {
    const f = e.target.files?.[0];
    if (!f) return;

    const isPdf =
      f.type === "application/pdf" || f.name.toLowerCase().endsWith(".pdf");
    if (!isPdf) {
      alert("กรุณาเลือกไฟล์ PDF เท่านั้น");
      if (inputRef.current) inputRef.current.value = "";
      return;
    }

    setFile(f);
  };

  const onClear = () => {
    setFile(null);
    if (inputRef.current) inputRef.current.value = "";
  };

  const onUploadKeyDown = (e) => {
    if (e.key === "Enter" || e.key === " ") {
      e.preventDefault();
      onPickFile();
    }
  };

  return (
    <div className="profile-page ai-page">
      <div className="profile-container">
        <Sidebar />

        <main className="profile-content">
          <div className="profile-shell">
            {/* ✅ Layout container */}
            <div className="ai-layout">
              {/* ✅ Source panel (LEFT) */}
              <aside className="ai-source-panel">
                <div className="ai-left-title">AI ช่วยสรุป</div>
                <div className="ai-left-sub">แหล่งข้อมูล</div>

                <input
                  ref={inputRef}
                  type="file"
                  accept=".pdf,application/pdf"
                  className="ai-hidden"
                  onChange={onFileChange}
                />

                {!file ? (
                  <div
                    className="ai-upload"
                    role="button"
                    tabIndex={0}
                    onClick={onPickFile}
                    onKeyDown={onUploadKeyDown}
                    aria-label="อัปโหลดไฟล์ PDF"
                    title="คลิกเพื่ออัปโหลดไฟล์ PDF"
                  >
                    <div className="ai-upload-icon">
                      <UploadIcon />
                    </div>
                    <div className="ai-upload-text">เพิ่มไฟล์ PDF</div>
                  </div>
                ) : (
                  <div className="ai-file">
                    <div className="ai-file-row">
                      <div className="ai-file-name" title={file.name}>
                        {file.name}
                      </div>
                    </div>

                    <div className="ai-file-actions">
                      <button className="ai-btn ai-btn-ghost" type="button" onClick={onPickFile}>
                        เปลี่ยนไฟล์
                      </button>
                      <button className="ai-btn ai-btn-danger" type="button" onClick={onClear}>
                        ลบไฟล์
                      </button>
                    </div>
                  </div>
                )}
              </aside>

              {/* ✅ Output panel (RIGHT) */}
              <section className="ai-output-panel">
                <div className="ai-greet">
                  <div className="ai-greet-icon">
                    <SparkleIcon />
                  </div>

                  <div className="ai-greet-text">
                    <div className="ai-greet-title">
                      สวัสดี, ฉันคือ AI ที่จะช่วยสรุปเนื้อหาของคุณ
                    </div>
                    <div className="ai-greet-sub">
                      อัปโหลดแหล่งข้อมูลแล้วเริ่มช่วยเหลือได้เลย
                    </div>
                  </div>
                </div>
              </section>
            </div>

          </div>
        </main>
      </div>
    </div>
  );
};

export default AISummary;
