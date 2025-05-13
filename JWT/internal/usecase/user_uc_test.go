package usecase_test

import (
	"JWT/internal/entity"
	"JWT/internal/usecase"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockUserRepository struct {
	Users []entity.User
}

// GetByEmail implements repository.UserRepository.
func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	panic("unimplemented")
}

// GetByUsername implements repository.UserRepository.
func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (entity.User, error) {
	panic("unimplemented")
}

func (m *MockUserRepository) GetUserById(ctx context.Context, id primitive.ObjectID) (entity.User, error) {
	for _, user := range m.Users {
		if user.ID == id {
			return user, nil
		}
	}
	return entity.User{}, errors.New("user not found")
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	user.ID = primitive.NewObjectID() // Assign an ID
	m.Users = append(m.Users, *user)
	return nil
}

func (m *MockUserRepository) Update(ctx context.Context, user entity.User) error {
	for i, u := range m.Users {
		if u.ID == user.ID {
			m.Users[i] = user
			return nil
		}
	}
	return errors.New("user not found")
}

func (m *MockUserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	for i, user := range m.Users {
		if user.ID == id {
			m.Users = append(m.Users[:i], m.Users[i+1:]...)
			return nil
		}
	}
	return errors.New("user not found")
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]entity.User, error) {
	return m.Users, nil
}

func TestGetUser(t *testing.T) {
	mockRepo := &MockUserRepository{
		Users: []entity.User{
			{ID: primitive.NewObjectID(), Username: "testuser", Email: "test@example.com", Password: "password"},
		},
	}
	useCase := usecase.NewUserUseCase(mockRepo)
	user, err := useCase.GetUser(mockRepo.Users[0].ID)
	assert.NoError(t, err)
	assert.Equal(t, mockRepo.Users[0].Username, user.Username)
}

func TestCreateUser(t *testing.T) {
	mockRepo := &MockUserRepository{}
	useCase := usecase.NewUserUseCase(mockRepo)
	user, err := useCase.CreateUser("newuser", "new@example.com", "securepassword")
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, "newuser", user.Username)
	assert.Len(t, mockRepo.Users, 1) //check if user was added to the repo
}

func TestCreateUser_ValidationFails(t *testing.T) {
	mockRepo := &MockUserRepository{}
	useCase := usecase.NewUserUseCase(mockRepo)
	_, err := useCase.CreateUser("", "", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "username, email, and password are required")
}

func TestUpdateUser(t *testing.T) {
	mockRepo := &MockUserRepository{
		Users: []entity.User{
			{ID: primitive.NewObjectID(), Username: "testuser", Email: "test@example.com", Password: "password"},
		},
	}
	useCase := usecase.NewUserUseCase(mockRepo)
	updatedUser := mockRepo.Users[0]
	updatedUser.Username = "updateduser"
	err := useCase.UpdateUser(updatedUser)
	assert.NoError(t, err)
	assert.Equal(t, "updateduser", mockRepo.Users[0].Username)
}

func TestDeleteUser(t *testing.T) {
	mockRepo := &MockUserRepository{
		Users: []entity.User{
			{ID: primitive.NewObjectID(), Username: "testuser", Email: "test@example.com", Password: "password"},
		},
	}
	useCase := usecase.NewUserUseCase(mockRepo)
	err := useCase.DeleteUser(mockRepo.Users[0].ID)
	assert.NoError(t, err)
	assert.Len(t, mockRepo.Users, 0)
}

func TestGetAllUsers(t *testing.T) {
	mockRepo := &MockUserRepository{
		Users: []entity.User{
			{ID: primitive.NewObjectID(), Username: "testuser1", Email: "test1@example.com", Password: "password"},
			{ID: primitive.NewObjectID(), Username: "testuser2", Email: "test2@example.com", Password: "password"},
		},
	}
	useCase := usecase.NewUserUseCase(mockRepo)
	users, err := useCase.GetAllUsers()
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestUpdateUser_UserNotFound(t *testing.T) {
	mockRepo := &MockUserRepository{}
	useCase := usecase.NewUserUseCase(mockRepo)
	err := useCase.UpdateUser(entity.User{ID: primitive.NewObjectID(), Username: "nonexistent"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")

}

func TestDeleteUser_UserNotFound(t *testing.T) {
	mockRepo := &MockUserRepository{}
	useCase := usecase.NewUserUseCase(mockRepo)
	err := useCase.DeleteUser(primitive.NewObjectID())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}
