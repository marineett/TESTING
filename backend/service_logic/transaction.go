package service_logic

import (
	"data_base_project/data_base"
	"data_base_project/types"
	"time"
)

type ITransactionService interface {
	GetTransactionsList(user_id int64, from int64, size int64) ([]types.ServiceTransaction, error)
	GetContractTransactionsList(contract_id int64, from int64, size int64) ([]types.ServiceTransaction, error)
	CreateContractPaymentTransaction(amount int64, user_id int64, contract_id int64) (int64, error)
	ChangeTransactionStatus(transaction_id int64, status types.TransactionStatus) error
	GetTransaction(transaction_id int64) (*types.ServiceTransaction, error)
	GetPendingContractPaymentTransaction() (*types.ServicePendingContractPaymentTransaction, error)
	ApproveTransaction(transaction_id int64) error
}

type TransactionService struct {
	transactionRepository data_base.ITransactionRepository
}

func CreateTransactionService(transactionRepository data_base.ITransactionRepository) *TransactionService {
	return &TransactionService{transactionRepository: transactionRepository}
}

func (s *TransactionService) CreateContractPaymentTransaction(amount int64, user_id int64, contract_id int64) (int64, error) {
	return s.transactionRepository.InsertTransaction(types.DBTransaction{
		UserID:     user_id,
		ContractID: contract_id,
		Amount:     amount,
		Status:     types.TransactionStatusPending,
		Type:       types.TransactionTypeContractPayment,
		CreatedAt:  time.Now(),
	})
}

func (s *TransactionService) GetTransactionsList(user_id int64, from int64, size int64) ([]types.ServiceTransaction, error) {
	transactions, err := s.transactionRepository.GetTransactionsList(user_id, from, size)
	if err != nil {
		return nil, err
	}
	serviceTransactions := make([]types.ServiceTransaction, len(transactions))
	for i, transaction := range transactions {
		serviceTransactions[i] = types.ServiceTransaction(transaction)
	}
	return serviceTransactions, nil
}

func (s *TransactionService) ChangeTransactionStatus(transaction_id int64, status types.TransactionStatus) error {
	return s.transactionRepository.UpdateTransactionStatus(transaction_id, status)
}

func (s *TransactionService) GetTransaction(transaction_id int64) (*types.ServiceTransaction, error) {
	transaction, err := s.transactionRepository.GetTransaction(transaction_id)
	if err != nil {
		return nil, err
	}
	return types.MapperTransactionDBToService(transaction), nil
}

func (s *TransactionService) GetPendingContractPaymentTransaction() (*types.ServicePendingContractPaymentTransaction, error) {
	transaction, err := s.transactionRepository.GetPendingContractPaymentTransaction()
	if err != nil {
		return nil, err
	}
	return types.MapperPendingContractPaymentTransactionDBToService(transaction), nil
}

func (s *TransactionService) ApproveTransaction(transaction_id int64) error {
	return s.transactionRepository.ApproveTransaction(transaction_id)
}

func (s *TransactionService) GetContractTransactionsList(contract_id int64, from int64, size int64) ([]types.ServiceTransaction, error) {
	transactions, err := s.transactionRepository.GetContractTransactionsList(contract_id, from, size)
	if err != nil {
		return nil, err
	}
	serviceTransactions := make([]types.ServiceTransaction, len(transactions))
	for i, transaction := range transactions {
		serviceTransactions[i] = types.ServiceTransaction(transaction)
	}
	return serviceTransactions, nil
}
