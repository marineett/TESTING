package service_logic

import (
	"data_base_project/data_base"
	tu "data_base_project/test_service_utility"
	"data_base_project/types"
	"database/sql"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type ClientSuite struct {
	suite.Suite
}

func TestRunClientSuite(t *testing.T) {
	suite.RunSuite(t, new(ClientSuite))
}

func (s *ClientSuite) TestCreateClientCorrectLondon(t provider.T) {
	var (
		pdRepo = tu.CreateTestPersonalDataRepository()
		aRepo  = tu.CreateTestAuthRepository()
		uRepo  = tu.CreateTestUserRepository()
		rRepo  = tu.CreateTestReviewRepository()
		cRepo  = tu.CreateTestClientRepository(pdRepo, aRepo, uRepo)
	)
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		err := CreateClientService(cRepo, pdRepo, uRepo, rRepo).CreateClient(tu.TestInitClientData)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		pd, err := pdRepo.GetPersonalData(1)
		sx.Assert().NoError(err)
		sx.Assert().Equal(tu.TestPD.TelephoneNumber, pd.TelephoneNumber)
		sx.Assert().Equal(tu.TestPD.Email, pd.Email)
		auth, err := aRepo.TestGetAuth(1)
		sx.Assert().NoError(err)
		sx.Assert().Equal(tu.TestAuth.Login, auth.Login)
		sx.Assert().Equal(tu.TestAuth.Password, auth.Password)
		cd, err := cRepo.GetClient(1)
		sx.Assert().NoError(err)
		sx.Assert().Equal(tu.TestSummaryRating, cd.SummaryRating)
	})
}

func (s *ClientSuite) TestCreateClientCorrectClassic(t provider.T) {
	var (
		db  *sql.DB
		mod *data_base.DataBaseModule
	)
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		var err error
		db, err = sql.Open("duckdb", ":memory:")
		sx.Assert().NoError(err)
		t.Cleanup(func() { _ = db.Close() })
		mod, err = tu.SetupModule(db)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		err := CreateClientService(mod.ClientRepository, mod.PersonalDataRepository, mod.UserRepository, mod.ReviewRepository).CreateClient(tu.TestInitClientData)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		res, err := mod.AuthRepository.Authorize(types.DBAuthData{Login: tu.TestAuth.Login, Password: tu.TestAuth.Password})
		sx.Assert().NoError(err)
		cd, err := mod.ClientRepository.GetClient(res.UserID)
		sx.Assert().NoError(err)
		sx.Assert().Equal(tu.TestSummaryRating, cd.SummaryRating)
	})
}

func (s *ClientSuite) TestGetClientDataCorrectLondon(t provider.T) {
	var (
		pdRepo = tu.CreateTestPersonalDataRepository()
		aRepo  = tu.CreateTestAuthRepository()
		uRepo  = tu.CreateTestUserRepository()
		rRepo  = tu.CreateTestReviewRepository()
		cRepo  = tu.CreateTestClientRepository(pdRepo, aRepo, uRepo)
	)
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		err := CreateClientService(cRepo, pdRepo, uRepo, rRepo).CreateClient(tu.TestInitClientData)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		_, err := CreateClientService(cRepo, pdRepo, uRepo, rRepo).GetClientProfile(1, 0, 10)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		cp, err := CreateClientService(cRepo, pdRepo, uRepo, rRepo).GetClientProfile(1, 0, 10)
		sx.Assert().NoError(err)
		sx.Assert().Equal(tu.TestMeanRating, cp.MeanRating)
		sx.Assert().Equal(tu.TestPD.FirstName, cp.FirstName)
		sx.Assert().Equal(tu.TestPD.LastName, cp.LastName)
		sx.Assert().Equal(tu.TestPD.MiddleName, cp.MiddleName)
		sx.Assert().Equal(tu.TestPD.TelephoneNumber, cp.TelephoneNumber)
		sx.Assert().Equal(tu.TestPD.Email, cp.Email)
	})
}

func (s *ClientSuite) TestGetClientDataCorrectClassic(t provider.T) {
	var (
		db  *sql.DB
		mod *data_base.DataBaseModule
		uid int64
	)
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		var err error
		db, err = sql.Open("duckdb", ":memory:")
		sx.Assert().NoError(err)
		t.Cleanup(func() { _ = db.Close() })
		mod, err = tu.SetupModule(db)
		sx.Assert().NoError(err)
		err = CreateClientService(mod.ClientRepository, mod.PersonalDataRepository, mod.UserRepository, mod.ReviewRepository).CreateClient(tu.TestInitClientData)
		sx.Assert().NoError(err)
		res, err := mod.AuthRepository.Authorize(types.DBAuthData{Login: tu.TestAuth.Login, Password: tu.TestAuth.Password})
		sx.Assert().NoError(err)
		uid = res.UserID
	})
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		_, err := CreateClientService(mod.ClientRepository, mod.PersonalDataRepository, mod.UserRepository, mod.ReviewRepository).GetClientProfile(uid, 0, 10)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		cp, err := CreateClientService(mod.ClientRepository, mod.PersonalDataRepository, mod.UserRepository, mod.ReviewRepository).GetClientProfile(uid, 0, 10)
		sx.Assert().NoError(err)
		sx.Assert().Equal(tu.TestMeanRating, cp.MeanRating)
		sx.Assert().Equal(tu.TestPD.FirstName, cp.FirstName)
		sx.Assert().Equal(tu.TestPD.LastName, cp.LastName)
		sx.Assert().Equal(tu.TestPD.MiddleName, cp.MiddleName)
		sx.Assert().Equal(tu.TestPD.TelephoneNumber, cp.TelephoneNumber)
		sx.Assert().Equal(tu.TestPD.Email, cp.Email)
	})
}

func (s *ClientSuite) TestGetClientDataIncorrectLondon(t provider.T) {
	var (
		pdRepo = tu.CreateTestPersonalDataRepository()
		aRepo  = tu.CreateTestAuthRepository()
		uRepo  = tu.CreateTestUserRepository()
		rRepo  = tu.CreateTestReviewRepository()
		cRepo  = tu.CreateTestClientRepository(pdRepo, aRepo, uRepo)
	)
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		err := CreateClientService(cRepo, pdRepo, uRepo, rRepo).CreateClient(tu.TestInitClientData)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Act+Assert", func(sx provider.StepCtx) {
		_, err := CreateClientService(cRepo, pdRepo, uRepo, rRepo).GetClientProfile(2, 0, 10)
		sx.Assert().Error(err)
	})
}

func (s *ClientSuite) TestGetClientDataIncorrectClassic(t provider.T) {
	var db *sql.DB
	var mod *data_base.DataBaseModule
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		var err error
		db, err = sql.Open("duckdb", ":memory:")
		sx.Assert().NoError(err)
		t.Cleanup(func() { _ = db.Close() })
		var e error
		mod, e = tu.SetupModule(db)
		sx.Assert().NoError(e)
	})
	t.WithNewStep("Act+Assert", func(sx provider.StepCtx) {
		_, err := CreateClientService(mod.ClientRepository, mod.PersonalDataRepository, mod.UserRepository, mod.ReviewRepository).GetClientProfile(1, 0, 10)
		sx.Assert().Error(err)
	})
}

func (s *ClientSuite) TestUpdateClientPersonalDataCorrectLondon(t provider.T) {
	var (
		pdRepo = tu.CreateTestPersonalDataRepository()
		aRepo  = tu.CreateTestAuthRepository()
		uRepo  = tu.CreateTestUserRepository()
		rRepo  = tu.CreateTestReviewRepository()
		cRepo  = tu.CreateTestClientRepository(pdRepo, aRepo, uRepo)
	)
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		err := CreateClientService(cRepo, pdRepo, uRepo, rRepo).CreateClient(tu.TestInitClientData)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		newPD := types.ServicePersonalData{FirstName: "Petr", LastName: "Petrov", MiddleName: "Petrovich", TelephoneNumber: "88005553536", Email: "test2@test.com"}
		err := CreateClientService(cRepo, pdRepo, uRepo, rRepo).UpdateClientPersonalData(1, newPD)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		pd, err := pdRepo.GetPersonalData(1)
		sx.Assert().NoError(err)
		sx.Assert().Equal("Petr", pd.FirstName)
		sx.Assert().Equal("Petrov", pd.LastName)
		sx.Assert().Equal("Petrovich", pd.MiddleName)
		sx.Assert().Equal("88005553536", pd.TelephoneNumber)
		sx.Assert().Equal("test2@test.com", pd.Email)
	})
}

func (s *ClientSuite) TestUpdateClientPersonalDataCorrectClassic(t provider.T) {
	var (
		db  *sql.DB
		mod *data_base.DataBaseModule
		uid int64
	)
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		var err error
		db, err = sql.Open("duckdb", ":memory:")
		sx.Assert().NoError(err)
		t.Cleanup(func() { _ = db.Close() })
		mod, err = tu.SetupModule(db)
		sx.Assert().NoError(err)
		err = CreateClientService(mod.ClientRepository, mod.PersonalDataRepository, mod.UserRepository, mod.ReviewRepository).CreateClient(tu.TestInitClientData)
		sx.Assert().NoError(err)
		res, err := mod.AuthRepository.Authorize(types.DBAuthData{Login: tu.TestAuth.Login, Password: tu.TestAuth.Password})
		sx.Assert().NoError(err)
		uid = res.UserID
	})
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		err := CreateClientService(mod.ClientRepository, mod.PersonalDataRepository, mod.UserRepository, mod.ReviewRepository).UpdateClientPersonalData(uid, types.ServicePersonalData{FirstName: "Petr", LastName: "Petrov", MiddleName: "Petrovich", TelephoneNumber: "88005553536", Email: "test2@test.com"})
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		pd, err := mod.PersonalDataRepository.GetPersonalData(1)
		sx.Assert().NoError(err)
		sx.Assert().Equal("Petr", pd.FirstName)
		sx.Assert().Equal("Petrov", pd.LastName)
		sx.Assert().Equal("Petrovich", pd.MiddleName)
		sx.Assert().Equal("88005553536", pd.TelephoneNumber)
		sx.Assert().Equal("test2@test.com", pd.Email)
	})
}

func (s *ClientSuite) TestUpdateClientPersonalDataIncorrectLondon(t provider.T) {
	var (
		pdRepo = tu.CreateTestPersonalDataRepository()
		aRepo  = tu.CreateTestAuthRepository()
		uRepo  = tu.CreateTestUserRepository()
		rRepo  = tu.CreateTestReviewRepository()
		cRepo  = tu.CreateTestClientRepository(pdRepo, aRepo, uRepo)
	)
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		err := CreateClientService(cRepo, pdRepo, uRepo, rRepo).CreateClient(tu.TestInitClientData)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Act+Assert", func(sx provider.StepCtx) {
		err := CreateClientService(cRepo, pdRepo, uRepo, rRepo).UpdateClientPersonalData(2, types.ServicePersonalData{FirstName: "Petr", LastName: "Petrov", MiddleName: "Petrovich", TelephoneNumber: "88005553536", Email: "test2@test.com"})
		sx.Assert().Error(err)
	})
}

func (s *ClientSuite) TestUpdateClientPersonalDataIncorrectClassic(t provider.T) {
	var db *sql.DB
	var mod *data_base.DataBaseModule
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		var err error
		db, err = sql.Open("duckdb", ":memory:")
		sx.Assert().NoError(err)
		t.Cleanup(func() { _ = db.Close() })
		mod, err = tu.SetupModule(db)
		sx.Assert().NoError(err)
		err = CreateClientService(mod.ClientRepository, mod.PersonalDataRepository, mod.UserRepository, mod.ReviewRepository).CreateClient(tu.TestInitClientData)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Act+Assert", func(sx provider.StepCtx) {
		err := CreateClientService(mod.ClientRepository, mod.PersonalDataRepository, mod.UserRepository, mod.ReviewRepository).UpdateClientPersonalData(1, types.ServicePersonalData{FirstName: "Petr", LastName: "Petrov", MiddleName: "Petrovich", TelephoneNumber: "88005553536", Email: "test2@test.com"})
		sx.Assert().Error(err)
	})
}

func (s *ClientSuite) TestUpdateClientPasswordCorrectLondon(t provider.T) {
	var (
		pdRepo = tu.CreateTestPersonalDataRepository()
		aRepo  = tu.CreateTestAuthRepository()
		uRepo  = tu.CreateTestUserRepository()
		rRepo  = tu.CreateTestReviewRepository()
		cRepo  = tu.CreateTestClientRepository(pdRepo, aRepo, uRepo)
	)
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		err := CreateClientService(cRepo, pdRepo, uRepo, rRepo).CreateClient(tu.TestInitClientData)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		err := CreateClientService(cRepo, pdRepo, uRepo, rRepo).UpdateClientPassword(1, tu.TestAuth, "test3")
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		auth, err := aRepo.TestGetAuth(1)
		sx.Assert().NoError(err)
		sx.Assert().Equal("test3", auth.Password)
	})
}

func (s *ClientSuite) TestUpdateClientPasswordCorrectClassic(t provider.T) {
	var (
		db  *sql.DB
		mod *data_base.DataBaseModule
		uid int64
	)
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		var err error
		db, err = sql.Open("duckdb", ":memory:")
		sx.Assert().NoError(err)
		t.Cleanup(func() { _ = db.Close() })
		mod, err = tu.SetupModule(db)
		sx.Assert().NoError(err)
		err = CreateClientService(mod.ClientRepository, mod.PersonalDataRepository, mod.UserRepository, mod.ReviewRepository).CreateClient(tu.TestInitClientData)
		sx.Assert().NoError(err)
		res, err := mod.AuthRepository.Authorize(types.DBAuthData{Login: tu.TestAuth.Login, Password: tu.TestAuth.Password})
		sx.Assert().NoError(err)
		uid = res.UserID
	})
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		err := CreateClientService(mod.ClientRepository, mod.PersonalDataRepository, mod.UserRepository, mod.ReviewRepository).UpdateClientPassword(uid, tu.TestAuth, "test3")
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		_, err := mod.AuthRepository.Authorize(types.DBAuthData{Login: tu.TestAuth.Login, Password: "test3"})
		sx.Assert().NoError(err)
	})
}

func (s *ClientSuite) TestUpdateClientPasswordIncorrectLondon(t provider.T) {
	var (
		pdRepo = tu.CreateTestPersonalDataRepository()
		aRepo  = tu.CreateTestAuthRepository()
		uRepo  = tu.CreateTestUserRepository()
		rRepo  = tu.CreateTestReviewRepository()
		cRepo  = tu.CreateTestClientRepository(pdRepo, aRepo, uRepo)
	)
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		err := CreateClientService(cRepo, pdRepo, uRepo, rRepo).CreateClient(tu.TestInitClientData)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Act+Assert", func(sx provider.StepCtx) {
		err := CreateClientService(cRepo, pdRepo, uRepo, rRepo).UpdateClientPassword(2, tu.TestAuth, "test3")
		sx.Assert().Error(err)
	})
}

func (s *ClientSuite) TestUpdateClientPasswordIncorrectClassic(t provider.T) {
	var (
		db  *sql.DB
		mod *data_base.DataBaseModule
	)
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		var err error
		db, err = sql.Open("duckdb", ":memory:")
		sx.Assert().NoError(err)
		t.Cleanup(func() { _ = db.Close() })
		mod, err = tu.SetupModule(db)
		sx.Assert().NoError(err)
		err = CreateClientService(mod.ClientRepository, mod.PersonalDataRepository, mod.UserRepository, mod.ReviewRepository).CreateClient(tu.TestInitClientData)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Act+Assert", func(sx provider.StepCtx) {
		err := CreateClientService(mod.ClientRepository, mod.PersonalDataRepository, mod.UserRepository, mod.ReviewRepository).UpdateClientPassword(1, tu.TestAuth, "test3")
		sx.Assert().Error(err)
	})
}
