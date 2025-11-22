package types

import "time"

type UserType int

const (
	Unauthorized UserType = iota
	Client
	Repetitor
	Moderator
	Admin
)

func (u UserType) String() string {
	switch u {
	case Client:
		return "client"
	case Repetitor:
		return "repetitor"
	case Moderator:
		return "moderator"
	case Admin:
		return "admin"
	default:
		return "guest"
	}
}

type DBUserData struct {
	ID               int64     `json:"id"`
	RegistrationDate time.Time `json:"registration_date"`
	LastLoginDate    time.Time `json:"last_login_date"`
	PersonalDataID   int64     `json:"personal_data_id"`
}

type DBAuthData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type DBAuthVerdict struct {
	UserID   int64    `json:"user_id"`
	UserType UserType `json:"user_type"`
}

type DBAuthInfo struct {
	ID       int64    `json:"id"`
	UserID   int64    `json:"user_id"`
	UserType UserType `json:"user_type"`
	Login    string   `json:"login"`
	Password string   `json:"password"`
}

type DBPassportData struct {
	PassportNumber   string    `json:"passport_number"`
	PassportSeries   string    `json:"passport_series"`
	PassportDate     time.Time `json:"passport_date"`
	PassportIssuedBy string    `json:"passport_issued_by"`
}

type DBPersonalData struct {
	ID              int64  `json:"id"`
	TelephoneNumber string `json:"telephone_number"`
	Email           string `json:"email"`
	DBPassportData
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
}
