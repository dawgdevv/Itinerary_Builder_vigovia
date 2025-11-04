package models

import "time"

// User represents a user account in the system
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Never include password in JSON responses
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SignupRequest represents the signup request payload
type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

// UserResponse represents a safe user response without sensitive data
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToUserResponse converts User to UserResponse
func (u *User) ToUserResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		FullName:  u.FullName,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
