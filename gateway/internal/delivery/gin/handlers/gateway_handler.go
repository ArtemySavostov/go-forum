package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"gateway/internal/proxy"
	"io"
	"log"
	"os"

	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/ArtemySavostov/chat-protos"
	chat_grpc "github.com/ArtemySavostov/chat-protos"
	pb "github.com/ArtemySavostov/chat-protos"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc/status"
)

type GatewayHandler struct {
	usersServiceProxy    *httputil.ReverseProxy
	articlesServiceProxy *httputil.ReverseProxy
	chatServiceClient    chat_grpc.ChatServiceClient
}

func NewGatewayHandler(usersServiceURL string, articlesServiceURL string, chatServiceClient chat_grpc.ChatServiceClient) (*GatewayHandler, error) {
	usersURL, err := url.Parse(usersServiceURL)
	if err != nil {
		return nil, err
	}

	articlesURL, err := url.Parse(articlesServiceURL)
	if err != nil {
		return nil, err
	}
	return &GatewayHandler{
		usersServiceProxy:    proxy.NewReverseProxy(usersURL),
		articlesServiceProxy: proxy.NewReverseProxy(articlesURL),
		chatServiceClient:    chatServiceClient,
	}, nil
}

// ProxyToUsersService godoc
// @Summary Proxy request to Users Service
// @Description Proxy POST request /login and /register, GET/POST request /users/*path to Users Service
// @Tags Users
// @Accept json
// @Produce json
// @Param   path     path    string     false    "Path после /users"
// @Success 200 {object} map[string]interface{} "Succses respons"
// @Failure 403 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Gateway server error"
// @Router /login [post]
// @Router /register [post]
// @Router /users/{path} [get]
// @Router /users/{path} [post]
func (g *GatewayHandler) ProxyToUsersService(c *gin.Context) {
	g.usersServiceProxy.ServeHTTP(c.Writer, c.Request)
}

// ProxyToArticlesService godoc
// @Summary Proxy request to Articles Service
// @Description Proxy GET/POST request /articles, GET/PUT/DELETE /article/:id, GET/POST /articles/:articleId/comments to Articles Service
// @Tags Articles
// @Accept json
// @Produce json
// @Param id path string false "ID статьи"
// @Param articleId path string false "ArticleId"
// @Success 200 {object} map[string]interface{} "Succses respons"
// @Failure 500 {object} map[string]interface{} "Gateway server error"
// @Router /articles [get]
// @Router /articles [post]
// @Router /article/{id} [get]
// @Router /article/{id} [put]
// @Router /article/{id} [delete]
// @Router /articles/{articleId}/comments [get]
// @Router /articles/{articleId}/comments [post]
func (g *GatewayHandler) ProxyToArticlesService(c *gin.Context) {
	g.articlesServiceProxy.Director(c.Request)

	log.Printf("Proxying to Articles Service: %s", c.Request.URL.String())

	g.articlesServiceProxy.ServeHTTP(c.Writer, c.Request)
}

type responseWriter struct {
	gin.ResponseWriter
	statusCode int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// func (g *GatewayHandler) ProxyToChatService(c *gin.Context) {

// 	switch c.Request.Method {
// 	case http.MethodPost:

// 		g.createRoomHandler(c)
// 	case http.MethodGet:
// 		g.subscribeToRoomHandler(c)
// 	default:
// 		c.AbortWithStatus(http.StatusMethodNotAllowed)
// 	}

// }
// func (g *GatewayHandler) createRoomHandler(c *gin.Context) {
// 	log.Println("createRoomHandler called")
// 	req := &chat.CreateRoomRequest{
// 		RoomName: c.Query("roomName"),
// 		ClientId: c.Query("clientId"),
// 	}

// 	resp, err := g.chatServiceClient.CreateRoom(context.Background(), req)
// 	if err != nil {

// 		st, ok := status.FromError(err)
// 		if ok {
// 			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 				"error": st.Message(),
// 				"code":  st.Code(),
// 			})
// 		} else {
// 			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("ошибка gRPC: %w", err))
// 		}
// 		log.Printf("ClientId: %s", req.ClientId)
// 		log.Printf("gRPC error: %v", err)

// 		return
// 	}
// 	log.Printf("gRPC response: %v", resp)

// 	c.JSON(http.StatusOK, resp)
// }

func (g *GatewayHandler) subscribeToRoomHandler(c *gin.Context) {
	fmt.Println("Calling SubscribeToRoom handler")

	c.String(http.StatusOK, "SubscribeToRoom не реализован, используйте Websockets или SSE")
}

func getEnv(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // TODO: Проверка Origin в продакшн-среде!
	},
}

// func (g *GatewayHandler) ProxyToChatService(c *gin.Context) {
// 	// ProxyToChatService теперь просто проксирует запрос в нужный обработчик
// 	switch c.Request.Method {
// 	case http.MethodPost:
// 		g.CreateRoomHandler(c)
// 	case http.MethodGet:
// 		g.GetRoomHandler(c)
// 	default:
// 		c.AbortWithStatus(http.StatusMethodNotAllowed)
// 	}
// }

func (g *GatewayHandler) CreateRoomHandler(c *gin.Context) {
	log.Println("createRoomHandler called")

	req := &chat.CreateRoomRequest{
		RoomName: c.Query("roomName"),
		ClientId: c.Query("clientId"),
	}
	roomName := c.Query("roomName")
	clientID := c.Query("clientId")

	log.Printf("Получен roomName: %s, clientId: %s\n", roomName, clientID)

	resp, err := g.chatServiceClient.CreateRoom(context.Background(), req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": st.Message(),
				"code":  st.Code(),
			})
		} else {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("ошибка gRPC: %w", err))
		}
		log.Printf("ClientId: %s", req.ClientId)
		log.Printf("gRPC error: %v", err)
		return
	}
	log.Printf("gRPC response: %v", resp)
	c.JSON(http.StatusOK, resp)
}

// func (h *GatewayHandler) CreateRoomHandler(c *gin.Context) {
// 	log.Println("createRoomHandler called (Simplified)")
// 	c.String(http.StatusOK, "OK")
// }

func (g *GatewayHandler) GetRoomHandler(c *gin.Context) {
	// Получаем roomID из параметров URL
	roomID := c.Param("roomID")
	if roomID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "roomID is required"})
		return
	}

	// Вызываем gRPC метод для получения информации о комнате
	roomResponse, err := g.chatServiceClient.GetRoom(context.Background(), &chat.GetRoomRequest{RoomId: roomID})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": st.Message(),
				"code":  st.Code(),
			})
		} else {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("ошибка gRPC: %w", err))
		}
		log.Printf("gRPC error: %v", err)
		return
	}

	// Создаем структуру для ответа
	response := gin.H{
		"RoomId":   roomResponse.RoomId,
		"RoomName": roomResponse.RoomName,
	}

	c.JSON(http.StatusOK, response)
}

func (g *GatewayHandler) HandleWebSocket(c *gin.Context) {
	log.Println("HandleWebSocket CALLED!!!")
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	clientID := c.Query("clientID")
	roomID := c.Query("roomID")
	log.Printf("Room: %s", roomID)

	if clientID == "" || roomID == "" {
		log.Println("clientID and roomID are required")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := g.chatServiceClient.SubscribeToRoom(ctx)
	if err != nil {
		log.Printf("Error establishing ChatStream: %v", err)
		ws.Close() // Close websocket on gRPC stream error
		return
	}
	defer func() {
		if err := ws.Close(); err != nil {
			log.Printf("Error closing WebSocket: %v", err)
		}
	}()

	log.Printf("Sending gRPC request: RoomID=%s, ClientID=%s", roomID, clientID)

	// Отправляем первое сообщение с RoomID и ClientID
	initialMessage := &pb.ChatMessage{
		RoomId:   roomID,
		ClientId: clientID,
		Text:     "", // Initial message should not have text
	}
	if err := stream.Send(initialMessage); err != nil {
		log.Printf("Failed to send initial message: %v", err)
		return
	}

	go g.readFromGrpcStream(ws, stream)
	go g.writeToGrpcStream(ws, stream, roomID, clientID)

	// Keep the handler alive until the context is cancelled (client disconnects)
	<-ctx.Done()
	log.Println("Client disconnected, closing streams")
}

func (g *GatewayHandler) readFromGrpcStream(ws *websocket.Conn, stream pb.ChatService_SubscribeToRoomClient) {
	defer func() {
		if err := ws.Close(); err != nil {
			log.Printf("Error closing WebSocket in readFromGrpcStream: %v", err)
		}
	}()
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			log.Println("gRPC stream closed")
			return
		}
		if err != nil {
			log.Printf("Error receiving message from gRPC stream: %v", err)
			return
		}

		err = ws.WriteJSON(msg)
		if err != nil {
			log.Printf("Error writing message to WebSocket: %v", err)
			return
		}
	}
}

func (g *GatewayHandler) writeToGrpcStream(ws *websocket.Conn, stream pb.ChatService_SubscribeToRoomClient, roomID string, clientID string) {
	defer func() {
		log.Println("Closing writeToGrpcStream")
		if err := stream.CloseSend(); err != nil {
			log.Printf("Error closing gRPC stream: %v", err)
		}
		if err := ws.Close(); err != nil {
			log.Printf("Error closing WebSocket in writeToGrpcStream: %v", err)
		}
	}()

	for {
		messageType, messageBytes, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading from WebSocket: %v", err)
			// Send error message to client
			if closeErr, ok := err.(*websocket.CloseError); ok {
				log.Printf("WebSocket closed with code %d and reason %s", closeErr.Code, closeErr.Text)
			}
			break
		}

		log.Printf("MessageType: %d", messageType)
		messageString := string(messageBytes)
		log.Printf("Received message from WebSocket: %s", messageString)

		if messageType == websocket.PongMessage {
			log.Println("Received pong message from client")
			continue
		}

		var message map[string]interface{}
		if err := json.Unmarshal([]byte(messageString), &message); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			errMsg := map[string]interface{}{"error": "Invalid message format"}
			if err := ws.WriteJSON(errMsg); err != nil {
				log.Printf("Failed to send error message to WebSocket: %v", err)
				break
			}
			continue
		}

		log.Printf("Unmarshaled message: %v", message)

		text, ok := message["Text"].(string)
		if !ok {
			log.Println("Text field missing in message")
			errMsg := map[string]interface{}{"error": "Text field is required"}
			if err := ws.WriteJSON(errMsg); err != nil {
				log.Printf("Failed to send error message to WebSocket: %v", err)
				break
			}
			continue
		}

		log.Printf("Text before sending to gRPC: %s", text)

		chatMessage := &pb.ChatMessage{
			RoomId:   roomID,
			ClientId: clientID,
			Text:     text,
		}

		err = stream.Send(chatMessage)
		if err != nil {
			log.Printf("Error sending message to gRPC stream: %v", err)
			break
		}
	}
}

// ProxyToChatService godoc
// @Summary Proxy request to Chat Service for messages
// @Description Proxy GET request to /channels/{roomId}/messages to Chat Service
// @Tags Chat
// @Accept json
// @Produce json
// @Param roomId path string true "Room ID"
// @Success 200 {array} chat_protos.ChatMessage "List of messages"
// @Failure 500 {object} map[string]interface{} "Gateway server error"
// @Router /channels/{roomId}/messages [get]
func (g *GatewayHandler) GetMessagesHandler(c *gin.Context) {
	roomID := c.Param("roomId")

	// 1. Подготовка gRPC запроса
	req := &chat_grpc.ListByChannelRequest{ // Предполагается, что у вас есть ListByChannelRequest
		ChannelId: roomID, // Предполагается, что у вас есть поле ChannelId
	}

	// 2. Вызов gRPC метода
	resp, err := g.chatServiceClient.ListByChannel(context.Background(), req) //  Вызов ListByChannel
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": st.Message(),
				"code":  st.Code(),
			})
		} else {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("gRPC error: %w", err))
		}
		log.Printf("gRPC error: %v", err)
		return
	}

	// 3. Возврат результата
	c.JSON(http.StatusOK, resp.Messages) // Предполагаем, что resp содержит поле Messages с массивом сообщений
}
