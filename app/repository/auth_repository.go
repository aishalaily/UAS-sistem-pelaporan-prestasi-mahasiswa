package repository

import (
	"context"
	"uas-go/app/model"
	"uas-go/database"
)

func FindUserByUsername(Username string) (*model.User, error) {
	user := new(model.User)

	query := `
        SELECT id, username, email, password_hash, role_id, created_at
        FROM users
        WHERE username=$1
        LIMIT 1
    `

	err := database.PostgresDB.
		QueryRow(context.Background(), query, Username).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.RoleID, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByID(userID string) (*model.User, error) {
	query := `
		SELECT id, username, email, full_name, role_id, is_active, created_at, updated_at 
		FROM users 
		WHERE id = $1 
		LIMIT 1
	`

	row := database.PostgresDB.QueryRow(context.Background(), query, userID)

	user := new(model.User)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FullName,
		&user.RoleID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}