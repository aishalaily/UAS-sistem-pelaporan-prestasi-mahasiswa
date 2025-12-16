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
	Permissions []string `json:"permissions"`
	
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

type RefreshTokenRequest struct {
	Token string `json:"token"`
}

type UserResponse struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	FullName    string   `json:"full_name"`
	Email       string   `json:"email"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

type CreateUserRequest struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	FullName     string `json:"full_name"`
	Password     string `json:"password"`
	Role         string `json:"role"`

	// mahasiswa
	StudentID     string `json:"student_id"`
	ProgramStudy  string `json:"program_study"`
	AcademicYear  string `json:"academic_year"`
	AdvisorID     string `json:"advisor_id"`

	// dosen wali
	NIDN       string `json:"nidn"`
	Department string `json:"department"`
}

type UpdateUserRequest struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	FullName *string `json:"full_name"`
	IsActive *bool   `json:"is_active"`
}