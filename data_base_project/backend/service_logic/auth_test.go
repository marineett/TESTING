package service_logic

import (
	tu "data_base_project/test_service_utility"
	"data_base_project/types"
	"database/sql"
	"testing"

	_ "github.com/marcboeker/go-duckdb"
)

func TestAuthorizeCorrectLondon(t *testing.T) {
	authRepository := tu.CreateTestAuthRepository()
	authService := CreateAuthService(authRepository)
	authRepository.InsertAuth(types.DBAuthInfo{
		UserID:   1,
		UserType: types.Admin,
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	verdict, err := authService.Authorize(tu.TestAuth)
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	if verdict.UserID != 1 {
		t.Fatalf("User id not updated: %v", verdict)
	}
	if verdict.UserType != types.Admin {
		t.Fatalf("User type not updated: %v", verdict)
	}
}

func TestAuthorizeCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up auth tables: %v", err)
	}
	authRepository := module.AuthRepository
	authService := CreateAuthService(authRepository)
	authRepository.InsertAuth(types.DBAuthInfo{
		UserID:   1,
		UserType: types.Admin,
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	verdict, err := authService.Authorize(tu.TestAuth)
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	if verdict.UserID != 1 {
		t.Fatalf("User id not updated: %v", verdict)
	}
	if verdict.UserType != types.Admin {
		t.Fatalf("User type not updated: %v", verdict)
	}
}

func TestAuthorizeIncorrectLoginLondon(t *testing.T) {
	authRepository := tu.CreateTestAuthRepository()
	authService := CreateAuthService(authRepository)
	authRepository.InsertAuth(types.DBAuthInfo{
		UserID:   1,
		UserType: types.Admin,
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	newAuth := tu.TestAuth
	newAuth.Login = "incorrect"
	verdict, err := authService.Authorize(newAuth)
	if err == nil {
		t.Fatalf("No error authorizing: %v", verdict)
	}
}

func TestAuthorizeIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up auth tables: %v", err)
	}
	authRepository := module.AuthRepository
	authService := CreateAuthService(authRepository)
	authRepository.InsertAuth(types.DBAuthInfo{
		UserID:   1,
		UserType: types.Admin,
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	newAuth := tu.TestAuth
	newAuth.Login = "incorrect"
	verdict, err := authService.Authorize(newAuth)
	if err == nil {
		t.Fatalf("No error authorizing: %v", verdict)
	}
}

func TestAuthorizeIncorrectPasswordLondon(t *testing.T) {
	authRepository := tu.CreateTestAuthRepository()
	authService := CreateAuthService(authRepository)
	authRepository.InsertAuth(types.DBAuthInfo{
		UserID:   1,
		UserType: types.Admin,
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	newAuth := tu.TestAuth
	newAuth.Password = "incorrect"
	verdict, err := authService.Authorize(newAuth)
	if err == nil {
		t.Fatalf("No error authorizing: %v", verdict)
	}
}

func TestAuthorizeIncorrectPasswordClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up auth tables: %v", err)
	}
	authRepository := module.AuthRepository
	authService := CreateAuthService(authRepository)
	authRepository.InsertAuth(types.DBAuthInfo{
		UserID:   1,
		UserType: types.Admin,
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	newAuth := tu.TestAuth
	newAuth.Password = "incorrect"
	verdict, err := authService.Authorize(newAuth)
	if err == nil {
		t.Fatalf("No error authorizing: %v", verdict)
	}
}

func TestCheckLoginCorrectLondon(t *testing.T) {
	authRepository := tu.CreateTestAuthRepository()
	authService := CreateAuthService(authRepository)
	authRepository.InsertAuth(types.DBAuthInfo{
		UserID:   1,
		UserType: types.Admin,
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	loginExists, err := authService.CheckLogin(tu.TestAuth.Login)
	if err != nil {
		t.Fatalf("Error checking login: %v", err)
	}
	if !loginExists {
		t.Fatalf("Login not found: %v", loginExists)
	}
}

func TestCheckLoginCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up auth tables: %v", err)
	}
	authRepository := module.AuthRepository
	authService := CreateAuthService(authRepository)
	authRepository.InsertAuth(types.DBAuthInfo{
		UserID:   1,
		UserType: types.Admin,
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	loginExists, err := authService.CheckLogin(tu.TestAuth.Login)
	if err != nil {
		t.Fatalf("Error checking login: %v", err)
	}
	if !loginExists {
		t.Fatalf("Login not found: %v", loginExists)
	}
}

func TestCheckLoginIncorrectLondon(t *testing.T) {
	authRepository := tu.CreateTestAuthRepository()
	authService := CreateAuthService(authRepository)
	loginExists, err := authService.CheckLogin("incorrect")
	if err != nil {
		t.Fatalf("Error checking login: %v", err)
	}
	if loginExists {
		t.Fatalf("Login found: %v", loginExists)
	}
}

func TestCheckLoginIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up auth tables: %v", err)
	}
	authRepository := module.AuthRepository
	authService := CreateAuthService(authRepository)
	loginExists, err := authService.CheckLogin("incorrect")
	if err != nil {
		t.Fatalf("Error checking login: %v", err)
	}
	if loginExists {
		t.Fatalf("Login found: %v", loginExists)
	}
}
