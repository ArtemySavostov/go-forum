package models

type RegisterRequest struct {
	Username string `json:"username" example:"john_doe"`
	Email    string `json:"email" example:"john.doe@example.com"`
	Password string `json:"password" example:"password123"`
}

type RegisterResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
}

type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Invalid request body"`
}

type ServerError struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"Internal Server Error"`
}
type StatusUnauthorized struct {
	Code    int    `json:"code" example:"401"`
	Message string `json:"message" example:"Status Unauthorized"`
}
type LoginRequest struct {
	Username string `json:"username" example:"john_doe"`
	Password string `json:"password" example:"password123"`
}

type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
}
