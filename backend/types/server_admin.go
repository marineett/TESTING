package types

type ServerInitAdminData struct {
	ServerInitUserData
	Salary int64 `json:"salary"`
}

type ServerAdminProfile struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	MiddleName      string `json:"middle_name"`
	TelephoneNumber string `json:"telephone_number"`
	Email           string `json:"email"`
	Salary          int64  `json:"salary"`
}

// V2 Auth generic response
type ServerAuthResponseV2 struct {
	Token  string `json:"token"`
	Role   string `json:"role"`
	UserID int64  `json:"user_id"`
}
