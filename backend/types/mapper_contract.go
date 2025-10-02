package types

func MapperContractDBToService(contract *DBContract) *ServiceContract {
	if contract == nil {
		return nil
	}
	return &ServiceContract{
		ID:                contract.ID,
		ClientID:          contract.ClientID,
		RepetitorID:       contract.RepetitorID,
		TransactionID:     contract.TransactionID,
		CreatedAt:         contract.CreatedAt,
		Description:       contract.Description,
		Status:            contract.Status,
		PaymentStatus:     contract.PaymentStatus,
		ReviewClientID:    contract.ReviewClientID,
		ReviewRepetitorID: contract.ReviewRepetitorID,
		Price:             contract.Price,
		Commission:        contract.Commission,
		StartDate:         contract.StartDate,
		EndDate:           contract.EndDate,
		IDCRChat:          contract.IDCRChat,
		IDCMRepChat:       contract.IDCMRepChat,
		IDRMRepChat:       contract.IDRMRepChat,
	}
}

func MapperContractServiceToDB(contract *ServiceContract) *DBContract {
	if contract == nil {
		return nil
	}
	return &DBContract{
		ID:                contract.ID,
		ClientID:          contract.ClientID,
		RepetitorID:       contract.RepetitorID,
		TransactionID:     contract.TransactionID,
		CreatedAt:         contract.CreatedAt,
		Description:       contract.Description,
		Status:            contract.Status,
		PaymentStatus:     contract.PaymentStatus,
		ReviewClientID:    contract.ReviewClientID,
		ReviewRepetitorID: contract.ReviewRepetitorID,
		Price:             contract.Price,
		Commission:        contract.Commission,
		StartDate:         contract.StartDate,
		EndDate:           contract.EndDate,
		IDCRChat:          contract.IDCRChat,
		IDCMRepChat:       contract.IDCMRepChat,
		IDRMRepChat:       contract.IDRMRepChat,
	}
}

func MapperContractServiceToServer(contract *ServiceContract) *ServerContract {
	if contract == nil {
		return nil
	}
	return &ServerContract{
		ID:                contract.ID,
		ClientID:          contract.ClientID,
		RepetitorID:       contract.RepetitorID,
		TransactionID:     contract.TransactionID,
		CreatedAt:         contract.CreatedAt,
		Description:       contract.Description,
		Status:            contract.Status,
		PaymentStatus:     contract.PaymentStatus,
		ReviewClientID:    contract.ReviewClientID,
		ReviewRepetitorID: contract.ReviewRepetitorID,
		Price:             contract.Price,
		Commission:        contract.Commission,
		StartDate:         contract.StartDate,
		EndDate:           contract.EndDate,
		IDCRChat:          contract.IDCRChat,
		IDCMRepChat:       contract.IDCMRepChat,
		IDRMRepChat:       contract.IDRMRepChat,
	}
}

func MapperContractServerToService(contract *ServerContract) *ServiceContract {
	if contract == nil {
		return nil
	}
	return &ServiceContract{
		ID:                contract.ID,
		ClientID:          contract.ClientID,
		RepetitorID:       contract.RepetitorID,
		TransactionID:     contract.TransactionID,
		CreatedAt:         contract.CreatedAt,
		Description:       contract.Description,
		Status:            contract.Status,
		PaymentStatus:     contract.PaymentStatus,
		ReviewClientID:    contract.ReviewClientID,
		ReviewRepetitorID: contract.ReviewRepetitorID,
		Price:             contract.Price,
		Commission:        contract.Commission,
		StartDate:         contract.StartDate,
		EndDate:           contract.EndDate,
		IDCRChat:          contract.IDCRChat,
		IDCMRepChat:       contract.IDCMRepChat,
		IDRMRepChat:       contract.IDRMRepChat,
	}
}

func MapperReviewServiceToServer(review *ServiceReview) *ServerReview {
	if review == nil {
		return nil
	}
	return &ServerReview{
		ClientID:    review.ClientID,
		RepetitorID: review.RepetitorID,
		Rating:      review.Rating,
		Comment:     review.Comment,
		CreatedAt:   review.CreatedAt,
	}
}

func MapperReviewServerToService(review *ServerReview) *ServiceReview {
	if review == nil {
		return nil
	}
	return &ServiceReview{
		ClientID:    review.ClientID,
		RepetitorID: review.RepetitorID,
		Rating:      review.Rating,
		Comment:     review.Comment,
		CreatedAt:   review.CreatedAt,
	}
}
