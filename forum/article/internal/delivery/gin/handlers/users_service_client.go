package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type UsersService interface {
	GetUserByID(id string, token string) (*User, error)
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"username"`
}

func NewUsersServiceClient(baseURL string) UsersService {
	return &usersServiceClient{
		client:  &http.Client{},
		baseURL: baseURL,
	}
}

type usersServiceClient struct {
	client  *http.Client
	baseURL string
}

func (c *usersServiceClient) GetUserByID(id string, token string) (*User, error) {
	url := fmt.Sprintf("%s/users/%s", c.baseURL, id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &user, nil
}
