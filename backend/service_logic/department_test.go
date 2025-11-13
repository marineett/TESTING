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

type DepartmentSuite struct {
	suite.Suite
}

func TestRunDepartmentSuite(t *testing.T) {
	suite.RunSuite(t, new(DepartmentSuite))
}

func (s *DepartmentSuite) TestCreateDepartmentCorrectLondon(t provider.T) {
	var (
		depRepo = tu.CreateTestDepartmentRepository()
		pdRepo  = tu.CreateTestPersonalDataRepository()
		aRepo   = tu.CreateTestAuthRepository()
		uRepo   = tu.CreateTestUserRepository()
		modRepo = tu.CreateTestModeratorRepository(pdRepo, aRepo, uRepo)
	)
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		_, err := CreateDepartmentService(depRepo, modRepo, uRepo, pdRepo).CreateDepartment(tu.TestInitDepartmentData)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		dep, err := depRepo.GetDepartment(1)
		sx.Assert().NoError(err)
		sx.Assert().Equal(tu.TestInitDepartmentData.Name, dep.Name)
		sx.Assert().Equal(tu.TestInitDepartmentData.HeadID, dep.HeadID)
	})
}

func (s *DepartmentSuite) TestCreateDepartmentCorrectClassic(t provider.T) {
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
		err = CreateAdminService(mod.AdminRepository, mod.UserRepository, mod.PersonalDataRepository).CreateAdmin(tu.TestInitAdminData)
		sx.Assert().NoError(err)
		res, err := mod.AuthRepository.Authorize(types.DBAuthData{Login: tu.TestAuth.Login, Password: tu.TestAuth.Password})
		sx.Assert().NoError(err)
		uid = res.UserID
	})
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		init := tu.TestInitDepartmentData
		init.HeadID = uid
		_, err := CreateDepartmentService(mod.DepartmentRepository, mod.ModeratorRepository, mod.UserRepository, mod.PersonalDataRepository).CreateDepartment(init)
		sx.Assert().NoError(err)
	})
}

func (s *DepartmentSuite) TestGetDepartmentsByHeadIDCorrectLondon(t provider.T) {
	var (
		depRepo = tu.CreateTestDepartmentRepository()
		pdRepo  = tu.CreateTestPersonalDataRepository()
		aRepo   = tu.CreateTestAuthRepository()
		uRepo   = tu.CreateTestUserRepository()
		modRepo = tu.CreateTestModeratorRepository(pdRepo, aRepo, uRepo)
	)
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		_, err := CreateDepartmentService(depRepo, modRepo, uRepo, pdRepo).CreateDepartment(tu.TestInitDepartmentData)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		list, err := CreateDepartmentService(depRepo, modRepo, uRepo, pdRepo).GetDepartmentsByHeadID(tu.TestInitDepartmentData.HeadID)
		sx.Assert().NoError(err)
		sx.Assert().Equal(1, len(list))
		sx.Assert().Equal(tu.TestInitDepartmentData.Name, list[0].Name)
	})
}

func (s *DepartmentSuite) TestGetDepartmentsByHeadIDCorrectClassic(t provider.T) {
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
		err = CreateAdminService(mod.AdminRepository, mod.UserRepository, mod.PersonalDataRepository).CreateAdmin(tu.TestInitAdminData)
		sx.Assert().NoError(err)
		res, err := mod.AuthRepository.Authorize(types.DBAuthData{Login: tu.TestAuth.Login, Password: tu.TestAuth.Password})
		sx.Assert().NoError(err)
		uid = res.UserID
	})
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		init := tu.TestInitDepartmentData
		init.HeadID = uid
		ds := CreateDepartmentService(mod.DepartmentRepository, mod.ModeratorRepository, mod.UserRepository, mod.PersonalDataRepository)
		_, err := ds.CreateDepartment(init)
		sx.Assert().NoError(err)
		_, err = ds.CreateDepartment(init)
		sx.Assert().NoError(err)
		list, err := ds.GetDepartmentsByHeadID(uid)
		sx.Assert().NoError(err)
		sx.Assert().Equal(2, len(list))
		sx.Assert().Equal(uid, list[0].HeadID)
		sx.Assert().Equal(uid, list[1].HeadID)
	})
}

func (s *DepartmentSuite) TestGetDepartmentCorrectLondon(t provider.T) {
	var (
		depRepo = tu.CreateTestDepartmentRepository()
		pdRepo  = tu.CreateTestPersonalDataRepository()
		aRepo   = tu.CreateTestAuthRepository()
		uRepo   = tu.CreateTestUserRepository()
		modRepo = tu.CreateTestModeratorRepository(pdRepo, aRepo, uRepo)
	)
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		_, err := CreateDepartmentService(depRepo, modRepo, uRepo, pdRepo).CreateDepartment(tu.TestInitDepartmentData)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		dep, err := CreateDepartmentService(depRepo, modRepo, uRepo, pdRepo).GetDepartment(1)
		sx.Assert().NoError(err)
		sx.Assert().Equal(tu.TestInitDepartmentData.Name, dep.Name)
		sx.Assert().Equal(tu.TestInitDepartmentData.HeadID, dep.HeadID)
	})
}

func (s *DepartmentSuite) TestGetDepartmentCorrectClassic(t provider.T) {
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
		sx.Assert().NoError(CreateAdminService(mod.AdminRepository, mod.UserRepository, mod.PersonalDataRepository).CreateAdmin(tu.TestInitAdminData))
		res, err := mod.AuthRepository.Authorize(types.DBAuthData{Login: tu.TestAuth.Login, Password: tu.TestAuth.Password})
		sx.Assert().NoError(err)
		uid = res.UserID
	})
	t.WithNewStep("Act+Assert", func(sx provider.StepCtx) {
		init := tu.TestInitDepartmentData
		init.HeadID = uid
		ds := CreateDepartmentService(mod.DepartmentRepository, mod.ModeratorRepository, mod.UserRepository, mod.PersonalDataRepository)
		_, err := ds.CreateDepartment(init)
		sx.Assert().NoError(err)
		list, err := ds.GetDepartmentsByHeadID(uid)
		sx.Assert().NoError(err)
		sx.Require().Equal(1, len(list))
		dep, err := ds.GetDepartment(list[0].ID)
		sx.Assert().NoError(err)
		sx.Assert().Equal(init.Name, dep.Name)
		sx.Assert().Equal(uid, dep.HeadID)
	})
}

func (s *DepartmentSuite) TestAssignAdminToDepartmentCorrectLondon(t provider.T) {
	var (
		depRepo = tu.CreateTestDepartmentRepository()
		pdRepo  = tu.CreateTestPersonalDataRepository()
		aRepo   = tu.CreateTestAuthRepository()
		uRepo   = tu.CreateTestUserRepository()
		modRepo = tu.CreateTestModeratorRepository(pdRepo, aRepo, uRepo)
	)
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		init := tu.TestInitDepartmentData
		init.HeadID = 0
		_, err := CreateDepartmentService(depRepo, modRepo, uRepo, pdRepo).CreateDepartment(init)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Act+Assert", func(sx provider.StepCtx) {
		ds := CreateDepartmentService(depRepo, modRepo, uRepo, pdRepo)
		sx.Assert().NoError(ds.AssignAdminToDepartment(1, 1))
		dep, err := depRepo.GetDepartment(1)
		sx.Assert().NoError(err)
		sx.Assert().Equal(int64(1), dep.HeadID)
	})
}

func (s *DepartmentSuite) TestAssignAdminToDepartmentCorrectClassic(t provider.T) {
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
		sx.Assert().NoError(CreateAdminService(mod.AdminRepository, mod.UserRepository, mod.PersonalDataRepository).CreateAdmin(tu.TestInitAdminData))
		res, err := mod.AuthRepository.Authorize(types.DBAuthData{Login: tu.TestAuth.Login, Password: tu.TestAuth.Password})
		sx.Assert().NoError(err)
		uid = res.UserID
	})
	t.WithNewStep("Act+Assert", func(sx provider.StepCtx) {
		ds := CreateDepartmentService(mod.DepartmentRepository, mod.ModeratorRepository, mod.UserRepository, mod.PersonalDataRepository)
		init := tu.TestInitDepartmentData
		init.HeadID = 0
		_, err := ds.CreateDepartment(init)
		sx.Assert().NoError(err)
		list, err := ds.GetDepartmentsByHeadID(0)
		sx.Assert().NoError(err)
		sx.Require().Equal(1, len(list))
		sx.Assert().NoError(ds.AssignAdminToDepartment(uid, list[0].ID))
		dep, err := mod.DepartmentRepository.GetDepartment(list[0].ID)
		sx.Assert().NoError(err)
		sx.Assert().Equal(uid, dep.HeadID)
	})
}

func (s *DepartmentSuite) TestAssignAdminToDepartmentIncorrectLondon(t provider.T) {
	var (
		depRepo = tu.CreateTestDepartmentRepository()
		pdRepo  = tu.CreateTestPersonalDataRepository()
		aRepo   = tu.CreateTestAuthRepository()
		uRepo   = tu.CreateTestUserRepository()
		modRepo = tu.CreateTestModeratorRepository(pdRepo, aRepo, uRepo)
	)
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		_, err := CreateDepartmentService(depRepo, modRepo, uRepo, pdRepo).CreateDepartment(tu.TestInitDepartmentData)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Act+Assert", func(sx provider.StepCtx) {
		err := CreateDepartmentService(depRepo, modRepo, uRepo, pdRepo).AssignAdminToDepartment(1, 2)
		sx.Assert().Error(err)
	})
}

func (s *DepartmentSuite) TestAssignAdminToDepartmentIncorrectClassic(t provider.T) {
	var db *sql.DB
	var mod *data_base.DataBaseModule
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		var err error
		db, err = sql.Open("duckdb", ":memory:")
		sx.Assert().NoError(err)
		t.Cleanup(func() { _ = db.Close() })
		mod, err = tu.SetupModule(db)
		sx.Assert().NoError(err)
		// create admin to ensure auth table is not empty (mimic original)
		sx.Assert().NoError(CreateAdminService(mod.AdminRepository, mod.UserRepository, mod.PersonalDataRepository).CreateAdmin(tu.TestInitAdminData))
	})
	t.WithNewStep("Act+Assert", func(sx provider.StepCtx) {
		err := CreateDepartmentService(mod.DepartmentRepository, mod.ModeratorRepository, mod.UserRepository, mod.PersonalDataRepository).AssignAdminToDepartment(1, 2)
		sx.Assert().Error(err)
	})
}

func (s *DepartmentSuite) TestFireAdminFromDepartmentCorrectLondon(t provider.T) {
	var (
		depRepo = tu.CreateTestDepartmentRepository()
		pdRepo  = tu.CreateTestPersonalDataRepository()
		aRepo   = tu.CreateTestAuthRepository()
		uRepo   = tu.CreateTestUserRepository()
		modRepo = tu.CreateTestModeratorRepository(pdRepo, aRepo, uRepo)
	)
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		_, err := CreateDepartmentService(depRepo, modRepo, uRepo, pdRepo).CreateDepartment(tu.TestInitDepartmentData)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Act+Assert", func(sx provider.StepCtx) {
		ds := CreateDepartmentService(depRepo, modRepo, uRepo, pdRepo)
		sx.Assert().NoError(ds.FireAdminFromDepartment(tu.TestInitDepartmentData.HeadID, 1))
		dep, err := depRepo.GetDepartment(1)
		sx.Assert().NoError(err)
		sx.Assert().Equal(int64(0), dep.HeadID)
	})
}

func (s *DepartmentSuite) TestFireAdminFromDepartmentCorrectClassic(t provider.T) {
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
		sx.Assert().NoError(CreateAdminService(mod.AdminRepository, mod.UserRepository, mod.PersonalDataRepository).CreateAdmin(tu.TestInitAdminData))
		res, err := mod.AuthRepository.Authorize(types.DBAuthData{Login: tu.TestAuth.Login, Password: tu.TestAuth.Password})
		sx.Assert().NoError(err)
		uid = res.UserID
	})
	t.WithNewStep("Act+Assert", func(sx provider.StepCtx) {
		ds := CreateDepartmentService(mod.DepartmentRepository, mod.ModeratorRepository, mod.UserRepository, mod.PersonalDataRepository)
		init := tu.TestInitDepartmentData
		init.HeadID = uid
		_, err := ds.CreateDepartment(init)
		sx.Assert().NoError(err)
		list, err := ds.GetDepartmentsByHeadID(uid)
		sx.Assert().NoError(err)
		sx.Require().Equal(1, len(list))
		sx.Assert().NoError(ds.FireAdminFromDepartment(uid, list[0].ID))
		dep, err := mod.DepartmentRepository.GetDepartment(list[0].ID)
		sx.Assert().NoError(err)
		sx.Assert().Equal(int64(0), dep.HeadID)
	})
}

func (s *DepartmentSuite) TestFireAdminFromDepartmentIncorrectLondon(t provider.T) {
	var (
		depRepo = tu.CreateTestDepartmentRepository()
		pdRepo  = tu.CreateTestPersonalDataRepository()
		aRepo   = tu.CreateTestAuthRepository()
		uRepo   = tu.CreateTestUserRepository()
		modRepo = tu.CreateTestModeratorRepository(pdRepo, aRepo, uRepo)
	)
	t.WithNewStep("Act+Assert", func(sx provider.StepCtx) {
		err := CreateDepartmentService(depRepo, modRepo, uRepo, pdRepo).FireAdminFromDepartment(tu.TestInitDepartmentData.HeadID+1, tu.TestInitDepartmentData.ID)
		sx.Assert().Error(err)
	})
}

func (s *DepartmentSuite) TestFireAdminFromDepartmentIncorrectClassic(t provider.T) {
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
		err := CreateDepartmentService(mod.DepartmentRepository, mod.ModeratorRepository, mod.UserRepository, mod.PersonalDataRepository).FireAdminFromDepartment(1, 2)
		sx.Assert().Error(err)
	})
}

func (s *DepartmentSuite) TestFireModeratorFromDepartmentCorrectLondon(t provider.T) {
	var (
		depRepo = tu.CreateTestDepartmentRepository()
		pdRepo  = tu.CreateTestPersonalDataRepository()
		aRepo   = tu.CreateTestAuthRepository()
		uRepo   = tu.CreateTestUserRepository()
		modRepo = tu.CreateTestModeratorRepository(pdRepo, aRepo, uRepo)
	)
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		_, err := CreateDepartmentService(depRepo, modRepo, uRepo, pdRepo).CreateDepartment(tu.TestInitDepartmentData)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Act+Assert", func(sx provider.StepCtx) {
		err := CreateDepartmentService(depRepo, modRepo, uRepo, pdRepo).FireModeratorFromDepartment(tu.TestInitDepartmentData.HeadID, 1)
		sx.Assert().NoError(err)
		dep, err := depRepo.GetDepartment(1)
		sx.Assert().NoError(err)
		sx.Assert().Equal(tu.TestInitDepartmentData.HeadID, dep.HeadID)
	})
}

func (s *DepartmentSuite) TestGetDepartmentUsersIDsIncorrectLondon(t provider.T) {
	var (
		depRepo = tu.CreateTestDepartmentRepository()
		pdRepo  = tu.CreateTestPersonalDataRepository()
		aRepo   = tu.CreateTestAuthRepository()
		uRepo   = tu.CreateTestUserRepository()
		modRepo = tu.CreateTestModeratorRepository(pdRepo, aRepo, uRepo)
	)
	t.WithNewStep("Act+Assert", func(sx provider.StepCtx) {
		_, err := CreateDepartmentService(depRepo, modRepo, uRepo, pdRepo).GetDepartmentUsersIDs(tu.TestInitDepartmentData.ID + 1)
		sx.Assert().Error(err)
	})
}
