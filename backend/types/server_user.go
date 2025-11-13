package types

import "time"

type ServerAuthData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type ServerPassportData struct {
	PassportNumber   string    `json:"passport_number"`
	PassportSeries   string    `json:"passport_series"`
	PassportDate     time.Time `json:"passport_date"`
	PassportIssuedBy string    `json:"passport_issued_by"`
}

type ServerUserData struct {
	ID               int64     `json:"id"`
	RegistrationDate time.Time `json:"registration_date"`
	LastLoginDate    time.Time `json:"last_login_date"`
	PersonalDataID   int64     `json:"personal_data_id"`
}

type ServerPersonalData struct {
	TelephoneNumber string `json:"telephone_number"`
	Email           string `json:"email"`
	ServerPassportData
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
}

type ServerInitUserData struct {
	ServerPersonalData
	ServerAuthData
}

type ServerRegistrationData struct {
	Username         string    `json:"username"`
	Password         string    `json:"password"`
	UserType         UserType  `json:"user_type"`
	RegistrationDate time.Time `json:"registration_date"`
}

type ServerVerdict struct {
	UserID   int64    `json:"user_id"`
	UserType UserType `json:"user_type"`
}

type ServerRegistrationDataV2 struct {
	Login           string `json:"login"`
	Password        string `json:"password"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	MiddleName      string `json:"middle_name"`
	Email           string `json:"email"`
	TelephoneNumber string `json:"telephone_number"`
	Role            string `json:"role"`
	Salary          int    `json:"salary"`
}
