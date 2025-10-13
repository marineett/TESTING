package integration

var (
	testClientAuthData = map[string]interface{}{
		"login":    "client@test.com",
		"password": "password123",
	}

	testRepetitorAuthData = map[string]interface{}{
		"login":    "repetitor@test.com",
		"password": "password123",
	}

	testModeratorAuthData = map[string]interface{}{
		"login":    "moderator@test.com",
		"password": "password123",
	}

	testClientData = map[string]interface{}{
		"login":        testClientAuthData["login"],
		"password":     testClientAuthData["password"],
		"first_name":   "Иван",
		"last_name":    "Иванов",
		"middle_name":  "Иванович",
		"passport":     "1234 567890",
		"phone_number": "+7-900-123-45-67",
	}

	testRepetitorData = map[string]interface{}{
		"login":        testRepetitorAuthData["login"],
		"password":     testRepetitorAuthData["password"],
		"first_name":   "Иван",
		"last_name":    "Иванов",
		"middle_name":  "Иванович",
		"passport":     "1234 567890",
		"phone_number": "+7-900-123-45-67",
	}

	testModeratorData = map[string]interface{}{
		"login":        testModeratorAuthData["login"],
		"password":     testModeratorAuthData["password"],
		"first_name":   "Иван",
		"last_name":    "Иванов",
		"middle_name":  "Иванович",
		"passport":     "1234 567890",
		"phone_number": "+7-900-123-45-67",
	}
)
