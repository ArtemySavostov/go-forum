package websocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"realtime-chat/internal/entity"
	"realtime-chat/pkg/auth"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	hub         *Hub
	authService auth.AuthServiceClient
	chatUC      ChatUsecase
}

func NewHandler(hub *Hub, authService auth.AuthServiceClient, chatUC ChatUsecase) *Handler {
	return &Handler{
		hub:         hub,
		authService: authService,
		chatUC:      chatUC,
	}
}

func (h *Handler) HandleWebSocket(c *gin.Context) {

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade connection to WebSocket:", err)
		return
	}

	clientID := c.Query("clientID")
	roomID := c.Query("roomID")
	senderName := c.Query("sender_name")

	client := &Client{
		conn:       conn,
		hub:        h.hub,
		send:       make(chan []byte, 256),
		clientID:   clientID,
		roomID:     roomID,
		senderName: senderName,
	}

	h.hub.Register <- client

	go client.writePump()
	go h.readPump(client)

	fmt.Println("Client connected")
}

func (h *Handler) readPump(client *Client) {
	defer func() {
		h.hub.Unregister <- client
		client.conn.Close()
		log.Println("Client disconnected (readPump)")
	}()

	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		log.Printf("Received raw message: %s", message)
		var msg entity.Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("error unmarshaling message: %v", err)
			continue
		}

		if msg.Type == "heartbeat" {
			log.Println("Received heartbeat message")
			continue
		}

		err = h.chatUC.CreateMessage(&msg)
		if err != nil {
			log.Printf("Error creating message: %v", err)
			continue
		}
		log.Printf("Received message: %+v", msg)
		h.broadcastMessage(client, &msg)
	}
}

type ResponseMessage struct {
	Message string `json:"message"`
}

func (h *Handler) broadcastMessage(client *Client, message *entity.Message) {

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return
	}

	for c := range h.hub.Clients {
		if c.roomID == client.roomID {
			select {
			case c.send <- jsonMessage:
			default:
				close(c.send)
				h.hub.Unregister <- c
			}
		}
	}
}

func (h *Handler) ListMessagesByChannel(c *gin.Context) {
	channelID := c.Param("room_id")

	if channelID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Channel ID is required"})
		return
	}

	messages, err := h.chatUC.ListByChannel(channelID)
	if err != nil {
		log.Printf("Error listing messages by channel: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
		return
	}
	if messages == nil {
		messages = []entity.Message{}
	}

	c.JSON(http.StatusOK, messages)
}
