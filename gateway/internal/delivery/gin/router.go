package gin

import (
	"gateway/internal/delivery/gin/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(gatewayHandler *handlers.GatewayHandler) *gin.Engine {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // Разрешите ваш клиентский домен
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Accept", "Accept-Language", "Content-Language", "Origin", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	r.POST("/login", gatewayHandler.ProxyToUsersService)
	r.POST("/register", gatewayHandler.ProxyToUsersService)

	r.GET("/users/*path", gatewayHandler.ProxyToUsersService)
	r.POST("/users/*path", gatewayHandler.ProxyToUsersService)

	r.GET("/articles", gatewayHandler.ProxyToArticlesService)
	r.POST("/articles", gatewayHandler.ProxyToArticlesService)

	r.GET("/article/:id", gatewayHandler.ProxyToArticlesService)
	r.PUT("/article/:id", gatewayHandler.ProxyToArticlesService)
	r.DELETE("/article/:id", gatewayHandler.ProxyToArticlesService)

	r.GET("/articles/:articleId/comments", gatewayHandler.ProxyToArticlesService)
	r.POST("/articles/:articleId/comments", gatewayHandler.ProxyToArticlesService)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// r.POST("/chat/createRoom", gatewayHandler.CreateRoomHandler) // POST для создания комнаты
	// r.GET("/chat/room/:roomID", gatewayHandler.GetRoomHandler)   // GET для получения комнаты
	//r.GET("/ws/chat", gatewayHandler.HandleWebSocket)
	//r.GET("/channels/{roomId}/messages", gatewayHandler.GetMessagesHandler)
	//r.GET("/channels/:roomId/messages", gatewayHandler.GetMessagesHandler)
	return r
}
