package repository

import (
	"context"
	"fmt"
	"log"

	"JWT/internal/entity"

	"time"

	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBUserRepository struct {
	collection *mongo.Collection
}

func NewMongoDBUserRepository(collection *mongo.Collection) UserRepository {
	return &MongoDBUserRepository{collection: collection}
}

func (r *MongoDBUserRepository) GetUserById(ctx context.Context, id primitive.ObjectID) (entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.User{}, fmt.Errorf("user not found")
		}
		log.Printf("failed to FindOne on db: %v", err)
		return entity.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *MongoDBUserRepository) Create(ctx context.Context, user *entity.User) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, user)

	if err != nil {
		log.Printf("failed to InsertOne on db: %v", err)
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *MongoDBUserRepository) Update(ctx context.Context, user entity.User) error {
	// todo: Implement the function update user
	return nil
}

func (r *MongoDBUserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	// todo: Implement the function delete user
	return nil
}

func (r *MongoDBUserRepository) GetByUsername(ctx context.Context, username string) (entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.User{}, fmt.Errorf("user not found")
		}
		return entity.User{}, fmt.Errorf("failed to get user by username: %w", err)
	}

	return user, nil
}
func (r *MongoDBUserRepository) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.User{}, fmt.Errorf("user not found")
		}
		return entity.User{}, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

func (r *MongoDBUserRepository) GetAll(ctx context.Context) ([]entity.User, error) {
	return nil, nil
}
