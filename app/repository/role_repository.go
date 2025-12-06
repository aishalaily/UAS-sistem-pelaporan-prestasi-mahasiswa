package repository

import (
	"context"
	"uas-go/database"
)

func GetRoleName(roleID string) (string, error) {
	query := `
	SELECT name 
	FROM roles 
	WHERE id=$1 
	LIMIT 1`

	var roleName string
	err := database.PgPool.QueryRow(context.Background(), query, roleID).Scan(&roleName)
	if err != nil {
		return "", err
	}

	return roleName, nil
}

func GetRoleIDByName(name string) (string, error) {
	query := `
		SELECT id
		FROM roles
		WHERE LOWER(name) = LOWER($1)
		LIMIT 1
	`

	var roleID string
	err := database.PgPool.QueryRow(context.Background(), query, name).Scan(&roleID)
	if err != nil {
		return "", err
	}

	return roleID, nil
}
