package service_logic

import (
	tu "data_base_project/test_service_utility"
	"data_base_project/types"
	"database/sql"
	"testing"
	"time"
)

func TestCreateLessonCorrectLondon(t *testing.T) {
	lessonRepository := tu.CreateTestLessonRepository()
	lessonService := CreateLessonService(lessonRepository)
	tu.TestLesson.ContractID = 1
	_, err := lessonService.CreateLesson(tu.TestLesson)
	if err != nil {
		t.Fatalf("error creating lesson: %v", err)
	}
	lesson, err := lessonRepository.GetLesson(1)
	if err != nil {
		t.Fatalf("error getting lesson: %v", err)
	}
	if lesson.Duration != tu.TestLesson.Duration {
		t.Fatalf("lesson duration is not correct: %v", lesson.Duration)
	}
}

func TestCreateLessonCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up lesson tables: %v", err)
	}
	lessonRepository := module.LessonRepository
	contractRepository := module.ContractRepository
	lessonService := CreateLessonService(lessonRepository)
	clientRepository := module.ClientRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	reviewRepository := module.ReviewRepository
	authRepository := module.AuthRepository
	clientService := CreateClientService(clientRepository, personalDataRepository, userRepository, reviewRepository)
	err = clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("error creating client: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("error authorizing: %v", err)
	}
	contractID, err := contractRepository.InsertContract(types.DBContract{
		ClientID:      result.UserID,
		Status:        types.ContractStatusActive,
		PaymentStatus: types.PaymentStatusPaid,
		CreatedAt:     time.Now(),
	})
	if err != nil {
		t.Fatalf("error inserting contract: %v", err)
	}
	tu.TestLesson.ContractID = contractID
	_, err = lessonService.CreateLesson(tu.TestLesson)
	if err != nil {
		t.Fatalf("error creating lesson: %v", err)
	}
	lesson, err := lessonRepository.GetLessons(contractID, 0, 10)
	if err != nil {
		t.Fatalf("error getting lesson: %v", err)
	}
	if lesson[0].Duration != tu.TestLesson.Duration {
		t.Fatalf("lesson duration is not correct: %v", lesson[0].Duration)
	}
}

func TestCreateLessonIncorrectLondon(t *testing.T) {
	lessonRepository := tu.CreateTestLessonRepository()
	lessonService := CreateLessonService(lessonRepository)
	tu.TestLesson.ContractID = 0
	_, err := lessonService.CreateLesson(tu.TestLesson)
	if err == nil {
		t.Fatalf("No error creating lesson: %v", err)
	}
}

func TestCreateLessonIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up lesson tables: %v", err)
	}
	lessonRepository := module.LessonRepository
	lessonService := CreateLessonService(lessonRepository)
	_, err = lessonService.CreateLesson(tu.TestLesson)
	if err == nil {
		t.Fatalf("No error creating lesson: %v", err)
	}
}

func TestGetLessonsCorrectLondon(t *testing.T) {
	lessonRepository := tu.CreateTestLessonRepository()
	lessonService := CreateLessonService(lessonRepository)
	lessons, err := lessonService.GetLessons(1, 0, 10)
	if err != nil {
		t.Fatalf("error getting lessons: %v", err)
	}
	if len(lessons) != 0 {
		t.Fatalf("lessons list is not correct: %v", lessons)
	}
	tu.TestLesson.ContractID = 1
	_, err = lessonRepository.InsertLesson(*types.MapperLessonServiceToDB(&tu.TestLesson))
	if err != nil {
		t.Fatalf("error inserting lesson: %v", err)
	}
	lessons, err = lessonService.GetLessons(1, 0, 10)
	if err != nil {
		t.Fatalf("error getting lessons: %v", err)
	}
	if len(lessons) != 1 {
		t.Fatalf("lessons list is not correct: %v", lessons)
	}
}

func TestGetLessonsCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up lesson tables: %v", err)
	}
	lessonRepository := module.LessonRepository
	contractRepository := module.ContractRepository
	lessonService := CreateLessonService(lessonRepository)
	clientRepository := module.ClientRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	reviewRepository := module.ReviewRepository
	authRepository := module.AuthRepository
	clientService := CreateClientService(clientRepository, personalDataRepository, userRepository, reviewRepository)
	err = clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("error creating client: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("error authorizing: %v", err)
	}
	contractID, err := contractRepository.InsertContract(types.DBContract{
		ClientID:      result.UserID,
		Status:        types.ContractStatusActive,
		PaymentStatus: types.PaymentStatusPaid,
		CreatedAt:     time.Now(),
	})
	if err != nil {
		t.Fatalf("error inserting contract: %v", err)
	}
	tu.TestLesson.ContractID = contractID
	_, err = lessonService.CreateLesson(tu.TestLesson)
	if err != nil {
		t.Fatalf("error creating lesson: %v", err)
	}
	lesson, err := lessonService.GetLessons(contractID, 0, 10)
	if err != nil {
		t.Fatalf("error getting lesson: %v", err)
	}
	if lesson[0].Duration != tu.TestLesson.Duration {
		t.Fatalf("lesson duration is not correct: %v", lesson[0].Duration)
	}
}
