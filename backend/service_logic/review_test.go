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
	defer db.Close()
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
	defer db.Close()
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

func TestGetReviewsByRepetitorIDCorrectLondon(t *testing.T) {
	reviewRepository := tu.CreateTestReviewRepository()
	reviewService := CreateReviewService(reviewRepository)
	firstReview := tu.TestReview
	firstReview.RepetitorID = 1
	firstReview.ClientID = 2
	secondReview := tu.TestReview
	secondReview.RepetitorID = 1
	secondReview.ClientID = 3
	thirdReview := tu.TestReview
	thirdReview.RepetitorID = 5
	thirdReview.ClientID = 4
	reviewRepository.InsertReview(firstReview)
	reviewRepository.InsertReview(secondReview)
	reviewRepository.InsertReview(thirdReview)
	reviews, err := reviewService.GetReviewsByRepetitorID(1, 0, 10)
	if err != nil {
		t.Fatalf("Error getting reviews: %v", err)
	}
	if len(reviews) != 2 {
		t.Fatalf("Reviews not updated: %v", reviews)
	}
	if reviews[0].RepetitorID != 1 || reviews[1].RepetitorID != 1 {
		t.Fatalf("Repetitor id not updated: %v", reviews[0])
	}
	if reviews[0].Rating != tu.TestReview.Rating || reviews[1].Rating != tu.TestReview.Rating {
		t.Fatalf("Rating not updated: %v", reviews[0])
	}
	if reviews[0].Comment != tu.TestReview.Comment || reviews[1].Comment != tu.TestReview.Comment {
		t.Fatalf("Comment not updated: %v", reviews[0])
	}
	reviews, err = reviewService.GetReviewsByRepetitorID(2, 0, 10)
	if err != nil {
		t.Fatalf("Error getting reviews: %v", err)
	}
	if len(reviews) != 0 {
		t.Fatalf("Reviews not updated: %v", reviews)
	}
	reviews, err = reviewService.GetReviewsByRepetitorID(1, 1, 1)
	if err != nil {
		t.Fatalf("Error getting reviews: %v", err)
	}
	if len(reviews) != 1 {
		t.Fatalf("Reviews not updated: %v", reviews)
	}
	if reviews[0].RepetitorID != 1 {
		t.Fatalf("Repetitor id not updated: %v", reviews[0])
	}
	if reviews[0].Rating != tu.TestReview.Rating {
		t.Fatalf("Rating not updated: %v", reviews[0])
	}
	if reviews[0].Comment != tu.TestReview.Comment {
		t.Fatalf("Comment not updated: %v", reviews[0])
	}
}

func TestGetReviewsByRepetitorIDCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
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
	reviews, err := reviewService.GetReviewsByRepetitorID(repetitorID, 0, 10)
	if err != nil {
		t.Fatalf("Error getting reviews: %v", err)
	}
	if len(reviews) != 2 {
		t.Fatalf("Reviews not updated: %v", reviews)
	}
}

func TestGetReviewsByClientIDCorrectLondon(t *testing.T) {
	reviewRepository := tu.CreateTestReviewRepository()
	reviewService := CreateReviewService(reviewRepository)
	firstReview := tu.TestReview
	firstReview.RepetitorID = 2
	firstReview.ClientID = 1
	secondReview := tu.TestReview
	secondReview.RepetitorID = 3
	secondReview.ClientID = 1
	thirdReview := tu.TestReview
	thirdReview.RepetitorID = 4
	thirdReview.ClientID = 5
	reviewRepository.InsertReview(firstReview)
	reviewRepository.InsertReview(secondReview)
	reviewRepository.InsertReview(thirdReview)
	reviews, err := reviewService.GetReviewsByClientID(1, 0, 10)
	if err != nil {
		t.Fatalf("Error getting reviews: %v", err)
	}
	if len(reviews) != 2 {
		t.Fatalf("Reviews not updated: %v", reviews)
	}
	if reviews[0].ClientID != 1 || reviews[1].ClientID != 1 {
		t.Fatalf("Client id not updated: %v", reviews[0])
	}
	if reviews[0].Rating != tu.TestReview.Rating || reviews[1].Rating != tu.TestReview.Rating {
		t.Fatalf("Rating not updated: %v", reviews[0])
	}
	if reviews[0].Comment != tu.TestReview.Comment || reviews[1].Comment != tu.TestReview.Comment {
		t.Fatalf("Comment not updated: %v", reviews[0])
	}
	reviews, err = reviewService.GetReviewsByClientID(2, 0, 10)
	if err != nil {
		t.Fatalf("Error getting reviews: %v", err)
	}
	if len(reviews) != 0 {
		t.Fatalf("Reviews not updated: %v", reviews)
	}
	reviews, err = reviewService.GetReviewsByClientID(1, 1, 1)
	if err != nil {
		t.Fatalf("Error getting reviews: %v", err)
	}
	if len(reviews) != 1 {
		t.Fatalf("Reviews not updated: %v", reviews)
	}
	if reviews[0].ClientID != 1 {
		t.Fatalf("Client id not updated: %v", reviews[0])
	}
	if reviews[0].Rating != tu.TestReview.Rating {
		t.Fatalf("Rating not updated: %v", reviews[0])
	}
	if reviews[0].Comment != tu.TestReview.Comment {
		t.Fatalf("Comment not updated: %v", reviews[0])
	}
}

func TestGetReviewsByClientIDCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
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
	reviews, err := reviewService.GetReviewsByClientID(clientID, 0, 10)
	if err != nil {
		t.Fatalf("Error getting reviews: %v", err)
	}
	if len(reviews) != 2 {
		t.Fatalf("Reviews not updated: %v", reviews)
	}
}
