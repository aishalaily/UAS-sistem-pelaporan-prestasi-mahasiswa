package repository

import (
	"context"
	"uas-go/app/model"
	"uas-go/database"
)

func CreateStudent(s model.Student) error {
	query := `
		INSERT INTO students (id, user_id, student_id, program_study, academic_year, advisor_id, created_at)
		VALUES ($1,$2,$3,$4,$5,$6, NOW())
	`
	_, err := database.PgPool.Exec(context.Background(), query,
		s.ID, s.UserID, s.StudentID, s.ProgramStudy, s.AcademicYear, s.AdvisorID,
	)
	return err
}
