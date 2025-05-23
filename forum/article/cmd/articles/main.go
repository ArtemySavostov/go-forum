package main

import (
	"article/internal/delivery/gin"
	"article/internal/delivery/gin/handlers"
	"article/internal/repository"
	"article/internal/usecase"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "article/docs"

	"github.com/ArtemySavostov/JWT-Golang-mongodb/pkg/database"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// @title User service
// @version 1.0
// @description AuthServer

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable not set")
	}

	client, err := database.ConnectMongoDB(context.Background(), mongoURI)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatalf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	articleRepo := repository.NewMongoDBArticleRepository(client.Database("forum_db").Collection("article"))
	articleUC := usecase.NewArticleUseCase(articleRepo)
	articleHandler := handlers.NewArticleHandler(articleUC)
	commentRepo := repository.NewMongoDBCommentRepository(client.Database("forum_db").Collection("comment"))
	commentUC := usecase.NewCommentUseCase(commentRepo)
	//userService := handlers.NewUsersServiceClient("http://localhost:8088")
	commentHandler := handlers.NewCommentHandler(commentUC)
	allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
	allowedOrigins := strings.Split(allowedOriginsStr, ",")

	router := gin.NewRouter(articleHandler, commentHandler)
	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	handler := c.Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server is running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
