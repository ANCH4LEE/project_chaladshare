-- CREATE TABLE IF NOT EXISTS usersRegis(
--     user_id SERIAL PRIMARY KEY,
--     user_email VARCHAR(255) UNIQUE NOT NULL,
--     user_name VARCHAR(255) NOT NULL,
--     user_password varchar(255) NOT NULL;
-- );

-- login/register
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,                -- รหัสผู้ใช้ (ไล่เลขอัตโนมัต
    email TEXT UNIQUE NOT NULL,           -- อีเมล (ต้องไม่ซ้ำ, ใช้ล็อกอิน)
    username TEXT NOT NULL,               -- ชื่อผู้ใช้ (แสดงผลในระบบ)
    password_hash TEXT NOT NULL,          -- รหัสผ่านแบบเข้ารหัส
    avatar_url TEXT,                      -- URL ของรูปโปรไฟล์ (optional)
    bio TEXT,                             -- ข้อความแนะนำตัว (optional)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- วันที่สมัคร
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP   -- แก้ไขโปรไฟล์ล่าสุด
);

-- file original


-- result summarize from gg colab




INSERT INTO users(email,username,password_hash) VALUES
('red@example.com', 'RED','12345'),
('blue@example.com', 'Bule','01234');

COMMIT;