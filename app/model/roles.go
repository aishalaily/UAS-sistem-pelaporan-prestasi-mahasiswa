package model

import (
	"time"
)

type Role struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type UpdateUserRoleRequest struct {
	Role string `json:"role"`
}
