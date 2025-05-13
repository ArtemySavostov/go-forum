package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDBClient(uri string) (*mongo.Database, error) {

	if uri == "" {
		uri = os.Getenv("MONGO_URI")
		if uri == "" {
			log.Println("MONGO_URI environment variable not set, using default localhost")
			uri = "mongodb://localhost:27017"
		}
	}

	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка при проверке подключения к MongoDB: %w", err)
	}

	fmt.Println("Успешно подключено к MongoDB!")

	return client.Database("chat_db"), nil
}
