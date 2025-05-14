package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"JWT/pkg/database"

	"github.com/joho/godotenv"
)

func main() {
	var (
		migrateUp   = flag.Bool("up", false, "Run migrations up")
		migrateDown = flag.Bool("down", false, "Run migrations down (rollback)")
	)
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	mysqlPass := os.Getenv("MYSQL_PASSWORD")
	if mysqlPass == "" {
		mysqlPass = "admin"
	}

	cfg := database.MySQLConfig{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: mysqlPass,
		Database: "forum_db",
	}

	db, err := database.ConnectMySQL(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	if *migrateUp {
		fmt.Println("Running migrations up...")
		database.RunMigrations(db)
		fmt.Println("Migrations completed.")
	} else if *migrateDown {
		fmt.Println("Rolling back migrations...")
		database.RollbackMigrations(db)
		fmt.Println("Rollback completed.")
	} else {
		fmt.Println("Usage: migrate -up | -down")
		os.Exit(1)
	}
}
