package usecase

import (
	"article/internal/entity"
	"article/internal/repository"
	"context"
)

type CommentUseCase interface {
	CreateComment(ctx context.Context, comment *entity.Comment) error
	GetCommentById(ctx context.Context, id string) (*entity.Comment, error)
	UpdateComment(ctx context.Context, comment *entity.Comment) error
	DeleteComment(ctx context.Context, id string) error
	GetAllComments(ctx context.Context) ([]*entity.Comment, error)
	GetCommentsByArticleID(ctx context.Context, id string) ([]*entity.Comment, error)
}

type commentUseCase struct {
	commentRepo repository.CommentRepository
}

// GetCommentsByArticleID implements CommentUseCase.
func (c *commentUseCase) GetCommentsByArticleID(ctx context.Context, articleID string) ([]*entity.Comment, error) {
	return c.commentRepo.GetCommentsByArticleID(ctx, articleID)
}

// CreateComment implements CommentUseCase.
func (c *commentUseCase) CreateComment(ctx context.Context, comment *entity.Comment) error {
	return c.commentRepo.CreateComment(ctx, comment)
}

// DeleteComment implements CommentUseCase.
func (c *commentUseCase) DeleteComment(ctx context.Context, id string) error {
	return c.commentRepo.DeleteComment(ctx, id)
}

// GetAllComments implements CommentUseCase.
func (c *commentUseCase) GetAllComments(ctx context.Context) ([]*entity.Comment, error) {
	return c.commentRepo.GetAllComments(ctx)
}

// GetCommentById implements CommentUseCase.
func (c *commentUseCase) GetCommentById(ctx context.Context, id string) (*entity.Comment, error) {
	return c.commentRepo.GetCommentById(ctx, id)
}

// UpdateComment implements CommentUseCase.
func (c *commentUseCase) UpdateComment(ctx context.Context, comment *entity.Comment) error {
	return c.commentRepo.UpdateComment(ctx, comment)
}

func NewCommentUseCase(commentRepo repository.CommentRepository) CommentUseCase {
	return &commentUseCase{commentRepo: commentRepo}
}
