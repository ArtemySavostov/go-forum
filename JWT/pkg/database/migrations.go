package database

import (
	"database/sql"
	"log"

	"JWT/migrations"
)

func RunMigrations(db *sql.DB) {
	err := migrations.Up_0001(db)
	if err != nil {
		log.Fatalf("Error running migration 0001: %v", err)
	}

	err = migrations.Up_0002(db)
	if err != nil {
		log.Fatalf("Error running migration 0002: %v", err)
	}
}

func RollbackMigrations(db *sql.DB) {
	err := migrations.Down_0002(db)
	if err != nil {
		log.Printf("Error rolling back migration 0002: %v", err)
	}

	err = migrations.Down_0001(db)
	if err != nil {
		log.Printf("Error rolling back migration 0001: %v", err)
	}
}
