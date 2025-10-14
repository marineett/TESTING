package types

import "time"

type ServerInitRepetitorData struct {
	ServerInitUserData
}

type ServerRepetitorProfile struct {
	FirstName         string         `json:"first_name"`
	LastName          string         `json:"last_name"`
	MiddleName        string         `json:"middle_name"`
	TelephoneNumber   string         `json:"telephone_number"`
	Email             string         `json:"email"`
	MeanRating        float64        `json:"mean_rating"`
	ResumeTitle       string         `json:"resume_title"`
	ResumeDescription string         `json:"resume_description"`
	ResumePrices      map[string]int `json:"resume_prices"`
	Reviews           []ServerReview `json:"reviews"`
}

type ServerResume struct {
	RepetitorID int64          `json:"repetitor_id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Prices      map[string]int `json:"price"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type ServerRepetitorView struct {
	FirstName  string  `json:"first_name"`
	MeanRating float64 `json:"mean_rating"`
}
