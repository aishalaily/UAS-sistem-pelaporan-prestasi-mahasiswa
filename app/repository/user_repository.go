package repository

import (
	"context"
	"uas-go/app/model"
	"uas-go/database"
)

func GetUserByUsername(username string) (*model.User, string, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users
		WHERE username = $1 
		LIMIT 1
	`

	row := database.PgPool.QueryRow(context.Background(), query, username)

	var user model.User
	var passwordHash string

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&passwordHash,
		&user.FullName,
		&user.RoleID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, "", err
	}

	return &user, passwordHash, nil
}

func GetUserByID(id string) (*model.User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users
		WHERE id = $1 LIMIT 1
	`

	row := database.PgPool.QueryRow(context.Background(), query, id)

	var user model.User
	var passwordHash string

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&passwordHash,
		&user.FullName,
		&user.RoleID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func CreateUser(u model.User, passwordHash string) (string, error) {
	query := `
		INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,true, NOW(), NOW())
	`

	_, err := database.PgPool.Exec(context.Background(), query,
		u.ID, u.Username, u.Email, passwordHash, u.FullName, u.RoleID,
	)

	if err != nil {
		return "", err
	}

	return u.ID, nil
}

func IsUsernameExists(username string) bool {
	query := `
		SELECT 1
		FROM users
		WHERE username = $1
		LIMIT 1
	`

	var exists int
	err := database.PgPool.QueryRow(context.Background(), query, username).Scan(&exists)

	return err == nil
}


