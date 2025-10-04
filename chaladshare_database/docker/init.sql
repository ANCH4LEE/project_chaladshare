-- ตาราง users สำหรับ login/register 172.20.10.2
create table if not exists users (
    user_id         serial primary key,                     -- id auto increment
    email           varchar(50) unique not null,            -- email ไม่ซ้ำ ห้ามว่าง
    username        varchar(50) unique not null,            -- username ไม่ซ้ำ ห้ามว่าง
    username_ci     varchar(50) generated always as 
                     (lower(username)) stored,              -- ทำ index คำเล็ก (case-insensitive)
    password_hash   varchar(255) not null,                  -- เก็บรหัสผ่านแบบ hash
    user_created_at timestamptz default now(),              -- เวลาสร้าง
    user_status     varchar(20) default 'active'            -- สถานะ เช่น active / inactive
);

-- สร้าง unique index สำหรับ username_ci กันซ้ำแบบ case-insensitive
create unique index if not exists users_username_ci_uq on users(username_ci);

-------------------------------------------------------------------------------
-- ตารางโปรไฟล์ผู้ใช้
create table if not exists user_profiles (
    profile_user_id integer primary key
        references users(user_id) on delete cascade,
    avatar_url      varchar(255),
    avatar_storage  varchar(255),
    bio             varchar(150),
    created_at      timestamptz default now(),
    updated_at      timestamptz default now()
);

-- ฟังก์ชันสำหรับอัปเดต updated_at
create or replace function set_updated_at()
returns trigger as $$
begin
    new.updated_at := now();
    return new;
end;
$$ language plpgsql;

-- ทริกเกอร์สำหรับ user_profiles
drop trigger if exists trg_user_profiles_updated_at on user_profiles;
create trigger trg_user_profiles_updated_at
before update on user_profiles
for each row
execute function set_updated_at();

-- ตารางเก็บ session การล็อกอิน (ใช้ refresh token)
create table if not exists auth_sessions (
    session_id         serial primary key,                         -- id auto increment
    session_user_id    integer references users(user_id) 
                       on delete cascade,                          -- ผูกกับ users ถ้าลบ user ก็ลบ session
    refresh_token_hash varchar(255) not null,                      -- เก็บค่า refresh token แบบ hash
    session_expires_at timestamptz not null,                       -- เวลาหมดอายุของ session
    revoked_at         timestamptz                                 -- เวลาเพิกถอน (เช่น logout หรือ revoke)
);

-- ตารางเก็บการ reset password (otp หรือโค้ดชั่วคราว)
create table if not exists password_resets (
    reset_pass_id         serial primary key,                      -- id auto increment
    reset_pass_user_id    integer references users(user_id) 
                          on delete cascade,                       -- ผูกกับ users
    otp_hash              varchar(255) not null,                   -- เก็บรหัส OTP แบบ hash
    reset_pass_expires_at timestamptz not null,                    -- เวลาหมดอายุของการ reset
    used_at               timestamptz                              -- เวลาใช้ reset ไปแล้ว
);

-------------------------------------------------------------------------------
-- เพิ่มตารางหัวข้อให้มาก่อน user_interests
create table if not exists topics (
    topic_id   serial primary key,
    topic_name varchar(20) not null unique
);

-- ตารางเก็บความสนใจของผู้ใช้ (mapping user ↔ topic)
create table if not exists user_interests (
    interest_user_id    integer references users(user_id) on delete cascade,   -- ผู้ใช้
    interest_topic_id   integer references topics(topic_id) on delete cascade, -- หัวข้อที่สนใจ
    interest_created_at timestamptz default now(),                             -- เวลาเพิ่มความสนใจ
    primary key (interest_user_id, interest_topic_id)                          -- กันซ้ำ user เลือก topic เดิม
);

-- ตารางเก็บการติดตาม (Follow)
create table if not exists follows (
    follower_user_id   integer references users(user_id) on delete cascade, -- คนที่กดติดตาม
    followee_user_id   integer references users(user_id) on delete cascade, -- คนที่ถูกติดตาม
    follow_created_at  timestamptz default now(),                           -- เวลา follow
    primary key (follower_user_id, followee_user_id)                        -- กันซ้ำ
);

-- enum สถานะคำขอเป็นเพื่อน (ใช้ DO-block กันกรณี IF NOT EXISTS ใช้ไม่ได้)
do $$
begin
  if not exists (select 1 from pg_type where typname = 'friend_request_status') then
    create type friend_request_status as enum ('pending','accepted','declined');
  end if;
end
$$ language plpgsql;

-- ตารางเก็บคำขอเป็นเพื่อน
create table if not exists friend_requests (
    request_id         serial primary key,                                      -- id auto increment
    requester_user_id  integer not null references users(user_id) on delete cascade, -- คนที่ส่งคำขอ
    addressee_user_id  integer not null references users(user_id) on delete cascade, -- คนที่ถูกส่งคำขอ
    request_status     friend_request_status not null default 'pending',        -- สถานะ
    request_created_at timestamptz default now(),                               -- เวลาเริ่มส่งคำขอ
    decided_at         timestamptz,                                             -- เวลาตอบรับ/ปฏิเสธ
    check (requester_user_id <> addressee_user_id)                              -- กันไม่ให้ส่งหาตัวเอง
);

-- index กัน pending ซ้ำทิศทาง (A → B, B → A)
create unique index if not exists uq_friend_requests_pending_pair
on friend_requests (
    least(requester_user_id, addressee_user_id),
    greatest(requester_user_id, addressee_user_id)
) where request_status = 'pending';

-- ตารางเพื่อน (friendships) ที่ถูก accept แล้ว
create table if not exists friendships (
    user_id    integer not null references users(user_id) on delete cascade, -- user
    friend_id  integer not null references users(user_id) on delete cascade, -- friend
    created_at timestamptz default now(),                                    -- เวลาเป็นเพื่อน
    primary key (user_id, friend_id),
    check (user_id <> friend_id),                                            -- กันไม่ให้เป็นเพื่อนกับตัวเอง
    check (user_id < friend_id)                                              -- เก็บทิศเดียว (user_id < friend_id)
);

-- ตารางเก็บไฟล์เอกสาร (documents)
create table if not exists documents (
    document_id       serial primary key,
    document_user_id  integer references users(user_id) on delete cascade, -- เจ้าของไฟล์
    document_name     varchar(255),                                        -- ชื่อไฟล์
    document_url      text not null,                                       -- URL ของไฟล์
    storage_provider  varchar(50),                                         -- เช่น s3, firebase, local
    uploaded_at       timestamptz default now()                            -- เวลาอัปโหลด
);

-- ตารางเก็บสรุป (summaries)
create table if not exists summaries (
    summary_id         serial primary key,
    summary_text       text not null,  -- ตัวสรุปข้อความ
    summary_html       text,           -- ถ้ามี highlight HTML
    summary_pdf_url    text,           -- ถ้ามี export PDF
    summary_created_at timestamptz default now() -- เวลา generate
);

-- ตารางโพสต์
create table if not exists posts (
    post_id             serial primary key,
    post_author_user_id integer references users(user_id) on delete cascade, -- ผู้โพสต์
    post_title          varchar(150) not null,                                -- หัวข้อโพสต์
    post_description    varchar(500),                                         -- คำอธิบาย
    post_visibility     varchar(10) check (post_visibility in ('public','friends')), -- การมองเห็น
    post_document_id    integer references documents(document_id),            -- อ้างไฟล์
    post_summary_id     integer references summaries(summary_id),             -- อ้างสรุป
    post_created_at     timestamptz default now(),
    post_updated_at     timestamptz default now()
);

-- ตารางแท็ก (tags)
create table if not exists tags (
    tag_id   serial primary key,
    tag_name varchar(50) unique not null -- ชื่อแท็กไม่ซ้ำ
);

-- ตารางเชื่อมโพสต์กับแท็ก (many-to-many)
create table if not exists post_tags (
    post_tag_post_id integer references posts(post_id) on delete cascade,
    post_tag_tag_id  integer references tags(tag_id) on delete cascade,
    primary key (post_tag_post_id, post_tag_tag_id)
);

-- ตารางเก็บการกดถูกใจ
create table if not exists likes (
    like_user_id    integer references users(user_id) on delete cascade, -- คนที่กดไลก์
    like_post_id    integer references posts(post_id) on delete cascade, -- โพสต์
    like_created_at timestamptz default now(),                           -- เวลาไลก์
    primary key (like_user_id, like_post_id)
);

-- ตารางเก็บการบันทึกโพสต์
create table if not exists saved_posts (
    save_user_id    integer references users(user_id) on delete cascade, -- คนที่บันทึก
    save_post_id    integer references posts(post_id) on delete cascade, -- โพสต์
    save_created_at timestamptz default now(),                           -- เวลาบันทึก
    primary key (save_user_id, save_post_id)
);

-- ตารางสถิติโพสต์
create table if not exists post_stats (
    post_stats_post_id    integer primary key references posts(post_id) on delete cascade, -- อิงโพสต์
    post_like_count       integer default 0,    -- จำนวนไลก์
    post_save_count       integer default 0,    -- จำนวนบันทึก
    post_last_activity_at timestamptz default now() -- เวลากิจกรรมล่าสุด
);

commit;