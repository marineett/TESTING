package integration_tests

import (
	"data_base_project/data_base"
	"data_base_project/service_logic"
	tu "data_base_project/test_service_utility"
	"data_base_project/types"
	"database/sql"
	"testing"
	"time"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type ChatSuite struct {
	suite.Suite
	db  *sql.DB
	mod *data_base.DataBaseModule
}

func (s *ChatSuite) BeforeAll(t provider.T) {
	var err error
	s.db, err = sql.Open("duckdb", ":memory:")
	t.Assert().NoError(err)
	s.mod, err = tu.SetupModule(s.db)
	t.Assert().NoError(err)
}

func (s *ChatSuite) AfterAll(t provider.T) {
	if s.db != nil {
		_ = s.db.Close()
	}
}

func TestRunChatSuite(t *testing.T) {
	suite.RunSuite(t, new(ChatSuite))
}

func (s *ChatSuite) TestCreateCRChatCorrect(t provider.T) {
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

	var chatID int64
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		cs := service_logic.CreateChatService(mod.ChatRepository, mod.MessageRepository)
		var err error
		chatID, err = cs.CreateCRChat(1, 2)
		sx.Assert().NoError(err)
	})

	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		sx.Assert().Equal(int64(1), chatID)
		chat, err := mod.ChatRepository.GetChat(chatID)
		sx.Assert().NoError(err)
		sx.Assert().Equal(int64(1), chat.ClientID)
		sx.Assert().Equal(int64(2), chat.RepetitorID)
		sx.Assert().Equal(int64(0), chat.ModeratorID)
	})
}

func (s *ChatSuite) TestCreateRMChatCorrect(t provider.T) {
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
	var chatID int64
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		cs := service_logic.CreateChatService(mod.ChatRepository, mod.MessageRepository)
		var err error
		chatID, err = cs.CreateRMChat(1, 2)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		sx.Assert().Equal(int64(1), chatID)
		chat, err := mod.ChatRepository.GetChat(chatID)
		sx.Assert().NoError(err)
		sx.Assert().Equal(int64(1), chat.RepetitorID)
		sx.Assert().Equal(int64(2), chat.ModeratorID)
		sx.Assert().Equal(int64(0), chat.ClientID)
	})
}

func (s *ChatSuite) TestCreateCMChatCorrect(t provider.T) {
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
	var chatID int64
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		cs := service_logic.CreateChatService(mod.ChatRepository, mod.MessageRepository)
		var err error
		chatID, err = cs.CreateCMChat(1, 2)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		sx.Assert().Equal(int64(1), chatID)
		chat, err := mod.ChatRepository.GetChat(chatID)
		sx.Assert().NoError(err)
		sx.Assert().Equal(int64(1), chat.ClientID)
		sx.Assert().Equal(int64(2), chat.ModeratorID)
		sx.Assert().Equal(int64(0), chat.RepetitorID)
	})
}

func (s *ChatSuite) TestGetChatListByClientIDCorrect(t provider.T) {
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
		_, _ = mod.ChatRepository.InsertChat(types.DBChat{ClientID: 1, ModeratorID: 2, RepetitorID: 0, CreatedAt: time.Now()})
		_, _ = mod.ChatRepository.InsertChat(types.DBChat{ClientID: 1, ModeratorID: 0, RepetitorID: 3, CreatedAt: time.Now()})
		_, _ = mod.ChatRepository.InsertChat(types.DBChat{ClientID: 2, ModeratorID: 4, RepetitorID: 0, CreatedAt: time.Now()})
	})
	var list []types.ServiceChat
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		cs := service_logic.CreateChatService(mod.ChatRepository, mod.MessageRepository)
		var err error
		list, err = cs.GetChatListByClientID(1, 0, 10)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		sx.Assert().Equal(2, len(list))
		sx.Assert().Equal(int64(1), list[0].ClientID)
		sx.Assert().Equal(int64(1), list[1].ClientID)
		list2, err := service_logic.CreateChatService(mod.ChatRepository, mod.MessageRepository).GetChatListByClientID(5, 0, 10)
		sx.Assert().NoError(err)
		sx.Assert().Equal(0, len(list2))
	})
}

func (s *ChatSuite) TestGetChatListByRepetitorIDCorrect(t provider.T) {
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
		_, _ = mod.ChatRepository.InsertChat(types.DBChat{ClientID: 1, ModeratorID: 0, RepetitorID: 3, CreatedAt: time.Now()})
		_, _ = mod.ChatRepository.InsertChat(types.DBChat{ClientID: 2, ModeratorID: 0, RepetitorID: 3, CreatedAt: time.Now()})
		_, _ = mod.ChatRepository.InsertChat(types.DBChat{ClientID: 2, ModeratorID: 0, RepetitorID: 4, CreatedAt: time.Now()})
	})
	var list []types.ServiceChat
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		var err error
		list, err = service_logic.CreateChatService(mod.ChatRepository, mod.MessageRepository).GetChatListByRepetitorID(3, 0, 10)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		sx.Assert().Equal(2, len(list))
		sx.Assert().Equal(int64(3), list[0].RepetitorID)
		sx.Assert().Equal(int64(3), list[1].RepetitorID)
		list2, err := service_logic.CreateChatService(mod.ChatRepository, mod.MessageRepository).GetChatListByRepetitorID(5, 0, 10)
		sx.Assert().NoError(err)
		sx.Assert().Equal(0, len(list2))
	})
}

func (s *ChatSuite) TestGetChatListByModeratorIDCorrect(t provider.T) {
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
		_, _ = mod.ChatRepository.InsertChat(types.DBChat{ClientID: 1, ModeratorID: 2, RepetitorID: 0, CreatedAt: time.Now()})
		_, _ = mod.ChatRepository.InsertChat(types.DBChat{ClientID: 0, ModeratorID: 2, RepetitorID: 4, CreatedAt: time.Now()})
		_, _ = mod.ChatRepository.InsertChat(types.DBChat{ClientID: 7, ModeratorID: 3, RepetitorID: 0, CreatedAt: time.Now()})
	})
	var list []types.ServiceChat
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		var err error
		list, err = service_logic.CreateChatService(mod.ChatRepository, mod.MessageRepository).GetChatListByModeratorID(2, 0, 10)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		sx.Assert().Equal(2, len(list))
		sx.Assert().Equal(int64(2), list[0].ModeratorID)
		sx.Assert().Equal(int64(2), list[1].ModeratorID)
		list2, err := service_logic.CreateChatService(mod.ChatRepository, mod.MessageRepository).GetChatListByModeratorID(5, 0, 10)
		sx.Assert().NoError(err)
		sx.Assert().Equal(0, len(list2))
	})
}

func (s *ChatSuite) TestGetChatCorrect(t provider.T) {
	var (
		db     *sql.DB
		mod    *data_base.DataBaseModule
		chatID int64
	)
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		var err error
		db, err = sql.Open("duckdb", ":memory:")
		sx.Assert().NoError(err)
		t.Cleanup(func() { _ = db.Close() })
		mod, err = tu.SetupModule(db)
		sx.Assert().NoError(err)
		chatID, err = mod.ChatRepository.InsertChat(types.DBChat{ClientID: 1, ModeratorID: 2, RepetitorID: 0, CreatedAt: time.Now()})
		sx.Assert().NoError(err)
	})
	var chat *types.ServiceChat
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		var err error
		chat, err = service_logic.CreateChatService(mod.ChatRepository, mod.MessageRepository).GetChat(chatID)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		sx.Assert().Equal(int64(1), chat.ClientID)
		sx.Assert().Equal(int64(2), chat.ModeratorID)
		sx.Assert().Equal(int64(0), chat.RepetitorID)
	})
}

func (s *ChatSuite) TestSendMessageCorrect(t provider.T) {
	var (
		db     *sql.DB
		mod    *data_base.DataBaseModule
		chatID int64
		userID int64
	)
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		var err error
		db, err = sql.Open("duckdb", ":memory:")
		sx.Assert().NoError(err)
		t.Cleanup(func() { _ = db.Close() })
		mod, err = tu.SetupModule(db)
		sx.Assert().NoError(err)
		cs := service_logic.CreateClientService(mod.ClientRepository, mod.PersonalDataRepository, mod.UserRepository, mod.ReviewRepository)
		err = cs.CreateClient(tu.TestInitClientData)
		sx.Assert().NoError(err)
		res, err := mod.AuthRepository.Authorize(types.DBAuthData{Login: tu.TestAuth.Login, Password: tu.TestAuth.Password})
		sx.Assert().NoError(err)
		userID = res.UserID
		chatID, err = mod.ChatRepository.InsertChat(types.DBChat{ClientID: userID, ModeratorID: 2, RepetitorID: 0, CreatedAt: time.Now()})
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		_, err := service_logic.CreateChatService(mod.ChatRepository, mod.MessageRepository).SendMessage(chatID, userID, "Hello")
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		msgs, err := mod.MessageRepository.GetMessages(chatID, 0, 10)
		sx.Assert().NoError(err)
		sx.Require().Equal(1, len(msgs))
		sx.Assert().Equal(userID, msgs[0].SenderID)
		sx.Assert().Equal("Hello", msgs[0].Content)
	})
}

func (s *ChatSuite) TestGetMessagesCorrect(t provider.T) {
	var (
		db     *sql.DB
		mod    *data_base.DataBaseModule
		chatID int64
		msgID  int64
		userID int64
	)
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		var err error
		db, err = sql.Open("duckdb", ":memory:")
		sx.Assert().NoError(err)
		t.Cleanup(func() { _ = db.Close() })
		mod, err = tu.SetupModule(db)
		sx.Assert().NoError(err)
		cs := service_logic.CreateClientService(mod.ClientRepository, mod.PersonalDataRepository, mod.UserRepository, mod.ReviewRepository)
		err = cs.CreateClient(tu.TestInitClientData)
		sx.Assert().NoError(err)
		res, err := mod.AuthRepository.Authorize(types.DBAuthData{Login: tu.TestAuth.Login, Password: tu.TestAuth.Password})
		sx.Assert().NoError(err)
		userID = res.UserID
		chatID, err = mod.ChatRepository.InsertChat(types.DBChat{ClientID: userID, ModeratorID: 2, RepetitorID: 0, CreatedAt: time.Now()})
		sx.Assert().NoError(err)
		msgID, err = mod.MessageRepository.InsertMessage(types.DBMessage{ChatID: chatID, SenderID: userID, Content: "Hello", CreatedAt: time.Now()})
		sx.Assert().NoError(err)
	})
	var msgs []types.ServiceMessage
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		var err error
		msgs, err = service_logic.CreateChatService(mod.ChatRepository, mod.MessageRepository).GetMessages(chatID, 0, 10)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		sx.Require().Equal(1, len(msgs))
		sx.Assert().Equal(msgID, msgs[0].ID)
		sx.Assert().Equal(userID, msgs[0].SenderID)
		sx.Assert().Equal("Hello", msgs[0].Content)
	})
}

func (s *ChatSuite) TestGetMessagesIncorrect(t provider.T) {
	var (
		db   *sql.DB
		mod  *data_base.DataBaseModule
		msgs []types.ServiceMessage
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
		var err error
		msgs, err = service_logic.CreateChatService(mod.ChatRepository, mod.MessageRepository).GetMessages(1, 0, 10)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		sx.Assert().Equal(0, len(msgs))
	})
}

func (s *ChatSuite) TestGetChatIdByCIDAndMIDCorrect(t provider.T) {
	var (
		db     *sql.DB
		mod    *data_base.DataBaseModule
		chatID int64
		resID  int64
	)
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		var err error
		db, err = sql.Open("duckdb", ":memory:")
		sx.Assert().NoError(err)
		t.Cleanup(func() { _ = db.Close() })
		mod, err = tu.SetupModule(db)
		sx.Assert().NoError(err)
		chatID, err = mod.ChatRepository.InsertChat(types.DBChat{ClientID: 1, ModeratorID: 2, RepetitorID: 0, CreatedAt: time.Now()})
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		var err error
		resID, err = service_logic.CreateChatService(mod.ChatRepository, mod.MessageRepository).GetChatIdByCIDAndMID(1, 2)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		sx.Assert().Equal(chatID, resID)
	})
}

func (s *ChatSuite) TestGetChatIdByCIDAndRIDCorrect(t provider.T) {
	var (
		db     *sql.DB
		mod    *data_base.DataBaseModule
		chatID int64
		resID  int64
	)
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		var err error
		db, err = sql.Open("duckdb", ":memory:")
		sx.Assert().NoError(err)
		t.Cleanup(func() { _ = db.Close() })
		mod, err = tu.SetupModule(db)
		sx.Assert().NoError(err)
		chatID, err = mod.ChatRepository.InsertChat(types.DBChat{ClientID: 1, ModeratorID: 0, RepetitorID: 2, CreatedAt: time.Now()})
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		var err error
		resID, err = service_logic.CreateChatService(mod.ChatRepository, mod.MessageRepository).GetChatIdByCIDAndRID(1, 2)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		sx.Assert().Equal(chatID, resID)
	})
}

func (s *ChatSuite) TestGetChatIdByMIDAndRIDCorrect(t provider.T) {
	var (
		db     *sql.DB
		mod    *data_base.DataBaseModule
		chatID int64
		resID  int64
	)
	t.WithNewStep("Arrange", func(sx provider.StepCtx) {
		var err error
		db, err = sql.Open("duckdb", ":memory:")
		sx.Assert().NoError(err)
		t.Cleanup(func() { _ = db.Close() })
		mod, err = tu.SetupModule(db)
		sx.Assert().NoError(err)
		chatID, err = mod.ChatRepository.InsertChat(types.DBChat{ClientID: 0, ModeratorID: 1, RepetitorID: 2, CreatedAt: time.Now()})
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Act", func(sx provider.StepCtx) {
		var err error
		resID, err = service_logic.CreateChatService(mod.ChatRepository, mod.MessageRepository).GetChatIdByMIDAndRID(1, 2)
		sx.Assert().NoError(err)
	})
	t.WithNewStep("Assert", func(sx provider.StepCtx) {
		sx.Assert().Equal(chatID, resID)
	})
}
