package repository

import (
	"article/internal/entity"
	"context"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *entity.Comment) error
	GetCommentById(ctx context.Context, id string) (*entity.Comment, error)
	UpdateComment(ctx context.Context, comment *entity.Comment) error
	DeleteComment(ctx context.Context, id string) error
	GetAllComments(ctx context.Context) ([]*entity.Comment, error)
	GetCommentsByArticleID(ctx context.Context, id string) ([]*entity.Comment, error)
}
