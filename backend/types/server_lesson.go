package types

import "time"

type ServerLesson struct {
	ContractID int64     `json:"contract_id"`
	Duration   int64     `json:"duration"`
	CreatedAt  time.Time `json:"created_at"`
}

// V2 Lesson payloads
type ServerLessonV2 struct {
	ID          int64     `json:"id"`
	ContractID  int64     `json:"contract_id"`
	DurationMin int64     `json:"duration_min"`
	Format      string    `json:"format"`
	CreatedAt   time.Time `json:"created_at"`
}

type ServerLessonCreateV2 struct {
	DurationMin int64  `json:"duration_min"`
	Format      string `json:"format"`
}

type ServerLessonPatchV2 struct {
	DurationMin *int64  `json:"duration_min,omitempty"`
	Format      *string `json:"format,omitempty"`
}
