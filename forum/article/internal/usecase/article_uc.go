package usecase

import (
	"article/internal/entity"
	"article/internal/repository"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ArticleUseCase interface {
	CreateArticle(ctx context.Context, article *entity.Article, userID string) error
	GetArticleByID(ctx context.Context, id string) (*entity.Article, error)
	UpdateArticle(ctx context.Context, articleID string, userID string, updatedArticle *entity.Article) error
	DeleteArticle(ctx context.Context, id string) error
	GetAllArticles(ctx context.Context) ([]*entity.Article, error)
}

type articleUseCase struct {
	articleRepo repository.ArticleRepository
}

func NewArticleUseCase(articleRepo repository.ArticleRepository) ArticleUseCase {
	return &articleUseCase{articleRepo: articleRepo}
}

// CreateArticle implements ArticleUseCase.
func (a *articleUseCase) CreateArticle(ctx context.Context, article *entity.Article, userID string) error {

	if article.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if article.Content == "" {
		return fmt.Errorf("content cannot be empty")
	}
	article.AuthorID = userID
	article.ArticleID = uuid.New().String()
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()
	return a.articleRepo.CreateArticle(ctx, article)
}

// DeleteArticle implements ArticleUseCase.
func (a *articleUseCase) DeleteArticle(ctx context.Context, id string) error {
	return a.articleRepo.DeleteArticle(ctx, id)
}

// GetAllArticles implements ArticleUseCase.
func (a *articleUseCase) GetAllArticles(ctx context.Context) ([]*entity.Article, error) {
	return a.articleRepo.GetAllArticles(ctx)
}

// GetArticleByID implements ArticleUseCase.
func (a *articleUseCase) GetArticleByID(ctx context.Context, id string) (*entity.Article, error) {
	return a.articleRepo.GetArticleById(ctx, id)
}

// UpdateArticle implements ArticleUseCase.
func (a *articleUseCase) UpdateArticle(ctx context.Context, articleID string, userID string, updatedArticle *entity.Article) error {
	existingArticle, err := a.articleRepo.GetArticleById(ctx, articleID)
	if err != nil {
		return err // Или более конкретную ошибку, например, ErrArticleNotFound
	}

	if existingArticle.AuthorID != userID {
		return fmt.Errorf("unauthorized to update this article") // Или ErrUnauthorized
	}

	// Обновляем поля
	existingArticle.Title = updatedArticle.Title
	existingArticle.Content = updatedArticle.Content
	existingArticle.UpdatedAt = time.Now()
	return a.articleRepo.UpdateArticle(ctx, existingArticle)
}
