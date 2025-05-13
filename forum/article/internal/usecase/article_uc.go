package usecase

import (
	"article/internal/entity"
	"article/internal/repository"
	"context"
)

type ArticleUseCase interface {
	CreateArticle(ctx context.Context, article *entity.Article) error
	GetArticleByID(ctx context.Context, id string) (*entity.Article, error)
	UpdateArticle(ctx context.Context, article *entity.Article) error
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
func (a *articleUseCase) CreateArticle(ctx context.Context, article *entity.Article) error {
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
func (a *articleUseCase) UpdateArticle(ctx context.Context, article *entity.Article) error {
	return a.articleRepo.UpdateArticle(ctx, article)
}
