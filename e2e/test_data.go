package integration

var (
	testClientAuthData    map[string]interface{}
	testRepetitorAuthData map[string]interface{}
	testModeratorAuthData map[string]interface{}

	testClientData    map[string]interface{}
	testRepetitorData map[string]interface{}
	testModeratorData map[string]interface{}
)

func init() {
	clLogin, clPass := makeCreds("client")
	rpLogin, rpPass := makeCreds("repetitor")
	mdLogin, mdPass := makeCreds("moderator")

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
