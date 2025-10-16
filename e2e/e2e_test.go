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

type APISuite struct {
	suite.Suite
	c *APIClient
}

func (s *APISuite) BeforeAll(t provider.T) {
	s.c = NewAPIClient("http://backend:8000")
}

func (s *APISuite) TestCreateUsersAndChats(t provider.T) {
	var (
		ctx         = context.Background()
		clientID    int
		repetitorID int
		moderatorID int
	)

	t.WithNewStep("Arrange", func(sx provider.StepCtx) {})

	t.WithNewStep("Act", func(sx provider.StepCtx) {
		resp, err := s.c.makeRequestWithBody(ctx, "POST", "/api/registration/client", testClientData)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusCreated, resp.StatusCode)

		resp, err = s.c.makeRequestWithBody(ctx, "POST", "/api/registration/repetitor", testRepetitorData)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusCreated, resp.StatusCode)

		resp, err = s.c.makeRequestWithBody(ctx, "POST", "/api/registration/moderator", testModeratorData)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusCreated, resp.StatusCode)

		resp, err = s.c.makeRequestWithBody(ctx, "POST", "/api/auth/authorize", testClientAuthData)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusOK, resp.StatusCode)
		b, err := io.ReadAll(resp.Body)
		sx.Require().NoError(err)
		var auth struct {
			ID int `json:"id"`
		}
		sx.Require().NoError(json.Unmarshal(b, &auth))
		clientID = auth.ID

		resp, err = s.c.makeRequestWithBody(ctx, "POST", "/api/auth/authorize", testRepetitorAuthData)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusOK, resp.StatusCode)
		b, err = io.ReadAll(resp.Body)
		sx.Require().NoError(err)
		sx.Require().NoError(json.Unmarshal(b, &auth))
		repetitorID = auth.ID

		resp, err = s.c.makeRequestWithBody(ctx, "POST", "/api/auth/authorize", testModeratorAuthData)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusOK, resp.StatusCode)
		b, err = io.ReadAll(resp.Body)
		sx.Require().NoError(err)
		sx.Require().NoError(json.Unmarshal(b, &auth))
		moderatorID = auth.ID

		resp, err = s.c.makeRequest(ctx, "POST", fmt.Sprintf("/api/chat/start_cr_chat?c_id=%d&r_id=%d", clientID, repetitorID), nil)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusOK, resp.StatusCode)

		resp, err = s.c.makeRequest(ctx, "POST", fmt.Sprintf("/api/chat/start_cm_chat?c_id=%d&m_id=%d", clientID, moderatorID), nil)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusOK, resp.StatusCode)

		resp, err = s.c.makeRequest(ctx, "POST", fmt.Sprintf("/api/chat/start_rm_chat?r_id=%d&m_id=%d", repetitorID, moderatorID), nil)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusOK, resp.StatusCode)

		resp, err = s.c.makeRequest(ctx, "PUT", fmt.Sprintf("/api/chat/clear_messages?id=%d", 0), nil)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusOK, resp.StatusCode)

		resp, err = s.c.makeRequest(ctx, "DELETE", fmt.Sprintf("/api/chat/delete_chat?id=%d", 0), nil)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusOK, resp.StatusCode)

		resp, err = s.c.makeRequest(ctx, "PUT", fmt.Sprintf("/api/chat/clear_messages?id=%d", 0), nil)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusOK, resp.StatusCode)

		resp, err = s.c.makeRequest(ctx, "PATCH", fmt.Sprintf("/api/chat/send_message?id=%d&message=test", 0), nil)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Assert().NotEqual(http.StatusOK, resp.StatusCode)
	})

	t.WithNewStep("Assert", func(sx provider.StepCtx) {})
}

func TestRunAPISuite(t *testing.T) {
	suite.RunSuite(t, new(APISuite))
}
