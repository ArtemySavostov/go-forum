package handlers

import (
	"article/internal/entity"
	"article/internal/usecase"
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"

	_ "article/internal/delivery/gin/handlers/models"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	articleUC usecase.ArticleUseCase
}

func NewArticleHandler(articleUC usecase.ArticleUseCase) *ArticleHandler {
	return &ArticleHandler{
		articleUC: articleUC,
	}
}

type CreateArticleRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// @Summary CreateArticle
// @Security ApiKeyAuth
// @Description CreateNewArticle
// @Tags Article
// @Accept json
// @Produce json
// @Param request body models.NewArticleRequest true "Create New Article"
// @Success 201 {object} models.CreateArticleResponse "Article created successfully"
// @Failure 400 {object} models.HTTPError "Bad Request"
// @Failure 401 {object} models.StatusUnauthorized "Unauthorized"
// @Failure 500 {object} models.ServerError "Failed to create article"
// @Router /articles [post]
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return

	}
	var article entity.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	article.AuthorID = userID
	article.ArticleID = uuid.New().String()
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()

	err := h.articleUC.CreateArticle(context.Background(), &article)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Article created successfully", "article_id": article.ArticleID})

}

// @Summary GetArticleByID
// @Security ApiKeyAuth
// @Description GetArticleByID
// @Tags Article
// @Accept json
// @Produce json
// @Param request body models.GetArticleByIDRequest true "Create New Article"
// @Success 201 {object} models.GetArticleByIDResponse "OK"
// @Failure 400 {object} models.HTTPError "Bad Request"
// @Failure 401 {object} models.StatusUnauthorized "Unauthorized"
// @Failure 500 {object} models.InternalServerError "Failed to get article"
// @Router /article/:id [get]
func (h *ArticleHandler) GetArticleByID(c *gin.Context) {
	articleID := c.Param("id")
	article, err := h.articleUC.GetArticleByID(context.Background(), articleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get article"})
		return
	}

	c.JSON(http.StatusOK, article)
}

// @Summary GetAllArticles
// @Description GetAllArticles
// @Tags Article
// @Accept json
// @Produce json
// @Success 201 {array} models.Article "artile"
// @Failure 400 {object} models.HTTPError "Bad Request"
// @Failure 500 {object} models.InternalServerError "Failed to get articles"
// @Router /articles [get]
func (h *ArticleHandler) GetAllArticles(c *gin.Context) {
	articles, err := h.articleUC.GetAllArticles(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get articles"})
		return
	}

	c.JSON(http.StatusOK, articles)
}

// @Summary UpdateArticle
// @Description UpdateArticle
// @Security ApiKeyAuth
// @Description UpdateArticle
// @Tags Article
// @Accept json
// @Produce json
// @Param request body models.UpdateArticleRequest true "Update srticle succesfully"
// @Success 201 {object} models.UpdateArticleResponse "Article updated successfully"
// @Failure 400 {object} models.HTTPError "Bad Request"
// @Failure 401 {object} models.StatusUnauthorized "Unauthorized"
// @Failure 500 {object} models.UpdateServerError "Failed to update article"
// @Router /articles/:id [post]
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	articleID := c.Param("id")
	userID := c.GetString("userID")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var article entity.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingArticle, err := h.articleUC.GetArticleByID(context.Background(), articleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}
	if existingArticle.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	existingArticle.Title = article.Title
	existingArticle.Content = article.Content
	existingArticle.UpdatedAt = time.Now()

	err = h.articleUC.UpdateArticle(context.Background(), existingArticle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article updated successfully"})
}

// @Summary DeleteArticle
// @Description DeleteArticle
// @Security ApiKeyAuth
// @Description UpdateArticle
// @Tags Article
// @Accept json
// @Produce json
// @Param request body models.DeleteArticleRequest true "Delete article succesfully"
// @Success 201 {object} models.DeleteArticleResponse "Article deleted successfully"
// @Failure 400 {object} models.HTTPError "Bad Request"
// @Failure 401 {object} models.StatusUnauthorized "Unauthorized"
// @Failure 403 {object} models.StatusForbidden "Forbidden"
// @Failure 500 {object} models.DeleteServerError "Failed to delete article"
// @Router /articles/:id [delete]
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	articleID := c.Param("id")
	userID := c.GetString("userID")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// existingArticle, err := h.articleUC.GetArticleByID(context.Background(), articleID)
	// if err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
	// 	return
	// }
	// if existingArticle.AuthorID != userID {
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
	// 	return
	// }

	err := h.articleUC.DeleteArticle(context.Background(), articleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}
