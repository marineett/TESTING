package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

type APIClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *APIClient) makeRequest(ctx context.Context, method string, query string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + query
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return c.httpClient.Do(req)
}

func (c *APIClient) makeRequestWithBody(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(jsonData)
	}
	return c.makeRequest(ctx, method, endpoint, bodyReader)
}

func TestCreateUsersAndChats(t *testing.T) {
	ctx := context.Background()
	c := NewAPIClient("http://backend:8000")

	clientRegistrationResp, err := c.makeRequestWithBody(ctx, "POST", "/api/registration/client", testClientData)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	clientRegistrationResp.Body.Close()

	if clientRegistrationResp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status %d for client creation, got %d", http.StatusCreated, clientRegistrationResp.StatusCode)
	}

	repetitorRegistrationResp, err := c.makeRequestWithBody(ctx, "POST", "/api/registration/repetitor", testRepetitorData)
	if err != nil {
		t.Fatalf("Failed to create repetitor: %v", err)
	}
	repetitorRegistrationResp.Body.Close()

	if repetitorRegistrationResp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status %d for repetitor creation, got %d", http.StatusCreated, repetitorRegistrationResp.StatusCode)
	}

	moderatorRegistrationResp, err := c.makeRequestWithBody(ctx, "POST", "/api/registration/moderator", testModeratorData)
	if err != nil {
		t.Fatalf("Failed to create moderator: %v", err)
	}
	moderatorRegistrationResp.Body.Close()

	if moderatorRegistrationResp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status %d for moderator creation, got %d", http.StatusCreated, moderatorRegistrationResp.StatusCode)
	}

	clientAuthLoginResp, err := c.makeRequestWithBody(ctx, "POST", "/api/auth/authorize", testClientAuthData)
	if err != nil {
		t.Fatalf("Failed to login client: %v", err)
	}
	defer clientAuthLoginResp.Body.Close()

	if clientAuthLoginResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d for client login, got %d", http.StatusOK, clientAuthLoginResp.StatusCode)
	}

	body, err := io.ReadAll(clientAuthLoginResp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var authResponse struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(body, &authResponse); err != nil {
		t.Fatalf("Failed to parse auth response: %v", err)
	}

	clientID := authResponse.ID

	repetitorAuthLoginResp, err := c.makeRequestWithBody(ctx, "POST", "/api/auth/authorize", testRepetitorAuthData)
	if err != nil {
		t.Fatalf("Failed to login repetitor: %v", err)
	}
	defer repetitorAuthLoginResp.Body.Close()

	if repetitorAuthLoginResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d for repetitor login, got %d", http.StatusOK, repetitorAuthLoginResp.StatusCode)
	}

	body, err = io.ReadAll(repetitorAuthLoginResp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if err := json.Unmarshal(body, &authResponse); err != nil {
		t.Fatalf("Failed to parse auth response: %v", err)
	}

	repetitorID := authResponse.ID

	moderatorAuthLoginResp, err := c.makeRequestWithBody(ctx, "POST", "/api/auth/authorize", testModeratorAuthData)
	if err != nil {
		t.Fatalf("Failed to login moderator: %v", err)
	}
	defer moderatorAuthLoginResp.Body.Close()

	if moderatorAuthLoginResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d for moderator login, got %d", http.StatusOK, moderatorAuthLoginResp.StatusCode)
	}

	body, err = io.ReadAll(moderatorAuthLoginResp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if err := json.Unmarshal(body, &authResponse); err != nil {
		t.Fatalf("Failed to parse auth response: %v", err)
	}

	moderatorID := authResponse.ID

	crChatResp, err := c.makeRequest(ctx, "POST", fmt.Sprintf("/api/chat/start_cr_chat?c_id=%d&r_id=%d", clientID, repetitorID), nil)
	if err != nil {
		t.Fatalf("Failed to create CR chat: %v", err)
	}
	defer crChatResp.Body.Close()

	if crChatResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d for CR chat creation, got %d", http.StatusOK, crChatResp.StatusCode)
	}

	cmChatResp, err := c.makeRequest(ctx, "POST", fmt.Sprintf("/api/chat/start_cm_chat?c_id=%d&m_id=%d", clientID, moderatorID), nil)
	if err != nil {
		t.Fatalf("Failed to create CM chat: %v", err)
	}
	defer cmChatResp.Body.Close()

	if cmChatResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d for CM chat creation, got %d", http.StatusOK, cmChatResp.StatusCode)
	}

	rmChatResp, err := c.makeRequest(ctx, "POST", fmt.Sprintf("/api/chat/start_rm_chat?r_id=%d&m_id=%d", repetitorID, moderatorID), nil)
	if err != nil {
		t.Fatalf("Failed to create RM chat: %v", err)
	}
	defer rmChatResp.Body.Close()

	if rmChatResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d for RM chat creation, got %d", http.StatusOK, rmChatResp.StatusCode)
	}

	clearMessagesResp, err := c.makeRequest(ctx, "PUT", fmt.Sprintf("/api/chat/clear_messages?id=%d", 0), nil)
	if err != nil {
		t.Fatalf("Failed to clear messages: %v", err)
	}
	defer clearMessagesResp.Body.Close()

	if clearMessagesResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d for clearing messages, got %d", http.StatusOK, clearMessagesResp.StatusCode)
	}

	deleteChatResp, err := c.makeRequest(ctx, "DELETE", fmt.Sprintf("/api/chat/delete_chat?id=%d", 0), nil)
	if err != nil {
		t.Fatalf("Failed to delete chat: %v", err)
	}
	defer deleteChatResp.Body.Close()

	if deleteChatResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d for deleting chat, got %d", http.StatusInternalServerError, deleteChatResp.StatusCode)
	}

	clearMessagesResp, err = c.makeRequest(ctx, "PUT", fmt.Sprintf("/api/chat/clear_messages?id=%d", 0), nil)
	if err != nil {
		t.Fatalf("Failed to clear messages: %v", err)
	}
	defer clearMessagesResp.Body.Close()

	if clearMessagesResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d for clearing messages, got %d", http.StatusInternalServerError, clearMessagesResp.StatusCode)
	}

	sendMessageResp, err := c.makeRequest(ctx, "PATCH", fmt.Sprintf("/api/chat/send_message?id=%d&message=test", 0), nil)
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}
	defer sendMessageResp.Body.Close()

	if sendMessageResp.StatusCode == http.StatusOK {
		t.Fatalf("Expected status %d for sending message, got %d", http.StatusInternalServerError, sendMessageResp.StatusCode)
	}
}
