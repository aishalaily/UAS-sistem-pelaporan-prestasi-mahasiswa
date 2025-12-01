package model

type Lecturer struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	NIDN       string `json:"nidn"`
	Department string `json:"department"`
}
