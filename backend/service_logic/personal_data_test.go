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

type PersonalDataSuite struct {
	suite.Suite
}

func TestRunPersonalDataSuite(t *testing.T) {
	suite.RunSuite(t, new(PersonalDataSuite))
}

func (s *PersonalDataSuite) TestGetPersonalDataLondon(t provider.T) {
	var repo = tu.CreateTestPersonalDataRepository()
	var id int64
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		pd := types.MapperPersonalDataServiceToDB(&tu.TestPD)
		var err error
		id, err = repo.InsertPersonalData(*pd)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Act+Assert", func(sx provider.StepCtx) {
		service := CreatePersonalDataService(repo)
		res, err := service.GetPersonalData(id)
		sx.Assert().NoError(err)
		sx.Assert().Equal(tu.TestPD.TelephoneNumber, res.TelephoneNumber)
		sx.Assert().Equal(tu.TestPD.Email, res.Email)
		sx.Assert().Equal(tu.TestPD.FirstName, res.FirstName)
		sx.Assert().Equal(tu.TestPD.LastName, res.LastName)
		sx.Assert().Equal(tu.TestPD.MiddleName, res.MiddleName)
		sx.Assert().Equal(tu.TestPD.ServicePassportData.PassportNumber, res.ServicePassportData.PassportNumber)
		sx.Assert().Equal(tu.TestPD.ServicePassportData.PassportSeries, res.ServicePassportData.PassportSeries)
		sx.Assert().Equal(tu.TestPD.ServicePassportData.PassportIssuedBy, res.ServicePassportData.PassportIssuedBy)
	})
}

func (s *PersonalDataSuite) TestGetPersonalDataClassic(t provider.T) {
	var db *sql.DB
	var mod *data_base.DataBaseModule
	var id int64
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		var err error
		db, err = sql.Open("duckdb", ":memory:")
		sx.Assert().NoError(err)
		t.Cleanup(func() { _ = db.Close() })
		mod, err = tu.SetupModule(db)
		sx.Assert().NoError(err)
		pd := types.MapperPersonalDataServiceToDB(&tu.TestPD)
		id, err = mod.PersonalDataRepository.InsertPersonalData(*pd)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Act+Assert", func(sx provider.StepCtx) {
		service := CreatePersonalDataService(mod.PersonalDataRepository)
		res, err := service.GetPersonalData(id)
		sx.Assert().NoError(err)
		sx.Assert().Equal(tu.TestPD.TelephoneNumber, res.TelephoneNumber)
		sx.Assert().Equal(tu.TestPD.Email, res.Email)
		sx.Assert().Equal(tu.TestPD.FirstName, res.FirstName)
		sx.Assert().Equal(tu.TestPD.LastName, res.LastName)
		sx.Assert().Equal(tu.TestPD.MiddleName, res.MiddleName)
		sx.Assert().Equal(tu.TestPD.ServicePassportData.PassportNumber, res.ServicePassportData.PassportNumber)
		sx.Assert().Equal(tu.TestPD.ServicePassportData.PassportSeries, res.ServicePassportData.PassportSeries)
		sx.Assert().Equal(tu.TestPD.ServicePassportData.PassportIssuedBy, res.ServicePassportData.PassportIssuedBy)
	})
}

func (s *PersonalDataSuite) TestGetPersonalDataIncorrectLondon(t provider.T) {
	var repo = tu.CreateTestPersonalDataRepository()
	var id int64
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		pd := types.MapperPersonalDataServiceToDB(&tu.TestPD)
		var err error
		id, err = repo.InsertPersonalData(*pd)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Act+Assert", func(sx provider.StepCtx) {
		_, err := CreatePersonalDataService(repo).GetPersonalData(id + 1)
		sx.Assert().Error(err)
	})
}

func (s *PersonalDataSuite) TestGetPersonalDataIncorrectClassic(t provider.T) {
	var db *sql.DB
	var mod *data_base.DataBaseModule
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		var err error
		db, err = sql.Open("duckdb", ":memory:")
		sx.Assert().NoError(err)
		t.Cleanup(func() { _ = db.Close() })
		mod, err = tu.SetupModule(db)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Act+Assert", func(sx provider.StepCtx) {
		_, err := CreatePersonalDataService(mod.PersonalDataRepository).GetPersonalData(1)
		sx.Assert().Error(err)
	})
}
