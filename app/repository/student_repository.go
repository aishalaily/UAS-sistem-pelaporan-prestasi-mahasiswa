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

func GetAllStudents(db *pgxpool.Pool) ([]model.StudentWithAdvisor, error) {
	rows, err := db.Query(context.Background(), `
		SELECT
			s.id AS student_pk,
			u.id AS user_id,
			u.username,
			u.full_name,
			u.email,
			s.student_id,
			s.program_study,
			s.academic_year,
			l.id AS advisor_id,
			u2.full_name AS advisor_name
		FROM student s
		JOIN users u ON u.id = s.user_id
		LEFT JOIN lecturers l ON l.id = s.advisor_id
		LEFT JOIN users u2 ON u2.id = l.user_id
		ORDER BY s.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.StudentWithAdvisor
	for rows.Next() {
		var s model.StudentWithAdvisor
		if err := rows.Scan(
			&s.StudentPK,
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
		result = append(result, s)
	}

	return result, nil
}

func GetStudentByID(db *pgxpool.Pool, studentID string) (*model.StudentWithAdvisor, error) {
	row := db.QueryRow(context.Background(), `
		SELECT
			s.id AS student_pk,
			u.id AS user_id,
			u.username,
			u.full_name,
			u.email,
			s.student_id,
			s.program_study,
			s.academic_year,
			l.id AS advisor_id,
			u2.full_name AS advisor_name
		FROM student s
		JOIN users u ON u.id = s.user_id
		LEFT JOIN lecturers l ON l.id = s.advisor_id
		LEFT JOIN users u2 ON u2.id = l.user_id
		WHERE s.id = $1
		LIMIT 1
	`, studentID)

	var s model.StudentWithAdvisor
	if err := row.Scan(
		&s.StudentPK,
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

	return &s, nil
}

func UpdateStudentAdvisor(db *pgxpool.Pool, studentID, advisorID string) error {
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

func GetStudentsByAdvisor(db *pgxpool.Pool, lecturerID string) ([]model.StudentAdviseeResponse, error) {
	rows, err := db.Query(context.Background(), `
		SELECT
			s.id,
			s.user_id,
			u.full_name,
			s.student_id,
			s.program_study,
			s.academic_year
		FROM student s
		JOIN users u ON u.id = s.user_id
		WHERE s.advisor_id = $1
		ORDER BY u.full_name
	`, lecturerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.StudentAdviseeResponse
	for rows.Next() {
		var s model.StudentAdviseeResponse
		if err := rows.Scan(
			&s.StudentPK,
			&s.UserID,
			&s.FullName,
			&s.StudentID,
			&s.ProgramStudy,
			&s.AcademicYear,
		); err != nil {
			return nil, err
		}
		res = append(res, s)
	}

	return res, nil
}

func GetStudentAchievements(db *pgxpool.Pool, studentPK string) ([]model.AchievementReference, error) {
	rows, err := db.Query(context.Background(), `
		SELECT id, student_id, mongo_achievement_id, status, created_at
		FROM achievement_references
		WHERE student_id = $1
		  AND is_deleted = false
		ORDER BY created_at DESC
	`, studentPK)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.AchievementReference
	for rows.Next() {
		var a model.AchievementReference
		if err := rows.Scan(
			&a.ID,
			&a.StudentID,
			&a.MongoAchievementID,
			&a.Status,
			&a.CreatedAt,
		); err != nil {
			return nil, err
		}
		res = append(res, a)
	}

	return res, nil
}

