package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var PostgresDB *pgxpool.Pool

func ConnectPostgres() {
	dsn := os.Getenv("POSTGRES_URL")

	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Failed to connect PostgreSQL: %v", err)
	}

	if err := db.Ping(context.Background()); err != nil {
		log.Fatalf("Postgres ping failed: %v", err)
	}

	log.Println("PostgreSQL connected")
	PostgresDB = db
}