package types

import "time"

type DBContract struct {
	ID                int64          `json:"id"`
	ClientID          int64          `json:"client_id"`
	RepetitorID       int64          `json:"repetitor_id"`
	TransactionID     int64          `json:"transaction_id"`
	CreatedAt         time.Time      `json:"created_at"`
	Description       string         `json:"description"`
	Status            ContractStatus `json:"status"`
	PaymentStatus     PaymentStatus  `json:"payment_status"`
	ReviewClientID    int64          `json:"review_client_id"`
	ReviewRepetitorID int64          `json:"review_repetitor_id"`
	Price             int64          `json:"price"`
	Commission        int64          `json:"commission"`
	StartDate         time.Time      `json:"start_date"`
	EndDate           time.Time      `json:"end_date"`
	IDCRChat          int64          `json:"id_client_repetitor_chat"`
	IDCMRepChat       int64          `json:"id_client_manager_chat"`
	IDRMRepChat       int64          `json:"id_repetitor_moderator_chat"`
}

type DBReview struct {
	ID          int64     `json:"id"`
	ContractID  int64     `json:"contract_id"`
	ClientID    int64     `json:"client_id"`
	RepetitorID int64     `json:"repetitor_id"`
	Rating      int       `json:"rating"`
	Comment     string    `json:"comment"`
	CreatedAt   time.Time `json:"created_at"`
}
