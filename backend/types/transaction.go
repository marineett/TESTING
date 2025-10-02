package types

type TransactionStatus int

const (
	TransactionStatusNull TransactionStatus = iota
	TransactionStatusPending
	TransactionStatusPaid
	TransactionStatusRefunded
	TransactionStatusFailed
)

type TransactionType int

const (
	TransactionTypeNull TransactionType = iota
	TransactionTypeContractPayment
	TransactionTypeLessonPayment
)
