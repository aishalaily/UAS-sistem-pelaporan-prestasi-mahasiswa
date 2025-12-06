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
