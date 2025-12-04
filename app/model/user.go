package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID           string    `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`
	FullName     string    `json:"full_name" db:"full_name"`
	RoleID       string    `json:"role_id" db:"role_id"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type JWTClaims struct {
	UserID   string `json:"id"`
	Username string `json:"username"`
	RoleName string `json:"name"`

	jwt.RegisteredClaims
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token        string       `json:"token"`
	RefreshToken string       `json:"refreshToken"`
	User         UserResponse `json:"user"`
}

type UserResponse struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	FullName    string   `json:"full_name"`
	Email       string   `json:"email"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

type ProfileResponse struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	FullName    string   `json:"full_name"`
	Email       string   `json:"email"`
	Role        string   `json:"role"`
	IsActive    bool      `json:"is_active"`
}
