import React from "react";
import Sidebar from "./Sidebar";

const AISummary = () => {
  return (
    <div className="profile-page">
      <div className="profile-container">
        <Sidebar />
        <main className="profile-content">
          <div className="profile-shell">
            <h2>AI ช่วยสรุป</h2>
            <p style={{ color: "#6b7280" }}>
              หน้านี้อยู่ระหว่างพัฒนา
            </p>
          </div>
        </main>
      </div>
    </div>
  );
};

export default AISummary;
