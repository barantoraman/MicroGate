package integrationtests

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"testing"
	"time"
)

const baseURL = "http://microgate_api_gateway-svc:8081/v1"

func randomSuffix() string {
	nBig, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), nBig.Int64())
}

func TestUserFlow(t *testing.T) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	email := fmt.Sprintf("testuser-%s@test.com", randomSuffix())

	// Helper: authorized request
	doAuthRequest := func(method, url, token string, body io.Reader) (*http.Response, error) {
		req, err := http.NewRequest(method, url, body)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		return client.Do(req)
	}

	// --- Happy Path ---

	// Sign Up
	signupBody := fmt.Sprintf(`{"user":{"email":"%s","password":"test1234test"}}`, email)
	resp, err := client.Post(baseURL+"/signup", "application/json", bytes.NewBufferString(signupBody))
	if err != nil {
		t.Fatalf("SignUp request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("SignUp failed: expected 200 OK, got %s", resp.Status)
	}

	// Log In
	loginBody := fmt.Sprintf(`{"user":{"email":"%s","password":"test1234test"}}`, email)
	resp, err = client.Post(baseURL+"/login", "application/json", bytes.NewBufferString(loginBody))
	if err != nil {
		t.Fatalf("Login request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Login failed: expected 200 OK, got %s", resp.Status)
	}

	var loginResp struct {
		Token struct {
			Plaintext string `json:"plaintext"`
		} `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		t.Fatalf("Failed to parse login response: %v", err)
	}
	token := loginResp.Token.Plaintext
	if token == "" {
		t.Fatal("Empty token received on login")
	}

	// Add Task
	taskBody := `{
        "task": {
            "title": "sample title",
            "description": "sample description",
            "status": "sample status",
            "end": "2022-04-02T09:24:05Z"
        }
    }`
	resp, err = doAuthRequest("POST", baseURL+"/task", token, bytes.NewBufferString(taskBody))
	if err != nil {
		t.Fatalf("Add Task request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		t.Fatalf("Add Task failed: expected 200 or 201, got %s", resp.Status)
	}

	var taskResp struct {
		TaskID string `json:"task_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&taskResp); err != nil {
		t.Fatalf("Failed to parse add task response: %v", err)
	}
	taskID := taskResp.TaskID
	if taskID == "" {
		t.Fatal("Task ID is empty in add task response")
	}

	// List Task
	resp, err = doAuthRequest("GET", baseURL+"/task", token, nil)
	if err != nil {
		t.Fatalf("List Task request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("List Task failed: expected 200 OK, got %s", resp.Status)
	}

	var listResp struct {
		Tasks []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"tasks"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		t.Fatalf("Failed to parse list task response: %v", err)
	}
	found := false
	for _, task := range listResp.Tasks {
		if task.ID == taskID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Added task with ID %s not found in list", taskID)
	}

	// Delete Task
	deleteURL := fmt.Sprintf("%s/task/%s", baseURL, taskID)
	resp, err = doAuthRequest("DELETE", deleteURL, token, nil)
	if err != nil {
		t.Fatalf("Delete Task request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Delete Task failed: expected 200 OK or 204 No Content, got %s", resp.Status)
	}

	// Log Out
	logoutBody := fmt.Sprintf(`{"token":{"plaintext":"%s"}}`, token)
	resp, err = client.Post(baseURL+"/logout", "application/json", bytes.NewBufferString(logoutBody))
	if err != nil {
		t.Fatalf("Logout request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Logout failed: expected 200 OK, got %s", resp.Status)
	}

	// --- Negative Cases ---

	// Sign Up with missing fields
	badSignupBody := `{"user":{"email":"","password":""}}`
	resp, err = client.Post(baseURL+"/signup", "application/json", bytes.NewBufferString(badSignupBody))
	if err != nil {
		t.Fatalf("Bad signup request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("Expected failure on bad signup, got 200 OK")
	}

	// Duplicate Signup (using same email as before)
	resp, err = client.Post(baseURL+"/signup", "application/json", bytes.NewBufferString(signupBody))
	if err != nil {
		t.Fatalf("Duplicate signup request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("Expected failure on duplicate signup, got 200 OK")
	}

	// Login with wrong password
	wrongLoginBody := fmt.Sprintf(`{"user":{"email":"%s","password":"wrongpassword"}}`, email)
	resp, err = client.Post(baseURL+"/login", "application/json", bytes.NewBufferString(wrongLoginBody))
	if err != nil {
		t.Fatalf("Login with wrong password request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("Expected failure on wrong password login, got 200 OK")
	}

	// Access protected resource without token
	req, err := http.NewRequest("GET", baseURL+"/task", nil)
	if err != nil {
		t.Fatalf("Failed to create unauthorized request: %v", err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Unauthorized request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("Expected failure on unauthorized access, got 200 OK")
	}

	// Use token after logout
	resp, err = doAuthRequest("GET", baseURL+"/task", token, nil)
	if err != nil {
		t.Fatalf("Authorized request after logout failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("Expected failure when using token after logout, got 200 OK")
	}
}
