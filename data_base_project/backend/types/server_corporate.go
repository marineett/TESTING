package types

type ServerDepartmentInitData struct {
	Name   string `json:"name"`
	HeadID int64  `json:"head_id"`
}

type ServerDepartment struct {
	Name       string                         `json:"name"`
	HeadID     int64                          `json:"head_id"`
	Moderators []ServerModeratorProfileWithID `json:"moderators"`
}
