package ui

import (
	"fmt"
	"net/url"
	"strings"
)

type bodyContractCreate struct {
	ClientID    int64  `json:"client_id"`
	RepetitorID int64  `json:"repetitor_id"`
	Subject     string `json:"subject"`
}

type bodyLessonCreate struct {
	ContractID int64  `json:"contract_id"`
	Topic      string `json:"topic"`
	Format     string `json:"format"`
}

func RunContractsMenu(baseURL string, token *string) {
	api := baseURL + "/api/v2"
	for {
		fmt.Println("Contracts menu:")
		fmt.Println("1) List contracts")
		fmt.Println("2) Get contract")
		fmt.Println("3) Create contract")
		fmt.Println("4) List contract lessons")
		fmt.Println("5) Create contract lesson")
		fmt.Println("6) Create contract payment transaction")
		fmt.Println("0) Back")
		ch := PromptIntRequired("choice")
		switch ch {
		case 0:
			return
		case 1:
			offset := promptInt("offset", 0)
			size := promptInt("size", 10)
			u := fmt.Sprintf("%s/contracts?offset=%d&size=%d", api, offset, size)
			resp, b, err := doJSONRequest("GET", u, *token, nil)
			if err != nil {
				fmt.Println(err)
				break
			}
			printResponse(resp, b)
		case 2:
			id := PromptInt64Required("contractId")
			u := fmt.Sprintf("%s/contracts/%d", api, id)
			resp, b, err := doJSONRequest("GET", u, *token, nil)
			if err != nil {
				fmt.Println(err)
				break
			}
			printResponse(resp, b)
		case 3:
			c := bodyContractCreate{
				ClientID:    promptInt64("client_id", 1),
				RepetitorID: promptInt64("repetitor_id", 1),
				Subject:     promptWithDefault("subject", "math"),
			}
			u := api + "/contracts"
			resp, b, err := doJSONRequest("POST", u, *token, c)
			if err != nil {
				fmt.Println(err)
				break
			}
			printResponse(resp, b)
		case 4:
			id := PromptInt64Required("contractId")
			offset := promptInt("offset", 0)
			size := promptInt("size", 10)
			q := url.Values{}
			q.Set("offset", fmt.Sprintf("%d", offset))
			q.Set("size", fmt.Sprintf("%d", size))
			u := fmt.Sprintf("%s/contracts/%d/lessons?%s", api, id, q.Encode())
			resp, b, err := doJSONRequest("GET", u, *token, nil)
			if err != nil {
				fmt.Println(err)
				break
			}
			printResponse(resp, b)
		case 5:
			l := bodyLessonCreate{
				ContractID: PromptInt64Required("contract_id"),
				Topic:      promptWithDefault("topic", "Intro"),
				Format:     promptWithDefault("format", "online"),
			}
			u := api + "/lessons"
			resp, b, err := doJSONRequest("POST", u, *token, l)
			if err != nil {
				fmt.Println(err)
				break
			}
			printResponse(resp, b)
		case 6:
			contractID := PromptInt64Required("contractId")
			amount := promptInt64("amount", 1000)
			u := fmt.Sprintf("%s/contracts/%d/transactions?amount=%d", api, contractID, amount)
			resp, b, err := doJSONRequest("POST", u, *token, nil)
			if err != nil {
				fmt.Println(err)
				break
			}
			printResponse(resp, b)
		default:
			fmt.Println("unknown choice")
		}
	}
}

func JoinPath(parts ...string) string {
	return strings.TrimRight(strings.Join(parts, "/"), "/")
}
