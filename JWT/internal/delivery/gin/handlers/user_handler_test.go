package handlers_test

import (
	"JWT/internal/entity"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUserHandler_GetUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name              string
		userID            string
		authToken         string
		expectedStatus    int
		expectedBody      string
		getUserMock       func(m *MockUserRepository, authService *MockJWTAuthService, userID primitive.ObjectID)
		validateTokenMock func(m *MockJWTAuthService, token string)
	}{
		{
			name:           "Successful GetUser",
			userID:         "649f0e648e37756407f299bc",
			authToken:      "validtoken",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":"649f0e648e37756407f299bc","username":"testuser","email":"test@example.com"}`,
			getUserMock: func(m *MockUserRepository, authService *MockJWTAuthService, userID primitive.ObjectID) {
				m.On("GetUserById", mock.Anything, userID).Return(entity.User{
					ID:       userID,
					Username: "testuser",
					Email:    "test@example.com",
					Password: "hashedpassword",
				}, nil).Once()
			},
			validateTokenMock: func(m *MockJWTAuthService, token string) {
				m.On("ValidateToken", token).Return("649f0e648e37756407f299bc", nil).Once()
			},
		},
		{
			name:           "Invalid Token",
			userID:         "649f0e648e37756407f299bc",
			authToken:      "invalidtoken",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"invalid token"}`,
			getUserMock: func(m *MockUserRepository, authService *MockJWTAuthService, userID primitive.ObjectID) {
			},
			validateTokenMock: func(m *MockJWTAuthService, token string) {
				m.On("ValidateToken", token).Return("", errors.New("invalid token")).Once()
			},
		},
		{
			name:           "User Not Found",
			userID:         "649f0e648e37756407f299bd",
			authToken:      "validtoken",
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"user not found"}`,
			getUserMock: func(m *MockUserRepository, authService *MockJWTAuthService, userID primitive.ObjectID) {
				m.On("GetUserById", mock.Anything, userID).Return(entity.User{}, errors.New("user not found")).Once()
			},
			validateTokenMock: func(m *MockJWTAuthService, token string) {
				m.On("ValidateToken", token).Return("649f0e648e37756407f299bd", nil).Once()
			},
		},
		{
			name:           "Invalid user ID",
			userID:         "invalid-object-id",
			authToken:      "validtoken",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid user ID"}`,
			getUserMock: func(m *MockUserRepository, authService *MockJWTAuthService, userID primitive.ObjectID) {
			},
			validateTokenMock: func(m *MockJWTAuthService, token string) {
				m.On("ValidateToken", token).Return("649f0e648e37756407f299bd", nil).Once()
			},
		},

		{
			name:           "Missing Token",
			userID:         "649f0e648e37756407f299bd",
			authToken:      "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"missing token"}`,
			getUserMock: func(m *MockUserRepository, authService *MockJWTAuthService, userID primitive.ObjectID) {

			},
			validateTokenMock: func(m *MockJWTAuthService, token string) {

			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockUserRepo := new(MockUserRepository)
			mockAuthService := new(MockJWTAuthService)

			userID, _ := primitive.ObjectIDFromHex(tc.userID)

			tc.getUserMock(mockUserRepo, mockAuthService, userID)
			tc.validateTokenMock(mockAuthService, tc.authToken)

			router := createTestEngineWithMocks(mockUserRepo, mockAuthService) // Используем функцию createTestEngineWithMocks

			req, _ := http.NewRequest("GET", "/users/"+tc.userID, nil)
			if tc.authToken != "" {
				req.Header.Set("Authorization", tc.authToken)
			}

			recorder := httptest.NewRecorder()

			// Act
			router.ServeHTTP(recorder, req)

			// Assert
			assert.Equal(t, tc.expectedStatus, recorder.Code)
			assert.JSONEq(t, tc.expectedBody, recorder.Body.String())

			mockUserRepo.AssertExpectations(t)
			mockAuthService.AssertExpectations(t)
		})
	}
}
