package repository

import (
	"context"
	"uas-go/app/model"
	"uas-go/database"

	"fmt"
)

type UserRepository struct{}

func (r UserRepository) FindByUsername(username string) (*model.User, error) {
	query := `
		SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at
		FROM users
		WHERE username = $1
		LIMIT 1
	`

	row := database.PgPool.QueryRow(context.Background(), query, username)

	var user model.User

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.RoleID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		fmt.Println("SCAN ERROR:", err)
		return nil, err
	}

	return &user, nil
}

func (r UserRepository) FindByID(id string) (*model.UserResponse, error) {
    query := `
        SELECT u.id, u.username, u.email, u.full_name, r.name
        FROM users u
		JOIN roles r ON u.role_id = r.id
        WHERE u.id = $1
        LIMIT 1
    `

    row := database.PgPool.QueryRow(context.Background(), query, id)

    var user model.	UserResponse

    err := row.Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.FullName,
		&user.Role,
    )

    if err != nil {
        return nil, err
    }

    return &user, nil
}

