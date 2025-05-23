package integrationtests_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

var globalToken string
var globalUserID string

func TestUserRegistrationAndCreateArticle(t *testing.T) {

	time.Sleep(35 * time.Second)
	regBody := map[string]string{
		"username": "testuser",
		"email":    "testemail@y.ru",
		"password": "testpass",
	}
	t.Logf("Starting user registration with username=%s, email=%s", regBody["username"], regBody["email"])

	regJSON, _ := json.Marshal(regBody)
	resp, err := http.Post("http://gateway:8000/register", "application/json", bytes.NewBuffer(regJSON))
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	loginBody := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}
	loginJSON, _ := json.Marshal(loginBody)

	resp, err = http.Post("http://gateway:8000/login", "application/json", bytes.NewBuffer(loginJSON))
	if err != nil {
		t.Fatalf("failed to login: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK on login, got %d", resp.StatusCode)
	}
	var loginResp struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		t.Fatalf("failed to decode login response: %v", err)
	}

	articleBody := map[string]string{
		"title":   "Test Article",
		"content": "This is a test article",
	}
	articleJSON, _ := json.Marshal(articleBody)

	req, err := http.NewRequest("POST", "http://gateway:8000/articles", bytes.NewBuffer(articleJSON))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)
	req.Header.Set("Content-Type", "application/json")
	globalToken = loginResp.Token
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("failed to send create article request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created on article creation, got %d", resp.StatusCode)
	}

	commentBody := map[string]string{
		"comment_text":   "Test comment",
		"comment_author": "1",
		"article_id":     "1",
	}
	commentJSON, _ := json.Marshal(commentBody)

	req, err = http.NewRequest("POST", "http://gateway:8000/articles/:articleId/comments", bytes.NewBuffer(commentJSON))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)
	req.Header.Set("Content-Type", "application/json")
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("failed to send create article request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created on comment creation, got %d", resp.StatusCode)
	}

}

func TestGetArticles(t *testing.T) {
	time.Sleep(1 * time.Second)
	url := "http://gateway:8000/articles"

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("failed to get articles: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	t.Logf("GetArticles response: %s", string(body))

	var articles []struct {
		ID       string `json:"id"`
		Title    string `json:"title"`
		Content  string `json:"content"`
		AuthorID string `json:"author_id"`
	}
	if err := json.Unmarshal(body, &articles); err != nil {
		t.Fatalf("failed to unmarshal articles: %v", err)
	}

	if len(articles) == 0 {
		t.Fatalf("expected at least one article")
	}
	globalUserID = articles[0].AuthorID
	t.Logf("Saved globalUserID: %s", globalUserID)

}

func TestGetComments(t *testing.T) {
	time.Sleep(1 * time.Second)
	articleID := "1"
	url := fmt.Sprintf("http://gateway:8000/articles/%s/comments", articleID)

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("failed to get comments: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	t.Logf("GetComments response: %s", string(body))

	var comments []struct {
		ID      string `json:"id"`
		Text    string `json:"comment_text"`
		Author  string `json:"comment_author"`
		Article string `json:"article_id"`
	}
	if err := json.Unmarshal(body, &comments); err != nil {
		t.Fatalf("failed to unmarshal comments: %v", err)
	}
}

func TestGetUser(t *testing.T) {
	time.Sleep(1 * time.Second)
	userID := globalUserID
	fmt.Println("userID:", userID)
	fmt.Println("userID.Hex():", userID)
	url := fmt.Sprintf("http://gateway:8000/users/%s", userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	token := globalToken
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	t.Logf("GetUser response: %s", string(body))

	var user struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	}
	if err := json.Unmarshal(body, &user); err != nil {
		t.Fatalf("failed to unmarshal user: %v", err)
	}
}
