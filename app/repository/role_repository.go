package repository

import (
	"context"
	"uas-go/app/model"
	"uas-go/database"
)

type RoleRepository struct{}

func (r RoleRepository) FindByID(roleID string) (*model.Role, error) {
	query := `
		SELECT id, name, description, created_at
		FROM roles
		WHERE id = $1
	`

	row := database.PgPool.QueryRow(context.Background(), query, roleID)

	var role model.Role
	err := row.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &role, nil
}
