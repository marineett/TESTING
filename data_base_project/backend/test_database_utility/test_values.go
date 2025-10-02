package test_database_utility

import (
	"data_base_project/types"
	"time"
)

var (
	TestPD = types.DBPersonalData{
		TelephoneNumber: "88005553535",
		Email:           "test@test.com",
		FirstName:       "Ivan",
		LastName:        "Ivanov",
		MiddleName:      "Ivanovich",
		DBPassportData: types.DBPassportData{
			PassportNumber:   "1234567890",
			PassportSeries:   "1234",
			PassportDate:     time.Now(),
			PassportIssuedBy: "Moscow",
		},
	}
	TestUser = types.DBUserData{
		RegistrationDate: time.Now(),
		LastLoginDate:    time.Now(),
	}
	TestAdmin = types.DBAdminData{
		Salary: 100000,
	}
	TestClient = types.DBClientData{
		SummaryRating: 0,
		ReviewsCount:  0,
	}
	TestRepetitor = types.DBRepetitorData{
		SummaryRating: 0,
		ReviewsCount:  0,
	}
	TestAuthData = types.DBAuthData{
		Login:    "test1",
		Password: "test2",
	}
	TestAuthInfo = types.DBAuthInfo{
		UserID:   1,
		UserType: types.Admin,
		Login:    "test1",
		Password: "test2",
	}
	TestSalary        = int64(100000)
	TestInitAdminData = types.DBAdminData{
		Salary: 100000,
	}
	TestModeratorData = types.DBModeratorData{
		Salary: 100000,
	}
	TestTransaction = types.DBTransaction{
		Amount:    100,
		Status:    types.TransactionStatusPending,
		Type:      types.TransactionTypeContractPayment,
		CreatedAt: time.Now(),
	}
	TestPendingContractPaymentTransaction = types.DBPendingContractPaymentTransaction{
		UserID:    1,
		Amount:    100,
		CreatedAt: time.Now(),
	}
	TestRating  = 5
	TestComment = "comment for good work"
	TestReview  = types.DBReview{
		Rating:    TestRating,
		Comment:   TestComment,
		CreatedAt: time.Now(),
	}
	TestMessage = types.DBMessage{
		Content:   "Hello",
		CreatedAt: time.Now(),
	}
	TestContract = types.DBContract{
		Description:       "Test contract",
		Status:            types.ContractStatusPending,
		PaymentStatus:     types.PaymentStatusNull,
		ReviewClientID:    0,
		ReviewRepetitorID: 0,
		Price:             100,
		Commission:        10,
		StartDate:         time.Now(),
		EndDate:           time.Now(),
		TransactionID:     0,
		IDCRChat:          0,
		IDCMRepChat:       0,
		IDRMRepChat:       0,
	}
	TestDepartment = types.DBDepartment{
		Name:   "Test Department",
		HeadID: 1,
	}
	TestHireInfo = types.DBHireInfo{
		DepartmentID: 1,
		UserID:       1,
	}
	TestLesson = types.DBLesson{
		ContractID: 1,
		Duration:   10,
		CreatedAt:  time.Now(),
	}
	TestResume = types.DBResume{
		RepetitorID: 1,
		Title:       "Test Resume",
		Description: "Test Description",
		Prices:      map[string]int{"Test": 100},
	}
	TestResumePrices = map[string]int{"Test": 100}
)
