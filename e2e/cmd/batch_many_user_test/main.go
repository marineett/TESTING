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
	defaultOutDir  = "/metrics/batch_many_user_test"
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
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fmt.Printf("ERROR: failed to create directory %s: %v\n", dir, err)
		os.Exit(1)
	}
	fmt.Printf("Created/verified directory: %s\n", dir)
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
	go func(size int) {
		for i := 0; i < size; i++ {
			<-counter
		}
		done <- true
	}(batchSize)
	select {
	case err := <-errChan:
		return err
	case <-done:
		return nil
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	baseURL := getenv("BACKEND_URL", defaultBaseURL)
	outputDir := defaultOutDir
	latencyFile := outputDir + "/latency_ms.csv"
	resultFile := outputDir + "/result.json"
	messagesSize := 150
	bs := os.Getenv("BATCH_SIZE") // чтение нагрузки из env
	sz, err := strconv.Atoi(bs)
	if bs == "" || err != nil || sz <= 0 {
		fmt.Printf("ERROR: invalid or empty BATCH_SIZE: %q\n", bs)
		os.Exit(1)
	}

	usersNeeded := sz * 2
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

	if info, err := os.Stat(outputDir); err != nil {
		fmt.Printf("ERROR: cannot stat output directory %s: %v\n", outputDir, err)
		os.Exit(1)
	} else if !info.IsDir() {
		fmt.Printf("ERROR: output path %s is not a directory\n", outputDir)
		os.Exit(1)
	}

	testFile := outputDir + "/.write_test"
	if err := os.WriteFile(testFile, []byte("test"), 0o644); err != nil {
		fmt.Printf("ERROR: cannot write to output directory %s: %v\n", outputDir, err)
		os.Exit(1)
	}
	_ = os.Remove(testFile)
	fmt.Printf("Verified write access to: %s\n", outputDir)

	tsPath := outputDir + "/timestamps.env"
	startTs := time.Now().Unix()
	if err := os.WriteFile(tsPath, []byte(fmt.Sprintf("START_TS=%d\n", startTs)), 0o644); err != nil {
		fmt.Printf("ERROR: failed to write timestamps file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created timestamps file: %s\n", tsPath)

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
	fmt.Printf("Created latency file: %s\n", latencyFile)

	if err := testBatch(sz, f, messagesSize); err != nil {
		fmt.Printf("ERROR: testBatch failed for size %d: %v\n", sz, err)
		os.Exit(1)
	}

	sleepSec := 11
	if v := os.Getenv("POST_TEST_SLEEP"); v != "" {
		if n, e := strconv.Atoi(v); e == nil && n >= 0 {
			sleepSec = n
		}
	}
	time.Sleep(time.Duration(sleepSec) * time.Second)

	endTs := time.Now().Unix()
	if fh, err := os.OpenFile(tsPath, os.O_WRONLY|os.O_APPEND, 0o644); err == nil {
		if _, err := fh.WriteString(fmt.Sprintf("START_TS=%d\nEND_TS=%d\n", startTs, endTs)); err != nil { //временные рамки для Prometheus
			fmt.Printf("WARNING: failed to append END_TS: %v\n", err)
		} else {
			fmt.Printf("Appended END_TS to timestamps file\n")
		}
		_ = fh.Close()
	} else {
		fmt.Printf("WARNING: failed to open timestamps file for append: %v\n", err) 
	}

	// Emit result
	res := result{SelectedIndex: 0, SelectedSize: sz, Timestamp: time.Now().Unix()}
	bsBytes, _ := json.Marshal(res)
	if err := os.WriteFile(resultFile, bsBytes, 0o644); err != nil {
		fmt.Printf("ERROR: failed to write result file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created result file: %s\n", resultFile)
	fmt.Printf("BATCH_RUN_SIZE=%d\n", sz)
}
