package repository

import (
	"context"
	"uas-go/app/model"
	"uas-go/database"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateLecturer(l model.Lecturer) error {
	query := `
		INSERT INTO lecturers (id, user_id, nidn, department)
		VALUES ($1,$2,$3,$4)
	`

	_, err := database.PgPool.Exec(context.Background(), query,
		l.ID, l.UserID, l.NIDN, l.Department,
	)

	return err
}

func GetLecturerByUserID(userID string) (*model.Lecturer, error) {
	row := database.PgPool.QueryRow(context.Background(), `
		SELECT id, user_id, nidn, department
		FROM lecturers
		WHERE user_id = $1
		LIMIT 1
	`)

	var l model.Lecturer
	err := row.Scan(
		&l.ID,
		&l.UserID,
		&l.NIDN,
		&l.Department,
	)
	if err != nil {
		return nil, err
	}

	return &l, nil
}

func GetAllLecturers(db *pgxpool.Pool) ([]model.LecturerResponse, error) {
	rows, err := db.Query(context.Background(), `
		SELECT
			l.id,
			l.user_id,
			u.full_name,
			l.nidn,
			l.department
		FROM lecturers l
		JOIN users u ON u.id = l.user_id
		ORDER BY u.full_name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.LecturerResponse
	for rows.Next() {
		var l model.LecturerResponse
		if err := rows.Scan(
			&l.ID,
			&l.UserID,
			&l.FullName,
			&l.NIDN,
			&l.Department,
		); err != nil {
			return nil, err
		}
		res = append(res, l)
	}

	return res, nil
}
