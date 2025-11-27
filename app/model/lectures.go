package model

type Lecturer struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	LecturerID string `json:"lecturer_id"`
	Department string `json:"department"`
}
