package app

import (
	"article/internal/entity"
	"article/internal/repository"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDbArticleRepository struct {
	collection *mongo.Collection
}

func NewMongoDBArticleRepository(collection *mongo.Collection) repository.ArticleRepository {
	return &MongoDbArticleRepository{collection: collection}
}
func (r *MongoDbArticleRepository) CreateArticle(ctx context.Context, article *entity.Article) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, article)
	if err != nil {
		log.Printf("failed to InsertOne on db: %v", err)
		return fmt.Errorf("failed to create article: %w", err)
	}
	return nil
}
func (r *MongoDbArticleRepository) GetArticleById(ctx context.Context, id string) (*entity.Article, error) {
	var article entity.Article
	err := r.collection.FindOne(ctx, bson.M{"article_id": id}).Decode(&article)
	if err != nil {
		return nil, err
	}
	return &article, nil
}
func (r *MongoDbArticleRepository) UpdateArticle(ctx context.Context, article *entity.Article) error {
	filter := bson.M{"article_id": article.ArticleID}
	update := bson.M{"$set": article}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err

}
func (r *MongoDbArticleRepository) DeleteArticle(ctx context.Context, id string) error {
	filter := bson.M{"article_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err

}
func (r *MongoDbArticleRepository) GetAllArticles(ctx context.Context) ([]*entity.Article, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var articles []*entity.Article
	for cursor.Next(ctx) {
		var article entity.Article
		if err := cursor.Decode(&article); err != nil {
			return nil, err
		}
		articles = append(articles, &article)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return articles, nil

}
