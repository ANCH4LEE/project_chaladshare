package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"chaladshare_backend/internal/auth/models"
	"chaladshare_backend/internal/auth/service"
)

type AuthHandler struct {
	authService service.AuthService
	cookieName  string
	secure      bool
}

func NewAuthHandler(authService service.AuthService, cookieName string, secure bool) *AuthHandler {
	return &AuthHandler{authService: authService, cookieName: cookieName, secure: secure}
}

// Get all user
func (h *AuthHandler) GetAllUsers(c *gin.Context) {
	users, err := h.authService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// Get user by ID
func (h *AuthHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.authService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// Register
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest //
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON format"})
		return
	}

	user, err := h.authService.Register(req.Email, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ออก token ให้ user ใหม่
	token, err := h.authService.IssueToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "issue token failed"})
		return
	}

	// set cookie ทับของเดิม
	c.SetCookie(h.cookieName, token, 0, "/", "", h.secure, true)

	resp := models.AuthResponse{
		ID: user.ID, Email: user.Email, Username: user.Username,
		CreatedAt: user.CreatedAt, Status: user.Status,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    resp,
	})
}

// Login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON format"})
		return
	}

	user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := h.authService.IssueToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "issue token failed"})
		return
	}

	c.SetCookie(h.cookieName, token, 0, "/", "", h.secure, true)

	resp := models.AuthResponse{
		ID: user.ID, Email: user.Email, Username: user.Username, CreatedAt: user.CreatedAt, Status: user.Status,
	}
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": resp})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie(h.cookieName, "", -1, "/", "", h.secure, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
