package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type loginData struct {
	Password string
	Role     string
	ID       int64
	Client   *APIClient
}

var logins []string = setupLogins(5000)

var authData map[string]*loginData

var (
	messagesSize = 20000
)

func setupLogins(size int) []string {
	logins := make([]string, size)
	for i := 0; i < size; i++ {
		logins[i] = createLogin(fmt.Sprintf("user%d", i), 8)
	}
	return logins
}

func loginAuthorize(login string, role string, client *APIClient) error {
	if authData[login] == nil {
		authData[login] = &loginData{}
	}
	authData[login].Role = role
	authData[login].Client = client
	authData[login].Password = createPassword(12)
	regBody := map[string]interface{}{
		"login":            login,
		"password":         authData[login].Password,
		"first_name":       "Иван",
		"last_name":        "Иванов",
		"middle_name":      "Иванович",
		"email":            login,
		"telephone_number": "+7-900-000-00-00",
		"role":             role,
	}
	ctx := context.Background()
	resp, err := client.makeRequestWithBody(ctx, "POST", "/api/v2/registration", regBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("status code is not 201")
	}

	resp, err = client.makeRequestWithBody(ctx, "POST", "/api/v2/auth/login", map[string]string{
		"login":    login,
		"password": authData[login].Password,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code is not 200: %d", resp.StatusCode)
	}
	var authResp struct {
		Token  string `json:"token"`
		Role   string `json:"role"`
		UserID int64  `json:"user_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return err
	}
	authData[login].Client.token = authResp.Token
	authData[login].ID = authResp.UserID
	return nil
}

func createChatAndSpam(
	login1 string,
	login2 string,
	messagesSize int,
	file *os.File,
	batchSize int,
) (int64, error) {
	client1 := authData[login1].Client
	ctx := context.Background()
	resp, err := client1.makeRequestWithBody(ctx, "POST", "/api/v2/chats", map[string]interface{}{
		"type":         "client_repetitor",
		"client_id":    authData[login1].ID,
		"repetitor_id": authData[login2].ID,
	})
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("status code is not 201")
	}
	var chatResp struct {
		ID int64 `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return 0, err
	}
	chatID := chatResp.ID
	for i := 0; i < messagesSize; i++ {
		start := time.Now()
		resp, err := client1.makeRequestWithBody(ctx, "POST", fmt.Sprintf("/api/v2/chats/%d/messages", chatID), map[string]interface{}{
			"senderId": authData[login1].ID,
			"content":  fmt.Sprintf("test %d", i),
		})
		elapsed := time.Since(start)

		if err != nil {
			return 0, err
		}
		if resp.StatusCode != http.StatusCreated {
			resp.Body.Close()
			return 0, fmt.Errorf("status code is not 201")
		}
		if file != nil {
			fmt.Fprintf(file, "%d,%d\n", batchSize, elapsed.Milliseconds())
		}
		resp.Body.Close()
	}
	return chatID, nil
}

func beforeAllBenchmarks(t provider.T) {
	authData = make(map[string]*loginData)
	sharedClient := NewAPIClient("http://backend:8000")
	for i, login := range logins {
		role := ""
		if i%2 == 0 {
			role = "client"
		} else {
			role = "repetitor"
		}
		err := loginAuthorize(login, role, sharedClient)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func Counter(c chan bool, done chan bool, size int) {
	for range size {
		<-c
	}
	done <- true
}

func testBatch(batchSize int, file *os.File) error {
	errChan := make(chan error)
	counter := make(chan bool)
	done := make(chan bool)
	for i := 0; i < batchSize; i++ {
		login1 := logins[i*2]
		login2 := logins[i*2+1]
		go func() {
			_, err := createChatAndSpam(login1, login2, messagesSize, file, batchSize)
			if err != nil {
				errChan <- err
			}
			counter <- false
		}()
	}
	go Counter(counter, done, batchSize)
	select {
	case err := <-errChan:
		return err
	case <-done:
		return nil
	}
}

func selectBatch(batchSizes []int, file *os.File) (int, error) {
	for i, batchSize := range batchSizes {
		err := testBatch(batchSize, file)
		time.Sleep(20 * time.Second)
		if err != nil {
			fmt.Printf("Error testing batch size %d: %v\n", batchSize, err)
			if i == 0 {
				return 0, fmt.Errorf("No batches passed without errors")
			}
			return i - 1, nil
		}
	}
	return len(batchSizes) - 1, nil
}

func CreateChatAndSpamTest(t provider.T) {
	outPath := os.Getenv("LATENCY_FILE")
	if outPath == "" {
		outPath = "createChatAndSpam.txt"
	}
	batchSizes := []int{1, 80, 85, 90, 95, 100, 105, 110, 115, 120, 125, 130, 135, 140, 145, 150, 155, 160, 165, 170, 175, 180, 185, 190, 195, 200}
	file, err := os.Create(outPath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	file.WriteString("batch_size,latency_ms\n")
	batchIndex, err := selectBatch(batchSizes, file)
	if err != nil {
		t.Fatal(err)
	}

	if dir := filepath.Dir(outPath); dir != "" && dir != "." {
		_ = os.MkdirAll(dir, 0o755)
	}
	fmt.Printf("Selected batch size: %d\n", batchSizes[batchIndex])
	time.Sleep(20 * time.Second)
	err = testBatch(batchSizes[batchIndex], file)
	if err != nil {
		t.Fatal(err)
	}
}

type BenchmarkSuite struct {
	suite.Suite
}

func (s *BenchmarkSuite) BeforeAll(t provider.T) {
	beforeAllBenchmarks(t)
}

func (s *BenchmarkSuite) TestCreateChatAndSpamBatch(t provider.T) {
	CreateChatAndSpamTest(t)
}

type CreateChatSuite struct {
	suite.Suite
}

func (s *CreateChatSuite) BeforeAll(t provider.T)                  { beforeAllBenchmarks(t) }
func (s *CreateChatSuite) TestCreateChatAndSpamBatch(t provider.T) { CreateChatAndSpamTest(t) }

type OneBigChatSuite struct {
	suite.Suite
}

func (s *OneBigChatSuite) BeforeAll(t provider.T) { beforeAllBenchmarks(t) }

func TestRunBenchmarkSuite(t *testing.T) {
	suite.RunSuite(t, new(BenchmarkSuite))
}

func TestRunCreateChatSuite(t *testing.T) {
	suite.RunSuite(t, new(CreateChatSuite))
}

func TestRunOneBigChatSuite(t *testing.T) {
	suite.RunSuite(t, new(OneBigChatSuite))
}
