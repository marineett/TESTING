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
	defer db.Close()
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
	defer db.Close()
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

func TestGetChatCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
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
	if chat.ID != chatID {
		t.Fatalf("Chat id not correct: %v", chat)
	}
	if chat.ClientID != clientID {
		t.Fatalf("Chat client id not correct: %v", chat)
	}
	if chat.ModeratorID != moderatorID {
		t.Fatalf("Chat moderator id not correct: %v", chat)
	}
	if chat.RepetitorID != 0 {
		t.Fatalf("Chat repetitor id not correct: %v", chat)
	}
}

func TestGetChatIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
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

func TestGetChatListByClientIDCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
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
	moderatorRepository := CreateSqlModeratorRepository(db, "personal_data", "users", "moderators", "auth", "sequence")
	moderatorID, err := moderatorRepository.InsertModerator(tu.TestModeratorData, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting repetitor: %v", err)
	}
	chatRepository := CreateSqlChatRepository(db, "chat", "sequence")
	_, err = chatRepository.InsertChat(types.DBChat{
		ClientID:    clientID,
		ModeratorID: moderatorID,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting chat: %v", err)
	}
	_, err = chatRepository.InsertChat(types.DBChat{
		ClientID:    clientID,
		RepetitorID: repetitorID,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting chat: %v", err)
	}
	_, err = chatRepository.InsertChat(types.DBChat{
		ModeratorID: moderatorID,
		RepetitorID: repetitorID,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting chat: %v", err)
	}
	chatList, err := chatRepository.GetChatListByClientID(clientID, 0, 10)
	if err != nil {
		t.Fatalf("Error getting chat list: %v", err)
	}
	if len(chatList) != 2 {
		t.Fatalf("Chat list not correct: %v", chatList)
	}
	if chatList[0].ClientID != clientID || chatList[1].ClientID != clientID {
		t.Fatalf("Client id not updated: %v", chatList)
	}
	chatList, err = chatRepository.GetChatListByClientID(clientID+1, 0, 10)
	if err != nil {
		t.Fatalf("Error getting chat list: %v", err)
	}
	if len(chatList) != 0 {
		t.Fatalf("Chat list not correct: %v", chatList)
	}
}

func TestGetChatListByModeratorIDCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
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
	moderatorRepository := CreateSqlModeratorRepository(db, "personal_data", "users", "moderators", "auth", "sequence")
	moderatorID, err := moderatorRepository.InsertModerator(tu.TestModeratorData, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting repetitor: %v", err)
	}
	chatRepository := CreateSqlChatRepository(db, "chat", "sequence")
	_, err = chatRepository.InsertChat(types.DBChat{
		ClientID:    clientID,
		ModeratorID: moderatorID,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting chat: %v", err)
	}
	_, err = chatRepository.InsertChat(types.DBChat{
		ClientID:    clientID,
		RepetitorID: repetitorID,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting chat: %v", err)
	}
	_, err = chatRepository.InsertChat(types.DBChat{
		ModeratorID: moderatorID,
		RepetitorID: repetitorID,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting chat: %v", err)
	}
	chatList, err := chatRepository.GetChatListByModeratorID(moderatorID, 0, 10)
	if err != nil {
		t.Fatalf("Error getting chat list: %v", err)
	}
	if len(chatList) != 2 {
		t.Fatalf("Chat list not correct: %v", chatList)
	}
	if chatList[0].ModeratorID != moderatorID || chatList[1].ModeratorID != moderatorID {
		t.Fatalf("Moderator id not updated: %v", chatList)
	}
	chatList, err = chatRepository.GetChatListByModeratorID(moderatorID+1, 0, 10)
	if err != nil {
		t.Fatalf("Error getting chat list: %v", err)
	}
	if len(chatList) != 0 {
		t.Fatalf("Chat list not correct: %v", chatList)
	}
}

func TestGetChatListByRepetitorIDCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
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
	moderatorRepository := CreateSqlModeratorRepository(db, "personal_data", "users", "moderators", "auth", "sequence")
	moderatorID, err := moderatorRepository.InsertModerator(tu.TestModeratorData, tu.TestPD, tu.TestAuthData)
	if err != nil {
		t.Fatalf("Error inserting repetitor: %v", err)
	}
	chatRepository := CreateSqlChatRepository(db, "chat", "sequence")
	_, err = chatRepository.InsertChat(types.DBChat{
		ClientID:    clientID,
		ModeratorID: moderatorID,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting chat: %v", err)
	}
	_, err = chatRepository.InsertChat(types.DBChat{
		ClientID:    clientID,
		RepetitorID: repetitorID,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting chat: %v", err)
	}
	_, err = chatRepository.InsertChat(types.DBChat{
		ModeratorID: moderatorID,
		RepetitorID: repetitorID,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting chat: %v", err)
	}
	chatList, err := chatRepository.GetChatListByRepetitorID(repetitorID, 0, 10)
	if err != nil {
		t.Fatalf("Error getting chat list: %v", err)
	}
	if len(chatList) != 2 {
		t.Fatalf("Chat list not correct: %v", chatList)
	}
	if chatList[0].RepetitorID != repetitorID || chatList[1].RepetitorID != repetitorID {
		t.Fatalf("Repetitor id not updated: %v", chatList)
	}
	chatList, err = chatRepository.GetChatListByRepetitorID(repetitorID+1, 0, 10)
	if err != nil {
		t.Fatalf("Error getting chat list: %v", err)
	}
	if len(chatList) != 0 {
		t.Fatalf("Chat list not correct: %v", chatList)
	}
}

func TestGetChatIdByCIDAndMIDCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
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
	chatID, err = chatRepository.GetChatIdByCIDAndMID(clientID+1, moderatorID)
	if err != nil {
		t.Fatalf("Error getting chat id: %v", err)
	}
	if chatID != 0 {
		t.Fatalf("Chat id is not 0: %v", chatID)
	}
}

func TestGetChatIdByCIDAndRIDCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
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
	chatID, err = chatRepository.GetChatIdByCIDAndRID(clientID+1, repetitorID)
	if err != nil {
		t.Fatalf("Error getting chat id: %v", err)
	}
	if chatID != 0 {
		t.Fatalf("Chat id is not 0: %v", chatID)
	}
}

func TestGetChatIdByMIDAndRIDCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
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
	chatID, err = chatRepository.GetChatIdByMIDAndRID(moderatorID+1, repetitorID)
	if err != nil {
		t.Fatalf("Error getting chat id: %v", err)
	}
	if chatID != 0 {
		t.Fatalf("Chat id is not 0: %v", chatID)
	}
}
