package ui

import "fmt"

type bodyAuthLogin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type bodyUserRegister struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	MiddleName      string `json:"middle_name"`
	Email           string `json:"email"`
	TelephoneNumber string `json:"telephone_number"`
	Login           string `json:"login"`
	Password        string `json:"password"`
	Role            string `json:"role"`
	Resume          string `json:"resume,omitempty"`
	Salary          int64  `json:"salary,omitempty"`
}

func RunAuthMenu(baseURL string, token *string) {
	api := baseURL + "/api/v2"
	for {
		fmt.Println("Auth menu:")
		fmt.Println("1) Login (get JWT)")
		fmt.Println("2) Register user (client/repetitor/moderator/admin)")
		fmt.Println("3) Set token manually")
		fmt.Println("4) Show current token (first 16 chars)")
		fmt.Println("0) Exit")
		ch := PromptIntRequired("choice")
		switch ch {
		case 0:
			return
		case 1:
			login := promptWithDefault("login", "user@test.com")
			password := promptWithDefault("password", "password123")
			url := api + "/auth/login"
			body := bodyAuthLogin{Login: login, Password: password}
			resp, b, err := doJSONRequest("POST", url, *token, body)
			if err != nil {
				fmt.Println(err)
				break
			}
			printResponse(resp, b)
			tok := ExtractJSONField(b, "token")
			if tok != "" {
				*token = tok
				fmt.Println("Token updated.")
			}
		case 2:
			role := promptWithDefault("role (client/repetitor/moderator/admin)", "client")
			first := promptWithDefault("first_name", "Иван")
			last := promptWithDefault("last_name", "Иванов")
			middle := promptWithDefault("middle_name", "Иванович")
			email := promptWithDefault("email", "user@test.com")
			phone := promptWithDefault("telephone_number", "+7-900-000-00-00")
			login := promptWithDefault("login", "user_login")
			password := promptWithDefault("password", "password123")
			reg := bodyUserRegister{FirstName: first, LastName: last, MiddleName: middle, Email: email, TelephoneNumber: phone, Login: login, Password: password, Role: role}
			if role == "repetitor" {
				reg.Resume = promptWithDefault("resume (text)", "Опыт 5 лет")
			}
			if role == "moderator" || role == "admin" {
				reg.Salary = promptInt64("salary", 50000)
			}
			url := api + "/registration"
			resp, b, err := doJSONRequest("POST", url, *token, reg)
			if err != nil {
				fmt.Println(err)
				break
			}
			printResponse(resp, b)
		case 3:
			*token = promptWithDefault("token", *token)
		case 4:
			if len(*token) == 0 {
				fmt.Println("<empty>")
			} else if len(*token) < 16 {
				fmt.Println(*token)
			} else {
				fmt.Println((*token)[0:16] + "...")
			}
		default:
			fmt.Println("unknown choice")
		}
	}
}

func ExtractJSONField(b []byte, key string) string {
	var m map[string]any
	if err := jsonUnmarshal(b, &m); err != nil {
		return ""
	}
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
