package app

import (
	"fmt"
	"log"
	"realtime-chat/internal/delivery/websocket"

	"github.com/gin-gonic/gin"
)

func NewApp(r *gin.Engine, handler *websocket.Handler) {

	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.GET("/ws/chat", handler.HandleWebSocket)
	r.GET("/channels/:room_id/messages", handler.ListMessagesByChannel)
	log.Println("Route /ws/chat registered")
	fmt.Println("Starting")

}
