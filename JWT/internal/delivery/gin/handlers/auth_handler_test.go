package handlers_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"JWT/internal/delivery/gin/handlers"
	"JWT/internal/entity"
	"JWT/internal/usecase"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username string             `bson:"username" json:"username"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"-"`
}

type UserRepository interface {
	GetUserById(ctx context.Context, id primitive.ObjectID) (User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetByUsername(ctx context.Context, username string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	GetAll(ctx context.Context) ([]User, error)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserById(ctx context.Context, id primitive.ObjectID) (entity.User, error) {
	args := m.Called(ctx, id)
	user, ok := args.Get(0).(entity.User)
	if !ok {
		return entity.User{}, args.Error(1)
	}
	return user, args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (entity.User, error) {
	args := m.Called(ctx, username)
	user, ok := args.Get(0).(entity.User)
	if !ok {
		return entity.User{}, args.Error(1)
	}
	return user, args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	args := m.Called(ctx, email)
	user, ok := args.Get(0).(entity.User)
	if !ok {
		return entity.User{}, args.Error(1)
	}
	return user, args.Error(1)
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]entity.User, error) {
	args := m.Called(ctx)
	users, ok := args.Get(0).([]entity.User)
	if !ok {
		return []entity.User{}, args.Error(1)
	}
	return users, args.Error(1)
}

type MockJWTAuthService struct {
	mock.Mock
}

func (m *MockJWTAuthService) GenerateToken(username, userID, email string) (string, error) {
	args := m.Called(username, userID, email)
	return args.String(0), args.Error(1)
}

func (m *MockJWTAuthService) ValidateToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

func AuthMiddleware(authService *MockJWTAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из заголовка Authorization
		token := c.GetHeader("Authorization")

		// Проверяем, что токен присутствует
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		// Вызываем ValidateToken из AuthService
		userID, err := authService.ValidateToken(token)

		// Проверяем, что токен валиден
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Добавляем ID пользователя в контекст
		c.Set("userID", userID)

		// Переходим к следующему обработчику
		c.Next()
	}
}

func createTestEngineWithMocks(userRepository *MockUserRepository, mockAuthService *MockJWTAuthService) *gin.Engine {
	router := gin.Default()
	router.Use(gin.Recovery())

	userUC := usecase.NewUserUseCase(userRepository)
	authUC := usecase.NewAuthUseCase(userRepository, mockAuthService)
	userHandler := handlers.NewUserHandler(userUC)
	authHandler := handlers.NewAuthHandler(authUC)

	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)
	router.GET("/users/:id", AuthMiddleware(mockAuthService), userHandler.GetUser) // Защищаем GetUser

	return router
}

type MockAuthUseCase struct {
	mock.Mock
}

func (m *MockAuthUseCase) ValidateToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

func (m *MockAuthUseCase) Register(ctx context.Context, username, email, password string) (string, error) {
	args := m.Called(ctx, username, email, password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthUseCase) Login(ctx context.Context, username, password string) (string, error) {
	args := m.Called(ctx, username, password)
	return args.String(0), args.Error(1)
}

func TestAuthHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedBody   string
		registerMock   func(m *MockAuthUseCase)
	}{
		{
			name:           "Missing Password",                                      // Изменили название
			requestBody:    `{"username": "testuser", "email": "test@example.com"}`, // Удалили поле "password"
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid request body"}`, // Ожидаем ошибку от BindJSON
			registerMock:   func(m *MockAuthUseCase) {},
		},
		{
			name:           "Hashing Error",
			requestBody:    `{"username": "testuser", "email": "test@example.com", "password": "password"}`, // Добавлены username и email
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"could not hash password"}`,
			registerMock: func(m *MockAuthUseCase) {
				m.On("Register", mock.Anything, mock.Anything, mock.Anything, "password").Return("", errors.New("could not hash password")).Once()
			},
		},
		{
			name:           "Successful Registration",
			requestBody:    `{"username": "testuser", "email": "test@example.com", "password": "password"}`,
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"token":"testtoken"}`,
			registerMock: func(m *MockAuthUseCase) {
				m.On("Register", mock.Anything, "testuser", "test@example.com", "password").Return("testtoken", nil).Once()
			},
		},
		{
			name:           "Invalid Request Body",
			requestBody:    `{"username": "testuser", "email": "test@example.com" "password": "password"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid request body"}`,
			registerMock:   func(m *MockAuthUseCase) {},
		},
		{
			name:           "Register Usecase Error",
			requestBody:    `{"username": "testuser", "email": "test@example.com", "password": "password"}`,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"register error"}`,
			registerMock: func(m *MockAuthUseCase) {
				m.On("Register", mock.Anything, "testuser", "test@example.com", "password").Return("", errors.New("register error")).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockAuthUC := new(MockAuthUseCase)
			tc.registerMock(mockAuthUC)
			handler := handlers.NewAuthHandler(mockAuthUC)

			router := gin.Default()
			router.POST("/register", handler.Register)

			req, _ := http.NewRequest("POST", "/register", bytes.NewBufferString(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			// Act
			router.ServeHTTP(recorder, req)

			// Assert
			assert.Equal(t, tc.expectedStatus, recorder.Code)
			assert.JSONEq(t, tc.expectedBody, recorder.Body.String())

			mockAuthUC.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedBody   string
		loginMock      func(m *MockAuthUseCase)
	}{
		{
			name:           "Successful Login",
			requestBody:    `{"username": "testuser", "password": "password"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"token":"testtoken"}`,
			loginMock: func(m *MockAuthUseCase) {
				m.On("Login", mock.Anything, "testuser", "password").Return("testtoken", nil)
			},
		},
		{
			name:           "Invalid Request Body",
			requestBody:    `{"username": "testuser" "password": "password"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid request body"}`,
			loginMock:      func(m *MockAuthUseCase) {},
		},
		{
			name:           "Login Usecase Error",
			requestBody:    `{"username": "testuser", "password": "password"}`,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"login error"}`,
			loginMock: func(m *MockAuthUseCase) {
				m.On("Login", mock.Anything, "testuser", "password").Return("", errors.New("login error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockAuthUC := new(MockAuthUseCase)
			tc.loginMock(mockAuthUC)
			handler := handlers.NewAuthHandler(mockAuthUC)

			router := gin.Default()
			router.POST("/login", handler.Login)

			req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			// Act
			router.ServeHTTP(recorder, req)

			// Assert
			assert.Equal(t, tc.expectedStatus, recorder.Code)
			assert.JSONEq(t, tc.expectedBody, recorder.Body.String())

			mockAuthUC.AssertExpectations(t)
		})
	}
}
