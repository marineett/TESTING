package service_logic

import (
	"data_base_project/data_base"
	"data_base_project/types"
)

type IRepetitorService interface {
	CreateRepetitor(initData types.ServiceInitRepetitorData) error
	GetRepetitorData(userID int64) (*types.ServiceRepetitorData, error)
	GetRepetitorProfile(userID int64, reviewsOffset int64, reviewsLimit int64) (*types.ServiceRepetitorProfile, error)
	GetRepetitors(repetitorsOffset int64, repetitorsLimit int64) ([]*types.ServiceRepetitorView, error)
	UpdateRepetitorPersonalData(userID int64, personalData types.ServicePersonalData) error
	UpdateRepetitorPassword(userID int64, authData types.ServiceAuthData, newPassword string) error
}

type RepetitorService struct {
	repetitorRepository    data_base.IRepetitorRepository
	personalDataRepository data_base.IPersonalDataRepository
	userRepository         data_base.IUserRepository
	reviewRepository       data_base.IReviewRepository
	resumeRepository       data_base.IResumeRepository
}

func CreateRepetitorService(
	repetitorRepository data_base.IRepetitorRepository,
	personalDataRepository data_base.IPersonalDataRepository,
	userRepository data_base.IUserRepository,
	reviewRepository data_base.IReviewRepository,
	resumeRepository data_base.IResumeRepository,
) IRepetitorService {
	return &RepetitorService{
		repetitorRepository:    repetitorRepository,
		personalDataRepository: personalDataRepository,
		userRepository:         userRepository,
		reviewRepository:       reviewRepository,
		resumeRepository:       resumeRepository,
	}
}
func (s *RepetitorService) transormInitRepetitorData(initData types.ServiceInitRepetitorData) (*types.DBRepetitorData, *types.DBPersonalData, *types.DBAuthData) {
	repetitor := types.DBRepetitorData{
		SummaryRating: 0,
		ReviewsCount:  0,
		ResumeID:      0,
	}
	return &repetitor, types.MapperPersonalDataServiceToDB(&initData.ServicePersonalData), types.MapperAuthDataServiceToDB(&initData.ServiceAuthData)
}

func (s *RepetitorService) CreateRepetitor(initData types.ServiceInitRepetitorData) error {
	repetitor, personalData, authData := s.transormInitRepetitorData(initData)
	_, err := s.repetitorRepository.InsertRepetitor(*repetitor, *personalData, *authData)
	if err != nil {
		return err
	}
	return nil
}

func (s *RepetitorService) GetRepetitorData(userID int64) (*types.ServiceRepetitorData, error) {
	repetitor, err := s.repetitorRepository.GetRepetitor(userID)
	if err != nil {
		return nil, err
	}
	meanRating := 0.0
	if repetitor.ReviewsCount > 0 {
		meanRating = repetitor.SummaryRating / float64(repetitor.ReviewsCount)
	}
	return &types.ServiceRepetitorData{
		ID:         repetitor.ID,
		MeanRating: meanRating,
		ResumeID:   repetitor.ResumeID,
	}, nil
}

func (s *RepetitorService) UpdateRepetitorPersonalData(userID int64, personalData types.ServicePersonalData) error {
	return s.repetitorRepository.UpdateRepetitorPersonalData(userID, *types.MapperPersonalDataServiceToDB(&personalData))
}

func (s *RepetitorService) UpdateRepetitorPassword(userID int64, authData types.ServiceAuthData, newPassword string) error {
	return s.repetitorRepository.UpdateRepetitorPassword(userID, *types.MapperAuthDataServiceToDB(&authData), newPassword)
}

func (s *RepetitorService) GetRepetitorProfile(userID int64, reviewsOffset int64, reviewsLimit int64) (*types.ServiceRepetitorProfile, error) {
	userData, err := s.userRepository.GetUser(userID)
	if err != nil {
		return nil, err
	}
	repetitor, err := s.GetRepetitorData(userID)
	if err != nil {
		return nil, err
	}
	reviews, err := s.reviewRepository.GetReviewsByRepetitorID(userID, reviewsOffset, reviewsLimit)
	if err != nil {
		return nil, err
	}
	serviceReviews := make([]types.ServiceReview, len(reviews))
	for i, dbReview := range reviews {
		serviceReviews[i] = *types.MapperReviewDBToService(&dbReview)
	}
	personalData, err := s.personalDataRepository.GetPersonalData(userData.PersonalDataID)
	if err != nil {
		return nil, err
	}
	resume, err := s.resumeRepository.GetResume(userID)
	if err != nil {
		resume = &types.DBResume{
			RepetitorID: userID,
			Title:       "Very good Title",
			Description: "Very good Description",
			Prices:      map[string]int{"go": 100, "sql": 200, "dBeaver": 300},
		}
	}
	return &types.ServiceRepetitorProfile{
		FirstName:         personalData.FirstName,
		LastName:          personalData.LastName,
		MiddleName:        personalData.MiddleName,
		TelephoneNumber:   personalData.TelephoneNumber,
		Email:             personalData.Email,
		MeanRating:        repetitor.MeanRating,
		ResumeTitle:       resume.Title,
		ResumeDescription: resume.Description,
		ResumePrices:      resume.Prices,
		Reviews:           serviceReviews,
	}, nil
}

func (s *RepetitorService) GetRepetitors(repetitorsOffset int64, repetitorsLimit int64) ([]*types.ServiceRepetitorView, error) {
	repetitorsIds, err := s.repetitorRepository.GetRepetitorsIds(repetitorsOffset, repetitorsLimit)
	if err != nil {
		return nil, err
	}
	repetitors := make([]*types.ServiceRepetitorView, len(repetitorsIds))
	for i, id := range repetitorsIds {
		repetitor, err := s.GetRepetitorData(id)
		if err != nil {
			return nil, err
		}
		personalData, err := s.personalDataRepository.GetPersonalData(id)
		if err != nil {
			return nil, err
		}
		repetitors[i] = &types.ServiceRepetitorView{
			FirstName:  personalData.FirstName,
			MeanRating: repetitor.MeanRating,
		}
	}
	return repetitors, nil
}
