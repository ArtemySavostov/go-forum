package auth

type AuthService interface {
	GenerateToken(username string, userID string, email string) (string, error)
	ValidateToken(tokenString string) (string, error)
}
