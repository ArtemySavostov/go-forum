package auth

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var ErrInvalidToken = errors.New("invalid token")

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

type JWTAuthService struct {
	SecretKey string
}

func NewJWTAuthService(secretKey string) *JWTAuthService {
	return &JWTAuthService{
		SecretKey: secretKey,
	}
}

func (s *JWTAuthService) GenerateToken(username string, userID string, email string, role string) (string, error) {
	loadEnv()

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = s.SecretKey
	}

	tokenDuration, err := time.ParseDuration(os.Getenv("JWT_DURATION"))
	if err != nil {
		tokenDuration = time.Hour * 24
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"id":       userID,
		"email":    email,
		"exp":      time.Now().Add(tokenDuration).Unix(),
		"role":     role,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("couldn't sign token: %w", err)
	}

	return tokenString, nil
}

func (s *JWTAuthService) ValidateToken(tokenString string) (string, error) {

	loadEnv()

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = s.SecretKey
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return "", fmt.Errorf("couldn't parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", ErrInvalidToken

	}

	userID, ok := claims["id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid user ID claim")
	}

	return userID, nil
}
