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

func GetAllStudents() ([]model.StudentWithAdvisor, error) {
	query := `
		SELECT 
			u.id, u.username, u.full_name, u.email, s.student_id, s.program_study, s.academic_year, l.id AS advisor_id, u2.full_name AS advisor_name
		FROM users u
		JOIN student s ON s.user_id = u.id
		LEFT JOIN lecturers l ON s.advisor_id = l.id
		LEFT JOIN users u2 ON l.user_id = u2.id
		JOIN roles r ON r.id = u.role_id
		WHERE r.name = 'mahasiswa'
		ORDER BY u.created_at DESC
	`

	rows, err := database.PgPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.StudentWithAdvisor
	for rows.Next() {
		var s model.StudentWithAdvisor
		if err := rows.Scan(
			&s.UserID, 
			&s.Username, 
			&s.FullName, 
			&s.Email, 
			&s.StudentID, 
			&s.ProgramStudy, 
			&s.AcademicYear, 
			&s.AdvisorID, 
			&s.AdvisorName,
		); err != nil {
			return nil, err
		}
		res = append(res, s)
	}

	return res, nil
}

func GetStudentByID(studentID string) (*model.StudentWithAdvisor, error) {
	query := `
		SELECT 
			u.id, u.username, u.full_name, u.email, s.student_id, s.program_study, s.academic_year, l.id AS advisor_id, u2.full_name AS advisor_name
		FROM users u
		JOIN student s ON s.user_id = u.id
		LEFT JOIN lecturers l ON s.advisor_id = l.id
		LEFT JOIN users u2 ON l.user_id = u2.id
		WHERE s.id = $1
		LIMIT 1
	`

	var res model.StudentWithAdvisor
	err := database.PgPool.QueryRow(context.Background(), query).Scan(
		&res.UserID,
		&res.Username,
		&res.FullName,
		&res.Email,
		&res.StudentID,
		&res.ProgramStudy,
		&res.AcademicYear,
		&res.AdvisorID,
		&res.AdvisorName,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func UpdateStudentAdvisor(studentID, advisorID string) error {
	query := `
		UPDATE student
		SET advisor_id = $2
		WHERE id = $1
	`

	_, err := database.PgPool.Exec(
		context.Background(),
		query,
		studentID,
		advisorID,
	)

	return err
}

func GetStudentsByAdvisor(advisorID string) ([]model.StudentAdviseeResponse, error) {
	query := `
		SELECT 
			s.id, u.id, u.full_name, s.student_id, s.program_study, s.academic_year
		FROM student s
		JOIN users u ON s.user_id = u.id
		WHERE s.advisor_id = $1
		ORDER BY u.full_name ASC
	`

	rows, err := database.PgPool.Query(context.Background(), query, advisorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.StudentAdviseeResponse
	for rows.Next() {
		var s model.StudentAdviseeResponse
		if err := rows.Scan(
			&s.StudentID,
			&s.UserID,
			&s.FullName,
			&s.StudentNumber,
			&s.ProgramStudy,
			&s.AcademicYear,
		); err != nil {
			return nil, err
		}
		res = append(res, s)
	}

	return res, nil
}
