package model

type Lecturer struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	NIDN       string `json:"nidn"`
	Department string `json:"department"`
}

type LecturerResponse struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	FullName   string `json:"full_name"`
	NIDN       string `json:"nidn"`
	Department string `json:"department"`
}
