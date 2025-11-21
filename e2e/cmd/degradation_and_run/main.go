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
	SelectedIndex int    `json:"selected_index"`
	SelectedSize  int    `json:"selected_size"`
	Timestamp     int64  `json:"timestamp"`
	TestType      string `json:"test_type"`
}

var (
	authMap map[string]*loginData
	logins  []string
)

const (
	defaultBaseURL = "http://backend:8000"
	defaultOutDir  = "/metrics/degradation_many_user_test"
)

var (
	hardcodedBatchSizes = []int{
		70, 75, 80, 85,
		90, 95, 100, 105,
		110, 115, 120, 125,
		130, 135, 140, 145,
		150, 155, 160, 165,
		170, 175, 180, 185,
		190, 195, 200,
	}
)

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

func createChatAndSpam(login1, login2 string, messagesSize int, file *os.File, label int) (int64, error) {
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
			fmt.Fprintf(file, "%d,%d\n", label, elapsed.Milliseconds())
		}
		resp.Body.Close()
	}
	return chatID, nil
}

func testBatchUsers(batchSize, messagesSize int, file *os.File) error {
	errChan := make(chan error, batchSize)
	counter := make(chan struct{}, batchSize)
	done := make(chan struct{})

	go func(size int) {
		for i := 0; i < size; i++ {
			<-counter
		}
		close(done)
	}(batchSize)

	for i := 0; i < batchSize; i++ {
		login1 := logins[i*2]
		login2 := logins[i*2+1]
		go func(l1, l2 string) {
			if _, err := createChatAndSpam(l1, l2, messagesSize, file, batchSize); err != nil {
				select {
				case errChan <- err:
				default:
				}
			}
			counter <- struct{}{}
		}(login1, login2)
	}

	timeout := time.After(3 * time.Minute)
	select {
	case err := <-errChan:
		return err
	case <-done:
		return nil
	case <-timeout:
		return fmt.Errorf("timeout waiting users batch")
	}
}

func selectBatchUsers(batchSizes []int, warmupMessages int) (int, error) {
	for i, batchSize := range batchSizes {
		if err := testBatchUsers(batchSize, warmupMessages, nil); err != nil {
			fmt.Printf("Error testing batch size %d: %v\n", batchSize, err)
			if i == 0 {
				return 0, fmt.Errorf("no batches passed without errors")
			}
			return i - 1, nil
		}
		time.Sleep(3 * time.Second)
	}
	return len(batchSizes) - 1, nil
}

func main() {
	baseURL := getenv("BACKEND_URL", defaultBaseURL)
	outputDir := getenv("OUTPUT_DIR", defaultOutDir)
	postSleep := getenvInt("POST_TEST_SLEEP", 11)

	latencyFile := outputDir + "/latency_ms.csv"
	resultFile := outputDir + "/result.json"

	client := integration.NewAPIClient(baseURL)
	authMap = make(map[string]*loginData)

	maxBatch := hardcodedBatchSizes[len(hardcodedBatchSizes)-1]
	usersNeeded := maxBatch * 2
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
	messagesSize := 150
	idx, err := selectBatchUsers(hardcodedBatchSizes, messagesSize)
	if err != nil {
		fmt.Printf("ERROR: selectBatchUsers failed: %v\n", err)
		os.Exit(1)
	}
	selected := hardcodedBatchSizes[idx]
	time.Sleep(50 * time.Second)

	mustMkdirAll(outputDir)
	tsPath := outputDir + "/timestamps.env"
	startTs := time.Now().Unix()
	_ = os.WriteFile(tsPath, []byte(fmt.Sprintf("START_TS=%d\n", startTs)), 0o644)

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
	if err := testBatchUsers(selected, messagesSize, f); err != nil {
		fmt.Printf("ERROR: testBatchUsers failed at selected size %d: %v\n", selected, err)
		os.Exit(1)
	}

	time.Sleep(time.Duration(postSleep) * time.Second)
	endTs := time.Now().Unix()
	if fh, err := os.OpenFile(tsPath, os.O_WRONLY|os.O_APPEND, 0o644); err == nil {
		_, _ = fh.WriteString(fmt.Sprintf("START_TS=%d\nEND_TS=%d\n", startTs, endTs))
		_ = fh.Close()
	}
	res := result{SelectedIndex: idx, SelectedSize: selected, Timestamp: time.Now().Unix(), TestType: "many_users"}
	bsBytes, _ := json.Marshal(res)
	_ = os.WriteFile(resultFile, bsBytes, 0o644)
	fmt.Printf("SELECTED_BATCH_SIZE=%d\n", selected)
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}
