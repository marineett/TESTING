package types

import "time"

type ServerTransaction struct {
	ContractID int64             `json:"contract_id"`
	Amount     int64             `json:"amount"`
	Status     TransactionStatus `json:"status"`
	Type       TransactionType   `json:"type"`
	CreatedAt  time.Time         `json:"created_at"`
}

type ServerPendingContractPaymentTransaction struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	ContractID    int64     `json:"contract_id"`
	Amount        int64     `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
	TransactionID int64     `json:"transaction_id"`
}

// V2 Transaction payloads
type ServerTransactionV2 struct {
	ID         int64     `json:"id"`
	ContractID int64     `json:"contract_id"`
	Amount     int64     `json:"amount"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type ServerTransactionCreateV2 struct {
	Amount int64 `json:"amount"`
}

type ServerTransactionApproveV2 struct {
	Approved bool `json:"approved"`
}
