package service_logic

import (
	tu "data_base_project/test_service_utility"
	"data_base_project/types"
	"database/sql"
	"testing"
)

func TestCreateRepetitorCorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	resumeRepository := tu.CreateTestResumeRepository()
	repetitorRepository := tu.CreateTestRepetiorRepository(
		personalDataRepository,
		authRepository,
		userRepository,
		resumeRepository,
	)
	repetitorService := CreateRepetitorService(
		repetitorRepository,
		personalDataRepository,
		userRepository,
		reviewRepository,
		resumeRepository,
	)
	err := repetitorService.CreateRepetitor(tu.TestInitRepetitorData)
	if err != nil {
		t.Fatalf("Error creating repetitor: %v", err)
	}
	personalData, err := personalDataRepository.GetPersonalData(1)
	if err != nil {
		t.Fatalf("Error getting personal data: %v", err)
	}
	if personalData.TelephoneNumber != tu.TestPD.TelephoneNumber {
		t.Fatalf("Personal data not updated: %v", personalData)
	}
	if personalData.Email != tu.TestPD.Email {
		t.Fatalf("Personal data not updated: %v", personalData)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestInitRepetitorData.ServiceAuthData.Login,
		Password: tu.TestInitRepetitorData.ServiceAuthData.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	authData, err := authRepository.TestGetAuth(result.UserID)
	if err != nil {
		t.Fatalf("Error getting auth data: %v", err)
	}
	if authData.Login != tu.TestInitRepetitorData.ServiceAuthData.Login {
		t.Fatalf("Auth data not updated: %v", authData)
	}
	if authData.Password != tu.TestInitRepetitorData.ServiceAuthData.Password {
		t.Fatalf("Auth data not updated: %v", authData)
	}
	repetitorData, err := repetitorRepository.GetRepetitor(1)
	if err != nil {
		t.Fatalf("Error getting repetitor data: %v", err)
	}
	if repetitorData.SummaryRating != float64(tu.TestSummaryRating) {
		t.Fatalf("Repetitor data not updated: %v", repetitorData)
	}
}

func TestCreateRepetitorDataCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up repetitor tables: %v", err)
	}
	authRepository := module.AuthRepository
	repetitorRepository := module.RepetitorRepository
	repetitorService := CreateRepetitorService(repetitorRepository, module.PersonalDataRepository, module.UserRepository, module.ReviewRepository, module.ResumeRepository)
	err = repetitorService.CreateRepetitor(tu.TestInitRepetitorData)
	if err != nil {
		t.Fatalf("Error creating repetitor: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestInitRepetitorData.ServiceAuthData.Login,
		Password: tu.TestInitRepetitorData.ServiceAuthData.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	repetitorData, err := repetitorRepository.GetRepetitor(result.UserID)
	if err != nil {
		t.Fatalf("Error getting repetitor data: %v", err)
	}
	if repetitorData.SummaryRating != float64(tu.TestSummaryRating) {
		t.Fatalf("Repetitor data not updated: %v", repetitorData)
	}
}
func TestGetRepetitorDataCorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	resumeRepository := tu.CreateTestResumeRepository()
	repetitorRepository := tu.CreateTestRepetiorRepository(
		personalDataRepository,
		authRepository,
		userRepository,
		resumeRepository,
	)
	repetitorService := CreateRepetitorService(
		repetitorRepository,
		personalDataRepository,
		userRepository,
		reviewRepository,
		resumeRepository,
	)
	err := repetitorService.CreateRepetitor(tu.TestInitRepetitorData)
	if err != nil {
		t.Fatalf("Error creating repetitor: %v", err)
	}
	repetitorData, err := repetitorService.GetRepetitorProfile(1, 0, 10)
	if err != nil {
		t.Fatalf("Error getting repetitor data: %v", err)
	}
	if repetitorData.MeanRating != tu.TestMeanRating {
		t.Fatalf("Repetitor mean rating not updated: %v", repetitorData)
	}
	if repetitorData.MeanRating != tu.TestMeanRating {
		t.Fatalf("Repetitor data not updated: %v", repetitorData)
	}
	if repetitorData.FirstName != tu.TestPD.FirstName {
		t.Fatalf("Repetitor data not updated: %v", repetitorData)
	}
	if repetitorData.LastName != tu.TestPD.LastName {
		t.Fatalf("Repetitor data not updated: %v", repetitorData)
	}
	if repetitorData.MiddleName != tu.TestPD.MiddleName {
		t.Fatalf("Repetitor data not updated: %v", repetitorData)
	}
	if repetitorData.TelephoneNumber != tu.TestPD.TelephoneNumber {
		t.Fatalf("Repetitor data not updated: %v", repetitorData)
	}
	if repetitorData.Email != tu.TestPD.Email {
		t.Fatalf("Repetitor data not updated: %v", repetitorData)
	}
}

func TestGetRepetitorDataCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up repetitor tables: %v", err)
	}
	repetitorRepository := module.RepetitorRepository
	authRepository := module.AuthRepository
	repetitorService := CreateRepetitorService(repetitorRepository, module.PersonalDataRepository, module.UserRepository, module.ReviewRepository, module.ResumeRepository)
	err = repetitorService.CreateRepetitor(tu.TestInitRepetitorData)
	if err != nil {
		t.Fatalf("Error creating repetitor: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestInitRepetitorData.ServiceAuthData.Login,
		Password: tu.TestInitRepetitorData.ServiceAuthData.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	repetitorData, err := repetitorService.GetRepetitorProfile(result.UserID, 0, 10)
	if err != nil {
		t.Fatalf("Error getting repetitor data: %v", err)
	}
	if repetitorData.MeanRating != tu.TestMeanRating {
		t.Fatalf("Repetitor mean rating not updated: %v", repetitorData)
	}
	if repetitorData.MeanRating != tu.TestMeanRating {
		t.Fatalf("Repetitor data not updated: %v", repetitorData)
	}
	if repetitorData.FirstName != tu.TestPD.FirstName {
		t.Fatalf("Repetitor data not updated: %v", repetitorData)
	}
	if repetitorData.LastName != tu.TestPD.LastName {
		t.Fatalf("Repetitor data not updated: %v", repetitorData)
	}
	if repetitorData.MiddleName != tu.TestPD.MiddleName {
		t.Fatalf("Repetitor data not updated: %v", repetitorData)
	}
	if repetitorData.TelephoneNumber != tu.TestPD.TelephoneNumber {
		t.Fatalf("Repetitor data not updated: %v", repetitorData)
	}
	if repetitorData.Email != tu.TestPD.Email {
		t.Fatalf("Repetitor data not updated: %v", repetitorData)
	}
}

func TestGetRepetitorDataIncorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	resumeRepository := tu.CreateTestResumeRepository()
	repetitorRepository := tu.CreateTestRepetiorRepository(
		personalDataRepository,
		authRepository,
		userRepository,
		resumeRepository,
	)
	repetitorService := CreateRepetitorService(
		repetitorRepository,
		personalDataRepository,
		userRepository,
		reviewRepository,
		resumeRepository,
	)
	err := repetitorService.CreateRepetitor(tu.TestInitRepetitorData)
	if err != nil {
		t.Fatalf("Error creating repetitor: %v", err)
	}
	repetitorData, err := repetitorService.GetRepetitorProfile(2, 0, 10)
	if err == nil {
		t.Fatalf("No error getting repetitor data: %v", repetitorData)
	}
}

func TestGetRepetitorDataIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up repetitor tables: %v", err)
	}
	repetitorRepository := module.RepetitorRepository
	authRepository := module.AuthRepository
	repetitorService := CreateRepetitorService(repetitorRepository, module.PersonalDataRepository, module.UserRepository, module.ReviewRepository, module.ResumeRepository)
	err = repetitorService.CreateRepetitor(tu.TestInitRepetitorData)
	if err != nil {
		t.Fatalf("Error creating repetitor: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestInitRepetitorData.ServiceAuthData.Login,
		Password: tu.TestInitRepetitorData.ServiceAuthData.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	repetitorData, err := repetitorService.GetRepetitorProfile(result.UserID+1, 0, 10)
	if err == nil {
		t.Fatalf("No error getting repetitor data: %v", repetitorData)
	}
}

func TestUpdateRepetitorPersonalDataCorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	resumeRepository := tu.CreateTestResumeRepository()
	repetitorRepository := tu.CreateTestRepetiorRepository(
		personalDataRepository,
		authRepository,
		userRepository,
		resumeRepository,
	)
	repetitorService := CreateRepetitorService(
		repetitorRepository,
		personalDataRepository,
		userRepository,
		reviewRepository,
		resumeRepository,
	)
	err := repetitorService.CreateRepetitor(tu.TestInitRepetitorData)
	if err != nil {
		t.Fatalf("Error creating repetitor: %v", err)
	}
	newPersonalData := types.ServicePersonalData{
		FirstName:       "Petr",
		LastName:        "Petrov",
		MiddleName:      "Petrovich",
		TelephoneNumber: "88005553536",
		Email:           "test2@test.com",
	}
	err = repetitorService.UpdateRepetitorPersonalData(1, newPersonalData)
	if err != nil {
		t.Fatalf("Error updating repetitor personal data: %v", err)
	}
	personalData, err := personalDataRepository.GetPersonalData(1)
	if err != nil {
		t.Fatalf("Error getting personal data: %v", err)
	}
	if personalData.FirstName != newPersonalData.FirstName {
		t.Fatalf("Repetitor personal data not updated: %v", personalData)
	}
	if personalData.LastName != newPersonalData.LastName {
		t.Fatalf("Repetitor personal data not updated: %v", personalData)
	}
	if personalData.MiddleName != newPersonalData.MiddleName {
		t.Fatalf("Repetitor personal data not updated: %v", personalData)
	}
	if personalData.TelephoneNumber != newPersonalData.TelephoneNumber {
		t.Fatalf("Repetitor personal data not updated: %v", personalData)
	}
	if personalData.Email != newPersonalData.Email {
		t.Fatalf("Repetitor personal data not updated: %v", personalData)
	}
}

func TestUpdateRepetitorPersonalDataCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up repetitor tables: %v", err)
	}
	repetitorRepository := module.RepetitorRepository
	authRepository := module.AuthRepository
	repetitorService := CreateRepetitorService(repetitorRepository, module.PersonalDataRepository, module.UserRepository, module.ReviewRepository, module.ResumeRepository)
	err = repetitorService.CreateRepetitor(tu.TestInitRepetitorData)
	if err != nil {
		t.Fatalf("Error creating repetitor: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestInitRepetitorData.ServiceAuthData.Login,
		Password: tu.TestInitRepetitorData.ServiceAuthData.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	err = repetitorService.UpdateRepetitorPersonalData(result.UserID, types.ServicePersonalData{
		FirstName:       "Petr",
		LastName:        "Petrov",
		MiddleName:      "Petrovich",
		TelephoneNumber: "88005553536",
		Email:           "test2@test.com",
	})
	if err != nil {
		t.Fatalf("Error updating repetitor personal data: %v", err)
	}
}

func TestUpdateRepetitorPersonalDataIncorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	resumeRepository := tu.CreateTestResumeRepository()
	repetitorRepository := tu.CreateTestRepetiorRepository(
		personalDataRepository,
		authRepository,
		userRepository,
		resumeRepository,
	)
	repetitorService := CreateRepetitorService(
		repetitorRepository,
		personalDataRepository,
		userRepository,
		reviewRepository,
		resumeRepository,
	)
	err := repetitorService.CreateRepetitor(tu.TestInitRepetitorData)
	if err != nil {
		t.Fatalf("Error creating repetitor: %v", err)
	}
	err = repetitorService.UpdateRepetitorPersonalData(2, types.ServicePersonalData{
		FirstName:       "Petr",
		LastName:        "Petrov",
		MiddleName:      "Petrovich",
		TelephoneNumber: "88005553536",
		Email:           "test2@test.com",
	})
	if err == nil {
		t.Fatalf("No error updating repetitor personal data: %v", err)
	}
}

func TestUpdateRepetitorPersonalDataIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up repetitor tables: %v", err)
	}
	repetitorRepository := module.RepetitorRepository
	authRepository := module.AuthRepository
	repetitorService := CreateRepetitorService(repetitorRepository, module.PersonalDataRepository, module.UserRepository, module.ReviewRepository, module.ResumeRepository)
	err = repetitorService.CreateRepetitor(tu.TestInitRepetitorData)
	if err != nil {
		t.Fatalf("Error creating repetitor: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestInitRepetitorData.ServiceAuthData.Login,
		Password: tu.TestInitRepetitorData.ServiceAuthData.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	err = repetitorService.UpdateRepetitorPersonalData(result.UserID+1, types.ServicePersonalData{
		FirstName:       "Petr",
		LastName:        "Petrov",
		MiddleName:      "Petrovich",
		TelephoneNumber: "88005553536",
		Email:           "test2@test.com",
	})
	if err == nil {
		t.Fatalf("No error updating repetitor personal data: %v", err)
	}
}

func TestUpdateRepetitorPasswordCorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	resumeRepository := tu.CreateTestResumeRepository()
	repetitorRepository := tu.CreateTestRepetiorRepository(
		personalDataRepository,
		authRepository,
		userRepository,
		resumeRepository,
	)
	repetitorService := CreateRepetitorService(
		repetitorRepository,
		personalDataRepository,
		userRepository,
		reviewRepository,
		resumeRepository,
	)
	err := repetitorService.CreateRepetitor(tu.TestInitRepetitorData)
	if err != nil {
		t.Fatalf("Error creating repetitor: %v", err)
	}
	newPassword := "test3"
	err = repetitorService.UpdateRepetitorPassword(1, tu.TestInitRepetitorData.ServiceAuthData, newPassword)
	if err != nil {
		t.Fatalf("Error updating repetitor password: %v", err)
	}
	_, err = authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestInitRepetitorData.ServiceAuthData.Login,
		Password: newPassword,
	})
	if err != nil {
		t.Fatalf("Password not updated: %v", err)
	}
}

func TestUpdateRepetitorPasswordCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up repetitor tables: %v", err)
	}
	repetitorRepository := module.RepetitorRepository
	authRepository := module.AuthRepository
	repetitorService := CreateRepetitorService(repetitorRepository, module.PersonalDataRepository, module.UserRepository, module.ReviewRepository, module.ResumeRepository)
	err = repetitorService.CreateRepetitor(tu.TestInitRepetitorData)
	if err != nil {
		t.Fatalf("Error creating repetitor: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestInitRepetitorData.ServiceAuthData.Login,
		Password: tu.TestInitRepetitorData.ServiceAuthData.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	err = repetitorService.UpdateRepetitorPassword(result.UserID, tu.TestInitRepetitorData.ServiceAuthData, "test3")
	if err != nil {
		t.Fatalf("Error updating repetitor password: %v", err)
	}
}

func TestUpdateRepetitorPasswordIncorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	reviewRepository := tu.CreateTestReviewRepository()
	resumeRepository := tu.CreateTestResumeRepository()
	repetitorRepository := tu.CreateTestRepetiorRepository(
		personalDataRepository,
		authRepository,
		userRepository,
		resumeRepository,
	)
	repetitorService := CreateRepetitorService(
		repetitorRepository,
		personalDataRepository,
		userRepository,
		reviewRepository,
		resumeRepository,
	)
	err := repetitorService.CreateRepetitor(tu.TestInitRepetitorData)
	if err != nil {
		t.Fatalf("Error creating repetitor: %v", err)
	}
	err = repetitorService.UpdateRepetitorPassword(2, tu.TestAuth, "test3")
	if err == nil {
		t.Fatalf("No error updating repetitor password: %v", err)
	}
}

func TestUpdateRepetitorPasswordIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up repetitor tables: %v", err)
	}
	repetitorRepository := module.RepetitorRepository
	authRepository := module.AuthRepository
	repetitorService := CreateRepetitorService(repetitorRepository, module.PersonalDataRepository, module.UserRepository, module.ReviewRepository, module.ResumeRepository)
	err = repetitorService.CreateRepetitor(tu.TestInitRepetitorData)
	if err != nil {
		t.Fatalf("Error creating repetitor: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestInitRepetitorData.ServiceAuthData.Login,
		Password: tu.TestInitRepetitorData.ServiceAuthData.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	err = repetitorService.UpdateRepetitorPassword(result.UserID+1, tu.TestInitRepetitorData.ServiceAuthData, "test3")
	if err == nil {
		t.Fatalf("No error updating repetitor password: %v", err)
	}
}
