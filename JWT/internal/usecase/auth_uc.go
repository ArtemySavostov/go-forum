package usecase

import (
	"JWT/internal/entity"
	"JWT/internal/repository"
	"JWT/pkg/auth"
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	Register(ctx context.Context, username, email, password string) (string, error)
	Login(ctx context.Context, username, password string) (string, error)
	ValidateToken(token string) (string, error)
}

type authUseCase struct {
	userRepo    repository.UserRepository
	authService auth.AuthService
}

func NewAuthUseCase(userRepo repository.UserRepository, authService auth.AuthService) AuthUseCase {
	return &authUseCase{
		userRepo:    userRepo,
		authService: authService,
	}
}

func (uc *authUseCase) Register(ctx context.Context, username, email, password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("could not hash password: %w", err)
	}

	user := &entity.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	err = uc.userRepo.Create(ctx, user)
	if err != nil {
		return "", fmt.Errorf("could not create user: %w", err)
	}

	token, err := uc.authService.GenerateToken(username, user.ID.Hex(), email, user.Role)
	if err != nil {
		return "", fmt.Errorf("could not generate token: %w", err)
	}

	return token, nil
}

func (uc *authUseCase) Login(ctx context.Context, username, password string) (string, error) {
	user, err := uc.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("could not get user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := uc.authService.GenerateToken(username, user.ID.Hex(), user.Email, user.Role)
	if err != nil {
		return "", fmt.Errorf("could not generate token: %w", err)
	}

	return token, nil
}

func (uc *authUseCase) ValidateToken(token string) (string, error) {
	username, err := uc.authService.ValidateToken(token)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	return username, nil
}
