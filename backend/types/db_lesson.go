package types

import "time"

type DBLesson struct {
	ID         int64     `json:"id"`
	ContractID int64     `json:"contract_id"`
	Duration   int64     `json:"duration"`
	Format     string    `json:"format"`
	CreatedAt  time.Time `json:"created_at"`
}
