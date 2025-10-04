package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"chaladshare_backend/internal/auth/models"
)

type AuthRepository interface {
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
}

type authRepository struct {
	db *sql.DB
}

// func สร้าง repository
func NewUserRepository(db *sql.DB) AuthRepository {
	return &authRepository{db: db}
}

// ดึงข้อมูลผู้ใช้จาก email
func (r *authRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT user_id, email, username, password_hash, user_created_at, user_status
		FROM users
		WHERE email = $1
	`
	var user models.User
	err := r.db.QueryRow(query, email).
		Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("ไม่พบบัญชีผู้ใช้")
		}
		return nil, fmt.Errorf("query error: %v", err)
	}
	return &user, nil
}

// สร้างผู้ใช้ใหม่
func (r *authRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (email, username, password_hash)
		VALUES ($1, $2, $3)
		RETURNING user_id, user_created_at
	`
	err := r.db.QueryRow(query, user.Email, user.Username, user.Password).
		Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return errors.New("อีเมลหรือชื่อผู้ใช้งานถูกใช้แล้ว")
		}
		return fmt.Errorf("ไม่สามารถบันทึกผู้ใช้ได้: %v", err)
	}
	return nil
}
