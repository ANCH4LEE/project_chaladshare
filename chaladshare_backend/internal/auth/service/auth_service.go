package service

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"chaladshare_backend/internal/auth/models"
	"chaladshare_backend/internal/auth/repository"
)

type AuthService interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id int) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	Register(email, username, password string) (*models.User, error)
	Login(email, password string) (*models.User, error)
}

type authService struct {
	userRepo repository.AuthRepository
}

func NewAuthService(userRepo repository.AuthRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// ผู้ใช้ทั้งหมด
func (s *authService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAllUsers()
}

// ผู้ใช้ตาม ID
func (s *authService) GetUserByID(id int) (*models.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}
	return s.userRepo.GetUserByID(id)
}

// ดึงผู้ใช้จากอีเมล
func (s *authService) GetUserByEmail(email string) (*models.User, error) {
	if strings.TrimSpace(email) == "" {
		return nil, errors.New("email is required")
	}
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// func register
func (s *authService) Register(email, username, password string) (*models.User, error) {
	if strings.TrimSpace(email) == "" || strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		return nil, errors.New("email, username and password are required")
	}
	if !strings.Contains(email, "@") {
		return nil, errors.New("invalid email format")
	}

	// เข้ารหัสรหัสผ่าน
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// สร้างผู้ใช้ใหม่
	user, err := s.userRepo.CreateUser(email, username, string(hashedPassword))
	if err != nil {
		return nil, fmt.Errorf("cannot create user: %v", err)
	}

	return user, nil
}

// func login
func (s *authService) Login(email, password string) (*models.User, error) {
	if strings.TrimSpace(email) == "" || strings.TrimSpace(password) == "" {
		return nil, errors.New("email and password are required")
	}

	// ดึงข้อมูลผู้ใช้จาก email
	user, err := s.userRepo.GetUserByEmail(email) // ❗ ถ้ายังไม่มีฟังก์ชันนี้ เธอใช้ GetAllUsers แล้ว loop เช็ค email ก็ได้
	if err != nil {
		return nil, errors.New("user not found")
	}

	// ตรวจสอบรหัสผ่าน
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}
