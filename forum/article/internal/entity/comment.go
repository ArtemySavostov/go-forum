package entity

import "time"

type Comment struct {
	CommentID         string    `bson:"comment_id" json:"comment_id"`
	CommentText       string    `bson:"comment_text" json:"comment_text"`
	CommentAuthorID   string    `bson:"comment_author" json:"comment_author"`
	CreatedCommentAt  time.Time `bson:"created_comment_at" json:"created_comment_at"`
	UpdatedCommentAt  time.Time `bson:"updated_comment_at" json:"updated_comment_at"`
	CommentAuthorName string    `bson:"comment_author_name" json:"comment_author_name"`
	ArticleID         string    `bson:"article_id" json:"article_id"`
}
