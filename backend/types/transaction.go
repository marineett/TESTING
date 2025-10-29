package types

type TransactionStatus int

const (
	TransactionStatusNull TransactionStatus = iota
	TransactionStatusPending
	TransactionStatusPaid
	TransactionStatusRefunded
	TransactionStatusFailed
)

func (s TransactionStatus) String() string {
	switch s {
	case TransactionStatusPending:
		return "pending"
	case TransactionStatusPaid:
		return "confirmed"
	case TransactionStatusRefunded:
		return "rejected"
	case TransactionStatusFailed:
		return "rejected"
	default:
		return "pending"
	}
}

type TransactionType int

const (
	TransactionTypeNull TransactionType = iota
	TransactionTypeContractPayment
	TransactionTypeLessonPayment
)
