package ui

import "fmt"

func RunRepetitorsMenu(baseURL string, token *string) {
	api := baseURL + "/api/v2"
	for {
		fmt.Println("Repetitors menu:")
		fmt.Println("1) List repetitors")
		fmt.Println("2) Assign contract to repetitor (PATCH)")
		fmt.Println("0) Back")
		ch := PromptIntRequired("choice")
		switch ch {
		case 0:
			return
		case 1:
			offset := promptInt("offset", 0)
			size := promptInt("size", 10)
			u := fmt.Sprintf("%s/repetitors?offset=%d&size=%d", api, offset, size)
			resp, b, err := doJSONRequest("GET", u, *token, nil)
			if err != nil {
				fmt.Println(err)
				break
			}
			printResponse(resp, b)
		case 2:
			repetitorID := PromptInt64Required("repetitorId")
			contractID := PromptInt64Required("contractId")
			u := fmt.Sprintf("%s/repetitors/%d?contractId=%d", api, repetitorID, contractID)
			resp, b, err := doJSONRequest("PATCH", u, *token, nil)
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
