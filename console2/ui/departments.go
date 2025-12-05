package ui

import "fmt"

func RunDepartmentsMenu(baseURL string, token *string) {
	api := baseURL + "/api/v2"
	for {
		fmt.Println("Departments menu:")
		fmt.Println("1) List departments")
		fmt.Println("0) Back")
		ch := PromptIntRequired("choice")
		switch ch {
		case 0:
			return
		case 1:
			offset := promptInt("offset", 0)
			size := promptInt("size", 10)
			u := fmt.Sprintf("%s/departments?offset=%d&size=%d", api, offset, size)
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
