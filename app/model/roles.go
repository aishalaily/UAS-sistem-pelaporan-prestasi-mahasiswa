package model

import (
	"time"
)

type Role struct {
	ID          string `json:"id"`
	RoleName        string `json:"name"`
	Description string `json:"description"`
	CreatedAt time.Time `json:"created_at"`
}
