package repository

import (
	"article/internal/entity"
	"context"
)

type ArticleRepository interface {
	CreateArticle(ctx context.Context, article *entity.Article) error
	GetArticleById(ctx context.Context, id string) (*entity.Article, error)
	UpdateArticle(ctx context.Context, article *entity.Article) error
	DeleteArticle(ctx context.Context, id string) error
	GetAllArticles(ctx context.Context) ([]*entity.Article, error)
}
