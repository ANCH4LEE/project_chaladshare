package service

import (
	"chaladshare_backend/internal/auth/models"
	"chaladshare_backend/internal/auth/repository"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(req *models.RegisterRequest) (*models.AuthResponse, error)
	Login(req *models.LoginRequest) (*models.AuthResponse, error)
}

type authService struct {
	userRepo repository.AuthRepository
}

func NewAuthService(userRepo repository.AuthRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// func register
func (s *authService) Register(req *models.RegisterRequest) (*models.AuthResponse, error) {
	// ตรวจสอบว่าผู้ใช้มีอยู่แล้วหรือไม่
	existingUser, _ := s.userRepo.GetUserByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("อีเมลนี้ถูกใช้แล้ว")
	}

	// เข้ารหัสผ่าน
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถเข้ารหัสผ่านได้: %v", err)
	}

	//สร้าง struct user เก็บลง DB
	user := &models.User{
		Email:    req.Email,
		Username: req.Username,
		Password: string(hashedPassword),
	}

	//บันทึกลง BD ผ่าน repository
	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	// สร้าง response กลับให้ client
	resp := &models.AuthResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		Status:    "active",
	}
	return resp, nil
}

// func login
func (s *authService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	// ดึงข้อมูลผู้ใช้จาก email
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("เกิดข้อผิดพลาดในการเข้าสู่ระบบ")
	}
	if user == nil {
		return nil, errors.New("ไม่พบบัญชีผู้ใช้")
	}

	// ตรวจสอบรหัสผ่าน
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("รหัสผ่านไม่ถูกต้อง")
	}

	// password correct send response back to client
	resp := &models.AuthResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		Status:    "active",
	}
	return resp, nil
}
