package integration

import (
	"bytes"
	"context"
	cryptoRand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type APIClient struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

var sharedHTTPClient *http.Client

func NewAPIClient(baseURL string) *APIClient {
	if sharedHTTPClient == nil {
		transport := &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			DisableKeepAlives:     false,
			MaxIdleConns:          2048,
			MaxIdleConnsPerHost:   512,
			MaxConnsPerHost:       512,
			IdleConnTimeout:       120 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			ForceAttemptHTTP2:     false,
		}
		sharedHTTPClient = &http.Client{Timeout: 60 * time.Second, Transport: transport}
	}
	return &APIClient{baseURL: baseURL, httpClient: sharedHTTPClient}
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

// DoJSON is an exported wrapper to allow external packages (e.g., runners) to perform JSON requests.
func (c *APIClient) DoJSON(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
	return c.makeRequestWithBody(ctx, method, endpoint, body)
}

type APISuite struct {
	suite.Suite
	c *APIClient
}

func (s *APISuite) BeforeAll(t provider.T) {
	s.c = NewAPIClient("http://backend:8000")
}

func randHex(nBytes int) string {
	b := make([]byte, nBytes)
	_, _ = cryptoRand.Read(b)
	return hex.EncodeToString(b)
}

func createLogin(prefix string, size int) string {
	sfx := randHex(size)
	return fmt.Sprintf("%s_%s@test.com", prefix, sfx)
}

func createPassword(size int) string {
	return randHex(size)
}

func makeCreds(prefix string) (string, string) {
	login := createLogin(prefix, 8)
	password := createPassword(12)
	return login, password
}
