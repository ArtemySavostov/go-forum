package handlers

import (
	"context"

	"net/http"

	_ "JWT/internal/delivery/gin/handlers/models"
	"JWT/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUC usecase.AuthUseCase
}

func NewAuthHandler(authUC usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

// @Summary registr a new user
// @Description Registers a new user and returns a JWT token.
// @Tags register
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration request body"
// @Success 201 {object} models.RegisterResponse "Successfully registered user"
// @Failure 400 {object} models.HTTPError "Invalid request body"
// @Failure 500 {object} models.HTTPError "Internal server error"
// @Router /register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	token, err := h.authUC.Register(context.Background(), req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token})
}

// @Summary Login
// @Description Login user
// @Tags login
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login request body"
// @Success 201 {object} models.LoginResponse "Successfully login user"
// @Failure 400 {object} models.HTTPError "Invalid request body"
// @Failure 401 {object} models.StatusUnauthorized "Status Unauthorized"
// @Failure 500 {object} models.ServerError "Internal server error"
// @Router /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	token, err := h.authUC.Login(context.Background(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
