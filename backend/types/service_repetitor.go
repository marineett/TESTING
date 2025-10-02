package types

import "time"

type ServiceInitRepetitorData struct {
	ServiceInitUserData
}

type ServiceRepetitorProfile struct {
	FirstName         string          `json:"first_name"`
	LastName          string          `json:"last_name"`
	MiddleName        string          `json:"middle_name"`
	TelephoneNumber   string          `json:"telephone_number"`
	Email             string          `json:"email"`
	MeanRating        float64         `json:"mean_rating"`
	ResumeTitle       string          `json:"resume_title"`
	ResumeDescription string          `json:"resume_description"`
	ResumePrices      map[string]int  `json:"resume_prices"`
	Reviews           []ServiceReview `json:"reviews"`
}

type ServiceRepetitorData struct {
	ID         int64   `json:"id"`
	MeanRating float64 `json:"mean_rating"`
	ResumeID   int64   `json:"resume_id"`
}

type ServiceResume struct {
	RepetitorID int64          `json:"repetitor_id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Prices      map[string]int `json:"price"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type ServiceRepetitorView struct {
	FirstName  string  `json:"first_name"`
	MeanRating float64 `json:"mean_rating"`
}
