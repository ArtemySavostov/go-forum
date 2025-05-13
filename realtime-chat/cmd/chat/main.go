package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"realtime-chat/internal/app"
	"realtime-chat/internal/delivery/websocket"
	"realtime-chat/internal/repository"

	"realtime-chat/pkg/auth"
	db_pkg "realtime-chat/pkg/database"

	"realtime-chat/internal/grpcserver"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	authServiceAddr := os.Getenv("AUTH_SERVICE_ADDR")
	if authServiceAddr == "" {
		authServiceAddr = "localhost:8081"
	}
	redisClient, err := db_pkg.NewRedisClient(redisAddr)
	if err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}
	defer redisClient.Close()

	messageRepo := repository.NewMessageRepository(redisClient)

	hub := websocket.NewHub()
	go hub.Run()

	authService, err := auth.NewAuthServiceClient(authServiceAddr)
	if err != nil {
		log.Fatalf("Failed to create auth service client: %v", err)
	}

	chatUC := websocket.NewChatUsecase(hub, messageRepo)

	handler := websocket.NewHandler(hub, authService, chatUC)

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // Разрешить запросы только с этого домена.  В production укажите ваш домен.
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"} // Разрешенные заголовки
	config.AllowCredentials = true                                                      // Разрешить передачу cookies
	router.Use(cors.New(config))

	app.NewApp(router, handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	addr := fmt.Sprintf(":%s", port)

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}
	grpcAddr := fmt.Sprintf(":%s", grpcPort)
	lis, err := net.Listen("tcp", grpcAddr)

	if err != nil {
		log.Fatalf("failed to listen for gRPC: %v", err)
	}

	go func() {
		if err := grpcserver.StartGRPCServer(lis); err != nil {
			log.Fatalf("failed to start gRPC server: %v", err)
		}
	}()

	//log.Printf("Запуск сервиса чата на %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
