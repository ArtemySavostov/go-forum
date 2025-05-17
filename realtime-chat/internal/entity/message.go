package entity

import "time"

type Message struct {
	ID         string    `bson:"_id,omitempty"`
	Channel    string    `bson:"channel"`
	SenderID   string    `bson:"sender_id"`
	Text       string    `bson:"text"`
	Timestamp  time.Time `bson:"timestamp"`
	Type       string    `bson:"type" `
	SenderName string    `bson:"sender_name"`
}
