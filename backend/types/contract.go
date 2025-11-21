package types

type ContractStatus int

const (
	ContractStatusNull ContractStatus = iota
	ContractStatusPending
	ContractStatusActive
	ContractStatusCompleted
	ContractStatusCancelled
	ContractStatusBanned
)

func (s ContractStatus) String() string {
	switch s {
	case ContractStatusPending:
		return "review"
	case ContractStatusActive:
		return "active"
	case ContractStatusCompleted:
		return "completed"
	case ContractStatusCancelled:
		return "cancelled"
	case ContractStatusBanned:
		return "banned"
	default:
		return "review"
	}
}

// ParseContractStatus converts string status (as in OpenAPI) to ContractStatus enum
func ParseContractStatus(status string) (ContractStatus, error) {
	switch status {
	case "review":
		return ContractStatusPending, nil
	case "active":
		return ContractStatusActive, nil
	case "completed":
		return ContractStatusCompleted, nil
	case "cancelled":
		return ContractStatusCancelled, nil
	case "banned":
		return ContractStatusBanned, nil
	default:
		return ContractStatusNull, nil
	}
}

type PaymentStatus int

const (
	PaymentStatusNull PaymentStatus = iota
	PaymentStatusPending
	PaymentStatusPaid
	PaymentStatusRefused
	PaymentStatusRefunded
)

type ContractCategory int

const (
	ContractCategoryNull ContractCategory = iota
	ContractCategoryTranslation
	ContractCategoryWriting
	ContractCategoryDesign
	ContractCategoryProgramming
	ContractCategoryOther
)

type ContractSubcategory int

const (
	ContractSubcategoryNull ContractSubcategory = iota
	ContractSubcategoryTutoring
	ContractSubcategoryTranslation
	ContractSubcategoryWriting
	ContractSubcategoryDesign
	ContractSubcategoryProgramming
	ContractSubcategoryOther
)

type ContractUpdateInfo struct {
	Status            ContractStatus `json:"status"`
	PriceChange       int64          `json:"price_change"`
	DescriptionChange string         `json:"description_change"`
}
