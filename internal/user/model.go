package user

import "time"

type User struct {
	ID        int64     `json:"id" example:"1" description:"Unique identifier for the user"`
	Username  string    `json:"username" example:"johndoe" description:"Username for login"`
	Email     string    `json:"email" example:"john@example.com" description:"User's email address"`
	Password  string    `json:"-" description:"User's password"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T12:00:00Z"`
	IsActive  bool      `json:"is_active" example:"true" description:"Whether the user account is active"`
}

type RegisterRequest struct {
	Username string `json:"username" example:"johndoe"`
	Email    string `json:"email" example:"john@example.com" binding:"required"`
	Password string `json:"password" example:"securePassword123" binding:"required" `
}

type LoginRequest struct {
	Email    string `json:"email" example:"john@example.com" binding:"required"`
	Password string `json:"password" example:"securePassword123" binding:"required"`
}

type TokenResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." description:"JWT token for authentication"`
}
