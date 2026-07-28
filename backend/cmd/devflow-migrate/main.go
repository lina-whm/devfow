package main

import (
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
	defer db.Close()

	migrationsDir := filepath.Join("internal", "infrastructure", "postgres", "migrations")
	if err := goose.Run(command, db, migrationsDir); err != nil {
		log.Fatalf("migration %s: %v", command, err)
	}

	fmt.Printf("migration %s completed successfully\n", command)
}
