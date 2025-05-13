package handlers

import (
	"article/internal/entity"
	"article/internal/usecase"
	"context"
	"encoding/json"
	"fmt"

	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommentHandler struct {
	commentUC usecase.CommentUseCase
}

func NewCommentHandler(commentUC usecase.CommentUseCase) *CommentHandler {
	return &CommentHandler{
		commentUC: commentUC,
	}
}

type CreateCommentRequest struct {
}

func GetUserName(ctx context.Context, id string) (string, error) {
	url := "http://localhost:8088/users/:" + id
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call user service: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("user service returned status: %s", resp.Status)
	}
	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", fmt.Errorf("failed to decode user response: %w", err)
	}

	return user.Name, nil

}

// @Summary CreateComment
// @Security ApiKeyAuth
// @Description CreateComment
// @Tags Comment
// @Accept json
// @Produce json
// @Param request body models.NewCreateCommentRequest true "Create New Article"
// @Success 201 {object} models.CreateCommentResponse "Comment created successfully"
// @Failure 400 {object} models.CreateCommentBadRequest "Article ID is required"
// @Failure 401 {object} models.StatusUnauthorized "Unauthorized"
// @Failure 500 {object} models.CreateCommentServerError "Failed to create comment"
// @Router /articles/:articleId/comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {

	articleID := c.Param("articleId")
	if articleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Article ID is required"})
		return
	}

	var comment entity.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("userID")
	log.Printf("UserID from context: %s", userID)

	log.Printf("CommentText from JSON: %s", comment.CommentText)

	if comment.CommentText == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Comment text cannot be empty"})
		return
	}
	username, exist := c.Get("username")
	if !exist {
		log.Println("Username not found in context!")
		return
	}

	comment.CommentAuthorID = userID
	comment.CommentID = uuid.New().String()
	comment.CreatedCommentAt = time.Now()
	comment.UpdatedCommentAt = time.Now()
	comment.ArticleID = articleID
	comment.CommentAuthorName = username.(string)

	err := h.commentUC.CreateComment(context.Background(), &comment)
	log.Printf("Result from CreateComment UC: %v", err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":             "Comment created successfully",
		"comment_id":          comment.CommentID,
		"comment_author_name": comment.CommentAuthorName,
		"comment_text":        comment.CommentText,
	})
}

// @Summary GetCommentById
// @Description GetCommentById
// @Tags Comment
// @Accept json
// @Produce json
// @Param request body models.GetCommentByIdRequest true "Create New Article"
// @Success 201 {object} models.GetCommentByIdResponse "Comment created successfully"
// @Failure 500 {object} models.GetCommentByIdServerError "Failed to get comment"
// @Router /comments/:id [get]
func (h *CommentHandler) GetCommentById(c *gin.Context) {
	commentID := c.Param("id")
	comment, err := h.commentUC.GetCommentById(context.Background(), commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comment"})
		return
	}
	c.JSON(http.StatusOK, comment)
}

// @Summary UpdateComment
// @Description UpdateComment
// @Security ApiKeyAuth
// @Tags Comment
// @Accept json
// @Produce json
// @Param request body models.UpdateCommentRequest true "Create New Article"
// @Success 201 {object} models.UpdateCommentResponse "Comment created successfully"
// @Failure 400 {object} models.UpdateCommentBadRequest "Article ID is required"
// @Failure 401 {object} models.StatusUnauthorized "Unauthorized"
// @Failure 403 {object} models.StatusForbidden "Forbidden"
// @Failure 404 {object} models.StatusCommentNotFound "Comment not found"
// @Failure 500 {object} models.UpdateCommentServerError "Failed to update comment"
// @Router /comments/:id [post]
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	commentID := c.Param("id")
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var comment entity.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	existingComment, err := h.commentUC.GetCommentById(context.Background(), commentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}
	if existingComment.CommentAuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	//existingComment.CommentAuthorName = comment.CommentAuthorName
	existingComment.CommentText = comment.CommentText
	existingComment.UpdatedCommentAt = time.Now()
	err = h.commentUC.UpdateComment(context.Background(), existingComment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comment updated successfully"})
}

// @Summary DeleteComment
// @Description DeleteComment
// @Security ApiKeyAuth
// @Tags Comment
// @Accept json
// @Produce json
// @Param request body models.DeleteCommentRequest true "Dlete comment"
// @Success 201 {object} models.DeleteCommentResponse "Comment delete successfully"
// @Failure 401 {object} models.StatusUnauthorized "Unauthorized"
// @Failure 403 {object} models.StatusForbidden "Forbidden"
// @Failure 404 {object} models.StatusCommentNotFound "Comment not found"
// @Failure 500 {object} models.DeleteCommentServerError "Failed to delete comment"
// @Router /comments/:id [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentID := c.Param("id")
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	existingComment, err := h.commentUC.GetCommentById(context.Background(), commentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}
	if existingComment.CommentAuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	err = h.commentUC.DeleteComment(context.Background(), commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

// @Summary GetAllComments
// @Description GetAllComments
// @Tags Comment
// @Accept json
// @Produce json
// @Success 201 {object} models.GetAllCommentsResponse "200OK"
// @Failure 400 {object} models.UpdateCommentBadRequest "Article ID is required"
// @Failure 401 {object} models.StatusUnauthorized "Unauthorized"
// @Failure 403 {object} models.StatusForbidden "Forbidden"
// @Failure 404 {object} models.StatusCommentNotFound "Comment not found"
// @Failure 500 {object} models.GetAllCommentsServerError "Failed to get comments"
// @Router /comments [get]
func (h *CommentHandler) GetAllComments(c *gin.Context) {
	// token := c.GetString("token") // Get token from context
	// if token == "" {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	// 	return
	// }

	comments, err := h.commentUC.GetAllComments(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
		return
	}

	// for _, comment := range comments {
	// 	user, err := h.usersService.GetUserByID(comment.CommentAuthorID, token)
	// 	if err != nil {
	// 		log.Printf("Error fetching username for comment %s: %v", comment.CommentID, err)

	// 		comment.CommentAuthorName = "Anonymous"
	// 		continue
	// 	}
	// 	comment.CommentAuthorName = user.Name
	// }

	c.JSON(http.StatusOK, comments)
}

// @Summary GetCommentsByArticleID
// @Description GetCommentsByArticleID
// @Tags Comment
// @Accept json
// @Produce json
// @Param request body models.GetCommentsByArticleIDRequest true "Dlete comment"
// @Success 201 {object} models.Comment "Comments"
// @Failure 500 {object} models.GetCommentsByArticleIDServerError "Failed to get comments"
// @Router /articles/:articleId/comments [get]
func (h *CommentHandler) GetCommentsByArticleID(c *gin.Context) {
	articleID := c.Param("articleId")
	if articleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Article ID is required"})
		return
	}
	comments, err := h.commentUC.GetCommentsByArticleID(c.Request.Context(), articleID)
	if err != nil {
		log.Printf("Failed to get comments by article ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
		return
	}
	c.JSON(http.StatusOK, comments)

}
