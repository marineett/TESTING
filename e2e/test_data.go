package integration

import "os"

var (
	testClientAuthData    map[string]interface{}
	testRepetitorAuthData map[string]interface{}
	testModeratorAuthData map[string]interface{}

	testClientData    map[string]interface{}
	testRepetitorData map[string]interface{}
	testModeratorData map[string]interface{}
)

func init() {
	// Используем IMAP_USER из env, если установлен, иначе генерируем случайный email
	imapUser := os.Getenv("IMAP_USER")
	
	var clLogin, rpLogin, mdLogin string
	if imapUser != "" {
		// Используем один и тот же email для всех регистраций (для чтения токенов из одного почтового ящика)
		clLogin = imapUser
		rpLogin = imapUser
		mdLogin = imapUser
	} else {
		// Генерируем случайные email адреса
		clLogin, _ = makeCreds("client")
		rpLogin, _ = makeCreds("repetitor")
		mdLogin, _ = makeCreds("moderator")
	}
	
	clPass := createPassword(12)
	rpPass := createPassword(12)
	mdPass := createPassword(12)

	testClientAuthData = map[string]interface{}{
		"login":    clLogin,
		"password": clPass,
	}
	testRepetitorAuthData = map[string]interface{}{
		"login":    rpLogin,
		"password": rpPass,
	}
	testModeratorAuthData = map[string]interface{}{
		"login":    mdLogin,
		"password": mdPass,
	}

	testClientData = map[string]interface{}{
		"login":            testClientAuthData["login"],
		"password":         testClientAuthData["password"],
		"first_name":       "Иван",
		"last_name":        "Иванов",
		"middle_name":      "Иванович",
		"email":            clLogin,
		"telephone_number": "+7-900-123-45-67",
		"role":             "client",
	}

	testRepetitorData = map[string]interface{}{
		"login":            testRepetitorAuthData["login"],
		"password":         testRepetitorAuthData["password"],
		"first_name":       "Иван",
		"last_name":        "Иванов",
		"middle_name":      "Иванович",
		"email":            rpLogin,
		"telephone_number": "+7-900-123-45-67",
		"role":             "repetitor",
	}

	testModeratorData = map[string]interface{}{
		"login":            testModeratorAuthData["login"],
		"password":         testModeratorAuthData["password"],
		"first_name":       "Иван",
		"last_name":        "Иванов",
		"middle_name":      "Иванович",
		"email":            mdLogin,
		"telephone_number": "+7-900-123-45-67",
		"role":             "moderator",
	}
}
