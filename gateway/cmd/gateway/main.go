package main

import (
	"gateway/internal/app"
	"log"
	"os"

	_ "gateway/docs"

	chat_grpc "github.com/ArtemySavostov/chat-protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// @title User service
// @version 1.0
// @description AuthServer

// @host localhost:8000
// @BasePath /
func main() {

	chatServiceAddr := os.Getenv("CHAT_SERVICE_ADDR")
	if chatServiceAddr == "" {
		chatServiceAddr = "localhost:50051"
	}

	conn, err := grpc.Dial(chatServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться к gRPC-серверу чата: %v", err)
	}
	defer conn.Close()

	chatServiceClient := chat_grpc.NewChatServiceClient(conn)

	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8000"
	}

	if err := app.Run(port, chatServiceClient); err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}
}
