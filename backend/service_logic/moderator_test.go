package service_logic

import (
	tu "data_base_project/test_service_utility"
	"data_base_project/types"
	"database/sql"
	"testing"
)

func TestCreateModeratorCorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	departmentRepository := tu.CreateTestDepartmentRepository()
	moderatorRepository := tu.CreateTestModeratorRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	moderatorService := CreateModeratorService(
		moderatorRepository,
		personalDataRepository,
		userRepository,
		departmentRepository,
	)
	err := moderatorService.CreateModerator(tu.TestInitModeratorData)
	if err != nil {
		t.Fatalf("Error creating moderator: %v", err)
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
	authData, err := authRepository.TestGetAuth(1)
	if err != nil {
		t.Fatalf("Error getting auth data: %v", err)
	}
	if authData.Login != tu.TestAuth.Login {
		t.Fatalf("Auth data not updated: %v", authData)
	}
	if authData.Password != tu.TestAuth.Password {
		t.Fatalf("Auth data not updated: %v", authData)
	}
	moderatorData, err := moderatorRepository.GetModerator(1)
	if err != nil {
		t.Fatalf("Error getting moderator data: %v", err)
	}
	if moderatorData.Salary != 0 {
		t.Fatalf("Moderator data not updated: %v", moderatorData)
	}
}

func TestCreateModeratorCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up moderator tables: %v", err)
	}
	moderatorRepository := module.ModeratorRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	departmentRepository := module.DepartmentRepository
	moderatorService := CreateModeratorService(
		moderatorRepository,
		personalDataRepository,
		userRepository,
		departmentRepository,
	)
	err = moderatorService.CreateModerator(tu.TestInitModeratorData)
	if err != nil {
		t.Fatalf("Error creating moderator: %v", err)
	}
}

func TestGetModeratorDataCorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	departmentRepository := tu.CreateTestDepartmentRepository()
	moderatorRepository := tu.CreateTestModeratorRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	moderatorService := CreateModeratorService(
		moderatorRepository,
		personalDataRepository,
		userRepository,
		departmentRepository,
	)
	err := moderatorService.CreateModerator(tu.TestInitModeratorData)
	if err != nil {
		t.Fatalf("Error creating moderator: %v", err)
	}
	moderatorData, err := moderatorService.GetModeratorData(1)
	if err != nil {
		t.Fatalf("Error getting moderator data: %v", err)
	}
	if moderatorData.Salary != 0 {
		t.Fatalf("Moderator data not updated: %v", moderatorData)
	}
	if len(moderatorData.Departments) != 0 {
		t.Fatalf("Moderator data not updated: %v", moderatorData)
	}
	if moderatorData.ID != 1 {
		t.Fatalf("Moderator data not updated: %v", moderatorData)
	}
}

func TestGetModeratorDataCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up moderator tables: %v", err)
	}
	moderatorRepository := module.ModeratorRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	departmentRepository := module.DepartmentRepository
	authRepository := module.AuthRepository
	moderatorService := CreateModeratorService(
		moderatorRepository,
		personalDataRepository,
		userRepository,
		departmentRepository,
	)
	err = moderatorService.CreateModerator(tu.TestInitModeratorData)
	if err != nil {
		t.Fatalf("Error creating moderator: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	moderatorData, err := moderatorService.GetModeratorData(result.UserID)
	if err != nil {
		t.Fatalf("Error getting moderator data: %v", err)
	}
	if moderatorData.Salary != 0 {
		t.Fatalf("Moderator data not updated: %v", moderatorData)
	}
	if len(moderatorData.Departments) != 0 {
		t.Fatalf("Moderator data not updated: %v", moderatorData)
	}
	if moderatorData.ID != result.UserID {
		t.Fatalf("Moderator data not updated: %v", moderatorData)
	}
}

func TestGetModeratorDataIncorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	departmentRepository := tu.CreateTestDepartmentRepository()
	moderatorRepository := tu.CreateTestModeratorRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	moderatorService := CreateModeratorService(
		moderatorRepository,
		personalDataRepository,
		userRepository,
		departmentRepository,
	)
	err := moderatorService.CreateModerator(tu.TestInitModeratorData)
	if err != nil {
		t.Fatalf("Error creating moderator: %v", err)
	}
	moderatorData, err := moderatorService.GetModeratorData(2)
	if err == nil {
		t.Fatalf("No error getting moderator data: %v", moderatorData)
	}
}

func TestGetModeratorDataIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up moderator tables: %v", err)
	}
	moderatorRepository := module.ModeratorRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	departmentRepository := module.DepartmentRepository
	authRepository := module.AuthRepository
	moderatorService := CreateModeratorService(
		moderatorRepository,
		personalDataRepository,
		userRepository,
		departmentRepository,
	)
	err = moderatorService.CreateModerator(tu.TestInitModeratorData)
	if err != nil {
		t.Fatalf("Error creating moderator: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	moderatorData, err := moderatorService.GetModeratorData(result.UserID + 1)
	if err == nil {
		t.Fatalf("No error getting moderator data: %v", moderatorData)
	}
}

func TestGetModeratorProfileCorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	departmentRepository := tu.CreateTestDepartmentRepository()
	moderatorRepository := tu.CreateTestModeratorRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	moderatorService := CreateModeratorService(
		moderatorRepository,
		personalDataRepository,
		userRepository,
		departmentRepository,
	)
	moderatorService.CreateModerator(tu.TestInitModeratorData)
	moderatorProfile, err := moderatorService.GetModeratorProfile(1)
	if err != nil {
		t.Fatalf("Error getting moderator profile: %v", err)
	}
	if moderatorProfile.Salary != 0 {
		t.Fatalf("Moderator profile not updated: %v", moderatorProfile)
	}
	if moderatorProfile.FirstName != tu.TestPD.FirstName {
		t.Fatalf("Moderator profile not updated: %v", moderatorProfile)
	}
	if moderatorProfile.LastName != tu.TestPD.LastName {
		t.Fatalf("Moderator profile not updated: %v", moderatorProfile)
	}
	if moderatorProfile.MiddleName != tu.TestPD.MiddleName {
		t.Fatalf("Moderator profile not updated: %v", moderatorProfile)
	}
	if moderatorProfile.TelephoneNumber != tu.TestPD.TelephoneNumber {
		t.Fatalf("Moderator profile not updated: %v", moderatorProfile)
	}
	if moderatorProfile.Email != tu.TestPD.Email {
		t.Fatalf("Moderator profile not updated: %v", moderatorProfile)
	}
}

func TestGetModeratorProfileCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up moderator tables: %v", err)
	}
	moderatorRepository := module.ModeratorRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	departmentRepository := module.DepartmentRepository
	authRepository := module.AuthRepository
	moderatorService := CreateModeratorService(
		moderatorRepository,
		personalDataRepository,
		userRepository,
		departmentRepository,
	)
	err = moderatorService.CreateModerator(tu.TestInitModeratorData)
	if err != nil {
		t.Fatalf("Error creating moderator: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	moderatorProfile, err := moderatorService.GetModeratorProfile(result.UserID)
	if err != nil {
		t.Fatalf("Error getting moderator profile: %v", err)
	}
	if moderatorProfile.Salary != 0 {
		t.Fatalf("Moderator profile not updated: %v", moderatorProfile)
	}
	if moderatorProfile.FirstName != tu.TestPD.FirstName {
		t.Fatalf("Moderator profile not updated: %v", moderatorProfile)
	}
	if moderatorProfile.LastName != tu.TestPD.LastName {
		t.Fatalf("Moderator profile not updated: %v", moderatorProfile)
	}
	if moderatorProfile.MiddleName != tu.TestPD.MiddleName {
		t.Fatalf("Moderator profile not updated: %v", moderatorProfile)
	}
	if moderatorProfile.TelephoneNumber != tu.TestPD.TelephoneNumber {
		t.Fatalf("Moderator profile not updated: %v", moderatorProfile)
	}
	if moderatorProfile.Email != tu.TestPD.Email {
		t.Fatalf("Moderator profile not updated: %v", moderatorProfile)
	}
}

func TestGetModeratorProfileIncorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	departmentRepository := tu.CreateTestDepartmentRepository()
	moderatorRepository := tu.CreateTestModeratorRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	moderatorService := CreateModeratorService(
		moderatorRepository,
		personalDataRepository,
		userRepository,
		departmentRepository,
	)
	err := moderatorService.CreateModerator(tu.TestInitModeratorData)
	if err != nil {
		t.Fatalf("Error creating moderator: %v", err)
	}
	moderatorProfile, err := moderatorService.GetModeratorProfile(2)
	if err == nil {
		t.Fatalf("No error getting moderator profile: %v", moderatorProfile)
	}
}

func TestGetModeratorProfileIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up moderator tables: %v", err)
	}
	moderatorRepository := module.ModeratorRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	departmentRepository := module.DepartmentRepository
	authRepository := module.AuthRepository
	moderatorService := CreateModeratorService(
		moderatorRepository,
		personalDataRepository,
		userRepository,
		departmentRepository,
	)
	err = moderatorService.CreateModerator(tu.TestInitModeratorData)
	if err != nil {
		t.Fatalf("Error creating moderator: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	moderatorProfile, err := moderatorService.GetModeratorProfile(result.UserID + 1)
	if err == nil {
		t.Fatalf("No error getting moderator profile: %v", moderatorProfile)
	}
}

func TestUpdateModeratorPersonalDataCorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	departmentRepository := tu.CreateTestDepartmentRepository()
	moderatorRepository := tu.CreateTestModeratorRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	moderatorService := CreateModeratorService(
		moderatorRepository,
		personalDataRepository,
		userRepository,
		departmentRepository,
	)
	err := moderatorService.CreateModerator(tu.TestInitModeratorData)
	if err != nil {
		t.Fatalf("Error creating moderator: %v", err)
	}
	newPersonalData := types.ServicePersonalData{
		FirstName:       "Petr",
		LastName:        "Petrov",
		MiddleName:      "Petrovich",
		TelephoneNumber: "88005553536",
		Email:           "test2@test.com",
	}
	err = moderatorService.UpdateModeratorPersonalData(1, newPersonalData)
	if err != nil {
		t.Fatalf("Error updating moderator personal data: %v", err)
	}
	personalData, err := personalDataRepository.GetPersonalData(1)
	if err != nil {
		t.Fatalf("Error getting personal data: %v", err)
	}
	if personalData.FirstName != newPersonalData.FirstName {
		t.Fatalf("Moderator personal data not updated: %v", personalData)
	}
	if personalData.LastName != newPersonalData.LastName {
		t.Fatalf("Moderator personal data not updated: %v", personalData)
	}
	if personalData.MiddleName != newPersonalData.MiddleName {
		t.Fatalf("Moderator personal data not updated: %v", personalData)
	}
	if personalData.TelephoneNumber != newPersonalData.TelephoneNumber {
		t.Fatalf("Moderator personal data not updated: %v", personalData)
	}
	if personalData.Email != newPersonalData.Email {
		t.Fatalf("Moderator personal data not updated: %v", personalData)
	}
}

func TestUpdateModeratorPersonalDataCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up moderator tables: %v", err)
	}
	moderatorRepository := module.ModeratorRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	departmentRepository := module.DepartmentRepository
	authRepository := module.AuthRepository
	moderatorService := CreateModeratorService(
		moderatorRepository,
		personalDataRepository,
		userRepository,
		departmentRepository,
	)
	err = moderatorService.CreateModerator(tu.TestInitModeratorData)
	if err != nil {
		t.Fatalf("Error creating moderator: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	err = moderatorService.UpdateModeratorPersonalData(result.UserID, types.ServicePersonalData{
		FirstName:       "Petr",
		LastName:        "Petrov",
		MiddleName:      "Petrovich",
		TelephoneNumber: "88005553536",
		Email:           "test2@test.com",
	})
	if err != nil {
		t.Fatalf("Error updating moderator personal data: %v", err)
	}
	personalData, err := personalDataRepository.GetPersonalData(1)
	if err != nil {
		t.Fatalf("Error getting personal data: %v", err)
	}
	if personalData.FirstName != "Petr" {
		t.Fatalf("Moderator personal data not updated: %v", personalData)
	}
	if personalData.LastName != "Petrov" {
		t.Fatalf("Moderator personal data not updated: %v", personalData)
	}
	if personalData.MiddleName != "Petrovich" {
		t.Fatalf("Moderator personal data not updated: %v", personalData)
	}
	if personalData.TelephoneNumber != "88005553536" {
		t.Fatalf("Moderator personal data not updated: %v", personalData)
	}
	if personalData.Email != "test2@test.com" {
		t.Fatalf("Moderator personal data not updated: %v", personalData)
	}
}

func TestUpdateModeratorPersonalDataIncorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	departmentRepository := tu.CreateTestDepartmentRepository()
	moderatorRepository := tu.CreateTestModeratorRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	moderatorService := CreateModeratorService(
		moderatorRepository,
		personalDataRepository,
		userRepository,
		departmentRepository,
	)
	err := moderatorService.CreateModerator(tu.TestInitModeratorData)
	if err != nil {
		t.Fatalf("Error creating moderator: %v", err)
	}
	err = moderatorService.UpdateModeratorPersonalData(2, types.ServicePersonalData{
		FirstName:       "Petr",
		LastName:        "Petrov",
		MiddleName:      "Petrovich",
		TelephoneNumber: "88005553536",
		Email:           "test2@test.com",
	})
	if err == nil {
		t.Fatalf("No error updating moderator personal data: %v", err)
	}
}

func TestUpdateModeratorPersonalDataIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up moderator tables: %v", err)
	}
	moderatorRepository := module.ModeratorRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	departmentRepository := module.DepartmentRepository
	authRepository := module.AuthRepository
	moderatorService := CreateModeratorService(
		moderatorRepository,
		personalDataRepository,
		userRepository,
		departmentRepository,
	)
	err = moderatorService.CreateModerator(tu.TestInitModeratorData)
	if err != nil {
		t.Fatalf("Error creating moderator: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	err = moderatorService.UpdateModeratorPersonalData(result.UserID+1, types.ServicePersonalData{
		FirstName:       "Petr",
		LastName:        "Petrov",
		MiddleName:      "Petrovich",
		TelephoneNumber: "88005553536",
		Email:           "test2@test.com",
	})
	if err == nil {
		t.Fatalf("No error updating moderator personal data: %v", err)
	}
}

func TestUpdateModeratorPasswordCorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	departmentRepository := tu.CreateTestDepartmentRepository()
	moderatorRepository := tu.CreateTestModeratorRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	moderatorService := CreateModeratorService(
		moderatorRepository,
		personalDataRepository,
		userRepository,
		departmentRepository,
	)
	err := moderatorService.CreateModerator(tu.TestInitModeratorData)
	if err != nil {
		t.Fatalf("Error creating moderator: %v", err)
	}
	newPassword := "test3"
	err = moderatorService.UpdateModeratorPassword(1, tu.TestAuth, newPassword)
	if err != nil {
		t.Fatalf("Error updating moderator password: %v", err)
	}
	authData, err := authRepository.TestGetAuth(1)
	if err != nil {
		t.Fatalf("Error getting auth data: %v", err)
	}
	if authData.Password != newPassword {
		t.Fatalf("Moderator password not updated: %v", authData)
	}
}

func TestUpdateModeratorPasswordCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up moderator tables: %v", err)
	}
	moderatorRepository := module.ModeratorRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	departmentRepository := module.DepartmentRepository
	authRepository := module.AuthRepository
	moderatorService := CreateModeratorService(
		moderatorRepository,
		personalDataRepository,
		userRepository,
		departmentRepository,
	)
	err = moderatorService.CreateModerator(tu.TestInitModeratorData)
	if err != nil {
		t.Fatalf("Error creating moderator: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	err = moderatorService.UpdateModeratorPassword(result.UserID, tu.TestAuth, "test3")
	if err != nil {
		t.Fatalf("Error updating moderator password: %v", err)
	}
	result, err = authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: "test3",
	})
	if err != nil {
		t.Fatalf("Error updating moderator password: %v", err)
	}
}

func TestUpdateModeratorPasswordIncorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	departmentRepository := tu.CreateTestDepartmentRepository()
	moderatorRepository := tu.CreateTestModeratorRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	moderatorService := CreateModeratorService(
		moderatorRepository,
		personalDataRepository,
		userRepository,
		departmentRepository,
	)
	err := moderatorService.CreateModerator(tu.TestInitModeratorData)
	if err != nil {
		t.Fatalf("Error creating moderator: %v", err)
	}
	err = moderatorService.UpdateModeratorPassword(2, tu.TestAuth, "test3")
	if err == nil {
		t.Fatalf("No error updating moderator password: %v", err)
	}
}

func TestUpdateModeratorPasswordIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up moderator tables: %v", err)
	}
	moderatorRepository := module.ModeratorRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	departmentRepository := module.DepartmentRepository
	authRepository := module.AuthRepository
	moderatorService := CreateModeratorService(
		moderatorRepository,
		personalDataRepository,
		userRepository,
		departmentRepository,
	)
	err = moderatorService.CreateModerator(tu.TestInitModeratorData)
	if err != nil {
		t.Fatalf("Error creating moderator: %v", err)
	}
	err = moderatorService.UpdateModeratorPassword(1, tu.TestAuth, "test3")
	if err == nil {
		t.Fatalf("No error updating moderator password: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	err = moderatorService.UpdateModeratorPassword(result.UserID+1, tu.TestAuth, "test3")
	if err == nil {
		t.Fatalf("No error updating moderator password: %v", err)
	}
}

func TestUpdateModeratorSalaryCorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	departmentRepository := tu.CreateTestDepartmentRepository()
	moderatorRepository := tu.CreateTestModeratorRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	moderatorService := CreateModeratorService(
		moderatorRepository,
		personalDataRepository,
		userRepository,
		departmentRepository,
	)
	err := moderatorService.CreateModerator(tu.TestInitModeratorData)
	if err != nil {
		t.Fatalf("Error creating moderator: %v", err)
	}
	newSalary := int64(150000)
	err = moderatorService.UpdateModeratorSalary(1, newSalary)
	if err != nil {
		t.Fatalf("Error updating moderator salary: %v", err)
	}
	moderatorData, err := moderatorService.GetModeratorData(1)
	if err != nil {
		t.Fatalf("Error getting moderator data: %v", err)
	}
	if moderatorData.Salary != newSalary {
		t.Fatalf("Moderator salary not updated: %v", moderatorData)
	}
}

func TestUpdateModeratorSalaryCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up moderator tables: %v", err)
	}
	moderatorRepository := module.ModeratorRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	departmentRepository := module.DepartmentRepository
	authRepository := module.AuthRepository
	moderatorService := CreateModeratorService(
		moderatorRepository,
		personalDataRepository,
		userRepository,
		departmentRepository,
	)
	err = moderatorService.CreateModerator(tu.TestInitModeratorData)
	if err != nil {
		t.Fatalf("Error creating moderator: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	err = moderatorService.UpdateModeratorSalary(result.UserID, 150000)
	if err != nil {
		t.Fatalf("Error updating moderator salary: %v", err)
	}
	moderatorData, err := moderatorService.GetModeratorData(result.UserID)
	if err != nil {
		t.Fatalf("Error getting moderator data: %v", err)
	}
	if moderatorData.Salary != 150000 {
		t.Fatalf("Moderator salary not updated: %v", moderatorData)
	}
}

func TestUpdateModeratorSalaryIncorrectIdLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	departmentRepository := tu.CreateTestDepartmentRepository()
	moderatorRepository := tu.CreateTestModeratorRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	moderatorService := CreateModeratorService(
		moderatorRepository,
		personalDataRepository,
		userRepository,
		departmentRepository,
	)
	err := moderatorService.CreateModerator(tu.TestInitModeratorData)
	if err != nil {
		t.Fatalf("Error creating moderator: %v", err)
	}
	err = moderatorService.UpdateModeratorSalary(2, 150000)
	if err == nil {
		t.Fatalf("No error updating moderator salary: %v", err)
	}
}

func TestUpdateModeratorSalaryIncorrectIdClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up moderator tables: %v", err)
	}
	moderatorRepository := module.ModeratorRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	departmentRepository := module.DepartmentRepository
	authRepository := module.AuthRepository
	moderatorService := CreateModeratorService(
		moderatorRepository,
		personalDataRepository,
		userRepository,
		departmentRepository,
	)
	err = moderatorService.CreateModerator(tu.TestInitModeratorData)
	if err != nil {
		t.Fatalf("Error creating moderator: %v", err)
	}
	err = moderatorService.UpdateModeratorSalary(1, 150000)
	if err == nil {
		t.Fatalf("No error updating moderator salary: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	err = moderatorService.UpdateModeratorSalary(result.UserID+1, 150000)
	if err == nil {
		t.Fatalf("No error updating moderator salary: %v", err)
	}
}

func TestUpdateModeratorSalaryIncorrectSalaryLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	departmentRepository := tu.CreateTestDepartmentRepository()
	moderatorRepository := tu.CreateTestModeratorRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	moderatorService := CreateModeratorService(
		moderatorRepository,
		personalDataRepository,
		userRepository,
		departmentRepository,
	)
	err := moderatorService.CreateModerator(tu.TestInitModeratorData)
	if err != nil {
		t.Fatalf("Error creating moderator: %v", err)
	}
	err = moderatorService.UpdateModeratorSalary(1, -500)
	if err == nil {
		t.Fatalf("No error updating moderator salary: %v", err)
	}
}

func TestUpdateModeratorSalaryIncorrectSalaryClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up moderator tables: %v", err)
	}
	moderatorRepository := module.ModeratorRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	departmentRepository := module.DepartmentRepository
	authRepository := module.AuthRepository
	moderatorService := CreateModeratorService(
		moderatorRepository,
		personalDataRepository,
		userRepository,
		departmentRepository,
	)
	err = moderatorService.CreateModerator(tu.TestInitModeratorData)
	if err != nil {
		t.Fatalf("Error creating moderator: %v", err)
	}
	err = moderatorService.UpdateModeratorSalary(1, -500)
	if err == nil {
		t.Fatalf("No error updating moderator salary: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	err = moderatorService.UpdateModeratorSalary(result.UserID, -500)
	if err == nil {
		t.Fatalf("No error updating moderator salary: %v", err)
	}
}
