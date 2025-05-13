package entity

import "time"

type Article struct {
	ArticleID string    `json:"article_id" bson:"article_id"`
	Title     string    `json:"title" bson:"title"`
	Content   string    `json:"content" bson:"content"`
	AuthorID  string    `json:"author_id" bson:"author_id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
