package service_logic

import (
	tu "data_base_project/test_service_utility"
	"data_base_project/types"
	"database/sql"
	"testing"

	_ "github.com/marcboeker/go-duckdb"
)

func TestCreateAdminCorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	adminRepository := tu.CreateTestAdminRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	err := adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
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
	adminData, err := adminRepository.GetAdmin(1)
	if err != nil {
		t.Fatalf("Error getting admin data: %v", err)
	}
	if adminData.Salary != tu.TestSalary {
		t.Fatalf("Admin data not updated: %v", adminData)
	}
}

func TestCreateAdminCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	repositoryModule, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up tables: %v", err)
	}
	adminRepository := repositoryModule.AdminRepository
	personalDataRepository := repositoryModule.PersonalDataRepository
	userRepository := repositoryModule.UserRepository
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	err = adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
}

func TestCreateAdminIncorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	adminRepository := tu.CreateTestAdminRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	err := adminService.CreateAdmin(types.ServiceInitAdminData{
		ServiceInitUserData: types.ServiceInitUserData{
			ServicePersonalData: tu.TestPD,
			ServiceAuthData:     tu.TestAuth,
		},
		Salary: -500,
	})
	if err == nil {
		t.Fatalf("No error creating admin with salary -500: %v", err)
	}
}

func TestCreateAdminIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	repositoryModule, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up tables: %v", err)
	}
	adminRepository := repositoryModule.AdminRepository
	personalDataRepository := repositoryModule.PersonalDataRepository
	userRepository := repositoryModule.UserRepository
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	err = adminService.CreateAdmin(types.ServiceInitAdminData{
		ServiceInitUserData: types.ServiceInitUserData{
			ServicePersonalData: tu.TestPD,
			ServiceAuthData:     tu.TestAuth,
		},
		Salary: -500,
	})
	if err == nil {
		t.Fatalf("No error creating admin with salary -500: %v", err)
	}
}

func TestGetAdminDataCorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	adminRepository := tu.CreateTestAdminRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	err := adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	adminData, err := adminService.GetAdminData(1)
	if err != nil {
		t.Fatalf("Error getting admin data: %v", err)
	}
	if adminData.Salary != tu.TestSalary {
		t.Fatalf("Admin data not updated: %v", adminData)
	}
	if adminData.DepartmentID != 0 {
		t.Fatalf("Admin data not updated: %v", adminData)
	}
	if adminData.ID != 1 {
		t.Fatalf("Admin data not updated: %v", adminData)
	}
}

func TestGetAdminDataCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	repositoryModule, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up tables: %v", err)
	}
	adminRepository := repositoryModule.AdminRepository
	personalDataRepository := repositoryModule.PersonalDataRepository
	authRepository := repositoryModule.AuthRepository
	userRepository := repositoryModule.UserRepository
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	err = adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	adminData, err := adminRepository.GetAdmin(result.UserID)
	if err != nil {
		t.Fatalf("Error getting admin data: %v", err)
	}
	if adminData.Salary != tu.TestSalary {
		t.Fatalf("Admin data not updated: %v", adminData)
	}
	if adminData.DepartmentID != 0 {
		t.Fatalf("Admin data not updated: %v", adminData)
	}
	if adminData.ID != result.UserID {
		t.Fatalf("Admin data not updated: %v", adminData)
	}
}

func TestGetAdminDataInCorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	adminRepository := tu.CreateTestAdminRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	err := adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	adminData, err := adminService.GetAdminData(2)
	if err == nil {
		t.Fatalf("No error getting admin data: %v", adminData)
	}
}

func TestGetAdminDataIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	repositoryModule, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up tables: %v", err)
	}
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := repositoryModule.AuthRepository
	userRepository := tu.CreateTestUserRepository()
	adminRepository := tu.CreateTestAdminRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	err = adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	adminData, err := adminRepository.GetAdmin(result.UserID + 1)
	if err == nil {
		t.Fatalf("No error getting admin data: %v", adminData)
	}
}

func TestGetAdminProfileCorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	adminRepository := tu.CreateTestAdminRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	adminService.CreateAdmin(tu.TestInitAdminData)
	adminProfile, err := adminService.GetAdminProfile(1)
	if err != nil {
		t.Fatalf("Error getting admin profile: %v", err)
	}
	if adminProfile.Salary != tu.TestSalary {
		t.Fatalf("Admin profile not updated: %v", adminProfile)
	}
	if adminProfile.FirstName != tu.TestPD.FirstName {
		t.Fatalf("Admin profile not updated: %v", adminProfile)
	}
	if adminProfile.LastName != tu.TestPD.LastName {
		t.Fatalf("Admin profile not updated: %v", adminProfile)
	}
	if adminProfile.MiddleName != tu.TestPD.MiddleName {
		t.Fatalf("Admin profile not updated: %v", adminProfile)
	}
	if adminProfile.TelephoneNumber != tu.TestPD.TelephoneNumber {
		t.Fatalf("Admin profile not updated: %v", adminProfile)
	}
	if adminProfile.Email != tu.TestPD.Email {
		t.Fatalf("Admin profile not updated: %v", adminProfile)
	}
}

func TestGetAdminProfileCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	repositoryModule, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up tables: %v", err)
	}
	adminRepository := repositoryModule.AdminRepository
	personalDataRepository := repositoryModule.PersonalDataRepository
	authRepository := repositoryModule.AuthRepository
	userRepository := repositoryModule.UserRepository
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	err = adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	adminProfile, err := adminService.GetAdminProfile(result.UserID)
	if err != nil {
		t.Fatalf("Error getting admin profile: %v", err)
	}
	if adminProfile.Salary != tu.TestSalary {
		t.Fatalf("Admin profile not updated: %v", adminProfile)
	}
	if adminProfile.FirstName != tu.TestPD.FirstName {
		t.Fatalf("Admin profile not updated: %v", adminProfile)
	}
	if adminProfile.LastName != tu.TestPD.LastName {
		t.Fatalf("Admin profile not updated: %v", adminProfile)
	}
	if adminProfile.MiddleName != tu.TestPD.MiddleName {
		t.Fatalf("Admin profile not updated: %v", adminProfile)
	}
	if adminProfile.TelephoneNumber != tu.TestPD.TelephoneNumber {
		t.Fatalf("Admin profile not updated: %v", adminProfile)
	}
	if adminProfile.Email != tu.TestPD.Email {
		t.Fatalf("Admin profile not updated: %v", adminProfile)
	}
}

func TestGetAdminProfileIncorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	adminRepository := tu.CreateTestAdminRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	err := adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	adminProfile, err := adminService.GetAdminProfile(2)
	if err == nil {
		t.Fatalf("No error getting admin profile: %v", adminProfile)
	}
}

func TestGetAdminProfileIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	repositoryModule, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up tables: %v", err)
	}
	adminRepository := repositoryModule.AdminRepository
	personalDataRepository := repositoryModule.PersonalDataRepository
	authRepository := repositoryModule.AuthRepository
	userRepository := repositoryModule.UserRepository
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	err = adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	adminProfile, err := adminService.GetAdminProfile(result.UserID + 1)
	if err == nil {
		t.Fatalf("No error getting admin profile: %v", adminProfile)
	}
}

func TestUpdateAdminPersonalDataCorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	adminRepository := tu.CreateTestAdminRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	err := adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	newPersonalData := types.ServicePersonalData{
		FirstName:       "Petr",
		LastName:        "Petrov",
		MiddleName:      "Petrovich",
		TelephoneNumber: "88005553536",
		Email:           "test2@test.com",
	}
	err = adminService.UpdateAdminPersonalData(1, newPersonalData)
	if err != nil {
		t.Fatalf("Error updating admin personal data: %v", err)
	}
	personalData, err := personalDataRepository.GetPersonalData(1)
	if err != nil {
		t.Fatalf("Error getting personal data: %v", err)
	}
	if personalData.FirstName != newPersonalData.FirstName {
		t.Fatalf("Admin personal data not updated: %v", personalData)
	}
	if personalData.LastName != newPersonalData.LastName {
		t.Fatalf("Admin personal data not updated: %v", personalData)
	}
	if personalData.MiddleName != newPersonalData.MiddleName {
		t.Fatalf("Admin personal data not updated: %v", personalData)
	}
	if personalData.TelephoneNumber != newPersonalData.TelephoneNumber {
		t.Fatalf("Admin personal data not updated: %v", personalData)
	}
	if personalData.Email != newPersonalData.Email {
		t.Fatalf("Admin personal data not updated: %v", personalData)
	}
}

func TestUpdateAdminPersonalDataCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	repositoryModule, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up tables: %v", err)
	}
	adminRepository := repositoryModule.AdminRepository
	personalDataRepository := repositoryModule.PersonalDataRepository
	userRepository := repositoryModule.UserRepository
	authRepository := repositoryModule.AuthRepository
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)

	err = adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	err = adminService.UpdateAdminPersonalData(result.UserID, tu.TestPD)
	if err != nil {
		t.Fatalf("Error updating admin personal data: %v", err)
	}
	personalData, err := personalDataRepository.GetPersonalData(1)
	if err != nil {
		t.Fatalf("Error getting personal data: %v", err)
	}
	if personalData.FirstName != tu.TestPD.FirstName {
		t.Fatalf("Admin personal data not updated: %v", personalData)
	}
	if personalData.LastName != tu.TestPD.LastName {
		t.Fatalf("Admin personal data not updated: %v", personalData)
	}
	if personalData.MiddleName != tu.TestPD.MiddleName {
		t.Fatalf("Admin personal data not updated: %v", personalData)
	}
	if personalData.TelephoneNumber != tu.TestPD.TelephoneNumber {
		t.Fatalf("Admin personal data not updated: %v", personalData)
	}
	if personalData.Email != tu.TestPD.Email {
		t.Fatalf("Admin personal data not updated: %v", personalData)
	}
}

func TestUpdateAdminPersonalDataIncorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	adminRepository := tu.CreateTestAdminRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	err := adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	err = adminService.UpdateAdminPersonalData(2, types.ServicePersonalData{
		FirstName:       "Petr",
		LastName:        "Petrov",
		MiddleName:      "Petrovich",
		TelephoneNumber: "88005553536",
		Email:           "test2@test.com",
	})
	if err == nil {
		t.Fatalf("No error updating admin personal data: %v", err)
	}
}

func TestUpdateAdminPersonalDataIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	repositoryModule, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up tables: %v", err)
	}
	adminRepository := repositoryModule.AdminRepository
	personalDataRepository := repositoryModule.PersonalDataRepository
	userRepository := repositoryModule.UserRepository
	authRepository := repositoryModule.AuthRepository
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	err = adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	err = adminService.UpdateAdminPersonalData(result.UserID+1, tu.TestPD)
	if err == nil {
		t.Fatalf("No error updating admin personal data: %v", err)
	}
}

func TestUpdateAdminPasswordCorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	adminRepository := tu.CreateTestAdminRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	err := adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	newPassword := "test3"
	err = adminService.UpdateAdminPassword(1, tu.TestAuth, newPassword)
	if err != nil {
		t.Fatalf("Error updating admin password: %v", err)
	}
	authData, err := authRepository.TestGetAuth(1)
	if err != nil {
		t.Fatalf("Error getting auth data: %v", err)
	}
	if authData.Password != newPassword {
		t.Fatalf("Admin password not updated: %v", authData)
	}
}

func TestUpdateAdminPasswordCorrectClassic(t *testing.T) {

	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	repositoryModule, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up tables: %v", err)
	}
	adminRepository := repositoryModule.AdminRepository
	personalDataRepository := repositoryModule.PersonalDataRepository
	userRepository := repositoryModule.UserRepository
	authRepository := repositoryModule.AuthRepository
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	err = adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	err = adminService.UpdateAdminPassword(result.UserID, tu.TestAuth, "test3")
	if err != nil {
		t.Fatalf("Error updating admin password: %v", err)
	}
	result, err = authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: "test3",
	})
	if err != nil {
		t.Fatalf("Password not updated: %v", err)
	}
}

func TestUpdateAdminPasswordIncorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	adminRepository := tu.CreateTestAdminRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	err := adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	err = adminService.UpdateAdminPassword(2, tu.TestAuth, "test3")
	if err == nil {
		t.Fatalf("No error updating admin password: %v", err)
	}
}

func TestUpdateAdminPasswordIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	repositoryModule, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up tables: %v", err)
	}
	adminRepository := repositoryModule.AdminRepository
	personalDataRepository := repositoryModule.PersonalDataRepository
	userRepository := repositoryModule.UserRepository
	authRepository := repositoryModule.AuthRepository
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	err = adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	err = adminService.UpdateAdminPassword(result.UserID+1, tu.TestAuth, "test3")
	if err == nil {
		t.Fatalf("No error updating admin password: %v", err)
	}
}

func TestUpdateAdminDepartmentCorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	adminRepository := tu.CreateTestAdminRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	err := adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	newDepartmentID := int64(5)
	err = adminService.UpdateAdminDepartment(1, newDepartmentID)
	if err != nil {
		t.Fatalf("Error updating admin department: %v", err)
	}
	adminData, err := adminService.GetAdminData(1)
	if err != nil {
		t.Fatalf("Error getting admin data: %v", err)
	}
	if adminData.DepartmentID != newDepartmentID {
		t.Fatalf("Admin department not updated: %v", adminData)
	}
}

func TestUpdateAdminDepartmentCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	repositoryModule, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up tables: %v", err)
	}
	adminRepository := repositoryModule.AdminRepository
	personalDataRepository := repositoryModule.PersonalDataRepository
	userRepository := repositoryModule.UserRepository
	authRepository := repositoryModule.AuthRepository
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	err = adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	err = adminService.UpdateAdminDepartment(result.UserID, 5)
	if err != nil {
		t.Fatalf("Error updating admin department: %v", err)
	}
	adminData, err := adminService.GetAdminData(result.UserID)
	if err != nil {
		t.Fatalf("Error getting admin data: %v", err)
	}
	if adminData.DepartmentID != 5 {
		t.Fatalf("Admin department not updated: %v", adminData)
	}
}

func TestUpdateAdminDepartmentIncorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	adminRepository := tu.CreateTestAdminRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	err := adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	err = adminService.UpdateAdminDepartment(2, 5)
	if err == nil {
		t.Fatalf("No error updating admin department: %v", err)
	}
}

func TestUpdateAdminDepartmentIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	repositoryModule, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up tables: %v", err)
	}
	adminRepository := repositoryModule.AdminRepository
	personalDataRepository := repositoryModule.PersonalDataRepository
	userRepository := repositoryModule.UserRepository
	authRepository := repositoryModule.AuthRepository
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	err = adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	err = adminService.UpdateAdminDepartment(result.UserID+1, 5)
	if err == nil {
		t.Fatalf("No error updating admin department: %v", err)
	}
}

func TestUpdateAdminSalaryCorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	adminRepository := tu.CreateTestAdminRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	err := adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	newSalary := int64(150000)
	err = adminService.UpdateAdminSalary(1, newSalary)
	if err != nil {
		t.Fatalf("Error updating admin salary: %v", err)
	}
	adminData, err := adminService.GetAdminData(1)
	if err != nil {
		t.Fatalf("Error getting admin data: %v", err)
	}
	if adminData.Salary != newSalary {
		t.Fatalf("Admin salary not updated: %v", adminData)
	}
}

func TestUpdateAdminSalaryCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	repositoryModule, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up tables: %v", err)
	}
	adminRepository := repositoryModule.AdminRepository
	personalDataRepository := repositoryModule.PersonalDataRepository
	userRepository := repositoryModule.UserRepository
	authRepository := repositoryModule.AuthRepository
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	err = adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	err = adminService.UpdateAdminSalary(result.UserID, 150000)
	if err != nil {
		t.Fatalf("Error updating admin salary: %v", err)
	}
	adminData, err := adminService.GetAdminData(result.UserID)
	if err != nil {
		t.Fatalf("Error getting admin data: %v", err)
	}
	if adminData.Salary != 150000 {
		t.Fatalf("Admin salary not updated: %v", adminData)
	}
}

func TestUpdateAdminSalaryIncorrectLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	adminRepository := tu.CreateTestAdminRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	err := adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	err = adminService.UpdateAdminSalary(2, 150000)
	if err == nil {
		t.Fatalf("No error updating admin salary: %v", err)
	}
}

func TestUpdateAdminSalaryIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	repositoryModule, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up tables: %v", err)
	}
	adminRepository := repositoryModule.AdminRepository
	personalDataRepository := repositoryModule.PersonalDataRepository
	userRepository := repositoryModule.UserRepository
	authRepository := repositoryModule.AuthRepository
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	err = adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	err = adminService.UpdateAdminSalary(result.UserID+1, 150000)
	if err == nil {
		t.Fatalf("No error updating admin salary: %v", err)
	}
}

func TestUpdateAdminSalaryIncorrectSalaryLondon(t *testing.T) {
	personalDataRepository := tu.CreateTestPersonalDataRepository()
	authRepository := tu.CreateTestAuthRepository()
	userRepository := tu.CreateTestUserRepository()
	adminRepository := tu.CreateTestAdminRepository(
		personalDataRepository,
		authRepository,
		userRepository,
	)
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	err := adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	err = adminService.UpdateAdminSalary(1, -500)
	if err == nil {
		t.Fatalf("No error updating admin salary: %v", err)
	}
	adminData, err := adminService.GetAdminData(1)
	if err != nil {
		t.Fatalf("Error getting admin data: %v", err)
	}
	if adminData.Salary != tu.TestSalary {
		t.Fatalf("Admin salary not updated: %v", adminData)
	}
}

func TestUpdateAdminSalaryIncorrectSalaryClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	repositoryModule, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up tables: %v", err)
	}
	adminRepository := repositoryModule.AdminRepository
	personalDataRepository := repositoryModule.PersonalDataRepository
	userRepository := repositoryModule.UserRepository
	authRepository := repositoryModule.AuthRepository
	adminService := CreateAdminService(
		adminRepository,
		userRepository,
		personalDataRepository,
	)
	err = adminService.CreateAdmin(tu.TestInitAdminData)
	if err != nil {
		t.Fatalf("Error creating admin: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})

	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	err = adminService.UpdateAdminSalary(result.UserID, -500)
	if err == nil {
		t.Fatalf("No error updating admin salary: %v", err)
	}
}
