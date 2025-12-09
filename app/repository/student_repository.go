package repository

import (
	"context"
	"uas-go/app/model"
	"uas-go/database"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateStudent(s model.Student) error {
	query := `
		INSERT INTO student (id, user_id, student_id, program_study, academic_year, advisor_id, created_at)
		VALUES ($1,$2,$3,$4,$5,$6, NOW())
	`
	_, err := database.PgPool.Exec(context.Background(), query,
		s.ID, s.UserID, s.StudentID, s.ProgramStudy, s.AcademicYear, s.AdvisorID,
	)
	return err
}

func GetStudentIDByUserID(db *pgxpool.Pool, userID string) (string, error) {
	var studentID string
	err := db.QueryRow(context.Background(),
		`SELECT id 
		FROM student 
		WHERE user_id = $1 LIMIT 1`, userID).
		Scan(&studentID)
	if err != nil {
		return "", err
	}
	return studentID, nil
}

func GetStudentsUnderAdvisor(db *pgxpool.Pool, advisorUserID string) ([]string, error) {
	rows, err := db.Query(context.Background(), `
		SELECT s.id
		FROM student s
		JOIN lecturers l ON s.advisor_id = l.id
		WHERE l.user_id = $1
	`, advisorUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []string
	for rows.Next() {
		var sid string
		if err := rows.Scan(&sid); err != nil {
			return nil, err
		}
		res = append(res, sid)
	}
	return res, nil
}

func IsStudentUnderAdvisor(db *pgxpool.Pool, advisorUserID, studentID string) (bool, error) {
	var exists bool
	err := db.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT 1 
			FROM student s
			JOIN lecturers l ON s.advisor_id = l.id
			WHERE l.user_id = $1 AND s.id = $2
		)
	`, advisorUserID, studentID).Scan(&exists)
	return exists, err
}
