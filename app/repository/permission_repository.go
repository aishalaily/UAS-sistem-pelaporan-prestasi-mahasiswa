package repository

import (
	"context"
	"uas-go/database"
)

func GetPermissionsByRole(roleID string) ([]string, error) {
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

	var perms []string
	for rows.Next() {
		var perm string
		rows.Scan(&perm)
		perms = append(perms, perm)
	}

	return perms, nil
}