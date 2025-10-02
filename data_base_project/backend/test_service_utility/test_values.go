package test_service_utility

import (
	"data_base_project/types"
	"time"
)

var (
	TestPD = types.ServicePersonalData{
		TelephoneNumber: "88005553535",
		Email:           "test@test.com",
		FirstName:       "Ivan",
		LastName:        "Ivanov",
		MiddleName:      "Ivanovich",
		ServicePassportData: types.ServicePassportData{
			PassportNumber:   "1234567890",
			PassportSeries:   "1234567890",
			PassportDate:     time.Now(),
			PassportIssuedBy: "test",
		},
	}
	TestAdmin = types.DBAdminData{
		ID:     1,
		Salary: 100000,
	}
	TestAuth = types.ServiceAuthData{
		Login:    "test1",
		Password: "test2",
	}
	TestSalary = int64(100000)

	TestInitAdminData = types.ServiceInitAdminData{
		ServiceInitUserData: types.ServiceInitUserData{
			ServicePersonalData: TestPD,
			ServiceAuthData:     TestAuth,
		},
		Salary: TestSalary,
	}
	TestInitClientData = types.ServiceInitClientData{
		ServiceInitUserData: types.ServiceInitUserData{
			ServicePersonalData: TestPD,
			ServiceAuthData:     TestAuth,
		},
	}
	TestSummaryRating     = int64(0)
	TestMeanRating        = float64(0)
	TestInitModeratorData = types.ServiceInitModeratorData{
		ServiceInitUserData: types.ServiceInitUserData{
			ServicePersonalData: TestPD,
			ServiceAuthData:     TestAuth,
		},
	}
	TestInitRepetitorData = types.ServiceInitRepetitorData{
		ServiceInitUserData: types.ServiceInitUserData{
			ServicePersonalData: TestPD,
			ServiceAuthData:     TestAuth,
		},
	}
	TestResume = types.ServiceResume{
		RepetitorID: 1,
		Title:       "just a title",
		Description: "just a description",
		Prices:      map[string]int{"Go": 100, "C++": 200, "Python": 100},
	}
	TestRating  = 5
	TestComment = "comment for good work"
	TestReview  = types.DBReview{
		Rating:  TestRating,
		Comment: TestComment,
	}
	TestTransaction = types.DBTransaction{
		UserID:    1,
		Amount:    100,
		Status:    types.TransactionStatusPending,
		Type:      types.TransactionTypeContractPayment,
		CreatedAt: time.Now(),
	}
	TestPendingContractPaymentTransaction = types.DBPendingContractPaymentTransaction{
		UserID:        1,
		Amount:        100,
		CreatedAt:     time.Now(),
		TransactionID: 1,
	}
	TestLesson = types.ServiceLesson{
		Duration:  100,
		CreatedAt: time.Now(),
	}
	TestServiceContractInitData = types.ServiceContractInitData{
		ClientID:              1,
		ContractCategory:      types.ContractCategoryTranslation,
		ContractSubcategories: []types.ContractSubcategory{types.ContractSubcategoryTranslation},
		Description:           "test description",
		Price:                 100,
		Commission:            10,
		StartDate:             time.Now(),
		Duration:              10,
	}
	TestServiceReview = types.ServiceReview{
		ClientID:    1,
		RepetitorID: 2,
		Rating:      TestRating,
		Comment:     TestComment,
		CreatedAt:   time.Now(),
	}
	TestInitDepartmentData = types.ServiceDepartmentInitData{
		Name:   "test department",
		HeadID: 1,
	}
)
