package ui

import "fmt"

type bodyModeratorAssign struct {
	ModeratorID int64 `json:"moderator_id"`
	ChatID      int64 `json:"chat_id"`
}

func RunModeratorsMenu(baseURL string, token *string) {
	api := baseURL + "/api/v2"
	for {
		fmt.Println("Moderators menu:")
		fmt.Println("1) List moderators")
		fmt.Println("0) Back")
		ch := PromptIntRequired("choice")
		switch ch {
		case 0:
			return
		case 1:
			offset := promptInt("offset", 0)
			size := promptInt("size", 10)
			u := fmt.Sprintf("%s/moderators?offset=%d&size=%d", api, offset, size)
			resp, b, err := doJSONRequest("GET", u, *token, nil)
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
