package service_logic

import (
	tu "data_base_project/test_service_utility"
	"data_base_project/types"
	"database/sql"
	"testing"
)

func TestGetReviewCorrectLondon(t *testing.T) {
	reviewRepository := tu.CreateTestReviewRepository()
	reviewService := CreateReviewService(reviewRepository)
	reviewID, err := reviewRepository.InsertReview(tu.TestReview)
	if err != nil {
		t.Fatalf("Error inserting review: %v", err)
	}
	reviewServiceData, err := reviewService.GetReview(reviewID)
	if err != nil {
		t.Fatalf("Error getting review: %v", err)
	}
	if reviewServiceData.Rating != tu.TestReview.Rating {
		t.Fatalf("Review rating not updated: %v", reviewServiceData)
	}
	if reviewServiceData.Comment != tu.TestReview.Comment {
		t.Fatalf("Review comment not updated: %v", reviewServiceData)
	}
}

func TestGetReviewCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up review tables: %v", err)
	}
	clientRepository := module.ClientRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	reviewRepository := module.ReviewRepository
	resumeRepository := module.ResumeRepository
	repetitorRepository := module.RepetitorRepository
	reviewService := CreateReviewService(reviewRepository)
	clientService := CreateClientService(clientRepository, personalDataRepository, userRepository, reviewRepository)
	repetitorService := CreateRepetitorService(repetitorRepository, personalDataRepository, userRepository, reviewRepository, resumeRepository)
	err = clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	authRepository := module.AuthRepository
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	clientID := result.UserID
	err = repetitorService.CreateRepetitor(tu.TestInitRepetitorData)
	if err != nil {
		t.Fatalf("Error creating repetitor: %v", err)
	}
	repetitorID := result.UserID
	tu.TestReview.RepetitorID = repetitorID
	tu.TestReview.ClientID = clientID
	reviewID, err := reviewRepository.InsertReview(tu.TestReview)
	if err != nil {
		t.Fatalf("Error inserting review: %v", err)
	}
	_, err = reviewService.GetReview(reviewID)
	if err != nil {
		t.Fatalf("Error getting review: %v", err)
	}
}
func TestGetReviewIncorrectLondon(t *testing.T) {
	reviewRepository := tu.CreateTestReviewRepository()
	reviewService := CreateReviewService(reviewRepository)
	reviewServiceData, err := reviewService.GetReview(1)
	if err == nil {
		t.Fatalf("No error getting review: %v", reviewServiceData)
	}
}

func TestGetReviewIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up review tables: %v", err)
	}
	reviewRepository := module.ReviewRepository
	reviewService := CreateReviewService(reviewRepository)
	reviewServiceData, err := reviewService.GetReview(1)
	if err == nil {
		t.Fatalf("No error getting review: %v", reviewServiceData)
	}
}

func CheckReviewsLength(t *testing.T, reviews []types.ServiceReview, length int) {
	if len(reviews) != length {
		t.Fatalf("Reviews not updated: %v", reviews)
	}
}

func TestGetReviewsByRepetitorIDCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up review tables: %v", err)
	}
	clientRepository := module.ClientRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	reviewRepository := module.ReviewRepository
	resumeRepository := module.ResumeRepository
	repetitorRepository := module.RepetitorRepository
	reviewService := CreateReviewService(reviewRepository)
	clientService := CreateClientService(clientRepository, personalDataRepository, userRepository, reviewRepository)
	repetitorService := CreateRepetitorService(repetitorRepository, personalDataRepository, userRepository, reviewRepository, resumeRepository)
	err = clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	authRepository := module.AuthRepository
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	clientID := result.UserID
	err = repetitorService.CreateRepetitor(tu.TestInitRepetitorData)
	if err != nil {
		t.Fatalf("Error creating repetitor: %v", err)
	}
	repetitorID := result.UserID
	tu.TestReview.RepetitorID = repetitorID
	tu.TestReview.ClientID = clientID
	_, err = reviewRepository.InsertReview(tu.TestReview)
	if err != nil {
		t.Fatalf("Error inserting review: %v", err)
	}
	_, err = reviewRepository.InsertReview(tu.TestReview)
	if err != nil {
		t.Fatalf("Error inserting review: %v", err)
	}
	_, err = reviewService.GetReviewsByRepetitorID(repetitorID, 0, 10)
	if err != nil {
		t.Fatalf("Error getting reviews: %v", err)
	}
}

func CheckReview(
	t *testing.T,
	Review *types.ServiceReview,
	ReviewID int64,
	Rating int,
	Comment string,
) {
	if Review.ID != ReviewID {
		t.Fatalf("Review id not correct: %v", Review.ID)
	}
	if Review.Rating != Rating {
		t.Fatalf("Review rating not correct: %v", Review.Rating)
	}
	if Review.Comment != Comment {
		t.Fatalf("Review comment not correct: %v", Review.Comment)
	}
}

func TestGetReviewsByClientIDCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up review tables: %v", err)
	}
	clientRepository := module.ClientRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	reviewRepository := module.ReviewRepository
	resumeRepository := module.ResumeRepository
	repetitorRepository := module.RepetitorRepository
	reviewService := CreateReviewService(reviewRepository)
	clientService := CreateClientService(clientRepository, personalDataRepository, userRepository, reviewRepository)
	repetitorService := CreateRepetitorService(repetitorRepository, personalDataRepository, userRepository, reviewRepository, resumeRepository)
	err = clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	authRepository := module.AuthRepository
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	clientID := result.UserID
	err = repetitorService.CreateRepetitor(tu.TestInitRepetitorData)
	if err != nil {
		t.Fatalf("Error creating repetitor: %v", err)
	}
	repetitorID := result.UserID
	tu.TestReview.RepetitorID = repetitorID
	tu.TestReview.ClientID = clientID
	_, err = reviewRepository.InsertReview(tu.TestReview)
	if err != nil {
		t.Fatalf("Error inserting review: %v", err)
	}
	_, err = reviewRepository.InsertReview(tu.TestReview)
	if err != nil {
		t.Fatalf("Error inserting review: %v", err)
	}
	_, err = reviewService.GetReviewsByClientID(clientID, 0, 10)
	if err != nil {
		t.Fatalf("Error getting reviews: %v", err)
	}
}
