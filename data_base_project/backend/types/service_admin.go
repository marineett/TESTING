package types

type ServiceInitAdminData struct {
	ServiceInitUserData
	Salary int64 `json:"salary"`
}

type ServiceAdminProfile struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	MiddleName      string `json:"middle_name"`
	TelephoneNumber string `json:"telephone_number"`
	Email           string `json:"email"`
	Salary          int64  `json:"salary"`
}

type ServiceAdminData struct {
	ID           int64 `json:"id"`
	Salary       int64 `json:"salary"`
	DepartmentID int64 `json:"department_id"`
}
