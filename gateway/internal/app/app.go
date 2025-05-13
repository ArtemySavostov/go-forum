package app

import (
	"fmt"
	"gateway/internal/config"
	"gateway/internal/delivery/gin"
	"gateway/internal/delivery/gin/handlers"
	"log"

	chat_grpc "github.com/ArtemySavostov/chat-protos"
)

func Run(port string, chatServiceClient chat_grpc.ChatServiceClient) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	gatewayHandler, err := handlers.NewGatewayHandler(cfg.UsersServiceURL, cfg.ArticlesServiceURL, chatServiceClient)
	if err != nil {
		return fmt.Errorf("failed to create gateway handler: %w", err)
	}
	router := gin.SetupRouter(gatewayHandler)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting gateway server on %s", addr)
	if err := router.Run(addr); err != nil {
		return fmt.Errorf("failed to run server: %w", err)
	}

	return nil

}
