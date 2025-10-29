package types

type ServerInitModeratorData struct {
	ServerInitUserData
}

type ServerModeratorProfile struct {
	FirstName       string   `json:"first_name"`
	LastName        string   `json:"last_name"`
	MiddleName      string   `json:"middle_name"`
	TelephoneNumber string   `json:"telephone_number"`
	Email           string   `json:"email"`
	Salary          int64    `json:"salary"`
	Departments     []string `json:"departments"`
}

type ServerModeratorProfileWithID struct {
	ID        int64                  `json:"id"`
	Moderator ServerModeratorProfile `json:"moderator"`
}

type ServerModeratorProfileWithIDV2 struct {
	ID        int64                  `json:"id"`
	Moderator ServerModeratorProfile `json:"moderator"`
}

type ServerDepartmentV2 struct {
	ID         int64                            `json:"id"`
	Name       string                           `json:"name"`
	HeadID     int64                            `json:"head_id"`
	Moderators []ServerModeratorProfileWithIDV2 `json:"moderators"`
}

type ServerDepartmentNameUpdateV2 struct {
	Name string `json:"name"`
}

type ServerDepartmentCreateV2 struct {
	Name   string `json:"name"`
	HeadID int64  `json:"head_id"`
}

type ModeratorSalaryPatch struct {
	Salary int64 `json:"salary"`
}
