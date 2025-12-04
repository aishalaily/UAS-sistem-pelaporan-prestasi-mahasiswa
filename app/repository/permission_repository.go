package repository

import (
	"context"
	"uas-go/database"
)

type PermissionRepository struct{}

func (r PermissionRepository) GetPermissionsByRole(roleID string) ([]string, error) {
	query := `
		SELECT p.name
		FROM permissions p
		JOIN role_permissions rp ON rp.permission_id = p.id
		WHERE rp.role_id = $1
	`

	rows, err := database.PgPool.Query(context.Background(), query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err == nil {
			permissions = append(permissions, name)
		}
	}

	return permissions, nil
}
