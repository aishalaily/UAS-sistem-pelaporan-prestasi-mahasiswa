package repository

import (
    "context"
    "uas-go/app/model"
    "uas-go/database"
)

func FindUserByEmail(email string) (*model.User, error) {
    user := new(model.User)

    query := `
        SELECT id, email, password, full_name, role_id
        FROM users
        WHERE email=$1
        LIMIT 1
    `

    err := database.PostgresDB.
        QueryRow(context.Background(), query, email).
        Scan(&user.ID, &user.Email, &user.Password, &user.FullName, &user.RoleID)

    if err != nil {
        return nil, err
    }

    return user, nil
}
