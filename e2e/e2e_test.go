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

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type APIClient struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *APIClient) WithToken(token string) *APIClient {
	clone := *c
	clone.token = token
	return &clone
}

func (c *APIClient) makeRequest(ctx context.Context, method string, query string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + query
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
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

type APISuite struct {
	suite.Suite
	c *APIClient
}

func (s *APISuite) BeforeAll(t provider.T) {
	s.c = NewAPIClient("http://backend:8000")
}

func (s *APISuite) TestCreateUsersAndChats(t provider.T) {
	var (
		ctx            = context.Background()
		clientID       int64
		repetitorID    int64
		moderatorID    int64
		clientToken    string
		repetitorToken string
		cClient        *APIClient
		cRep           *APIClient
	)

	t.WithNewStep("Arrange", func(sx provider.StepCtx) {})

	t.WithNewStep("Act", func(sx provider.StepCtx) {
		var authResp struct {
			Token  string `json:"token"`
			Role   string `json:"role"`
			UserID int64  `json:"user_id"`
		}

		resp, err := s.c.makeRequestWithBody(ctx, "POST", "/api/v2/registration", testClientData)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusCreated, resp.StatusCode)
		b, err := io.ReadAll(resp.Body)
		sx.Require().NoError(err)
		sx.Require().NoError(json.Unmarshal(b, &authResp))
		clientID = authResp.UserID
		clientToken = authResp.Token

		resp, err = s.c.makeRequestWithBody(ctx, "POST", "/api/v2/registration", testRepetitorData)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusCreated, resp.StatusCode)
		b, err = io.ReadAll(resp.Body)
		sx.Require().NoError(err)
		sx.Require().NoError(json.Unmarshal(b, &authResp))
		repetitorID = authResp.UserID
		repetitorToken = authResp.Token

		resp, err = s.c.makeRequestWithBody(ctx, "POST", "/api/v2/registration", testModeratorData)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusCreated, resp.StatusCode)
		b, err = io.ReadAll(resp.Body)
		sx.Require().NoError(err)
		sx.Require().NoError(json.Unmarshal(b, &authResp))
		moderatorID = authResp.UserID

		cClient = s.c.WithToken(clientToken)
		cRep = s.c.WithToken(repetitorToken)

		body := map[string]interface{}{
			"type":         "client_repetitor",
			"client_id":    clientID,
			"repetitor_id": repetitorID,
			"moderator_id": 0,
		}
		resp, err = cClient.makeRequestWithBody(ctx, "POST", "/api/v2/chats", body)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusOK, resp.StatusCode)
		var crChatID int64
		b, err = io.ReadAll(resp.Body)
		sx.Require().NoError(err)
		sx.Require().NoError(json.Unmarshal(b, &crChatID))

		body = map[string]interface{}{
			"type":         "client_moderator",
			"client_id":    clientID,
			"moderator_id": moderatorID,
		}
		resp, err = cClient.makeRequestWithBody(ctx, "POST", "/api/v2/chats", body)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusOK, resp.StatusCode)
		var cmChatID int64
		b, err = io.ReadAll(resp.Body)
		sx.Require().NoError(err)
		sx.Require().NoError(json.Unmarshal(b, &cmChatID))

		body = map[string]interface{}{
			"type":         "repetitor_moderator",
			"repetitor_id": repetitorID,
			"moderator_id": moderatorID,
		}
		resp, err = cRep.makeRequestWithBody(ctx, "POST", "/api/v2/chats", body)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusOK, resp.StatusCode)
		var rmChatID int64
		b, err = io.ReadAll(resp.Body)
		sx.Require().NoError(err)
		sx.Require().NoError(json.Unmarshal(b, &rmChatID))

		resp, err = cClient.makeRequest(ctx, "PUT", fmt.Sprintf("/api/v2/chats/%d", crChatID), nil)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusOK, resp.StatusCode)

		resp, err = cClient.makeRequest(ctx, "DELETE", fmt.Sprintf("/api/v2/chats/%d", crChatID), nil)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusOK, resp.StatusCode)

		msg := map[string]interface{}{
			"senderId": clientID,
			"content":  "test",
		}
		resp, err = cClient.makeRequestWithBody(ctx, "POST", fmt.Sprintf("/api/v2/chats/%d/messages?offset=%d&limit=%d", crChatID, 0, 10), msg)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Assert().NotEqual(http.StatusOK, resp.StatusCode)
	})

	t.WithNewStep("Assert", func(sx provider.StepCtx) {})
}

func TestRunAPISuite(t *testing.T) {
	suite.RunSuite(t, new(APISuite))
}
