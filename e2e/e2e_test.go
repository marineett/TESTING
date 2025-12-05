package integration

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type emailReader struct {
	imapHost     string
	imapPort     string
	imapUser     string
	imapPassword string
}

func createEmailReader(imapHost, imapPort, imapUser, imapPassword string) *emailReader {
	return &emailReader{
		imapHost:     imapHost,
		imapPort:     imapPort,
		imapUser:     imapUser,
		imapPassword: imapPassword,
	}
}

func (r *emailReader) GetLastMessage() (string, error) {
	addr := r.imapHost + ":" + r.imapPort

	port, err := strconv.Atoi(r.imapPort)
	if err != nil {
		return "", fmt.Errorf("invalid IMAP_PORT: %v", err)
	}

	var conn net.Conn
	if port == 993 {
		tlsConfig := &tls.Config{
			ServerName:         r.imapHost,
			InsecureSkipVerify: false,
		}
		conn, err = tls.DialWithDialer(
			&net.Dialer{Timeout: 10 * time.Second},
			"tcp",
			addr,
			tlsConfig,
		)
	} else {
		conn, err = net.DialTimeout("tcp", addr, 10*time.Second)
	}

	if err != nil {
		return "", fmt.Errorf("failed to connect to IMAP server: %w", err)
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	buffer := make([]byte, 4096)
	_, err = conn.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("failed to read IMAP greeting: %w", err)
	}

	loginCmd := fmt.Sprintf("a001 LOGIN %s %s\r\n", r.imapUser, r.imapPassword)
	_, err = conn.Write([]byte(loginCmd))
	if err != nil {
		return "", fmt.Errorf("failed to send LOGIN command: %w", err)
	}

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	response, err := r.readIMAPResponse(conn, "a001")
	if err != nil {
		return "", fmt.Errorf("failed to read LOGIN response: %w", err)
	}

	if !strings.Contains(response, "a001 OK") {
		return "", fmt.Errorf("IMAP login failed: %s", response)
	}

	selectCmd := "a002 SELECT INBOX\r\n"
	_, err = conn.Write([]byte(selectCmd))
	if err != nil {
		return "", fmt.Errorf("failed to send SELECT command: %w", err)
	}

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	_, err = r.readIMAPResponse(conn, "a002")
	if err != nil {
		return "", fmt.Errorf("failed to read SELECT response: %w", err)
	}

	searchCmd := "a003 UID SEARCH ALL\r\n"
	_, err = conn.Write([]byte(searchCmd))
	if err != nil {
		return "", fmt.Errorf("failed to send SEARCH command: %w", err)
	}

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	searchResponse, err := r.readIMAPResponse(conn, "a003")
	if err != nil {
		return "", fmt.Errorf("failed to read SEARCH response: %w", err)
	}

	lastUID := "1"
	if strings.Contains(searchResponse, "SEARCH") {
		lines := strings.Split(searchResponse, "\r\n")
		for _, line := range lines {
			if strings.Contains(line, "SEARCH") && strings.HasPrefix(strings.TrimSpace(line), "*") {
				parts := strings.Fields(line)
				for i, part := range parts {
					if part == "SEARCH" && i+1 < len(parts) {
						uids := parts[i+1:]
						if len(uids) > 0 {
							lastUID = strings.TrimSpace(uids[len(uids)-1])
							lastUID = strings.Trim(lastUID, " \r\n\t")
						}
						break
					}
				}
				break
			}
		}
	}

	if lastUID == "" || lastUID == "0" {
		return "", fmt.Errorf("no messages found in mailbox")
	}

	fetchCmd := fmt.Sprintf("a004 UID FETCH %s (BODY[TEXT])\r\n", lastUID)
	_, err = conn.Write([]byte(fetchCmd))
	if err != nil {
		return "", fmt.Errorf("failed to send FETCH command: %w", err)
	}

	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	fetchResponse, err := r.readIMAPResponse(conn, "a004")
	if err != nil {
		return "", fmt.Errorf("failed to read FETCH response: %w", err)
	}

	body := fetchResponse
	startIdx := strings.Index(body, "{")
	if startIdx == -1 {
		return "", fmt.Errorf("failed to parse message body: %s", body[:min(200, len(body))])
	}

	endIdx := strings.Index(body[startIdx:], "}")
	if endIdx == -1 {
		return "", fmt.Errorf("failed to parse message body size")
	}

	sizeStr := body[startIdx+1 : startIdx+endIdx]
	size, err := strconv.Atoi(strings.TrimSpace(sizeStr))
	if err != nil {
		return "", fmt.Errorf("failed to parse message size: %w", err)
	}

	bodyStart := strings.Index(body[startIdx+endIdx:], "\r\n")
	if bodyStart == -1 {
		return "", fmt.Errorf("failed to find message body start")
	}
	bodyStart = startIdx + endIdx + bodyStart + 2

	if bodyStart+size > len(body) {
		remaining := size - (len(body) - bodyStart)
		messageText := body[bodyStart:]

		if remaining > 0 {
			conn.SetReadDeadline(time.Now().Add(10 * time.Second))
			remainingBuf := make([]byte, remaining)
			n, err := conn.Read(remainingBuf)
			if err != nil {
				return "", fmt.Errorf("failed to read remaining message: %w", err)
			}
			messageText += string(remainingBuf[:n])
		}

		return strings.TrimSpace(messageText), nil
	}

	messageText := body[bodyStart : bodyStart+size]
	messageText = strings.TrimSpace(messageText)

	logoutCmd := "a005 LOGOUT\r\n"
	conn.Write([]byte(logoutCmd))

	return messageText, nil
}

func (r *emailReader) readIMAPResponse(conn net.Conn, tag string) (string, error) {
	var response strings.Builder
	buffer := make([]byte, 4096)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if response.Len() > 0 {
				return response.String(), nil
			}
			return "", fmt.Errorf("failed to read IMAP response: %w", err)
		}

		response.Write(buffer[:n])
		responseStr := response.String()

		if strings.Contains(responseStr, tag+" OK") ||
			strings.Contains(responseStr, tag+" NO") ||
			strings.Contains(responseStr, tag+" BAD") {
			break
		}

		if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
			break
		}
	}

	return response.String(), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func extractTokenFromEmail(emailBody string) string {
	prefix := "Token:"
	idx := strings.Index(emailBody, prefix)
	if idx == -1 {
		return ""
	}
	token := strings.TrimSpace(emailBody[idx+len(prefix):])
	token = strings.Trim(token, " \r\n\t")
	return token
}

func getTokenFromEmail() (string, error) {
	imapHost := os.Getenv("IMAP_HOST")
	imapPort := os.Getenv("IMAP_PORT")
	imapUser := os.Getenv("IMAP_USER")
	imapPassword := os.Getenv("IMAP_PASSWORD")

	if imapHost == "" || imapPort == "" || imapUser == "" || imapPassword == "" {
		return "", fmt.Errorf("IMAP environment variables are not set (IMAP_HOST, IMAP_PORT, IMAP_USER, IMAP_PASSWORD)")
	}

	reader := createEmailReader(imapHost, imapPort, imapUser, imapPassword)

	time.Sleep(3 * time.Second)

	message, err := reader.GetLastMessage()
	if err != nil {
		return "", fmt.Errorf("failed to read email: %w", err)
	}

	token := extractTokenFromEmail(message)
	if token == "" {
		return "", fmt.Errorf("token not found in email message. Message content: %s", message[:min(200, len(message))])
	}

	return token, nil
}

func (c *APIClient) applyToken(ctx context.Context, token, login string) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v2/apply_token?token=%s&login=%s", c.baseURL, token, login)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.httpClient.Do(req)
}

func (s *APISuite) TestCreateUsers(t provider.T) {
	var (
		ctx      = context.Background()
		clientID int64
		login    string
		password string
	)

	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		imapUser := os.Getenv("IMAP_USER")
		sx.Require().NotEmpty(imapUser)

		login = createLogin("client", 8)
		password = createPassword(12)
	})

	t.WithNewStep("Act", func(sx provider.StepCtx) {
		var authResp struct {
			Token  string `json:"token"`
			Role   string `json:"role"`
			UserID int64  `json:"user_id"`
		}

		body := map[string]interface{}{
			"login":            login,
			"password":         password,
			"first_name":       "Иван",
			"last_name":        "Иванов",
			"middle_name":      "Иванович",
			"email":            os.Getenv("IMAP_USER"),
			"telephone_number": "+7-900-123-45-67",
			"role":             "client",
		}

		resp, err := s.c.makeRequestWithBody(ctx, "POST", "/api/v2/registration", body)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusCreated, resp.StatusCode)

		clientAuthToken, err := getTokenFromEmail()
		sx.Require().NoError(err, "Failed to get token from email for client")
		resp, err = s.c.applyToken(ctx, clientAuthToken, login)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusOK, resp.StatusCode)

		b, err := io.ReadAll(resp.Body)
		sx.Require().NoError(err)
		sx.Require().NoError(json.Unmarshal(b, &authResp))

		clientID = authResp.UserID

		sx.Require().True(clientID > 0, "clientID must be > 0")
	})

	t.WithNewStep("Assert", func(sx provider.StepCtx) {})
}

func (s *APISuite) TestWrongAccess(t provider.T) {
	ctx := context.Background()
	var login, password string

	t.WithNewStep("Arrange: register client", func(sx provider.StepCtx) {
		imapUser := os.Getenv("IMAP_USER")
		sx.Require().NotEmpty(imapUser)

		login = createLogin("client_wrong", 8)
		password = createPassword(12)

		body := map[string]interface{}{
			"login":            login,
			"password":         password,
			"first_name":       "Иван",
			"last_name":        "Иванов",
			"middle_name":      "Иванович",
			"email":            imapUser,
			"telephone_number": "+7-900-123-45-67",
			"role":             "client",
		}

		resp, err := s.c.makeRequestWithBody(ctx, "POST", "/api/v2/registration", body)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusCreated, resp.StatusCode)
	})

	t.WithNewStep("Act: 3 wrong token attempts then check lock", func(sx provider.StepCtx) {
		wrongToken := "definitely_wrong_token"

		for i := 0; i < 3; i++ {
			resp, err := s.c.applyToken(ctx, wrongToken, login)
			sx.Require().NoError(err)
			defer resp.Body.Close()
			sx.Require().Equal(http.StatusBadRequest, resp.StatusCode, "wrong attempt #%d must return 400", i+1)
		}

		resp, err := s.c.applyToken(ctx, wrongToken, login)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusTooManyRequests, resp.StatusCode, "4th wrong attempt must return 429")

		clientAuthToken, err := getTokenFromEmail()
		sx.Require().NoError(err, "Failed to get token from email for client in lock test")
		resp, err = s.c.applyToken(ctx, clientAuthToken, login)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusTooManyRequests, resp.StatusCode, "correct token must also return 429 after lock")
	})

	t.WithNewStep("Assert", func(sx provider.StepCtx) {})
}

func (s *APISuite) TestUpdateTokenFlow(t provider.T) {
	ctx := context.Background()
	var login, password string

	t.WithNewStep("Arrange: register client", func(sx provider.StepCtx) {
		imapUser := os.Getenv("IMAP_USER")
		sx.Require().NotEmpty(imapUser)

		login = createLogin("client_update", 8)
		password = createPassword(12)

		body := map[string]interface{}{
			"login":            login,
			"password":         password,
			"first_name":       "Иван",
			"last_name":        "Иванов",
			"middle_name":      "Иванович",
			"email":            imapUser,
			"telephone_number": "+7-900-123-45-67",
			"role":             "client",
		}

		resp, err := s.c.makeRequestWithBody(ctx, "POST", "/api/v2/registration", body)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusCreated, resp.StatusCode)
	})

	t.WithNewStep("Act: 3 wrong password login attempts", func(sx provider.StepCtx) {
		wrongToken := "definitely_wrong_token"

		for i := 0; i < 3; i++ {
			resp, err := s.c.applyToken(ctx, wrongToken, login)
			sx.Require().NoError(err)
			defer resp.Body.Close()
			sx.Require().Equal(http.StatusBadRequest, resp.StatusCode, "wrong attempt #%d must return 400", i+1)
		}
	})

	t.WithNewStep("Act: update token via endpoint", func(sx provider.StepCtx) {
		body := map[string]interface{}{
			"login":    login,
			"password": password,
		}
		resp, err := s.c.makeRequestWithBody(ctx, "POST", "/api/v2/update_token", body)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusOK, resp.StatusCode)
	})

	t.WithNewStep("Act: apply new token successfully", func(sx provider.StepCtx) {
		newToken, err := getTokenFromEmail()
		sx.Require().NoError(err, "Failed to get updated token from email")

		resp, err := s.c.applyToken(ctx, newToken, login)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusOK, resp.StatusCode)
	})

	t.WithNewStep("Assert", func(sx provider.StepCtx) {})
}

func (s *APISuite) TestTokenExpiration(t provider.T) {
	ctx := context.Background()
	var login, password string

	t.WithNewStep("Arrange: register client", func(sx provider.StepCtx) {
		imapUser := os.Getenv("IMAP_USER")
		sx.Require().NotEmpty(imapUser)

		login = createLogin("client_expire", 8)
		password = createPassword(12)

		body := map[string]interface{}{
			"login":            login,
			"password":         password,
			"first_name":       "Иван",
			"last_name":        "Иванов",
			"middle_name":      "Иванович",
			"email":            imapUser,
			"telephone_number": "+7-900-123-45-67",
			"role":             "client",
		}

		resp, err := s.c.makeRequestWithBody(ctx, "POST", "/api/v2/registration", body)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(http.StatusCreated, resp.StatusCode)
	})

	t.WithNewStep("Act: wait and apply expired token", func(sx provider.StepCtx) {
		time.Sleep(15 * time.Second)

		expiredToken, err := getTokenFromEmail()
		sx.Require().NoError(err, "Failed to get token from email for expiration test")

		resp, err := s.c.applyToken(ctx, expiredToken, login)
		sx.Require().NoError(err)
		defer resp.Body.Close()
		sx.Require().Equal(430, resp.StatusCode, "expired token must return 430")
	})

	t.WithNewStep("Assert", func(sx provider.StepCtx) {})
}

func TestRunAPISuite(t *testing.T) {
	suite.RunSuite(t, new(APISuite))
}
