package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var PgPool *pgxpool.Pool

func ConnectPostgres() error {
    dsn := os.Getenv("DB_DSN")

    db, err := pgxpool.New(context.Background(), dsn)
    if err != nil {
        return fmt.Errorf("failed to connect PostgreSQL: %w", err)
    }

    PgPool = db

    if err := PgPool.Ping(context.Background()); err != nil {
        return fmt.Errorf("PostgreSQL ping error: %w", err)
    }

    fmt.Println("PostgreSQL connected!")
    return nil
}
