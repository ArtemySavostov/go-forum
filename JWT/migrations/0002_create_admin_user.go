package migrations

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Замените этими значениями
const (
	AdminUsername = "admin"
	AdminEmail    = "admin@y.ru"
	AdminPassword = "admin"
)

func Up_0002(db *sql.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing admin password: %v", err)
		return err
	}

	query := `
  INSERT INTO users (id, username, email, password, role)
  VALUES (?, ?, ?, ?, 'admin')
 `
	id := uuid.New().String()

	_, err = db.Exec(query, id, AdminUsername, AdminEmail, hashedPassword)
	if err != nil {
		log.Printf("Error creating admin user: %v", err)
		return err
	}

	log.Println("Migration 0002 Up executed successfully: Admin user created")
	return nil
}

func Down_0002(db *sql.DB) error {
	query := `
  DELETE FROM users WHERE email = ?
 `

	_, err := db.Exec(query, AdminEmail)
	if err != nil {
		log.Printf("Error deleting admin user: %v", err)
		return err
	}

	log.Println("Migration 0002 Down executed successfully: Admin user deleted")
	return nil
}
