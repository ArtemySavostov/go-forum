package repository

import (
	"article/internal/entity"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *entity.Comment) error
	GetCommentById(ctx context.Context, id string) (*entity.Comment, error)
	UpdateComment(ctx context.Context, comment *entity.Comment) error
	DeleteComment(ctx context.Context, id string) error
	GetAllComments(ctx context.Context) ([]*entity.Comment, error)
	GetCommentsByArticleID(ctx context.Context, id string) ([]*entity.Comment, error)
}

type MongoDbCommentRepository struct {
	collection *mongo.Collection
}

func NewMongoDBCommentRepository(collection *mongo.Collection) CommentRepository {
	return &MongoDbCommentRepository{collection: collection}
}

// CreateComment implements repository.CommentRepository.
//
//	func (m *MongoDbCommentRepository) CreateComment(ctx context.Context, comment *entity.Comment) error {
//		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
//		defer cancel()
//		_, err := m.collection.InsertOne(ctx, comment)
//		if err != nil {
//			log.Printf("failed to InsertOne on db: %v", err)
//			return fmt.Errorf("failed to create comment: %w", err)
//		}
//		return nil
//	}
func (m *MongoDbCommentRepository) CreateComment(ctx context.Context, comment *entity.Comment) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := m.collection.InsertOne(ctx, comment)
	if err != nil {
		log.Printf("failed to InsertOne on db: %v", err)
		return fmt.Errorf("failed to create comment: %w", err)
	}

	// Log the inserted ID
	log.Printf("Inserted comment with ID: %v", result.InsertedID)

	// Verify the insertion
	var insertedComment entity.Comment
	err = m.collection.FindOne(context.TODO(), bson.M{"_id": result.InsertedID}).Decode(&insertedComment)
	if err != nil {
		log.Printf("Failed to find inserted comment: %v", err)
		return fmt.Errorf("failed to verify comment creation: %w", err)
	}

	log.Printf("Successfully created and verified comment: %+v", insertedComment)
	return nil
}

// DeleteComment implements repository.CommentRepository.
func (m *MongoDbCommentRepository) DeleteComment(ctx context.Context, id string) error {
	filter := bson.M{"comment_id": id}
	_, err := m.collection.DeleteOne(ctx, filter)
	return err
}

// GetAllComments implements repository.CommentRepository.
func (m *MongoDbCommentRepository) GetAllComments(ctx context.Context) ([]*entity.Comment, error) {
	cursor, err := m.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var comments []*entity.Comment
	for cursor.Next(ctx) {
		var comment entity.Comment
		if err := cursor.Decode(&comment); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

// GetCommentById implements repository.CommentRepository.
func (m *MongoDbCommentRepository) GetCommentById(ctx context.Context, id string) (*entity.Comment, error) {
	var comment entity.Comment
	err := m.collection.FindOne(ctx, bson.M{"comment_id": id}).Decode(&comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetCommentsByArticleID implements repository.CommentRepository.
func (m *MongoDbCommentRepository) GetCommentsByArticleID(ctx context.Context, articleID string) ([]*entity.Comment, error) {
	filter := bson.M{"article_id": articleID}
	cursor, err := m.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find comments: %w", err)
	}
	defer cursor.Close(ctx)
	var comments []*entity.Comment
	if err := cursor.All(ctx, &comments); err != nil {
		return nil, fmt.Errorf("failed to decode comments: %w", err)
	}
	return comments, nil
}

// UpdateComment implements repository.CommentRepository.
func (m *MongoDbCommentRepository) UpdateComment(ctx context.Context, comment *entity.Comment) error {
	filter := bson.M{"comment_id": comment.CommentID}
	update := bson.M{"$set": comment}
	_, err := m.collection.UpdateOne(ctx, filter, update)
	return err
}
