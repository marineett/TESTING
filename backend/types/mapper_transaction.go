package types

func MapperTransactionDBToService(transaction *DBTransaction) *ServiceTransaction {
	if transaction == nil {
		return nil
	}
	return &ServiceTransaction{
		ID:        transaction.ID,
		UserID:    transaction.UserID,
		Amount:    transaction.Amount,
		Status:    transaction.Status,
		Type:      transaction.Type,
		CreatedAt: transaction.CreatedAt,
	}
}

func MapperTransactionServiceToDB(transaction *ServiceTransaction) *DBTransaction {
	if transaction == nil {
		return nil
	}
	return &DBTransaction{
		ID:        transaction.ID,
		UserID:    transaction.UserID,
		Amount:    transaction.Amount,
		Status:    transaction.Status,
		Type:      transaction.Type,
		CreatedAt: transaction.CreatedAt,
	}
}

func MapperPendingContractPaymentTransactionDBToService(transaction *DBPendingContractPaymentTransaction) *ServicePendingContractPaymentTransaction {
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

func MapperPendingContractPaymentTransactionServiceToDB(transaction *ServicePendingContractPaymentTransaction) *DBPendingContractPaymentTransaction {
	if transaction == nil {
		return nil
	}
	return &DBPendingContractPaymentTransaction{
		ID:            transaction.ID,
		UserID:        transaction.UserID,
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
		Amount:    transaction.Amount,
		Status:    transaction.Status,
		Type:      transaction.Type,
		CreatedAt: transaction.CreatedAt,
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
