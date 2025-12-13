package repository

import (
	"context"
	"uas-go/app/model"
	"uas-go/database"
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
