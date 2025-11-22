package types

import "time"

type ServerContractInitData struct {
	ClientID              int64                 `json:"client_id"`
	ContractCategory      ContractCategory      `json:"contract_category"`
	ContractSubcategories []ContractSubcategory `json:"contract_subcategories"`
	Description           string                `json:"description"`
	Price                 int64                 `json:"price"`
	Commission            int64                 `json:"commission"`
	StartDate             time.Time             `json:"start_date"`
	Duration              int64                 `json:"duration"`
}

type ServerContract struct {
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

type ServerReview struct {
	ContractID  int64     `json:"contract_id"`
	ClientID    int64     `json:"client_id"`
	RepetitorID int64     `json:"repetitor_id"`
	Rating      int       `json:"rating"`
	Comment     string    `json:"comment"`
	CreatedAt   time.Time `json:"created_at"`
}

type ServerContractV2 struct {
	ID          int64     `json:"id"`
	ClientID    int64     `json:"client_id"`
	RepetitorID *int64    `json:"repetitor_id,omitempty"`
	Description string    `json:"description"`
	Rate        int64     `json:"rate"`
	Format      string    `json:"format"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type ServerContractCreateV2 struct {
	ClientID    int64  `json:"client_id"`
	Description string `json:"description"`
	Rate        int64  `json:"rate"`
	Format      string `json:"format"`
}

type ServerContractStatusPatchV2 struct {
	Status string `json:"status"`
}

type ServerReviewV2 struct {
	ID         int64     `json:"id"`
	ContractID int64     `json:"contract_id"`
	FromUserID int64     `json:"from_user_id"`
	ToUserID   int64     `json:"to_user_id"`
	Score      int       `json:"score"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"created_at"`
}

type ServerReviewCreateV2 struct {
	Score    int    `json:"score"`
	Text     string `json:"text"`
	SenderID int64  `json:"senderId"`
}
