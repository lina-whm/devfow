package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer func() { _ = db.Close() }()

	migrationsDir := filepath.Join("internal", "infrastructure", "postgres", "migrations")
	ctx := context.Background()
	if err := goose.RunContext(ctx, command, db, migrationsDir); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	fmt.Printf("migration %s completed successfully\n", command)
}
