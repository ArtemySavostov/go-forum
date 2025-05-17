package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"realtime-chat/internal/entity"
	"realtime-chat/internal/repository"

	"time"
)

type ChatUsecase interface {
	ProcessMessage(userID string, message []byte)
	ListByChannel(channelID string) ([]entity.Message, error)
	CreateMessage(message *entity.Message) error
}
type chatUsecase struct {
	hub         *Hub
	messageRepo repository.MessageRepository
}

type IncomingMessage struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func NewChatUsecase(hub *Hub, messageRepo repository.MessageRepository) ChatUsecase {
	return &chatUsecase{hub: hub, messageRepo: messageRepo}
}

func (c *chatUsecase) ProcessMessage(userID string, message []byte) {
	var incomingMessage IncomingMessage
	err := json.Unmarshal(message, &incomingMessage)
	if err != nil {
		log.Printf("Ошибка при разборе JSON: %v", err)
		return
	}

	if incomingMessage.Channel == "" || incomingMessage.Text == "" {
		log.Println("Некорректное сообщение: отсутствует канал или текст")
		return
	}

	newMessage := &entity.Message{
		Channel:   incomingMessage.Channel,
		SenderID:  userID,
		Text:      incomingMessage.Text,
		Timestamp: time.Now(),
	}
	ctx := context.Background()
	err = c.messageRepo.Create(ctx, newMessage)
	if err != nil {
		log.Printf("Ошибка при сохранении сообщения: %v", err)
		return
	}

	responseMessage, err := json.Marshal(newMessage)
	if err != nil {
		log.Printf("Ошибка при сериализации JSON: %v", err)
		return
	}

	fullMessage := []byte(fmt.Sprintf("%s: ", userID))
	fullMessage = append(fullMessage, responseMessage...)

}

func (c *chatUsecase) ListByChannel(channelID string) ([]entity.Message, error) {

	ctx := context.Background()
	messages, err := c.messageRepo.ListByChannel(ctx, channelID)
	if err != nil {
		log.Printf("Error getting messages from repository: %v", err)
		return nil, err
	}
	fmt.Println(messages)
	return messages, nil
}

func (c *chatUsecase) CreateMessage(message *entity.Message) error {
	ctx := context.Background()
	err := c.messageRepo.Create(ctx, message)
	if err != nil {
		log.Printf("Failed to create message: %v", err)
	}
	return nil
}

// func (s *chatUsecase) DeleteOldMessages(olderThan time.Time) {
// 	if s.messageRepo != nil {
// 		s.messageRepo.DeleteOldMessages(context.Background(), olderThan)
// 	}
// }
