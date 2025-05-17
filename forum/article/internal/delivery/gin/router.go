package gin

import (
	"article/internal/delivery/gin/handlers"
	"log"
	"net/http"
	"strings"

	"github.com/ArtemySavostov/JWT-Golang-mongodb/pkg/auth"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(articleHandler *handlers.ArticleHandler, commentHandler *handlers.CommentHandler) *gin.Engine {
	router := gin.Default()

	router.Use(CORSMiddleware())

	router.GET("/articles", articleHandler.GetAllArticles)

	router.GET("/article/:id", articleHandler.GetArticleByID)
	router.POST("/articles", AuthMiddleware(), articleHandler.CreateArticle)
	router.GET("/comments", commentHandler.GetAllComments)
	router.GET("/comments/:id", commentHandler.GetCommentById)
	router.POST("/articles/:articleId/comments", AuthMiddleware(), commentHandler.CreateComment)
	router.DELETE("/comments/:id", AuthMiddleware(), commentHandler.DeleteComment)
	router.POST("/comments/:id", AuthMiddleware(), commentHandler.UpdateComment)
	router.GET("/articles/:articleId/comments", commentHandler.GetCommentsByArticleID)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	protected := router.Group("/articles")
	protected.Use(AuthMiddleware())
	{
		protected.POST("/", articleHandler.CreateArticle)
		// protected.PUT("/:id", articleHandler.UpdateArticle)
		// protected.DELETE("/:id", articleHandler.DeleteArticle)
	}
	admin := router.Group("/admin")
	admin.Use(AuthMiddleware())
	admin.Use(RoleMiddleware("admin"))

	{
		//admin.PUT("/:id", articleHandler.UpdateArticle)
		admin.DELETE("/:id", articleHandler.DeleteArticle)
	}

	return router
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func GetUserInfo(c *gin.Context) (string, uint, error) {
	username, ok := c.Get("username")
	if !ok {
		return "", 0, gin.Error{
			Err:  nil,
			Type: 0,
			Meta: "username not found in context",
		}
	}

	usernameStr, ok := username.(string)

	if !ok {
		return "", 0, gin.Error{
			Err:  nil,
			Type: 0,
			Meta: "username is not a string",
		}
	}

	userID, ok := c.Get("user_id")
	if !ok {
		return "", 0, gin.Error{
			Err:  nil,
			Type: 0,
			Meta: "user_id not found in context",
		}
	}

	userIDUint, ok := userID.(uint)

	if !ok {
		userIDInt, ok := userID.(int)

		if !ok {
			return "", 0, gin.Error{
				Err:  nil,
				Type: 0,
				Meta: "user_id is not a uint",
			}
		}

		userIDUint = uint(userIDInt)
	}
	return usernameStr, userIDUint, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		username, userID, role, err := auth.ValidateToken(tokenString)
		if err != nil {
			log.Printf("Token validation failed: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		c.Set("role", role)
		c.Set("username", username)

		c.Set("userID", userID)
		c.Next()
	}
}

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		if role != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}

		c.Next()
	}
}
