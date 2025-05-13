package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"JWT/internal/app"
	"JWT/internal/delivery/gin"
	"JWT/internal/delivery/gin/handlers"
	"JWT/internal/usecase"
	"JWT/pkg/auth"
	"JWT/pkg/database"

	"github.com/joho/godotenv"
	"github.com/rs/cors"

	_ "JWT/docs"
)

// @title User service
// @version 1.0
// @description AuthServer

// @host localhost:8088
// @BasePath /

// @securityDefinitions.basic BasicAuth

func main() {
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

	fmt.Println("Successfully connected to MySQL!")
	userRepo := app.NewMySQLUserRepository(db)

	// mongoURI := os.Getenv("MONGO_URI")
	// if mongoURI == "" {
	// 	log.Fatal("MONGO_URI environment variable not set")
	// }

	// client, err := database.ConnectMongoDB(context.Background(), mongoURI)
	// if err != nil {
	// 	log.Fatalf("Error connecting to MongoDB: %v", err)
	// }
	// defer func() {
	// 	if err := client.Disconnect(context.Background()); err != nil {
	// 		log.Fatalf("Error disconnecting from MongoDB: %v", err)
	// 	}
	// }()
	jwtSecretKey := os.Getenv("JWT_SECRET")
	//userRepo := app.NewMongoDBUserRepository(client.Database("forum_db").Collection("users"))
	authService := auth.NewJWTAuthService(jwtSecretKey)
	userUC := usecase.NewUserUseCase(userRepo)
	authUC := usecase.NewAuthUseCase(userRepo, authService)
	userHandler := handlers.NewUserHandler(userUC)
	authHandler := handlers.NewAuthHandler(authUC)
	allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
	allowedOrigins := strings.Split(allowedOriginsStr, ",")

	router := gin.SetupRouter(authHandler, userHandler, authService)

	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	handler := c.Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8088"
	}
	fmt.Printf("Server is running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
