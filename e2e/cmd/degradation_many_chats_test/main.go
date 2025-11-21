package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	integration "data_base_project_integration"
)

type loginData struct {
	Password string
	Role     string
	ID       int64
	Client   *integration.APIClient
}

type result struct {
	SelectedIndex int   `json:"selected_index"`
	SelectedSize  int   `json:"selected_size"`
	Timestamp     int64 `json:"timestamp"`
}

var (
	authMap map[string]*loginData
	logins  []string
)

const (
	defaultBaseURL = "http://backend:8000"
	defaultOutDir  = "/metrics/degradation_many_user_test"
)

var hardcodedChatSizes = []int{
	5000, 10000, 15000, 20000,
}

func generateLogins(n int) []string {
	res := make([]string, n)
	now := time.Now().UnixNano()
	for i := 0; i < n; i++ {
		res[i] = fmt.Sprintf("user_%d_%d@test.com", now, i)
	}
	return res
}

func mustMkdirAll(dir string) {
	_ = os.MkdirAll(dir, 0o755)
}

func loginAuthorize(login, role string, client *integration.APIClient) error {
	if authMap[login] == nil {
		authMap[login] = &loginData{}
	}
	authMap[login].Role = role
	authMap[login].Client = client
	authMap[login].Password = fmt.Sprintf("p_%d", time.Now().UnixNano())
	regBody := map[string]interface{}{
		"login":            login,
		"password":         authMap[login].Password,
		"first_name":       "Иван",
		"last_name":        "Иванов",
		"middle_name":      "Иванович",
		"email":            login,
		"telephone_number": "+7-900-000-00-00",
		"role":             role,
	}
	ctx := context.Background()
	resp, err := client.DoJSON(ctx, "POST", "/api/v2/registration", regBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		return fmt.Errorf("status code is not 201")
	}
	resp, err = client.DoJSON(ctx, "POST", "/api/v2/auth/login", map[string]string{
		"login":    login,
		"password": authMap[login].Password,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
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
	authMap[login].Client = authMap[login].Client.WithToken(authResp.Token)
	authMap[login].ID = authResp.UserID
	return nil
}

func createChatAndSpam(login1, login2 string, messagesSize int, file *os.File, batchSize int) (int64, error) {
	client1 := authMap[login1].Client
	ctx := context.Background()
	resp, err := client1.DoJSON(ctx, "POST", "/api/v2/chats", map[string]interface{}{
		"type":         "client_repetitor",
		"client_id":    authMap[login1].ID,
		"repetitor_id": authMap[login2].ID,
	})
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
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
		resp, err := client1.DoJSON(ctx, "POST", fmt.Sprintf("/api/v2/chats/%d/messages", chatID), map[string]interface{}{
			"senderId": authMap[login1].ID,
			"content":  fmt.Sprintf("test %d", i),
		})
		elapsed := time.Since(start)
		if err != nil {
			return 0, err
		}
		if resp.StatusCode != 201 {
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

func testBatch(batchSize int, file *os.File, messagesSize int) error {
	// guard against deadlocks: buffered error channel, buffered counter, timeout
	errChan := make(chan error, batchSize)
	counter := make(chan struct{}, batchSize)
	done := make(chan struct{})

	// start watcher first
	go func(size int) {
		for i := 0; i < size; i++ {
			<-counter
		}
		close(done)
	}(batchSize)

	// launch workers
	for i := 0; i < batchSize; i++ {
		go func() {
			// single-chat load for many-chats model: reuse first two users
			if _, err := createChatAndSpam(logins[0], logins[1], messagesSize, file, batchSize); err != nil {
				select {
				case errChan <- err:
				default:
				}
			}
			counter <- struct{}{}
		}()
	}

	timeout := time.After(2 * time.Minute)
	select {
	case err := <-errChan:
		return err
	case <-done:
		return nil
	case <-timeout:
		return fmt.Errorf("timeout waiting batch completion (possible deadlock)")
	}
}

func selectBatch(batchSizes []int, file *os.File, messagesSize int) (int, error) {
	for i, batchSize := range batchSizes {
		err := testBatch(batchSize, file, messagesSize)
		time.Sleep(20 * time.Second)
		if err != nil {
			fmt.Printf("Error testing batch size %d: %v\n", batchSize, err)
			if i == 0 {
				return 0, fmt.Errorf("no batches passed without errors")
			}
			return i - 1, nil
		}
	}
	return len(batchSizes) - 1, nil
}

func main() {
	baseURL := defaultBaseURL
	outputDir := defaultOutDir
	latencyFile := outputDir + "/latency_ms.csv"
	resultFile := outputDir + "/result.json"
	chatSizes := hardcodedChatSizes
	messagesSize := 300

	if cs := os.Getenv("CHAT_SIZE"); cs != "" {
		sz, err := strconv.Atoi(cs)
		if err != nil || sz <= 0 {
			fmt.Printf("ERROR: invalid CHAT_SIZE: %q\n", cs)
			os.Exit(1)
		}
		client := integration.NewAPIClient(baseURL)
		authMap = make(map[string]*loginData)
		logins = generateLogins(2)
		for i := 0; i < 2; i++ {
			role := "client"
			if i%2 == 1 {
				role = "repetitor"
			}
			if err := loginAuthorize(logins[i], role, client); err != nil {
				fmt.Printf("ERROR: loginAuthorize failed for %s: %v\n", logins[i], err)
				os.Exit(1)
			}
		}

		mustMkdirAll(outputDir)
		f, err := os.Create(latencyFile)
		if err != nil {
			fmt.Printf("ERROR: create latency file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		if _, err := f.WriteString("batch_size,latency_ms\n"); err != nil {
			fmt.Printf("ERROR: write header: %v\n", err)
			os.Exit(1)
		}

		if err := testBatch(sz, f, messagesSize); err != nil {
			fmt.Printf("ERROR: testBatch failed for size %d: %v\n", sz, err)
			os.Exit(1)
		}
		res := result{SelectedIndex: 0, SelectedSize: sz, Timestamp: time.Now().Unix()}
		bs, _ := json.Marshal(res)
		_ = os.WriteFile(resultFile, bs, 0o644)
		fmt.Printf("BATCH_RUN_SIZE=%d\n", sz)
		return
	}

	maxChatSize := chatSizes[len(chatSizes)-1]
	usersNeeded := maxChatSize * 2
	client := integration.NewAPIClient(baseURL)
	authMap = make(map[string]*loginData)
	logins = generateLogins(usersNeeded)
	for i := 0; i < usersNeeded; i++ {
		role := "client"
		if i%2 == 1 {
			role = "repetitor"
		}
		if err := loginAuthorize(logins[i], role, client); err != nil {
			fmt.Printf("ERROR: loginAuthorize failed for %s: %v\n", logins[i], err)
			os.Exit(1)
		}
	}
	mustMkdirAll(outputDir)
	f, err := os.Create(latencyFile)
	if err != nil {
		fmt.Printf("ERROR: create latency file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	if _, err := f.WriteString("batch_size,latency_ms\n"); err != nil {
		fmt.Printf("ERROR: write header: %v\n", err)
		os.Exit(1)
	}
	idx, err := selectBatch(chatSizes, f, messagesSize)
	if err != nil {
		fmt.Printf("ERROR: selectBatch failed: %v\n", err)
		os.Exit(1)
	}
	selected := chatSizes[idx]
	res := result{SelectedIndex: idx, SelectedSize: selected, Timestamp: time.Now().Unix()}
	bsBytes, _ := json.Marshal(res)
	_ = os.WriteFile(resultFile, bsBytes, 0o644)
	fmt.Printf("DEGRADATION_BATCH_SIZE=%d\n", selected)
}
