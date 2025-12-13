package model

import (
	"time"
)

type Student struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	StudentID    string `json:"student_id"`
	ProgramStudy string `json:"program_study"`
	AcademicYear string `json:"academic_year"`
	AdvisorID    string `json:"advisor_id"`
	CreatedAt time.Time `json:"created_at"`
}

type StudentWithAdvisor struct {
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`

	StudentID    string `json:"student_id"`
	ProgramStudy string `json:"program_study"`
	AcademicYear string `json:"academic_year"`

	AdvisorID    string `json:"advisor_id"`
	AdvisorName  string `json:"advisor_name"`
}

type StudentAdviseeResponse struct {
	StudentID      string `json:"student_id"`
	UserID         string `json:"user_id"`
	FullName       string `json:"full_name"`
	StudentNumber  string `json:"student_number"`
	ProgramStudy   string `json:"program_study"`
	AcademicYear   string `json:"academic_year"`
}
