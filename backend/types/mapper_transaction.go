package types

import "time"

func MapperTransactionDBToService(transaction *DBTransaction) *ServiceTransaction {
	if transaction == nil {
		return nil
	}
	return &ServiceTransaction{
		ID:         transaction.ID,
		UserID:     transaction.UserID,
		ContractID: transaction.ContractID,
		Amount:     transaction.Amount,
		Status:     transaction.Status,
		Type:       transaction.Type,
		CreatedAt:  transaction.CreatedAt,
	}
}

func MapperTransactionServiceToDB(transaction *ServiceTransaction) *DBTransaction {
	if transaction == nil {
		return nil
	}
	return &DBTransaction{
		ID:         transaction.ID,
		UserID:     transaction.UserID,
		ContractID: transaction.ContractID,
		Amount:     transaction.Amount,
		Status:     transaction.Status,
		Type:       transaction.Type,
		CreatedAt:  transaction.CreatedAt,
	}
}

func MapperPendingContractPaymentTransactionDBToService(transaction *DBPendingContractPaymentTransaction) *ServicePendingContractPaymentTransaction {
	if transaction == nil {
		return nil
	}
	return &ServicePendingContractPaymentTransaction{
		ID:            transaction.ID,
		UserID:        transaction.UserID,
		ContractID:    transaction.ContractID,
		Amount:        transaction.Amount,
		CreatedAt:     transaction.CreatedAt,
		TransactionID: transaction.TransactionID,
	}
}

func MapperPendingContractPaymentTransactionServiceToDB(transaction *ServicePendingContractPaymentTransaction) *DBPendingContractPaymentTransaction {
	if transaction == nil {
		return nil
	}
	return &DBPendingContractPaymentTransaction{
		ID:            transaction.ID,
		UserID:        transaction.UserID,
		ContractID:    transaction.ContractID,
		Amount:        transaction.Amount,
		CreatedAt:     transaction.CreatedAt,
		TransactionID: transaction.TransactionID,
	}
}

func MapperTransactionServiceToServer(transaction *ServiceTransaction) *ServerTransaction {
	if transaction == nil {
		return nil
	}
	return &ServerTransaction{
		ContractID: transaction.ContractID,
		Amount:     transaction.Amount,
		Status:     transaction.Status,
		Type:       transaction.Type,
		CreatedAt:  transaction.CreatedAt,
	}
}

func MapperTransactionServiceToServerV2(tx *ServiceTransaction) *ServerTransactionV2 {
	if tx == nil {
		return nil
	}
	return &ServerTransactionV2{
		ID:         tx.ID,
		ContractID: tx.ContractID,
		Amount:     tx.Amount,
		Status:     tx.Status.String(),
		CreatedAt:  tx.CreatedAt,
	}
}

func MapperTransactionCreateV2ServerToService(contractID int64, req *ServerTransactionCreateV2) *ServiceTransaction {
	if req == nil {
		return nil
	}
	return &ServiceTransaction{
		ContractID: contractID,
		Amount:     req.Amount,
		Status:     TransactionStatusPending,
		CreatedAt:  time.Now(),
	}
}

func MapperTransactionServerToService(transaction *ServerTransaction) *ServiceTransaction {
	if transaction == nil {
		return nil
	}
	return &ServiceTransaction{
		Amount:    transaction.Amount,
		Status:    transaction.Status,
		Type:      transaction.Type,
		CreatedAt: transaction.CreatedAt,
	}
}

func MapperPendingContractPaymentTransactionServiceToServer(transaction *ServicePendingContractPaymentTransaction) *ServerPendingContractPaymentTransaction {
	if transaction == nil {
		return nil
	}
	return &ServerPendingContractPaymentTransaction{
		ID:            transaction.ID,
		UserID:        transaction.UserID,
		Amount:        transaction.Amount,
		CreatedAt:     transaction.CreatedAt,
		TransactionID: transaction.TransactionID,
	}
}

func MapperPendingContractPaymentTransactionServerToService(transaction *ServerPendingContractPaymentTransaction) *ServicePendingContractPaymentTransaction {
	if transaction == nil {
		return nil
	}
	return &ServicePendingContractPaymentTransaction{
		ID:            transaction.ID,
		UserID:        transaction.UserID,
		Amount:        transaction.Amount,
		CreatedAt:     transaction.CreatedAt,
		TransactionID: transaction.TransactionID,
	}
}
