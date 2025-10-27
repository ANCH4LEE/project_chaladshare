package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	Status       string    `json:"status"`
}

//register
type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//response ส่งกลับให้ client
type AuthResponse struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
	// Token 	  string 	`json:"token,omitempty"`
}
