package integrationtests

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"flag"
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

func doAuthRequest(t *testing.T, client *http.Client, method, url, token string, body any) *http.Response {
	t.Helper()
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal body: %v", err)
		}
		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(context.Background(), method, url, reader)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	return resp
}

func decodeJSON[T any](t *testing.T, r io.Reader) T {
	t.Helper()
	var v T
	if err := json.NewDecoder(r).Decode(&v); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}
	return v
}

func TestUserFlow(t *testing.T) {
	flag.Parse()
	client := &http.Client{Timeout: 15 * time.Second}
	email := fmt.Sprintf("testuser-%s@test.com", randomSuffix())
	var token, taskID string

	// --- Positive Cases ---
	t.Run("PositiveCases", func(t *testing.T) {
		t.Run("SignUp", func(t *testing.T) {
			body := map[string]any{"user": map[string]string{"email": email, "password": "test1234test"}}
			resp := doAuthRequest(t, client, "POST", baseURL+"/signup", "", body)
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected 200, got %d", resp.StatusCode)
			}
		})

		t.Run("Login", func(t *testing.T) {
			body := map[string]any{"user": map[string]string{"email": email, "password": "test1234test"}}
			resp := doAuthRequest(t, client, "POST", baseURL+"/login", "", body)
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected 200, got %d", resp.StatusCode)
			}
			loginResp := decodeJSON[struct {
				Token struct {
					Plaintext string `json:"plaintext"`
				} `json:"token"`
			}](t, resp.Body)
			token = loginResp.Token.Plaintext
			if token == "" {
				t.Fatal("empty token")
			}
		})

		t.Run("AddTask", func(t *testing.T) {
			body := map[string]any{"task": map[string]string{
				"title":       "sample title",
				"description": "sample description",
				"status":      "sample status",
				"end":         "2022-04-02T09:24:05Z",
			}}
			resp := doAuthRequest(t, client, "POST", baseURL+"/task", token, body)
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
				t.Fatalf("expected 200 or 201, got %d", resp.StatusCode)
			}
			taskResp := decodeJSON[struct {
				TaskID string `json:"task_id"`
			}](t, resp.Body)
			taskID = taskResp.TaskID
		})

		t.Run("ListTask", func(t *testing.T) {
			resp := doAuthRequest(t, client, "GET", baseURL+"/task", token, nil)
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected 200, got %d", resp.StatusCode)
			}
			listResp := decodeJSON[struct {
				Tasks []struct {
					ID string `json:"id"`
				} `json:"tasks"`
			}](t, resp.Body)
			found := false
			for _, task := range listResp.Tasks {
				if task.ID == taskID {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("task %s not found", taskID)
			}
		})

		t.Run("DeleteTask", func(t *testing.T) {
			resp := doAuthRequest(t, client, "DELETE", fmt.Sprintf("%s/task/%s", baseURL, taskID), token, nil)
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
				t.Fatalf("expected 200 or 204, got %d", resp.StatusCode)
			}
		})

		t.Run("Logout", func(t *testing.T) {
			body := map[string]any{"token": map[string]string{"plaintext": token}}
			resp := doAuthRequest(t, client, "POST", baseURL+"/logout", "", body)
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected 200, got %d", resp.StatusCode)
			}
		})
	})

	// --- Negative Cases ---
	t.Run("NegativeCases", func(t *testing.T) {
		t.Run("BadSignup", func(t *testing.T) {
			body := map[string]any{"user": map[string]string{"email": "", "password": ""}}
			resp := doAuthRequest(t, client, "POST", baseURL+"/signup", "", body)
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				t.Fatalf("expected failure, got 200")
			}
		})

		t.Run("DuplicateSignup", func(t *testing.T) {
			body := map[string]any{"user": map[string]string{"email": email, "password": "test1234test"}}
			resp := doAuthRequest(t, client, "POST", baseURL+"/signup", "", body)
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				t.Fatalf("expected failure, got 200")
			}
		})

		t.Run("WrongPasswordLogin", func(t *testing.T) {
			body := map[string]any{"user": map[string]string{"email": email, "password": "wrongpassword"}}
			resp := doAuthRequest(t, client, "POST", baseURL+"/login", "", body)
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				t.Fatalf("expected failure, got 200")
			}
		})

		t.Run("NoTokenAccess", func(t *testing.T) {
			resp := doAuthRequest(t, client, "GET", baseURL+"/task", "", nil)
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				t.Fatalf("expected failure, got 200")
			}
		})

		t.Run("TokenAfterLogout", func(t *testing.T) {
			resp := doAuthRequest(t, client, "GET", baseURL+"/task", token, nil)
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				t.Fatalf("expected failure, got 200")
			}
		})
	})
}
