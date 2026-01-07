package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"chaladshare_backend/internal/auth/models"
)

type AuthRepository interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id int) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(email, username, passwordHash string) (*models.User, error)
}

type authRepository struct {
	db *sql.DB
}

// func สร้าง repository
func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db: db}
}

// GET ผู้ใช้ทั้งหมด เรียง id
func (r *authRepository) GetAllUsers() ([]models.User, error) {
	rows, err := r.db.Query(`
		SELECT user_id, email, username, user_created_at, user_status
		FROM users
		ORDER BY user_id
	`)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถดึงข้อมูลผู้ใช้ทั้งหมดได้: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.ID, &u.Email, &u.Username,
			&u.CreatedAt, &u.Status,
		); err != nil {
			return nil, fmt.Errorf("อ่านข้อมูลผู้ใช้ไม่สำเร็จ: %w", err)
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("เกิดข้อผิดพลาดระหว่างอ่านข้อมูล: %w", err)
	}
	return users, nil
}

// ดึงข้อมูลผู้ใช้จาก id
func (r *authRepository) GetUserByID(id int) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(`
		SELECT user_id, email, username, user_created_at, user_status
		FROM users
		WHERE user_id = $1
	`, id).Scan(
		&u.ID, &u.Email, &u.Username, &u.CreatedAt, &u.Status,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("ไม่พบผู้ใช้")
	} else if err != nil {
		return nil, fmt.Errorf("เกิดข้อผิดพลาดในการดึงข้อมูลผู้ใช้: %w", err)
	}
	return &u, nil
}

// ผู้ใช้ตาม email
func (r *authRepository) GetUserByEmail(email string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(`
		SELECT user_id, email, username, password_hash, user_created_at, user_status
		FROM users
		WHERE LOWER(email) = LOWER($1)
	`, email).Scan(
		&u.ID, &u.Email, &u.Username, &u.PasswordHash,
		&u.CreatedAt, &u.Status,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("ไม่พบบัญชีผู้ใช้")
	} else if err != nil {
		return nil, fmt.Errorf("เกิดข้อผิดพลาด: %w", err)
	}
	return &u, nil
}

// สร้างผู้ใช้ใหม่
func (r *authRepository) CreateUser(email, username, passwordHash string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(`
		INSERT INTO users (email, username, password_hash)
		VALUES ($1, $2, $3)
		RETURNING user_id, email, username, user_created_at, user_status
	`, email, username, passwordHash).Scan(
		&u.ID, &u.Email, &u.Username,
		&u.CreatedAt, &u.Status,
	)

	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถสร้างผู้ใช้ใหม่ได้: %w", err)
	}

	return &u, nil
}
