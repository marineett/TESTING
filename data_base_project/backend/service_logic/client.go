package service_logic

import (
	"data_base_project/data_base"
	"data_base_project/types"
)

type IClientService interface {
	CreateClient(initData types.ServiceInitClientData) error
	//GetClientData(UserID int64) (*types.ServiceClientData, error)
	GetClientProfile(UserID int64, reviewsOffset int64, reviewsLimit int64) (*types.ServiceClientProfile, error)
	UpdateClientPersonalData(UserID int64, personalData types.ServicePersonalData) error
	UpdateClientPassword(UserID int64, authData types.ServiceAuthData, newPassword string) error
}

type ClientService struct {
	clientRepository       data_base.IClientRepository
	personalDataRepository data_base.IPersonalDataRepository
	userRepository         data_base.IUserRepository
	reviewRepository       data_base.IReviewRepository
}

func CreateClientService(
	clientRepository data_base.IClientRepository,
	personalDataRepository data_base.IPersonalDataRepository,
	userRepository data_base.IUserRepository,
	reviewRepository data_base.IReviewRepository,
) IClientService {
	return &ClientService{
		clientRepository:       clientRepository,
		personalDataRepository: personalDataRepository,
		userRepository:         userRepository,
		reviewRepository:       reviewRepository,
	}
}

func (s *ClientService) transormInitClientData(initData types.ServiceInitClientData) (*types.DBClientData, *types.DBPersonalData, *types.DBAuthData) {
	client := types.DBClientData{
		SummaryRating: 0,
		ReviewsCount:  0,
	}
	return &client, &types.DBPersonalData{
			FirstName:       initData.FirstName,
			LastName:        initData.LastName,
			MiddleName:      initData.MiddleName,
			TelephoneNumber: initData.TelephoneNumber,
			Email:           initData.Email,
		}, &types.DBAuthData{
			Login:    initData.Login,
			Password: initData.Password,
		}
}

func (s *ClientService) CreateClient(initData types.ServiceInitClientData) error {
	clientData, personalData, authData := s.transormInitClientData(initData)
	_, err := s.clientRepository.InsertClient(*clientData, *personalData, *authData)
	if err != nil {
		return err
	}
	return nil
}

func (s *ClientService) GetClientData(userID int64) (*types.DBClientData, error) {
	return s.clientRepository.GetClient(userID)
}

func (s *ClientService) UpdateClientPersonalData(userID int64, personalData types.ServicePersonalData) error {
	return s.clientRepository.UpdateClientPersonalData(userID, *types.MapperPersonalDataServiceToDB(&personalData))
}

func (s *ClientService) UpdateClientPassword(userID int64, authData types.ServiceAuthData, newPassword string) error {
	return s.clientRepository.UpdateClientPassword(userID, *types.MapperAuthDataServiceToDB(&authData), newPassword)
}

func (s *ClientService) GetClientProfile(userID int64, reviewsOffset int64, reviewsLimit int64) (*types.ServiceClientProfile, error) {
	user, err := s.userRepository.GetUser(userID)
	if err != nil {
		return nil, err
	}
	client, err := s.GetClientData(userID)
	if err != nil {
		return nil, err
	}
	personalData, err := s.personalDataRepository.GetPersonalData(user.PersonalDataID)
	if err != nil {
		return nil, err
	}
	reviews, err := s.reviewRepository.GetReviewsByClientID(userID, reviewsOffset, reviewsLimit)
	if err != nil {
		return nil, err
	}
	meanRating := 0.0
	if client.ReviewsCount > 0 {
		meanRating = float64(client.SummaryRating) / float64(client.ReviewsCount)
	}
	serviceReviews := make([]types.ServiceReview, len(reviews))
	for i, review := range reviews {
		serviceReviews[i] = *types.MapperReviewDBToService(&review)
	}
	return &types.ServiceClientProfile{
		FirstName:       personalData.FirstName,
		LastName:        personalData.LastName,
		MiddleName:      personalData.MiddleName,
		TelephoneNumber: personalData.TelephoneNumber,
		Email:           personalData.Email,
		MeanRating:      meanRating,
		Reviews:         serviceReviews,
	}, nil
}
