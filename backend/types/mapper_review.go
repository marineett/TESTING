package types

func MapperReviewDBToService(review *DBReview) *ServiceReview {
	if review == nil {
		return nil
	}
	return &ServiceReview{
		ID:          review.ID,
		ContractID:  review.ContractID,
		ClientID:    review.ClientID,
		RepetitorID: review.RepetitorID,
		Rating:      review.Rating,
		Comment:     review.Comment,
		CreatedAt:   review.CreatedAt,
	}
}

func MapperReviewServiceToDB(review *ServiceReview) *DBReview {
	if review == nil {
		return nil
	}
	return &DBReview{
		ID:          review.ID,
		ContractID:  review.ContractID,
		ClientID:    review.ClientID,
		RepetitorID: review.RepetitorID,
		Rating:      review.Rating,
		Comment:     review.Comment,
		CreatedAt:   review.CreatedAt,
	}
}
