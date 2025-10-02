package types

type DBClientData struct {
	ID            int64 `json:"id"`
	SummaryRating int64 `json:"summary_rating"`
	ReviewsCount  int64 `json:"reviews_count"`
}
