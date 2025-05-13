package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	UsersServiceURL    string
	ArticlesServiceURL string
	ChatServiceURL     string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	port := os.Getenv("GATEWAY_PORT")
	usersServiceURL := os.Getenv("USERS_SERVICE_URL")
	articlesServiceURL := os.Getenv("ARTICLES_SERVICE_URL")
	//chatServiceURL := os.Getenv("CHAT_SERVICE_URL")
	if usersServiceURL == "" || articlesServiceURL == "" /*|| chatServiceURL == ""*/ {
		log.Fatal("USERS_SERVICE_URL, ARTICLES_SERVICE_URL and CHAT_SERVICE_URL must be set in env variables or .env file")
	}

	return &Config{
		Port:               port,
		UsersServiceURL:    usersServiceURL,
		ArticlesServiceURL: articlesServiceURL,
		//ChatServiceURL:     chatServiceURL,
	}, nil

}
