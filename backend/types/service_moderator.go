package types

type ServiceInitModeratorData struct {
	ServiceInitUserData
	Salary int
}

type ServiceModeratorProfile struct {
	FirstName       string   `json:"first_name"`
	LastName        string   `json:"last_name"`
	MiddleName      string   `json:"middle_name"`
	TelephoneNumber string   `json:"telephone_number"`
	Email           string   `json:"email"`
	Salary          int64    `json:"salary"`
	Departments     []string `json:"departments"`
}

type ServiceModeratorData struct {
	ID          int64    `json:"id"`
	Salary      int64    `json:"salary"`
	Departments []string `json:"departments"`
}

type ServiceModeratorProfileWithID struct {
	ID int64 `json:"id"`
	ServiceModeratorProfile
}
