package types

type ServiceInitClientData struct {
	ServiceInitUserData
}

type ServiceClientProfile struct {
	FirstName       string          `json:"first_name"`
	LastName        string          `json:"last_name"`
	MiddleName      string          `json:"middle_name"`
	TelephoneNumber string          `json:"telephone_number"`
	Email           string          `json:"email"`
	MeanRating      float64         `json:"mean_rating"`
	Reviews         []ServiceReview `json:"reviews"`
}
