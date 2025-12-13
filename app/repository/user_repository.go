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


func GetAllUsers() ([]model.User, error) {
	rows, err := database.PgPool.Query(context.Background(), `
		SELECT id, username, email, full_name, role_id, is_active, created_at
		FROM users
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.Email,
			&u.FullName,
			&u.RoleID,
			&u.IsActive,
			&u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func GetStudentByUserID(userID string) (*model.Student, error) {
	row := database.PgPool.QueryRow(context.Background(), `
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id
		FROM student
		WHERE user_id = $1
		LIMIT 1
	`)

	var s model.Student
	err := row.Scan(
		&s.ID,
		&s.UserID,
		&s.StudentID,
		&s.ProgramStudy,
		&s.AcademicYear,
		&s.AdvisorID,
	)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func UpdateUser(u model.User) error {
	query := `
		UPDATE users
		SET username = $2,
		    email = $3,
		    full_name = $4,
		    is_active = $5,
		    updated_at = NOW()
		WHERE id = $1
	`

	_, err := database.PgPool.Exec(
		context.Background(),
		query,
		u.ID,
		u.Username,
		u.Email,
		u.FullName,
		u.IsActive,
	)

	return err
}

func UpdateUserRole(userID, roleID string) error {
	query := `
		UPDATE users
		SET role_id = $2,
		    updated_at = NOW()
		WHERE id = $1
	`

	_, err := database.PgPool.Exec(context.Background(), query, userID, roleID)
	return err
}

func DeactivateUser(userID string) error {
	query := `
		UPDATE users
		SET is_active = false,
		    updated_at = NOW()
		WHERE id = $1
	`

	_, err := database.PgPool.Exec(context.Background(), query, userID)
	return err
}
