package handlers_test

import (
	"JWT/internal/delivery/gin/handlers"
	"JWT/pkg/auth"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) GenerateToken(username, userID, email, role string) (string, error) {
	args := m.Called(username, userID, email)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) ValidateToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name              string
		authHeader        string
		expectedStatus    int
		expectedBody      string
		validateTokenMock func(m *MockAuthService, token string)
	}{
		{
			name:           "Missing Authorization Header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"authorization header required"}`,
			validateTokenMock: func(m *MockAuthService, token string) {

			},
		},
		{
			name:           "Invalid Authorization Header Format",
			authHeader:     "InvalidFormat",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"invalid authorization header format"}`,
			validateTokenMock: func(m *MockAuthService, token string) {

			},
		},
		{
			name:           "Invalid Token",
			authHeader:     "Bearer invalidtoken",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"invalid token"}`,
			validateTokenMock: func(m *MockAuthService, token string) {
				m.On("ValidateToken", token).Return("", auth.ErrInvalidToken).Once()
			},
		},
		{
			name:           "Valid Token",
			authHeader:     "Bearer validtoken",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"OK"}`, // Изменим тело ответа, чтобы отличать от ошибок
			validateTokenMock: func(m *MockAuthService, token string) {
				m.On("ValidateToken", token).Return("user123", nil).Once()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockAuthService := new(MockAuthService)
			tc.validateTokenMock(mockAuthService, extractToken(tc.authHeader))

			router := gin.New()
			router.Use(handlers.AuthMiddleware(mockAuthService)) // Применяем middleware

			// Определяем тестовый handler, который будет вызван после middleware
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "OK"})
			})

			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", tc.authHeader)
			recorder := httptest.NewRecorder()

			// Act
			router.ServeHTTP(recorder, req)

			// Assert
			assert.Equal(t, tc.expectedStatus, recorder.Code)
			assert.JSONEq(t, tc.expectedBody, recorder.Body.String())

			mockAuthService.AssertExpectations(t)
		})
	}
}

func extractToken(authHeader string) string {
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1]
	}
	return ""
}
