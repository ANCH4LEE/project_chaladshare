package handlers

import (
	"context"
	"net/http"

	"chaladshare_backend/db"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	UserEmail string `json:"userEmail"`
	Password  string `json:"password"`
}

func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ถูกต้อง"})
		return
	}

	var passwordHash string
	err := db.Pool.QueryRow(context.Background(),
		"SELECT password_hash FROM users WHERE email=$1", req.UserEmail).Scan(&passwordHash)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ไม่พบบัญชีผู้ใช้"})
		return
	}

	if req.Password != passwordHash {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "รหัสผ่านไม่ถูกต้อง"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "เข้าสู่ระบบสำเร็จ"})
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterHandler(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ข้อมูลไม่ถูกต้อง"})
		return
	}

	// เช็กอีเมลซ้ำ
	var emailExists bool
	err := db.Pool.QueryRow(context.Background(),
		"SELECT EXISTS (SELECT 1 FROM users WHERE email=$1)", req.Email).Scan(&emailExists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "เกิดข้อผิดพลาดขณะตรวจสอบอีเมล"})
		return
	}
	if emailExists {
		c.JSON(http.StatusConflict, gin.H{"error": "อีเมลนี้ถูกใช้แล้ว"})
		return
	}

	// เช็กชื่อผู้ใช้ซ้ำ (ไม่สนตัวพิมพ์เล็ก/ใหญ่)
	var usernameExists bool
	err = db.Pool.QueryRow(context.Background(),
		"SELECT EXISTS (SELECT 1 FROM users WHERE LOWER(username) = LOWER($1))", req.Username).Scan(&usernameExists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "เกิดข้อผิดพลาดขณะตรวจสอบชื่อผู้ใช้"})
		return
	}
	if usernameExists {
		c.JSON(http.StatusConflict, gin.H{"error": "ชื่อผู้ใช้นี้ถูกใช้แล้ว"})
		return
	}

	// สมัครสมาชิก
	_, err = db.Pool.Exec(context.Background(),
		"INSERT INTO users (email, username, password_hash) VALUES ($1, $2, $3)",
		req.Email, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถสร้างผู้ใช้ได้"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "สมัครสมาชิกสำเร็จ"})
}
