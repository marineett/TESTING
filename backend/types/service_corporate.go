package types

type ServiceDepartment struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	HeadID int64  `json:"head_id"`
}

type ServiceDepartmentInitData struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	HeadID int64  `json:"head_id"`
}
