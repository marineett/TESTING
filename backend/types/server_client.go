package types

type ServerInitClientData struct {
	ServerInitUserData
}

type ServerClientProfile struct {
	FirstName       string         `json:"first_name"`
	LastName        string         `json:"last_name"`
	MiddleName      string         `json:"middle_name"`
	TelephoneNumber string         `json:"telephone_number"`
	Email           string         `json:"email"`
	MeanRating      float64        `json:"mean_rating"`
	Reviews         []ServerReview `json:"reviews"`
}

type ServerClientProfileV2 struct {
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	MiddleName      string  `json:"middle_name"`
	TelephoneNumber string  `json:"telephone_number"`
	Email           string  `json:"email"`
	Raiting         float64 `json:"raiting"`
}

// V2 wrappers with IDs when required by API
type ServerClientProfileWithIDV2 struct {
	ID     int64                 `json:"id"`
	Client ServerClientProfileV2 `json:"client"`
}
