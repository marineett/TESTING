package types

import "time"

type ServiceAuthData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type ServicePassportData struct {
	PassportNumber   string    `json:"passport_number"`
	PassportSeries   string    `json:"passport_series"`
	PassportDate     time.Time `json:"passport_date"`
	PassportIssuedBy string    `json:"passport_issued_by"`
}

type ServiceUserData struct {
	ID               int64     `json:"id"`
	RegistrationDate time.Time `json:"registration_date"`
	LastLoginDate    time.Time `json:"last_login_date"`
	PersonalDataID   int64     `json:"personal_data_id"`
}

type ServicePersonalData struct {
	TelephoneNumber string `json:"telephone_number"`
	Email           string `json:"email"`
	ServicePassportData
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
}

type ServiceInitUserData struct {
	ServicePersonalData
	ServiceAuthData
}

type ServiceAuthVerdict struct {
	UserID   int64    `json:"user_id"`
	UserType UserType `json:"user_type"`
}
