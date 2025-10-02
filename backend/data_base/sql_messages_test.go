package data_base

import (
	tu "data_base_project/test_database_utility"
	"data_base_project/types"
	"database/sql"
	"fmt"
	"testing"
	"time"
)

func setupMessagesTables(db *sql.DB) error {
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
	err = CreateSqlChatTable(db, "chat", "users")
	if err != nil {
		return fmt.Errorf("error creating chat table: %v", err)
	}
	err = CreateSqlMessageTable(db, "messages", "chat", "users")
	if err != nil {
		return fmt.Errorf("error creating message table: %v", err)
	}
	return nil
}

func TestCreateSqlMessageTable(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupMessagesTables(db)
	if err != nil {
		t.Fatalf("Error setting up messages tables: %v", err)
	}
	messageRepository := CreateSqlMessageRepository(db, "messages", "sequence")
	if messageRepository == nil {
		t.Fatalf("Error creating message repository: %v", err)
	}
}

func TestInsertMessageCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupMessagesTables(db)
	if err != nil {
		t.Fatalf("Error setting up messages tables: %v", err)
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
	chatID, err := chatRepository.InsertChat(types.DBChat{
		ClientID:    clientID,
		RepetitorID: repetitorID,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting chat: %v", err)
	}
	messageRepository := CreateSqlMessageRepository(db, "messages", "sequence")
	tu.TestMessage.ChatID = chatID
	tu.TestMessage.SenderID = clientID
	_, err = messageRepository.InsertMessage(tu.TestMessage)
	if err != nil {
		t.Fatalf("Error inserting message: %v", err)
	}
}

func TestInsertMessageIncorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupMessagesTables(db)
	if err != nil {
		t.Fatalf("Error setting up messages tables: %v", err)
	}
	messageRepository := CreateSqlMessageRepository(db, "messages", "sequence")
	_, err = messageRepository.InsertMessage(tu.TestMessage)
	if err == nil {
		t.Fatalf("No error inserting message: %v", err)
	}
}

func TestGetMessagesCorrect(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	err = setupMessagesTables(db)
	if err != nil {
		t.Fatalf("Error setting up messages tables: %v", err)
	}
	messageRepository := CreateSqlMessageRepository(db, "messages", "sequence")
	messages, err := messageRepository.GetMessages(1, 0, 10)
	if err != nil {
		t.Fatalf("Error getting messages: %v", err)
	}
	if len(messages) != 0 {
		t.Fatalf("Messages not found: %v", messages)
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
	chatID, err := chatRepository.InsertChat(types.DBChat{
		ClientID:    clientID,
		RepetitorID: repetitorID,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting chat: %v", err)
	}
	_, err = messageRepository.InsertMessage(types.DBMessage{
		ChatID:    chatID,
		SenderID:  clientID,
		Content:   "Hello2",
		CreatedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting message: %v", err)
	}
	_, err = messageRepository.InsertMessage(types.DBMessage{
		ChatID:    chatID,
		SenderID:  clientID,
		Content:   "Hello3",
		CreatedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting message: %v", err)
	}
	_, err = messageRepository.InsertMessage(types.DBMessage{
		ChatID:    chatID,
		SenderID:  clientID,
		Content:   "Hello4",
		CreatedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting message: %v", err)
	}
	messages, err = messageRepository.GetMessages(chatID, 0, 10)
	if err != nil {
		t.Fatalf("Error getting messages: %v", err)
	}
	if len(messages) != 3 {
		t.Fatalf("Messages are not in correct amount: %v", messages)
	}
}
