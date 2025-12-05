package types

import "time"

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
		ContractID:  review.ContractID,
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
		ContractID:  review.ContractID,
		ClientID:    review.ClientID,
		RepetitorID: review.RepetitorID,
		Rating:      review.Rating,
		Comment:     review.Comment,
		CreatedAt:   review.CreatedAt,
	}
}

func MapperContractServiceToServerV2(contract *ServiceContract) *ServerContractV2 {
	if contract == nil {
		return nil
	}
	var repID *int64
	if contract.RepetitorID != 0 {
		repID = &contract.RepetitorID
	}
	return &ServerContractV2{
		ID:          contract.ID,
		ClientID:    contract.ClientID,
		RepetitorID: repID,
		Description: contract.Description,
		Rate:        contract.Price,
		Format:      "online",
		Status:      contract.Status.String(),
		CreatedAt:   contract.CreatedAt,
	}
}

func MapperContractDBToServerV2(contract *DBContract) *ServerContractV2 {
	if contract == nil {
		return nil
	}
	var repID *int64
	if contract.RepetitorID != 0 {
		repID = &contract.RepetitorID
	}
	return &ServerContractV2{
		ID:          contract.ID,
		ClientID:    contract.ClientID,
		RepetitorID: repID,
		Description: contract.Description,
		Rate:        contract.Price,
		Format:      "online",
		Status:      contract.Status.String(),
		CreatedAt:   contract.CreatedAt,
	}
}

func MapperContractCreateV2ServerToService(req *ServerContractCreateV2) *ServiceContract {
	if req == nil {
		return nil
	}
	return &ServiceContract{
		ClientID:    req.ClientID,
		Description: req.Description,
		Price:       req.Rate,
		// default fields left zero; will be filled by service
	}
}

func MapperContractCreateV2ServerToServiceInit(req *ServerContractCreateV2) *ServiceContractInitData {
	if req == nil {
		return nil
	}
	return &ServiceContractInitData{
		ClientID:    req.ClientID,
		Description: req.Description,
		Price:       req.Rate,
		Commission:  0,
		StartDate:   time.Now(),
		Duration:    0,
	}
}

func MapperReviewServiceToServerV2(review *ServiceReview) *ServerReviewV2 {
	if review == nil {
		return nil
	}
	return &ServerReviewV2{
		ID:         review.ID,
		ContractID: review.ContractID,
		FromUserID: review.ClientID,
		ToUserID:   review.RepetitorID,
		Score:      review.Rating,
		Text:       review.Comment,
		CreatedAt:  review.CreatedAt,
	}
}

func MapperReviewCreateV2ServerToService(req *ServerReviewCreateV2, contractID int64, fromUserID int64, toUserID int64) *ServiceReview {
	if req == nil {
		return nil
	}
	return &ServiceReview{
		ContractID:  contractID,
		ClientID:    fromUserID,
		RepetitorID: toUserID,
		Rating:      req.Score,
		Comment:     req.Text,
		CreatedAt:   time.Now(),
	}
}
