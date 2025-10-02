package types

import "time"

type DBRepetitorData struct {
	ID            int64   `json:"id"`
	SummaryRating float64 `json:"summary_rating"`
	ReviewsCount  int64   `json:"reviews_count"`
	ResumeID      int64   `json:"resume_id"`
}

type DBResume struct {
	ID          int64          `json:"id"`
	RepetitorID int64          `json:"repetitor_id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Prices      map[string]int `json:"price"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
