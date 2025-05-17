package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"realtime-chat/internal/entity"
	"time"

	"github.com/redis/go-redis/v9"
)

type MessageRepository interface {
	Create(ctx context.Context, message *entity.Message) error
	ListByChannel(ctx context.Context, channel string) ([]entity.Message, error)
}

type messageRepository struct {
	rdb *redis.Client
}

func NewMessageRepository(rdb *redis.Client) *messageRepository {
	return &messageRepository{
		rdb: rdb,
	}
}

const messageTTL = 24 * time.Hour

func (r *messageRepository) Create(ctx context.Context, message *entity.Message) error {
	message.ID = generateID() // Generate a unique ID
	message.Timestamp = time.Now()
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message to JSON: %w", err)
	}

	key := fmt.Sprintf("chat:%s:%s", message.Channel, message.ID) // Construct the Redis key
	err = r.rdb.Set(ctx, key, messageJSON, 0).Err()               // Store the message in Redis
	if err != nil {
		return fmt.Errorf("failed to set message in Redis: %w", err)
	}

	// Add to a list of messages for the channel
	channelKey := fmt.Sprintf("channel:%s", message.Channel)
	err = r.rdb.RPush(ctx, channelKey, key).Err()
	if err != nil {
		return fmt.Errorf("failed to push message key to channel list: %w", err)
	}
	fmt.Println("MESSAGE!!!!", message)
	return nil
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
func (r *messageRepository) ListByChannel(ctx context.Context, channel string) ([]entity.Message, error) {
	channelKey := fmt.Sprintf("channel:%s", channel)
	messageKeys, err := r.rdb.LRange(ctx, channelKey, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get message keys from channel list: %w", err)
	}

	var messages []entity.Message
	for _, key := range messageKeys {
		messageJSON, err := r.rdb.Get(ctx, key).Result()
		if err != nil {
			log.Printf("failed to get message from Redis: %v", err)
			continue
		}

		var message entity.Message
		err = json.Unmarshal([]byte(messageJSON), &message)
		if err != nil {
			log.Printf("failed to unmarshal message from JSON: %v", err)
			continue
		}

		messages = append(messages, message)
	}

	return messages, nil
}
