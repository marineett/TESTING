package data_base

import (
	tu "data_base_project/test_database_utility"
	"data_base_project/types"
	"database/sql"
	"fmt"
	"testing"
	"time"
)

func setupChatTables(db *sql.DB) error {
	err := CreateSqlSequence(db, "sequence")
	if err != nil {
		return fmt.Errorf("error creating sequence: %v", err)
	}
	err = CreateSqlPersonalDataTable(db, "personal_data", "sequence")
	if err != nil {
		return fmt.Errorf("error creating personal data table: %v", err)
	}
	err = CreateSqlUserTable(db, "users", "personal_data", "sequence")
	if err != nil {
		return fmt.Errorf("error creating user table: %v", err)
	}
	err = CreateSqlAuthTable(db, "auth", "users", "sequence")
	if err != nil {
		return fmt.Errorf("error creating auth table: %v", err)
	}
	err = CreateSqlClientTable(db, "clients", "users", "sequence")
	if err != nil {
		return fmt.Errorf("error creating client table: %v", err)
	}
	err = CreateSqlRepetitorTable(db, "repetitors", "users", "sequence")
	if err != nil {
		return fmt.Errorf("error creating repetitor table: %v", err)
	}
	err = CreateSqlModeratorTable(db, "moderators", "users")
	if err != nil {
		return fmt.Errorf("error creating moderator table: %v", err)
	}
	err = CreateSqlChatTable(db, "chat", "users")
	if err != nil {
		return fmt.Errorf("error creating chat table: %v", err)
	}
	return nil
}

func TestCreateSqlChatTable(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err = setupChatTables(db)
	if err != nil {
		t.Fatalf("Error setting up chat tables: %v", err)
	}
	chatRepository := CreateSqlChatRepository(db, "chat", "sequence")
	if chatRepository == nil {
		t.Fatalf("Error creating chat repository: %v", err)
	}
}

func TestInsertChatCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err = setupChatTables(db)
	if err != nil {
		t.Fatalf("Error setting up chat tables: %v", err)
	}
	clientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	clientID, err := clientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	repetitorRepository := CreateSqlRepetitorRepository(db, "personal_data", "users", "repetitors", "auth", "resume", "review", "sequence")
	repetitorID, err := repetitorRepository.InsertRepetitor(tu.TestRepetitor, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting repetitor: %v", err)
	}

	chatRepository := CreateSqlChatRepository(db, "chat", "sequence")
	_, err = chatRepository.InsertChat(types.DBChat{
		ClientID:    clientID,
		RepetitorID: repetitorID,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting chat: %v", err)
	}
}

func CheckChat(
	t *testing.T,
	Chat *types.DBChat,
	ChatID int64,
	ClientID int64,
	ModeratorID int64,
	RepetitorID int64,
) {
	if Chat.ID != ChatID {
		t.Fatalf("Chat id not correct: %v", Chat)
	}
	if Chat.ClientID != ClientID {
		t.Fatalf("Chat client id not correct: %v", Chat)
	}
	if Chat.ModeratorID != ModeratorID {
		t.Fatalf("Chat moderator id not correct: %v", Chat)
	}
	if Chat.RepetitorID != RepetitorID {
		t.Fatalf("Chat repetitor id not correct: %v", Chat)
	}
}

func TestGetChatCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err = setupChatTables(db)
	if err != nil {
		t.Fatalf("Error setting up chat tables: %v", err)
	}
	clientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	clientID, err := clientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	moderatorRepository := CreateSqlModeratorRepository(db, "personal_data", "users", "moderators", "auth", "sequence")
	moderatorID, err := moderatorRepository.InsertModerator(tu.TestModeratorData, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting repetitor: %v", err)
	}
	chatRepository := CreateSqlChatRepository(db, "chat", "sequence")
	chatID, err := chatRepository.InsertChat(types.DBChat{
		ClientID:    clientID,
		ModeratorID: moderatorID,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting chat: %v", err)
	}
	chat, err := chatRepository.GetChat(chatID)
	if err != nil {
		t.Fatalf("Error getting chat: %v", err)
	}
	CheckChat(t, chat, chatID, clientID, moderatorID, 0)
}

func TestGetChatIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err = setupChatTables(db)
	if err != nil {
		t.Fatalf("Error setting up chat tables: %v", err)
	}
	chatRepository := CreateSqlChatRepository(db, "chat", "sequence")
	_, err = chatRepository.GetChat(1)
	if err == nil {
		t.Fatalf("No error getting chat: %v", err)
	}
}

func CheckChatsLength(t *testing.T, chats []types.DBChat, length int) {
	if len(chats) != length {
		t.Fatalf("Chats not updated: %v", chats)
	}
}

func TestGetChatIdByCIDAndMIDCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err = setupChatTables(db)
	if err != nil {
		t.Fatalf("Error setting up chat tables: %v", err)
	}
	clientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	clientID, err := clientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	moderatorRepository := CreateSqlModeratorRepository(db, "personal_data", "users", "moderators", "auth", "sequence")
	moderatorID, err := moderatorRepository.InsertModerator(tu.TestModeratorData, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting repetitor: %v", err)
	}
	chatRepository := CreateSqlChatRepository(db, "chat", "sequence")
	insertedChatID, err := chatRepository.InsertChat(types.DBChat{
		ClientID:    clientID,
		ModeratorID: moderatorID,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting chat: %v", err)
	}
	chatID, err := chatRepository.GetChatIdByCIDAndMID(clientID, moderatorID)
	if err != nil {
		t.Fatalf("Error getting chat id: %v", err)
	}
	if chatID != insertedChatID {
		t.Fatalf("Chat id not correct: %v", chatID)
	}
	_, err = chatRepository.GetChatIdByCIDAndMID(clientID+1, moderatorID)
	if err != nil {
		t.Fatalf("Error getting chat id: %v", err)
	}
}

func TestGetChatIdByCIDAndRIDCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err = setupChatTables(db)
	if err != nil {
		t.Fatalf("Error setting up chat tables: %v", err)
	}
	clientRepository := CreateSqlClientRepository(db, "personal_data", "users", "clients", "auth", "sequence")
	clientID, err := clientRepository.InsertClient(tu.TestClient, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	repetitorRepository := CreateSqlRepetitorRepository(db, "personal_data", "users", "repetitors", "auth", "resume", "review", "sequence")
	repetitorID, err := repetitorRepository.InsertRepetitor(tu.TestRepetitor, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting repetitor: %v", err)
	}
	chatRepository := CreateSqlChatRepository(db, "chat", "sequence")
	insertedChatID, err := chatRepository.InsertChat(types.DBChat{
		ClientID:    clientID,
		RepetitorID: repetitorID,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting chat: %v", err)
	}
	chatID, err := chatRepository.GetChatIdByCIDAndRID(clientID, repetitorID)
	if err != nil {
		t.Fatalf("Error getting chat id: %v", err)
	}
	if chatID != insertedChatID {
		t.Fatalf("Chat id not correct: %v", chatID)
	}
	_, err = chatRepository.GetChatIdByCIDAndRID(clientID+1, repetitorID)
	if err != nil {
		t.Fatalf("Error getting chat id: %v", err)
	}
}

func TestGetChatIdByMIDAndRIDCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			t.Fatalf("Error closing database: %v", err)
		}
	}()
	err = setupChatTables(db)
	if err != nil {
		t.Fatalf("Error setting up chat tables: %v", err)
	}
	repetitorRepository := CreateSqlRepetitorRepository(db, "personal_data", "users", "repetitors", "auth", "resume", "review", "sequence")
	repetitorID, err := repetitorRepository.InsertRepetitor(tu.TestRepetitor, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting repetitor: %v", err)
	}
	moderatorRepository := CreateSqlModeratorRepository(db, "personal_data", "users", "moderators", "auth", "sequence")
	moderatorID, err := moderatorRepository.InsertModerator(tu.TestModeratorData, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting repetitor: %v", err)
	}
	chatRepository := CreateSqlChatRepository(db, "chat", "sequence")
	insertedChatID, err := chatRepository.InsertChat(types.DBChat{
		RepetitorID: repetitorID,
		ModeratorID: moderatorID,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting chat: %v", err)
	}
	chatID, err := chatRepository.GetChatIdByMIDAndRID(moderatorID, repetitorID)
	if err != nil {
		t.Fatalf("Error getting chat id: %v", err)
	}
	if chatID != insertedChatID {
		t.Fatalf("Chat id not correct: %v", chatID)
	}
	_, err = chatRepository.GetChatIdByMIDAndRID(moderatorID+1, repetitorID)
	if err != nil {
		t.Fatalf("Error getting chat id: %v", err)
	}
}
