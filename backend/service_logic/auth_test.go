package service_logic

import (
	"data_base_project/data_base"
	tu "data_base_project/test_service_utility"
	"data_base_project/types"
	"database/sql"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"

	_ "github.com/marcboeker/go-duckdb"
)

type AuthSuite struct {
	suite.Suite
}

func TestRunAuthSuite(t *testing.T) {
	suite.RunSuite(t, new(AuthSuite))
}

func (s *AuthSuite) TestAuthorizeCorrectLondon(t provider.T) {
	var (
		repo    data_base.IAuthRepository
		verdict types.ServiceAuthVerdict
	)
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		repo = tu.CreateTestAuthRepository()
		repo.InsertAuth(types.DBAuthInfo{UserID: 1, UserType: types.Admin, Login: tu.TestAuth.Login, Password: tu.TestAuth.Password})
	})
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		var err error
		verdict, err = CreateAuthService(repo).Authorize(types.ServiceAuthData{Login: tu.TestAuth.Login, Password: tu.TestAuth.Password})
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		sx.Assert().Equal(int64(1), verdict.UserID)
		sx.Assert().Equal(types.Admin, verdict.UserType)
	})
}

func (s *AuthSuite) TestAuthorizeCorrectClassic(t provider.T) {
	var (
		db      *sql.DB
		repo    data_base.IAuthRepository
		verdict types.ServiceAuthVerdict
	)
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		var err error
		db, err = sql.Open("duckdb", ":memory:")
		sx.Assert().NoError(err)
		t.Cleanup(func() { _ = db.Close() })
		mod, err := tu.SetupModule(db)
		sx.Assert().NoError(err)
		repo = mod.AuthRepository
		repo.InsertAuth(types.DBAuthInfo{UserID: 1, UserType: types.Admin, Login: tu.TestAuth.Login, Password: tu.TestAuth.Password})
	})
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		var err error
		verdict, err = CreateAuthService(repo).Authorize(tu.TestAuth)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		sx.Assert().Equal(int64(1), verdict.UserID)
		sx.Assert().Equal(types.Admin, verdict.UserType)
	})
}

func (s *AuthSuite) TestAuthorizeIncorrectLoginLondon(t provider.T) {
	var repo data_base.IAuthRepository
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		repo = tu.CreateTestAuthRepository()
		repo.InsertAuth(types.DBAuthInfo{UserID: 1, UserType: types.Admin, Login: tu.TestAuth.Login, Password: tu.TestAuth.Password})
	})
	t.WithNewStep("Act+Assert", func(sx provider.StepCtx) {
		bad := tu.TestAuth
		bad.Login = "incorrect"
		_, err := CreateAuthService(repo).Authorize(bad)
		sx.Assert().Error(err)
	})
}

func (s *AuthSuite) TestAuthorizeIncorrectClassic(t provider.T) {
	var db *sql.DB
	var repo data_base.IAuthRepository
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		var err error
		db, err = sql.Open("duckdb", ":memory:")
		sx.Assert().NoError(err)
		t.Cleanup(func() { _ = db.Close() })
		mod, err := tu.SetupModule(db)
		sx.Assert().NoError(err)
		repo = mod.AuthRepository
		repo.InsertAuth(types.DBAuthInfo{UserID: 1, UserType: types.Admin, Login: tu.TestAuth.Login, Password: tu.TestAuth.Password})
	})
	t.WithNewStep("Act+Assert", func(sx provider.StepCtx) {
		bad := tu.TestAuth
		bad.Login = "incorrect"
		_, err := CreateAuthService(repo).Authorize(bad)
		sx.Assert().Error(err)
	})
}

func (s *AuthSuite) TestAuthorizeIncorrectPasswordLondon(t provider.T) {
	var repo data_base.IAuthRepository
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		repo = tu.CreateTestAuthRepository()
		repo.InsertAuth(types.DBAuthInfo{UserID: 1, UserType: types.Admin, Login: tu.TestAuth.Login, Password: tu.TestAuth.Password})
	})
	t.WithNewStep("Act+Assert", func(sx provider.StepCtx) {
		bad := tu.TestAuth
		bad.Password = "incorrect"
		_, err := CreateAuthService(repo).Authorize(bad)
		sx.Assert().Error(err)
	})
}

func (s *AuthSuite) TestAuthorizeIncorrectPasswordClassic(t provider.T) {
	var db *sql.DB
	var repo data_base.IAuthRepository
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		var err error
		db, err = sql.Open("duckdb", ":memory:")
		sx.Assert().NoError(err)
		t.Cleanup(func() { _ = db.Close() })
		mod, err := tu.SetupModule(db)
		sx.Assert().NoError(err)
		repo = mod.AuthRepository
		repo.InsertAuth(types.DBAuthInfo{UserID: 1, UserType: types.Admin, Login: tu.TestAuth.Login, Password: tu.TestAuth.Password})
	})
	t.WithNewStep("Act+Assert", func(sx provider.StepCtx) {
		bad := tu.TestAuth
		bad.Password = "incorrect"
		_, err := CreateAuthService(repo).Authorize(bad)
		sx.Assert().Error(err)
	})
}

func (s *AuthSuite) TestCheckLoginCorrectLondon(t provider.T) {
	var repo data_base.IAuthRepository
	var exists bool
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		repo = tu.CreateTestAuthRepository()
		repo.InsertAuth(types.DBAuthInfo{UserID: 1, UserType: types.Admin, Login: tu.TestAuth.Login, Password: tu.TestAuth.Password})
	})
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		var err error
		exists, err = CreateAuthService(repo).CheckLogin(tu.TestAuth.Login)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		sx.Assert().True(exists)
	})
}

func (s *AuthSuite) TestCheckLoginCorrectClassic(t provider.T) {
	var db *sql.DB
	var repo data_base.IAuthRepository
	var exists bool
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		var err error
		db, err = sql.Open("duckdb", ":memory:")
		sx.Assert().NoError(err)
		t.Cleanup(func() { _ = db.Close() })
		mod, err := tu.SetupModule(db)
		sx.Assert().NoError(err)
		repo = mod.AuthRepository
		repo.InsertAuth(types.DBAuthInfo{UserID: 1, UserType: types.Admin, Login: tu.TestAuth.Login, Password: tu.TestAuth.Password})
	})
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		var err error
		exists, err = CreateAuthService(repo).CheckLogin(tu.TestAuth.Login)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		sx.Assert().True(exists)
	})
}

func (s *AuthSuite) TestCheckLoginIncorrectLondon(t provider.T) {
	var exists bool
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		var err error
		exists, err = CreateAuthService(tu.CreateTestAuthRepository()).CheckLogin("incorrect")
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		sx.Assert().False(exists)
	})
}

func (s *AuthSuite) TestCheckLoginIncorrectClassic(t provider.T) {
	var db *sql.DB
	var exists bool
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		var err error
		db, err = sql.Open("duckdb", ":memory:")
		sx.Assert().NoError(err)
		t.Cleanup(func() { _ = db.Close() })
	})
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		mod, err := tu.SetupModule(db)
		sx.Assert().NoError(err)
		service := CreateAuthService(mod.AuthRepository)
		exists, err = service.CheckLogin("incorrect")
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		sx.Assert().False(exists)
	})
}
