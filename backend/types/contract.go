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
