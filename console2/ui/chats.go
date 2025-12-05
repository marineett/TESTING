package ui

import "fmt"

type bodyChatCreate struct {
	ClientID    int64 `json:"client_id"`
	RepetitorID int64 `json:"repetitor_id"`
}

type bodyMessageCreate struct {
	ChatID  int64  `json:"chat_id"`
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

func RunChatsMenu(baseURL string, token *string) {
	api := baseURL + "/api/v2"
	for {
		fmt.Println("Chats menu:")
		fmt.Println("1) List chats")
		fmt.Println("2) Create chat")
		fmt.Println("3) Send message")
		fmt.Println("0) Back")
		ch := PromptIntRequired("choice")
		switch ch {
		case 0:
			return
		case 1:
			offset := promptInt("offset", 0)
			size := promptInt("size", 10)
			u := fmt.Sprintf("%s/chats?offset=%d&size=%d", api, offset, size)
			resp, b, err := doJSONRequest("GET", u, *token, nil)
			if err != nil {
				fmt.Println(err)
				break
			}
			printResponse(resp, b)
		case 2:
			body := bodyChatCreate{ClientID: PromptInt64Required("client_id"), RepetitorID: PromptInt64Required("repetitor_id")}
			u := api + "/chats"
			resp, b, err := doJSONRequest("POST", u, *token, body)
			if err != nil {
				fmt.Println(err)
				break
			}
			printResponse(resp, b)
		case 3:
			body := bodyMessageCreate{ChatID: PromptInt64Required("chat_id"), Sender: promptWithDefault("sender", "client"), Content: promptWithDefault("content", "hi")}
			u := api + "/messages"
			resp, b, err := doJSONRequest("POST", u, *token, body)
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
