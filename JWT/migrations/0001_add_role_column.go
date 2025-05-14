package migrations

import (
	"database/sql"
	"log"
)

func Up_0001(db *sql.DB) error {

	alterIDQuery := `
  ALTER TABLE users
  MODIFY COLUMN id VARCHAR(36) NOT NULL;
 `

	_, err := db.Exec(alterIDQuery)
	if err != nil {
		log.Printf("Error altering id column: %v", err)
		return err
	}

	query := `
  ALTER TABLE users
  ADD COLUMN role ENUM('user', 'admin') NOT NULL DEFAULT 'user';
 `

	_, err = db.Exec(query)
	if err != nil {
		log.Printf("Error executing migration 0001 Up (adding role): %v", err)
		return err
	}

	log.Println("Migration 0001 Up executed successfully")
	return nil
}

func Down_0001(db *sql.DB) error {
	query := `
  ALTER TABLE users
  DROP COLUMN role;
 `

	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Error executing migration 0001 Down (dropping role): %v", err)
		return err
	}

	alterIDQuery := `
  ALTER TABLE users
  MODIFY COLUMN id VARCHAR(24) NOT NULL;
 `
	_, err = db.Exec(alterIDQuery)

	if err != nil {
		log.Printf("Error altering id column back: %v", err)
		return err
	}

	log.Println("Migration 0001 Down executed successfully")
	return nil
}
