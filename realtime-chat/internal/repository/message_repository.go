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
	DeleteOldMessages(ctx context.Context, olderThan time.Time) error
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
	message.ID = generateID()
	message.Timestamp = time.Now()
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message to JSON: %w", err)
	}

	key := fmt.Sprintf("chat:%s:%s", message.Channel, message.ID)
	err = r.rdb.Set(ctx, key, messageJSON, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to set message in Redis: %w", err)
	}

	channelKey := fmt.Sprintf("channel:%s", message.Channel)
	err = r.rdb.RPush(ctx, channelKey, key).Err()
	if err != nil {
		return fmt.Errorf("failed to push message key to channel list: %w", err)
	}

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

func (r *messageRepository) DeleteOldMessages(ctx context.Context, olderThan time.Time) error {
	var cursor uint64 = 0
	batchSize := int64(100)

	for {
		var keys []string
		var err error
		keys, cursor, err = r.rdb.Scan(ctx, cursor, "chat:*", batchSize).Result()
		if err != nil && err != redis.Nil {
			return fmt.Errorf("failed to scan Redis keys: %w", err)
		}

		for _, key := range keys {
			messageJSON, err := r.rdb.Get(ctx, key).Result()
			if err != nil && err != redis.Nil {
				log.Printf("failed to get message from Redis: %v", err)
				continue
			}

			if err == redis.Nil {
				log.Printf("message key not found in Redis: %s", key)
				continue
			}
			var message entity.Message
			err = json.Unmarshal([]byte(messageJSON), &message)
			if err != nil {
				log.Printf("failed to unmarshal message from JSON: %v", err)
				continue
			}

			if message.Timestamp.Before(olderThan) {

				err := r.rdb.Del(ctx, key).Err()
				if err != nil {
					log.Printf("failed to delete message from Redis: %v", err)
					continue
				}
				log.Printf("Deleted old message: %s", key)

				channelKey := fmt.Sprintf("channel:%s", message.Channel)
				err = r.rdb.LRem(ctx, channelKey, 0, key).Err()
				if err != nil {
					log.Printf("failed to remove message key from channel list: %v", err)
					continue
				}
			}
		}

		if cursor == 0 {

			break
		}
	}
	return nil
}
