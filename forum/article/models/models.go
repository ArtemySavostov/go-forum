package models

import "time"

type NewArticleRequest struct {
	Title    string `json:"title" bson:"title" example:"New title"`
	Content  string `json:"content" bson:"content" example:"Content"`
	AuthorID string `json:"author_id" bson:"author_id" example:"67f64f84963c9e7f74634284"`
	Token    string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
}

type CreateArticleResponse struct {
	ArticleID string    `json:"article_id" bson:"article_id" example:"7a2be7e3-b34b-43d6-8e62-daa9f615f037"`
	Title     string    `json:"title" bson:"title" example:"New title"`
	Content   string    `json:"content" bson:"content" example:"Content"`
	AuthorID  string    `json:"author_id" bson:"author_id" example:"67f64f84963c9e7f74634284"`
	CreatedAt time.Time `json:"created_at" bson:"created_at" example:"2025-04-15T20:27:56.805+00:00"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at" example:"2025-04-15T20:27:56.805+00:00"`
}
type StatusUnauthorized struct {
	Code    int    `json:"code" example:"401"`
	Message string `json:"message" example:"Unauthorized"`
}
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Invalid request body"`
}
type ServerError struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"Failed to create article"`
}

type GetArticleByIDRequest struct {
	ArticleID string `json:"article_id" bson:"article_id" example:"7a2be7e3-b34b-43d6-8e62-daa9f615f037"`
}
type GetArticleByIDResponse struct {
	ArticleID string    `json:"article_id" bson:"article_id" example:"7a2be7e3-b34b-43d6-8e62-daa9f615f037"`
	Title     string    `json:"title" bson:"title" example:"New title"`
	Content   string    `json:"content" bson:"content" example:"Content"`
	AuthorID  string    `json:"author_id" bson:"author_id" example:"67f64f84963c9e7f74634284"`
	CreatedAt time.Time `json:"created_at" bson:"created_at" example:"2025-04-15T20:27:56.805+00:00"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at" example:"2025-04-15T20:27:56.805+00:00"`
}

type InternalServerError struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"Failed to get article"`
}

type Article struct {
	ArticleID string    `json:"article_id" bson:"article_id" example:"7a2be7e3-b34b-43d6-8e62-daa9f615f037"`
	Title     string    `json:"title" bson:"title" example:"New title"`
	Content   string    `json:"content" bson:"content" example:"Content"`
	AuthorID  string    `json:"author_id" bson:"author_id" example:"67f64f84963c9e7f74634284"`
	CreatedAt time.Time `json:"created_at" bson:"created_at" example:"2025-04-15T20:27:56.805+00:00"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at" example:"2025-04-15T20:27:56.805+00:00"`
}

type UpdateArticleRequest struct {
	ArticleID string    `json:"article_id" bson:"article_id" example:"7a2be7e3-b34b-43d6-8e62-daa9f615f037"`
	Title     string    `json:"title" bson:"title" example:"New title"`
	Content   string    `json:"content" bson:"content" example:"Content"`
	AuthorID  string    `json:"author_id" bson:"author_id" example:"67f64f84963c9e7f74634284"`
	CreatedAt time.Time `json:"created_at" bson:"created_at" example:"2025-04-15T20:27:56.805+00:00"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at" example:"2025-04-15T20:27:56.805+00:00"`
}
type UpdateArticleResponse struct {
	ArticleID string    `json:"article_id" bson:"article_id" example:"7a2be7e3-b34b-43d6-8e62-daa9f615f037"`
	Title     string    `json:"title" bson:"title" example:"New title"`
	Content   string    `json:"content" bson:"content" example:"Content"`
	AuthorID  string    `json:"author_id" bson:"author_id" example:"67f64f84963c9e7f74634284"`
	CreatedAt time.Time `json:"created_at" bson:"created_at" example:"2025-04-15T20:27:56.805+00:00"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at" example:"2025-04-15T20:27:56.805+00:00"`
}
type UpdateServerError struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"Failed to update article"`
}

type DeleteArticleRequest struct {
	ArticleID string `json:"article_id" bson:"article_id" example:"7a2be7e3-b34b-43d6-8e62-daa9f615f037"`
}

type StatusForbidden struct {
	Code    int    `json:"code" example:"403"`
	Message string `json:"message" example:"Forbiden"`
}
type DeleteServerError struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"Failed to delete article"`
}

type DeleteArticleResponse struct {
	Code    int    `json:"code" example:"201"`
	Message string `json:"message" example:"Delete article succesfully"`
}

type Comment struct {
	CommentID         string    `bson:"comment_id" json:"comment_id" example:"87bcf5bd-9725-4aef-937b-f2e651258c97"`
	CommentText       string    `bson:"comment_text" json:"comment_text" example:"Comment content"`
	CommentAuthorID   string    `bson:"comment_author" json:"comment_author" example:"6805e865276e7a3ca751f2b1"`
	CreatedCommentAt  time.Time `bson:"created_comment_at" json:"created_comment_at" example:"2025-04-21T10:44:39.309+00:00"`
	UpdatedCommentAt  time.Time `bson:"updated_comment_at" json:"updated_comment_at example:"2025-04-21T10:44:39.309+00:00"`
	CommentAuthorName string    `bson:"comment_author_name" json:"comment_author_name" example:""`
	ArticleID         string    `bson:"article_id" json:"article_id example:"7a2be7e3-b34b-43d6-8e62-daa9f615f037"`
}
type NewCreateCommentRequest struct {
	CommentText string `bson:"comment_text" json:"comment_text" example:"Comment content"`
}
type CreateCommentResponse struct {
	Code    int    `json:"code" example:"201"`
	Message string `json:"message" example:"Create comment succesfully"`
}
type CreateCommentBadRequest struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Article ID is required"`
}
type CreateCommentServerError struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"Failed to create comment"`
}

type GetCommentByIdRequest struct {
	CommentID string `bson:"comment_id" json:"comment_id" example:"87bcf5bd-9725-4aef-937b-f2e651258c97"`
}

type GetCommentByIdResponse struct {
	Code              int       `json:"code" example:"201"`
	Message           string    `json:"message" example:"Get comment succesfully"`
	CommentID         string    `bson:"comment_id" json:"comment_id" example:"87bcf5bd-9725-4aef-937b-f2e651258c97"`
	CommentText       string    `bson:"comment_text" json:"comment_text" example:"Comment content"`
	CommentAuthorID   string    `bson:"comment_author" json:"comment_author" example:"6805e865276e7a3ca751f2b1"`
	CreatedCommentAt  time.Time `bson:"created_comment_at" json:"created_comment_at" example:"2025-04-21T10:44:39.309+00:00"`
	UpdatedCommentAt  time.Time `bson:"updated_comment_at" json:"updated_comment_at example:"2025-04-21T10:44:39.309+00:00"`
	CommentAuthorName string    `bson:"comment_author_name" json:"comment_author_name" example:""`
	ArticleID         string    `bson:"article_id" json:"article_id example:"7a2be7e3-b34b-43d6-8e62-daa9f615f037"`
}

type GetCommentByIdServerError struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"Failed to get comment"`
}
type UpdateCommentRequest struct {
	Code             int       `json:"code" example:"201"`
	Message          string    `json:"message" example:"Get comment succesfully"`
	CommentID        string    `bson:"comment_id" json:"comment_id" example:"87bcf5bd-9725-4aef-937b-f2e651258c97"`
	CommentText      string    `bson:"comment_text" json:"comment_text" example:"Comment content"`
	CommentAuthorID  string    `bson:"comment_author" json:"comment_author" example:"6805e865276e7a3ca751f2b1"`
	CreatedCommentAt time.Time `bson:"created_comment_at" json:"created_comment_at" example:"2025-04-21T10:44:39.309+00:00"`
	UpdatedCommentAt time.Time `bson:"updated_comment_at" json:"updated_comment_at example:"2025-04-21T10:44:39.309+00:00"`
}
type UpdateCommentResponse struct {
	Code    int    `json:"code" example:"201"`
	Message string `json:"message" example:"Comment updated succesfully"`
}
type UpdateCommentBadRequest struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Article ID is required"`
}
type StatusCommentNotFound struct {
	Code    int    `json:"code" example:"404"`
	Message string `json:"message" example:"Comment not founf"`
}
type UpdateCommentServerError struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"Failed to update comment"`
}

type DeleteCommentRequest struct {
	CommentID string `bson:"comment_id" json:"comment_id" example:"87bcf5bd-9725-4aef-937b-f2e651258c97"`
}
type DeleteCommentResponse struct {
	Code    int    `json:"code" example:"201"`
	Message string `json:"message" example:"Comment deleted succesfully"`
}
type DeleteCommentServerError struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"Failed to delete comment"`
}
type GetAllCommentsServerError struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"Failed to get comments"`
}
type GetAllCommentsResponse struct {
	Code    int    `json:"code" example:"201"`
	Message string `json:"message" example:"Comments list..."`
}
type GetCommentsByArticleIDRequest struct {
	ArticleID string `bson:"article_id" json:"article_id example:"7a2be7e3-b34b-43d6-8e62-daa9f615f037"`
}
type GetCommentsByArticleIDServerError struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"Failed to get comments"`
}
