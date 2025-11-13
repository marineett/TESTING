package ui

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var stdin = bufio.NewReader(os.Stdin)

func PromptWithDefault(label string, def string) string { return promptWithDefault(label, def) }
func PromptInt64(label string, def int64) int64         { return promptInt64(label, def) }
func PromptInt(label string, def int) int               { return promptInt(label, def) }
func PromptIntRequired(label string) int                { return promptIntNoDefault(label) }
func PromptInt64Required(label string) int64            { return promptInt64NoDefault(label) }

func promptWithDefault(label string, def string) string {
	if label != "" {
		fmt.Printf("%s [%s]: ", label, def)
	}
	line, _ := stdin.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return def
	}
	return line
}

func promptInt64(label string, def int64) int64 {
	for {
		s := promptWithDefault(label, fmt.Sprintf("%d", def))
		v, err := strconv.ParseInt(s, 10, 64)
		if err == nil {
			return v
		}
		fmt.Println("enter integer value")
	}
}

func promptInt(label string, def int) int {
	for {
		s := promptWithDefault(label, fmt.Sprintf("%d", def))
		v, err := strconv.Atoi(s)
		if err == nil {
			return v
		}
		fmt.Println("enter integer value")
	}
}

func promptInt64NoDefault(label string) int64 {
	for {
		if label != "" {
			fmt.Printf("%s: ", label)
		}
		line, _ := stdin.ReadString('\n')
		line = strings.TrimSpace(line)
		v, err := strconv.ParseInt(line, 10, 64)
		if err == nil {
			return v
		}
		fmt.Println("enter integer value")
	}
}

func promptIntNoDefault(label string) int {
	for {
		if label != "" {
			fmt.Printf("%s: ", label)
		}
		line, _ := stdin.ReadString('\n')
		line = strings.TrimSpace(line)
		v, err := strconv.Atoi(line)
		if err == nil {
			return v
		}
		fmt.Println("enter integer value")
	}
}

func doJSONRequest(method, url, token string, body any) (*http.Response, []byte, error) {
	var rdr io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, nil, err
		}
		rdr = bytes.NewBuffer(b)
	}
	req, err := http.NewRequest(method, url, rdr)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	b, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	return resp, b, nil
}

func printResponse(resp *http.Response, b []byte) {
	fmt.Println("Status:", resp.Status)
	if json.Valid(b) {
		var out bytes.Buffer
		_ = json.Indent(&out, b, "", "  ")
		fmt.Println(out.String())
	} else {
		fmt.Println(string(b))
	}
}

func jsonUnmarshal(b []byte, v any) error { return json.Unmarshal(b, v) }
