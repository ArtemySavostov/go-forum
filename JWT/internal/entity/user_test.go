package entity_test

import (
	"JWT/internal/entity"
	"encoding/json"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUserMarshalJSON(t *testing.T) {
	user := entity.User{
		ID:       primitive.NewObjectID(),
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Unexpected error marshaling User: %v", err)
	}

	expectedJSON := `{"id":"` + user.ID.Hex() + `","username":"testuser","email":"test@example.com"}`
	if string(jsonData) != expectedJSON {
		t.Errorf("Unexpected JSON output:\nGot: %s\nExpected: %s", string(jsonData), expectedJSON)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	jsonData := `{"id":"65167b252e2a8d67d66b4a0f","username":"testuser","email":"test@example.com"}`
	var user entity.User
	err := json.Unmarshal([]byte(jsonData), &user)
	if err != nil {
		t.Fatalf("Unexpected error unmarshaling User: %v", err)
	}

	expectedID, _ := primitive.ObjectIDFromHex("65167b252e2a8d67d66b4a0f")
	if user.ID != expectedID || user.Username != "testuser" || user.Email != "test@example.com" {
		t.Errorf("Unexpected User data after unmarshaling")
	}
}
