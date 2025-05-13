package app

import (
	"JWT/internal/entity"
	"JWT/internal/repository"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MySQLUserRepository struct {
	db *sql.DB
}

func NewMySQLUserRepository(db *sql.DB) repository.UserRepository {
	return &MySQLUserRepository{db: db}
}

func (r *MySQLUserRepository) GetUserById(ctx context.Context, id primitive.ObjectID) (entity.User, error) {
	var user entity.User
	var idStr string
	idHex := id.Hex()
	query := "SELECT id, username, email, password FROM users WHERE id = ?"
	err := r.db.QueryRowContext(ctx, query, idHex).Scan(&idStr, &user.Username, &user.Email, &user.Password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, fmt.Errorf("user not found")
		}
		log.Printf("failed to QueryRowContext: %v", err)
		return entity.User{}, fmt.Errorf("failed to get user: %w", err)
	}
	user.ID, err = primitive.ObjectIDFromHex(idStr)
	if err != nil {
		log.Printf("failed to convert hex to ObjectID: %v", err)
		return entity.User{}, fmt.Errorf("failed to convert hex to ObjectID: %w", err)
	}

	return user, nil
}
func (r *MySQLUserRepository) Create(ctx context.Context, user *entity.User) error {

	objectID := primitive.NewObjectID()

	user.ID = objectID
	idString := objectID.Hex()

	query := "INSERT INTO users (id, username, email, password) VALUES (?, ?, ?, ?)"

	_, err := r.db.ExecContext(ctx, query, idString, user.Username, user.Email, user.Password)
	if err != nil {
		log.Printf("failed to ExecContext: %v", err)
		return fmt.Errorf("failed to create user: %w", err)
	}

	//user.ID = objectID
	return nil
}

func (r *MySQLUserRepository) Update(ctx context.Context, user entity.User) error {
	query := "UPDATE users SET username = ?, email = ?, password = ? WHERE id = ?"
	result, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.Password, user.ID)
	if err != nil {
		log.Printf("failed to ExecContext: %v", err)
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("failed to get rows affected: %v", err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found or no changes applied")
	}

	return nil
}

func (r *MySQLUserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	query := "DELETE FROM users WHERE id = ?"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("failed to ExecContext: %v", err)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("failed to get rows affected: %v", err)
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *MySQLUserRepository) GetByUsername(ctx context.Context, username string) (entity.User, error) {
	var user entity.User
	var idString string
	query := "SELECT id, username, email, password FROM users WHERE username = ?"
	err := r.db.QueryRowContext(ctx, query, username).Scan(&idString, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, fmt.Errorf("user not found")
		}
		log.Printf("failed to QueryRowContext: %v", err)
		return entity.User{}, fmt.Errorf("failed to get user by username: %w", err)
	}
	objectID, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		log.Printf("failed to convert id from string to ObjectID: %v", err)
		return entity.User{}, fmt.Errorf("failed to convert id to ObjectID: %w", err)
	}

	user.ID = objectID
	return user, nil
}

func (r *MySQLUserRepository) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User
	query := "SELECT id, username, email, password FROM users WHERE email = ?"
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, fmt.Errorf("user not found")
		}
		log.Printf("failed to QueryRowContext: %v", err)
		return entity.User{}, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

func (r *MySQLUserRepository) GetAll(ctx context.Context) ([]entity.User, error) {
	return nil, fmt.Errorf("GetAll not implemented")
}
