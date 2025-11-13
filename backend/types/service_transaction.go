package types

import "time"

type ServiceTransaction struct {
	ID         int64             `json:"id"`
	UserID     int64             `json:"user_id"`
	ContractID int64             `json:"contract_id"`
	Amount     int64             `json:"amount"`
	Status     TransactionStatus `json:"status"`
	Type       TransactionType   `json:"type"`
	CreatedAt  time.Time         `json:"created_at"`
}

type ServicePendingContractPaymentTransaction struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	ContractID    int64     `json:"contract_id"`
	Amount        int64     `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
	TransactionID int64     `json:"transaction_id"`
}
