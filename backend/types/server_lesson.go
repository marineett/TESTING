package types

import "time"

type ServerLesson struct {
	ContractID int64     `json:"contract_id"`
	Duration   int64     `json:"duration"`
	CreatedAt  time.Time `json:"created_at"`
}
