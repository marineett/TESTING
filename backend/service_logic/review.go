package service_logic

import (
	"data_base_project/data_base"
	"data_base_project/types"
)

type IReviewService interface {
	GetReview(id int64) (*types.ServiceReview, error)
	GetReviewsByRepetitorID(repetitorID int64, from int64, size int64) ([]types.ServiceReview, error)
	GetReviewsByClientID(clientID int64, from int64, size int64) ([]types.ServiceReview, error)
}

type ReviewService struct {
	reviewRepository data_base.IReviewRepository
}

func CreateReviewService(reviewRepository data_base.IReviewRepository) IReviewService {
	return &ReviewService{reviewRepository: reviewRepository}
}

func (s *ReviewService) GetReview(id int64) (*types.ServiceReview, error) {
	dbReview, err := s.reviewRepository.GetReview(id)
	if err != nil {
		return nil, err
	}
	return &types.ServiceReview{
		ID:          dbReview.ID,
		ContractID:  dbReview.ContractID,
		ClientID:    dbReview.ClientID,
		RepetitorID: dbReview.RepetitorID,
		Rating:      dbReview.Rating,
		Comment:     dbReview.Comment,
		CreatedAt:   dbReview.CreatedAt,
	}, nil
}

func (s *ReviewService) GetReviewsByRepetitorID(repetitorID int64, from int64, size int64) ([]types.ServiceReview, error) {
	dbReviews, err := s.reviewRepository.GetReviewsByRepetitorID(repetitorID, from, size)
	if err != nil {
		return nil, err
	}
	serviceReviews := make([]types.ServiceReview, len(dbReviews))
	for i, dbReview := range dbReviews {
		serviceReviews[i] = *types.MapperReviewDBToService(&dbReview)
	}
	return serviceReviews, nil
}

func (s *ReviewService) GetReviewsByClientID(clientID int64, from int64, size int64) ([]types.ServiceReview, error) {
	dbReviews, err := s.reviewRepository.GetReviewsByClientID(clientID, from, size)
	if err != nil {
		return nil, err
	}
	serviceReviews := make([]types.ServiceReview, len(dbReviews))
	for i, dbReview := range dbReviews {
		serviceReviews[i] = *types.MapperReviewDBToService(&dbReview)
	}
	return serviceReviews, nil
}
