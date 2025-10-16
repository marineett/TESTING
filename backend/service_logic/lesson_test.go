package service_logic

import (
	"data_base_project/data_base"
	tu "data_base_project/test_service_utility"
	"data_base_project/types"
	"database/sql"
	"testing"
	"time"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type LessonSuite struct {
	suite.Suite
}

func TestRunLessonSuite(t *testing.T) {
	suite.RunSuite(t, new(LessonSuite))
}

func (s *LessonSuite) TestCreateLessonCorrectLondon(t provider.T) {
	var repo = tu.CreateTestLessonRepository()
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		ts := CreateLessonService(repo)
		tu.TestLesson.ContractID = 1
		_, err := ts.CreateLesson(tu.TestLesson)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		lesson, err := repo.GetLesson(1)
		sx.Assert().NoError(err)
		sx.Assert().Equal(tu.TestLesson.Duration, lesson.Duration)
	})
}

func (s *LessonSuite) TestCreateLessonCorrectClassic(t provider.T) {
	var (
		db  *sql.DB
		mod *data_base.DataBaseModule
		cid int64
	)
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		var err error
		db, err = sql.Open("duckdb", ":memory:")
		sx.Assert().NoError(err)
		t.Cleanup(func() { _ = db.Close() })
		mod, err = tu.SetupModule(db)
		sx.Assert().NoError(err)
		cs := CreateClientService(mod.ClientRepository, mod.PersonalDataRepository, mod.UserRepository, mod.ReviewRepository)
		sx.Assert().NoError(cs.CreateClient(tu.TestInitClientData))
		res, err := mod.AuthRepository.Authorize(types.DBAuthData{Login: tu.TestAuth.Login, Password: tu.TestAuth.Password})
		sx.Assert().NoError(err)
		cID, err := mod.ContractRepository.InsertContract(types.DBContract{ClientID: res.UserID, Status: types.ContractStatusActive, PaymentStatus: types.PaymentStatusPaid, CreatedAt: time.Now()})
		sx.Assert().NoError(err)
		cid = cID
	})
	t.WithNewStep("Act+Assert", func(sx provider.StepCtx) {
		ls := CreateLessonService(mod.LessonRepository)
		tu.TestLesson.ContractID = cid
		_, err := ls.CreateLesson(tu.TestLesson)
		sx.Assert().NoError(err)
		lessons, err := mod.LessonRepository.GetLessons(cid, 0, 10)
		sx.Assert().NoError(err)
		sx.Require().True(len(lessons) >= 1)
		sx.Assert().Equal(tu.TestLesson.Duration, lessons[0].Duration)
	})
}

func (s *LessonSuite) TestCreateLessonIncorrectLondon(t provider.T) {
	var repo = tu.CreateTestLessonRepository()
	t.WithNewStep("Act+Assert", func(sx provider.StepCtx) {
		ls := CreateLessonService(repo)
		tu.TestLesson.ContractID = 0
		_, err := ls.CreateLesson(tu.TestLesson)
		sx.Assert().Error(err)
	})
}

func (s *LessonSuite) TestCreateLessonIncorrectClassic(t provider.T) {
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
		_, err := CreateLessonService(mod.LessonRepository).CreateLesson(tu.TestLesson)
		sx.Assert().Error(err)
	})
}

func (s *LessonSuite) TestGetLessonsCorrectLondon(t provider.T) {
	var repo = tu.CreateTestLessonRepository()
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		_, err := CreateLessonService(repo).GetLessons(1, 0, 10)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		tu.TestLesson.ContractID = 1
		_, err := repo.InsertLesson(*types.MapperLessonServiceToDB(&tu.TestLesson))
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		lessons, err := CreateLessonService(repo).GetLessons(1, 0, 10)
		sx.Assert().NoError(err)
		sx.Assert().Equal(1, len(lessons))
	})
}

func (s *LessonSuite) TestGetLessonsCorrectClassic(t provider.T) {
	var (
		db  *sql.DB
		mod *data_base.DataBaseModule
		cid int64
	)
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		var err error
		db, err = sql.Open("duckdb", ":memory:")
		sx.Assert().NoError(err)
		t.Cleanup(func() { _ = db.Close() })
		mod, err = tu.SetupModule(db)
		sx.Assert().NoError(err)
		cs := CreateClientService(mod.ClientRepository, mod.PersonalDataRepository, mod.UserRepository, mod.ReviewRepository)
		sx.Assert().NoError(cs.CreateClient(tu.TestInitClientData))
		res, err := mod.AuthRepository.Authorize(types.DBAuthData{Login: tu.TestAuth.Login, Password: tu.TestAuth.Password})
		sx.Assert().NoError(err)
		cID, err := mod.ContractRepository.InsertContract(types.DBContract{ClientID: res.UserID, Status: types.ContractStatusActive, PaymentStatus: types.PaymentStatusPaid, CreatedAt: time.Now()})
		sx.Assert().NoError(err)
		cid = cID
	})
	t.WithNewStep("Act+Assert", func(sx provider.StepCtx) {
		ts := CreateLessonService(mod.LessonRepository)
		tu.TestLesson.ContractID = cid
		_, err := ts.CreateLesson(tu.TestLesson)
		sx.Assert().NoError(err)
		lessons, err := ts.GetLessons(cid, 0, 10)
		sx.Assert().NoError(err)
		sx.Require().True(len(lessons) >= 1)
		sx.Assert().Equal(tu.TestLesson.Duration, lessons[0].Duration)
	})
}
