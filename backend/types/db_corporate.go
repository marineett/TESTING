package types

type DBDepartment struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	HeadID int64  `json:"head_id"`
}

type DBHireInfo struct {
	ID           int64 `json:"id"`
	UserID       int64 `json:"user_id"`
	DepartmentID int64 `json:"department_id"`
}
