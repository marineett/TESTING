package service_logic

import (
	"data_base_project/data_base"
	"data_base_project/types"
	"time"
)

type IContractService interface {
	CreateContract(initInfo types.ServiceContractInitData) (int64, error)
	GetContract(contractID int64) (*types.ServiceContract, error)
	UpdateContractStatus(contractID int64, status types.ContractStatus) error
	UpdateContractPaymentStatus(contractID int64, paymentStatus types.PaymentStatus) error
	CreateContractReviewClient(contractID int64, review types.ServiceReview) (int64, error)
	CreateContractReviewRepetitor(contractID int64, review types.ServiceReview) (int64, error)
	UpdateContractReviewClient(contractID int64, review types.ServiceReview) error
	UpdateContractReviewRepetitor(contractID int64, review types.ServiceReview) error
	UpdateContractRepetitorID(contractID int64, repetitorID int64) error
	GetClientContractList(clientID int64, from int64, size int64, status types.ContractStatus) ([]types.ServiceContract, error)
	GetRepetitorContractList(repetitorID int64, from int64, size int64, status types.ContractStatus) ([]types.ServiceContract, error)
	GetContractList(from int64, size int64, status types.ContractStatus) ([]types.ServiceContract, error)
	GetAllContracts(from int64, size int64) ([]types.ServiceContract, error)
	GetContracts(clientID int64, repetitorID int64, from int64, size int64) ([]types.ServerContractV2, error)
}

type ContractService struct {
	contractRepository data_base.IContractRepository
	reviewRepository   data_base.IReviewRepository
}

func CreateContractService(contractRepository data_base.IContractRepository, reviewRepository data_base.IReviewRepository) IContractService {
	return &ContractService{
		contractRepository: contractRepository,
		reviewRepository:   reviewRepository,
	}
}

func (s *ContractService) transormInitContractData(initInfo types.ServiceContractInitData) (*types.ServiceContract, error) {
	return &types.ServiceContract{
		ClientID:          initInfo.ClientID,
		Description:       initInfo.Description,
		Price:             int64(initInfo.Price),
		Commission:        initInfo.Commission,
		Status:            types.ContractStatusPending,
		PaymentStatus:     types.PaymentStatusNull,
		ReviewClientID:    0,
		ReviewRepetitorID: 0,
		StartDate:         initInfo.StartDate,
		EndDate:           initInfo.StartDate.AddDate(0, 0, int(initInfo.Duration)),
	}, nil
}

func (s *ContractService) CreateContract(initInfo types.ServiceContractInitData) (int64, error) {
	contract, err := s.transormInitContractData(initInfo)
	if err != nil {
		return 0, err
	}
	contract.CreatedAt = time.Now()
	contractID, err := s.contractRepository.InsertContract(*types.MapperContractServiceToDB(contract))
	if err != nil {
		return 0, err
	}
	contract.ID = contractID
	return contractID, nil
}

func (s *ContractService) GetContract(contractID int64) (*types.ServiceContract, error) {
	dbContract, err := s.contractRepository.GetContract(contractID)
	if err != nil {
		return nil, err
	}
	return types.MapperContractDBToService(dbContract), nil
}

func (s *ContractService) UpdateContractStatus(contractID int64, status types.ContractStatus) error {
	return s.contractRepository.UpdateContractStatus(contractID, status)
}

func (s *ContractService) UpdateContractPaymentStatus(contractID int64, paymentStatus types.PaymentStatus) error {
	return s.contractRepository.UpdateContractPaymentStatus(contractID, paymentStatus)
}

func (s *ContractService) CreateContractReviewClient(contractID int64, review types.ServiceReview) (int64, error) {
	tx, err := s.contractRepository.BeginTx()
	if err != nil {
		return 0, err
	}
	if tx != nil {
		defer tx.Rollback()
	}
	reviewID, err := s.reviewRepository.InsertReviewInSeq(tx, *types.MapperReviewServiceToDB(&review))
	if err != nil {
		return 0, err
	}
	err = s.contractRepository.UpdateContractReviewClientIDInSeq(tx, contractID, reviewID)
	if err != nil {
		return 0, err
	}
	if tx != nil {
		err = tx.Commit()
		if err != nil {
			return 0, err
		}
	}
	return reviewID, nil
}

func (s *ContractService) CreateContractReviewRepetitor(contractID int64, review types.ServiceReview) (int64, error) {
	tx, err := s.contractRepository.BeginTx()
	if err != nil {
		return 0, err
	}
	if tx != nil {
		defer tx.Rollback()
	}
	reviewID, err := s.reviewRepository.InsertReviewInSeq(tx, *types.MapperReviewServiceToDB(&review))
	if err != nil {
		return 0, err
	}
	err = s.contractRepository.UpdateContractReviewRepetitorID(contractID, reviewID)
	if err != nil {
		return 0, err
	}
	if tx != nil {
		err = tx.Commit()
		if err != nil {
			return 0, err
		}
	}
	return reviewID, nil
}

func (s *ContractService) UpdateContractReviewClient(contractID int64, review types.ServiceReview) error {
	return s.reviewRepository.UpdateReview(*types.MapperReviewServiceToDB(&review))
}

func (s *ContractService) UpdateContractReviewRepetitor(contractID int64, review types.ServiceReview) error {
	return s.reviewRepository.UpdateReview(*types.MapperReviewServiceToDB(&review))
}

func (s *ContractService) GetRepetitorContractList(repetitorID int64, from int64, size int64, status types.ContractStatus) ([]types.ServiceContract, error) {
	dbContracts, err := s.contractRepository.GetContractsByRepetitorID(repetitorID, from, size, status)
	if err != nil {
		return nil, err
	}
	serviceContracts := make([]types.ServiceContract, len(dbContracts))
	for i, dbContract := range dbContracts {
		serviceContracts[i] = *types.MapperContractDBToService(&dbContract)
	}
	return serviceContracts, nil
}

func (s *ContractService) GetClientContractList(clientID int64, from int64, size int64, status types.ContractStatus) ([]types.ServiceContract, error) {
	dbContracts, err := s.contractRepository.GetContractsByClientID(clientID, from, size, status)
	if err != nil {
		return nil, err
	}
	serviceContracts := make([]types.ServiceContract, len(dbContracts))
	for i, dbContract := range dbContracts {
		serviceContracts[i] = *types.MapperContractDBToService(&dbContract)
	}
	return serviceContracts, nil
}

func (s *ContractService) UpdateContractRepetitorID(contractID int64, repetitorID int64) error {
	return s.contractRepository.UpdateContractRepetitorID(contractID, repetitorID)
}

func (s *ContractService) GetContractList(from int64, size int64, status types.ContractStatus) ([]types.ServiceContract, error) {
	dbContracts, err := s.contractRepository.GetContractList(from, size, status)
	if err != nil {
		return nil, err
	}
	serviceContracts := make([]types.ServiceContract, len(dbContracts))
	for i, dbContract := range dbContracts {
		serviceContracts[i] = *types.MapperContractDBToService(&dbContract)
	}
	return serviceContracts, nil
}

func (s *ContractService) GetAllContracts(from int64, size int64) ([]types.ServiceContract, error) {
	dbContracts, err := s.contractRepository.GetAllContracts(from, size)
	if err != nil {
		return nil, err
	}
	serviceContracts := make([]types.ServiceContract, len(dbContracts))
	for i, dbContract := range dbContracts {
		serviceContracts[i] = *types.MapperContractDBToService(&dbContract)
	}
	return serviceContracts, nil
}

func (s *ContractService) GetContracts(clientID int64, repetitorID int64, from int64, size int64) ([]types.ServerContractV2, error) {
	dbContracts, err := s.contractRepository.GetContracts(clientID, repetitorID, from, size)
	if err != nil {
		return nil, err
	}
	serverContracts := make([]types.ServerContractV2, len(dbContracts))
	for i, dbContract := range dbContracts {
		conv := types.MapperContractDBToServerV2(&dbContract)
		if conv != nil {
			serverContracts[i] = *conv
		}
	}
	return serverContracts, nil
}
