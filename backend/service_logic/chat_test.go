package service_logic

import (
	tu "data_base_project/test_service_utility"
	"data_base_project/types"
	"database/sql"
	"testing"
	"time"
)

func TestCreateCRChatCorrectLondon(t *testing.T) {
	chatRepository := tu.CreateTestChatRepository()
	messageRepository := tu.CreateTestMessageRepository()
	chatService := CreateChatService(chatRepository, messageRepository)
	chatID, err := chatService.CreateCRChat(1, 2)
	if err != nil {
		t.Fatalf("Error creating chat: %v", err)
	}
	if chatID != 1 {
		t.Fatalf("Chat id not updated: %v", chatID)
	}
	chat, err := chatRepository.GetChat(chatID)
	if err != nil {
		t.Fatalf("Error getting chat: %v", err)
	}
	if chat.ClientID != 1 {
		t.Fatalf("Client id not updated: %v", chat)
	}
	if chat.RepetitorID != 2 {
		t.Fatalf("Repetitor id not updated: %v", chat)
	}
	if chat.ModeratorID != 0 {
		t.Fatalf("Moderator id not null: %v", chat)
	}
}

func TestCreateCRChatCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up chat tables: %v", err)
	}
	chatRepository := module.ChatRepository
	messageRepository := module.MessageRepository
	chatService := CreateChatService(chatRepository, messageRepository)
	chatID, err := chatService.CreateCRChat(1, 2)
	if err != nil {
		t.Fatalf("Error creating chat: %v", err)
	}
	if chatID != 1 {
		t.Fatalf("Chat id not updated: %v", chatID)
	}
	chat, err := chatRepository.GetChat(chatID)
	if err != nil {
		t.Fatalf("Error getting chat: %v", err)
	}
	if chat.ClientID != 1 {
		t.Fatalf("Client id not updated: %v", chat)
	}
	if chat.RepetitorID != 2 {
		t.Fatalf("Repetitor id not updated: %v", chat)
	}
	if chat.ModeratorID != 0 {
		t.Fatalf("Moderator id not null: %v", chat)
	}
}

func TestCreateRMChatCorrectLondon(t *testing.T) {
	chatRepository := tu.CreateTestChatRepository()
	messageRepository := tu.CreateTestMessageRepository()
	chatService := CreateChatService(chatRepository, messageRepository)
	chatID, err := chatService.CreateRMChat(1, 2)
	if err != nil {
		t.Fatalf("Error creating chat: %v", err)
	}
	if chatID != 1 {
		t.Fatalf("Chat id not updated: %v", chatID)
	}
	chat, err := chatRepository.GetChat(chatID)
	if err != nil {
		t.Fatalf("Error getting chat: %v", err)
	}
	if chat.RepetitorID != 1 {
		t.Fatalf("Repetitor id not updated: %v", chat)
	}
	if chat.ModeratorID != 2 {
		t.Fatalf("Moderator id not updated: %v", chat)
	}
	if chat.ClientID != 0 {
		t.Fatalf("Client id not null: %v", chat)
	}
}

func TestCreateRMChatCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up chat tables: %v", err)
	}
	chatRepository := module.ChatRepository
	messageRepository := module.MessageRepository
	chatService := CreateChatService(chatRepository, messageRepository)
	chatID, err := chatService.CreateRMChat(1, 2)
	if err != nil {
		t.Fatalf("Error creating chat: %v", err)
	}
	if chatID != 1 {
		t.Fatalf("Chat id not updated: %v", chatID)
	}
	chat, err := chatRepository.GetChat(chatID)
	if err != nil {
		t.Fatalf("Error getting chat: %v", err)
	}
	if chat.RepetitorID != 1 {
		t.Fatalf("Repetitor id not updated: %v", chat)
	}
	if chat.ModeratorID != 2 {
		t.Fatalf("Moderator id not updated: %v", chat)
	}
	if chat.ClientID != 0 {
		t.Fatalf("Client id not null: %v", chat)
	}
}

func TestCreateCMChatCorrectLondon(t *testing.T) {
	chatRepository := tu.CreateTestChatRepository()
	messageRepository := tu.CreateTestMessageRepository()
	chatService := CreateChatService(chatRepository, messageRepository)
	chatID, err := chatService.CreateCMChat(1, 2)
	if err != nil {
		t.Fatalf("Error creating chat: %v", err)
	}
	if chatID != 1 {
		t.Fatalf("Chat id not updated: %v", chatID)
	}
	chat, err := chatRepository.GetChat(chatID)
	if err != nil {
		t.Fatalf("Error getting chat: %v", err)
	}
	if chat.ClientID != 1 {
		t.Fatalf("Client id not updated: %v", chat)
	}
	if chat.ModeratorID != 2 {
		t.Fatalf("Moderator id not updated: %v", chat)
	}
	if chat.RepetitorID != 0 {
		t.Fatalf("Repetitor id not null: %v", chat)
	}
}

func TestCreateCMChatCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up chat tables: %v", err)
	}
	chatRepository := module.ChatRepository
	messageRepository := module.MessageRepository
	chatService := CreateChatService(chatRepository, messageRepository)
	chatID, err := chatService.CreateCMChat(1, 2)
	if err != nil {
		t.Fatalf("Error creating chat: %v", err)
	}
	if chatID != 1 {
		t.Fatalf("Chat id not updated: %v", chatID)
	}
	chat, err := chatRepository.GetChat(chatID)
	if err != nil {
		t.Fatalf("Error getting chat: %v", err)
	}
	if chat.ClientID != 1 {
		t.Fatalf("Client id not updated: %v", chat)
	}
	if chat.ModeratorID != 2 {
		t.Fatalf("Moderator id not updated: %v", chat)
	}
	if chat.RepetitorID != 0 {
		t.Fatalf("Repetitor id not null: %v", chat)
	}
}

func TestGetChatListByClientIDCorrectLondon(t *testing.T) {
	chatRepository := tu.CreateTestChatRepository()
	messageRepository := tu.CreateTestMessageRepository()
	chatService := CreateChatService(chatRepository, messageRepository)
	chatRepository.InsertChat(types.DBChat{
		ClientID:    1,
		ModeratorID: 2,
		RepetitorID: 0,
		CreatedAt:   time.Now(),
	})
	chatRepository.InsertChat(types.DBChat{
		ClientID:    1,
		ModeratorID: 0,
		RepetitorID: 3,
		CreatedAt:   time.Now(),
	})
	chatRepository.InsertChat(types.DBChat{
		ClientID:    2,
		ModeratorID: 4,
		RepetitorID: 0,
		CreatedAt:   time.Now(),
	})
	chatList, err := chatService.GetChatListByClientID(1, 0, 10)
	if err != nil {
		t.Fatalf("Error getting chat list: %v", err)
	}
	if len(chatList) != 2 {
		t.Fatalf("Chat list not updated: %v", chatList)
	}
	if chatList[0].ClientID != 1 || chatList[1].ClientID != 1 {
		t.Fatalf("Client id not updated: %v", chatList[0])
	}
	chatList, err = chatService.GetChatListByClientID(5, 0, 10)
	if err != nil {
		t.Fatalf("Error getting chat list: %v", err)
	}
	if len(chatList) != 0 {
		t.Fatalf("Chat list not updated: %v", chatList)
	}
}

func TestGetChatListByClientIDCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up chat tables: %v", err)
	}
	chatRepository := module.ChatRepository
	messageRepository := module.MessageRepository
	chatService := CreateChatService(chatRepository, messageRepository)
	chatRepository.InsertChat(types.DBChat{
		ClientID:    1,
		ModeratorID: 2,
		RepetitorID: 0,
		CreatedAt:   time.Now(),
	})
	chatRepository.InsertChat(types.DBChat{
		ClientID:    1,
		ModeratorID: 0,
		RepetitorID: 3,
		CreatedAt:   time.Now(),
	})
	chatRepository.InsertChat(types.DBChat{
		ClientID:    2,
		ModeratorID: 4,
		RepetitorID: 0,
		CreatedAt:   time.Now(),
	})
	chatList, err := chatService.GetChatListByClientID(1, 0, 10)
	if err != nil {
		t.Fatalf("Error getting chat list: %v", err)
	}
	if len(chatList) != 2 {
		t.Fatalf("Chat list not updated: %v", chatList)
	}
	if chatList[0].ClientID != 1 || chatList[1].ClientID != 1 {
		t.Fatalf("Client id not updated: %v", chatList[0])
	}
	chatList, err = chatService.GetChatListByClientID(5, 0, 10)
	if err != nil {
		t.Fatalf("Error getting chat list: %v", err)
	}
	if len(chatList) != 0 {
		t.Fatalf("Chat list not updated: %v", chatList)
	}
}

func TestGetChatListByRepetitorIDCorrectLondon(t *testing.T) {
	chatRepository := tu.CreateTestChatRepository()
	messageRepository := tu.CreateTestMessageRepository()
	chatService := CreateChatService(chatRepository, messageRepository)
	chatRepository.InsertChat(types.DBChat{
		ClientID:    1,
		ModeratorID: 0,
		RepetitorID: 3,
		CreatedAt:   time.Now(),
	})
	chatRepository.InsertChat(types.DBChat{
		ClientID:    2,
		ModeratorID: 0,
		RepetitorID: 3,
		CreatedAt:   time.Now(),
	})
	chatRepository.InsertChat(types.DBChat{
		ClientID:    2,
		ModeratorID: 0,
		RepetitorID: 4,
		CreatedAt:   time.Now(),
	})
	chatList, err := chatService.GetChatListByRepetitorID(3, 0, 10)
	if err != nil {
		t.Fatalf("Error getting chat list: %v", err)
	}
	if len(chatList) != 2 {
		t.Fatalf("Chat list not updated: %v", chatList)
	}
	if chatList[0].RepetitorID != 3 || chatList[1].RepetitorID != 3 {
		t.Fatalf("Repetitor id not updated: %v", chatList[0])
	}
	chatList, err = chatService.GetChatListByRepetitorID(5, 0, 10)
	if err != nil {
		t.Fatalf("Error getting chat list: %v", err)
	}
	if len(chatList) != 0 {
		t.Fatalf("Chat list not updated: %v", chatList)
	}
}

func TestGetChatListByRepetitorIDCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up chat tables: %v", err)
	}
	chatRepository := module.ChatRepository
	messageRepository := module.MessageRepository
	chatService := CreateChatService(chatRepository, messageRepository)
	chatRepository.InsertChat(types.DBChat{
		ClientID:    1,
		ModeratorID: 0,
		RepetitorID: 3,
		CreatedAt:   time.Now(),
	})
	chatRepository.InsertChat(types.DBChat{
		ClientID:    2,
		ModeratorID: 0,
		RepetitorID: 3,
		CreatedAt:   time.Now(),
	})
	chatRepository.InsertChat(types.DBChat{
		ClientID:    2,
		ModeratorID: 0,
		RepetitorID: 4,
		CreatedAt:   time.Now(),
	})
	chatList, err := chatService.GetChatListByRepetitorID(3, 0, 10)
	if err != nil {
		t.Fatalf("Error getting chat list: %v", err)
	}
	if len(chatList) != 2 {
		t.Fatalf("Chat list not updated: %v", chatList)
	}
	if chatList[0].RepetitorID != 3 || chatList[1].RepetitorID != 3 {
		t.Fatalf("Repetitor id not updated: %v", chatList[0])
	}
	chatList, err = chatService.GetChatListByRepetitorID(5, 0, 10)
	if err != nil {
		t.Fatalf("Error getting chat list: %v", err)
	}
	if len(chatList) != 0 {
		t.Fatalf("Chat list not updated: %v", chatList)
	}
}

func TestGetChatListByModeratorIDCorrectLondon(t *testing.T) {
	chatRepository := tu.CreateTestChatRepository()
	messageRepository := tu.CreateTestMessageRepository()
	chatService := CreateChatService(chatRepository, messageRepository)
	chatRepository.InsertChat(types.DBChat{
		ClientID:    1,
		ModeratorID: 2,
		RepetitorID: 0,
		CreatedAt:   time.Now(),
	})
	chatRepository.InsertChat(types.DBChat{
		ClientID:    0,
		ModeratorID: 2,
		RepetitorID: 4,
		CreatedAt:   time.Now(),
	})
	chatRepository.InsertChat(types.DBChat{
		ClientID:    7,
		ModeratorID: 3,
		RepetitorID: 0,
		CreatedAt:   time.Now(),
	})
	chatList, err := chatService.GetChatListByModeratorID(2, 0, 10)
	if err != nil {
		t.Fatalf("Error getting chat list: %v", err)
	}
	if len(chatList) != 2 {
		t.Fatalf("Chat list not updated: %v", chatList)
	}
	if chatList[0].ModeratorID != 2 || chatList[1].ModeratorID != 2 {
		t.Fatalf("Moderator id not updated: %v", chatList[0])
	}
	chatList, err = chatService.GetChatListByModeratorID(5, 0, 10)
	if err != nil {
		t.Fatalf("Error getting chat list: %v", err)
	}
	if len(chatList) != 0 {
		t.Fatalf("Chat list not updated: %v", chatList)
	}
}

func TestGetChatListByModeratorIDCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up chat tables: %v", err)
	}
	chatRepository := module.ChatRepository
	messageRepository := module.MessageRepository
	chatService := CreateChatService(chatRepository, messageRepository)
	chatRepository.InsertChat(types.DBChat{
		ClientID:    1,
		ModeratorID: 2,
		RepetitorID: 0,
		CreatedAt:   time.Now(),
	})
	chatRepository.InsertChat(types.DBChat{
		ClientID:    0,
		ModeratorID: 2,
		RepetitorID: 4,
		CreatedAt:   time.Now(),
	})
	chatRepository.InsertChat(types.DBChat{
		ClientID:    7,
		ModeratorID: 3,
		RepetitorID: 0,
		CreatedAt:   time.Now(),
	})
	chatList, err := chatService.GetChatListByModeratorID(2, 0, 10)
	if err != nil {
		t.Fatalf("Error getting chat list: %v", err)
	}
	if len(chatList) != 2 {
		t.Fatalf("Chat list not updated: %v", chatList)
	}
	if chatList[0].ModeratorID != 2 || chatList[1].ModeratorID != 2 {
		t.Fatalf("Moderator id not updated: %v", chatList[0])
	}
	chatList, err = chatService.GetChatListByModeratorID(5, 0, 10)
	if err != nil {
		t.Fatalf("Error getting chat list: %v", err)
	}
	if len(chatList) != 0 {
		t.Fatalf("Chat list not updated: %v", chatList)
	}
}

func TestGetChatCorrectLondon(t *testing.T) {
	chatRepository := tu.CreateTestChatRepository()
	messageRepository := tu.CreateTestMessageRepository()
	chatService := CreateChatService(chatRepository, messageRepository)
	chatID, err := chatRepository.InsertChat(types.DBChat{
		ClientID:    1,
		ModeratorID: 2,
		RepetitorID: 0,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting chat: %v", err)
	}
	chat, err := chatService.GetChat(chatID)
	if err != nil {
		t.Fatalf("Error getting chat: %v", err)
	}
	if chat.ClientID != 1 {
		t.Fatalf("Client id not updated: %v", chat)
	}
	if chat.ModeratorID != 2 {
		t.Fatalf("Moderator id not updated: %v", chat)
	}
	if chat.RepetitorID != 0 {
		t.Fatalf("Repetitor id not null: %v", chat)
	}
}

func TestGetChatCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up chat tables: %v", err)
	}
	chatRepository := module.ChatRepository
	messageRepository := module.MessageRepository
	chatService := CreateChatService(chatRepository, messageRepository)
	chatID, err := chatRepository.InsertChat(types.DBChat{
		ClientID:    1,
		ModeratorID: 2,
		RepetitorID: 0,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting chat: %v", err)
	}
	chat, err := chatService.GetChat(chatID)
	if err != nil {
		t.Fatalf("Error getting chat: %v", err)
	}
	if chat.ClientID != 1 {
		t.Fatalf("Client id not updated: %v", chat)
	}
	if chat.ModeratorID != 2 {
		t.Fatalf("Moderator id not updated: %v", chat)
	}
	if chat.RepetitorID != 0 {
		t.Fatalf("Repetitor id not null: %v", chat)
	}
}

func TestSendMessageCorrectLondon(t *testing.T) {
	chatRepository := tu.CreateTestChatRepository()
	messageRepository := tu.CreateTestMessageRepository()
	chatService := CreateChatService(chatRepository, messageRepository)
	chatID, err := chatRepository.InsertChat(types.DBChat{
		ClientID:    1,
		ModeratorID: 2,
		RepetitorID: 0,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting chat: %v", err)
	}
	_, err = chatService.SendMessage(chatID, 1, "Hello")
	if err != nil {
		t.Fatalf("Error sending message: %v", err)
	}
	message, err := messageRepository.GetMessages(chatID, 0, 10)
	if err != nil {
		t.Fatalf("Error getting messages: %v", err)
	}
	if len(message) != 1 {
		t.Fatalf("Message not updated: %v", message)
	}
	if message[0].SenderID != 1 {
		t.Fatalf("Sender id not updated: %v", message[0])
	}
	if message[0].Content != "Hello" {
		t.Fatalf("Content not updated: %v", message[0])
	}
}

func TestSendMessageCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up chat tables: %v", err)
	}
	chatRepository := module.ChatRepository
	messageRepository := module.MessageRepository
	chatService := CreateChatService(chatRepository, messageRepository)
	clientRepository := module.ClientRepository
	personalDataRepository := module.PersonalDataRepository
	userRepository := module.UserRepository
	reviewRepository := module.ReviewRepository
	authRepository := module.AuthRepository
	clientService := CreateClientService(clientRepository, personalDataRepository, userRepository, reviewRepository)

	if err != nil {
		t.Fatalf("Error inserting auth: %v", err)
	}
	err = clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	chatID, err := chatRepository.InsertChat(types.DBChat{
		ClientID:    result.UserID,
		ModeratorID: 2,
		RepetitorID: 0,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting chat: %v", err)
	}
	_, err = chatService.SendMessage(chatID, result.UserID, "Hello")
	if err != nil {
		t.Fatalf("Error sending message: %v", err)
	}
	message, err := messageRepository.GetMessages(chatID, 0, 10)
	if err != nil {
		t.Fatalf("Error getting messages: %v", err)
	}
	if len(message) != 1 {
		t.Fatalf("Message not updated: %v", message)
	}
	if message[0].SenderID != result.UserID {
		t.Fatalf("Sender id not updated: %v", message[0])
	}
	if message[0].Content != "Hello" {
		t.Fatalf("Content not updated: %v", message[0])
	}
}

func TestGetMessagesCorrectLondon(t *testing.T) {
	chatRepository := tu.CreateTestChatRepository()
	messageRepository := tu.CreateTestMessageRepository()
	chatService := CreateChatService(chatRepository, messageRepository)
	chatID, err := chatRepository.InsertChat(types.DBChat{
		ClientID:    1,
		ModeratorID: 2,
		RepetitorID: 0,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting chat: %v", err)
	}
	messageID, err := messageRepository.InsertMessage(types.DBMessage{
		ChatID:    chatID,
		SenderID:  1,
		Content:   "Hello",
		CreatedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting message: %v", err)
	}
	message, err := chatService.GetMessages(chatID, 0, 10)
	if err != nil {
		t.Fatalf("Error getting messages: %v", err)
	}
	if len(message) != 1 {
		t.Fatalf("Message not updated: %v", message)
	}
	if message[0].ID != messageID {
		t.Fatalf("Message id not updated: %v", message[0])
	}
	if message[0].SenderID != 1 {
		t.Fatalf("Sender id not updated: %v", message[0])
	}
	if message[0].Content != "Hello" {
		t.Fatalf("Content not updated: %v", message[0])
	}
}

func TestGetMessagesCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up chat tables: %v", err)
	}
	chatRepository := module.ChatRepository
	messageRepository := module.MessageRepository
	clientRepository := module.ClientRepository
	personalDataRepository := module.PersonalDataRepository
	authRepository := module.AuthRepository
	userRepository := module.UserRepository
	reviewRepository := module.ReviewRepository
	chatService := CreateChatService(chatRepository, messageRepository)
	clientService := CreateClientService(clientRepository, personalDataRepository, userRepository, reviewRepository)
	err = clientService.CreateClient(tu.TestInitClientData)
	if err != nil {
		t.Fatalf("Error inserting client: %v", err)
	}
	result, err := authRepository.Authorize(types.DBAuthData{
		Login:    tu.TestAuth.Login,
		Password: tu.TestAuth.Password,
	})
	if err != nil {
		t.Fatalf("Error authorizing: %v", err)
	}
	chatID, err := chatRepository.InsertChat(types.DBChat{
		ClientID:    result.UserID,
		ModeratorID: 2,
		RepetitorID: 0,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting chat: %v", err)
	}
	messageID, err := messageRepository.InsertMessage(types.DBMessage{
		ChatID:    chatID,
		SenderID:  result.UserID,
		Content:   "Hello",
		CreatedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting message: %v", err)
	}
	message, err := chatService.GetMessages(chatID, 0, 10)
	if err != nil {
		t.Fatalf("Error getting messages: %v", err)
	}
	if len(message) != 1 {
		t.Fatalf("Message not updated: %v", message)
	}
	if message[0].ID != messageID {
		t.Fatalf("Message id not updated: %v", message[0])
	}
	if message[0].SenderID != result.UserID {
		t.Fatalf("Sender id not updated: %v", message[0])
	}
	if message[0].Content != "Hello" {
		t.Fatalf("Content not updated: %v", message[0])
	}
}

func TestGetMessagesIncorrectLondon(t *testing.T) {
	chatRepository := tu.CreateTestChatRepository()
	messageRepository := tu.CreateTestMessageRepository()
	chatService := CreateChatService(chatRepository, messageRepository)
	message, err := chatService.GetMessages(1, 0, 10)
	if err != nil {
		t.Fatalf("Error getting messages: %v", err)
	}
	if len(message) != 0 {
		t.Fatalf("Message not updated: %v", message)
	}
}

func TestGetMessagesIncorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up chat tables: %v", err)
	}
	chatRepository := module.ChatRepository
	messageRepository := module.MessageRepository
	chatService := CreateChatService(chatRepository, messageRepository)
	message, err := chatService.GetMessages(1, 0, 10)
	if err != nil {
		t.Fatalf("Error getting messages: %v", err)
	}
	if len(message) != 0 {
		t.Fatalf("Message not updated: %v", message)
	}
}

func TestGetChatIdByCIDAndMIDCorrectLondon(t *testing.T) {
	chatRepository := tu.CreateTestChatRepository()
	messageRepository := tu.CreateTestMessageRepository()
	chatService := CreateChatService(chatRepository, messageRepository)
	chatID, err := chatRepository.InsertChat(types.DBChat{
		ClientID:    1,
		ModeratorID: 2,
		RepetitorID: 0,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting chat: %v", err)
	}
	resultChatID, err := chatService.GetChatIdByCIDAndMID(1, 2)
	if err != nil {
		t.Fatalf("Error getting chat id: %v", err)
	}
	if chatID != resultChatID {
		t.Fatalf("Chat id not updated: %v", chatID)
	}
}

func TestGetChatIdByCIDAndMIDCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up chat tables: %v", err)
	}
	chatRepository := module.ChatRepository
	messageRepository := module.MessageRepository
	chatService := CreateChatService(chatRepository, messageRepository)
	chatID, err := chatRepository.InsertChat(types.DBChat{
		ClientID:    1,
		ModeratorID: 2,
		RepetitorID: 0,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting chat: %v", err)
	}
	resultChatID, err := chatService.GetChatIdByCIDAndMID(1, 2)
	if err != nil {
		t.Fatalf("Error getting chat id: %v", err)
	}
	if chatID != resultChatID {
		t.Fatalf("Chat id not updated: %v", chatID)
	}
}

func TestGetChatIdByCIDAndRIDCorrectLondon(t *testing.T) {
	chatRepository := tu.CreateTestChatRepository()
	messageRepository := tu.CreateTestMessageRepository()
	chatService := CreateChatService(chatRepository, messageRepository)
	chatID, err := chatRepository.InsertChat(types.DBChat{
		ClientID:    1,
		ModeratorID: 0,
		RepetitorID: 2,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting chat: %v", err)
	}
	resultChatID, err := chatService.GetChatIdByCIDAndRID(1, 2)
	if err != nil {
		t.Fatalf("Error getting chat id: %v", err)
	}
	if chatID != resultChatID {
		t.Fatalf("Chat id not updated: %v", chatID)
	}
}

func TestGetChatIdByCIDAndRIDCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up chat tables: %v", err)
	}
	chatRepository := module.ChatRepository
	messageRepository := module.MessageRepository
	chatService := CreateChatService(chatRepository, messageRepository)
	chatID, err := chatRepository.InsertChat(types.DBChat{
		ClientID:    1,
		ModeratorID: 0,
		RepetitorID: 2,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting chat: %v", err)
	}
	resultChatID, err := chatService.GetChatIdByCIDAndRID(1, 2)
	if err != nil {
		t.Fatalf("Error getting chat id: %v", err)
	}
	if chatID != resultChatID {
		t.Fatalf("Chat id not updated: %v", chatID)
	}
}

func TestGetChatIdByMIDAndRIDCorrectLondon(t *testing.T) {
	chatRepository := tu.CreateTestChatRepository()
	messageRepository := tu.CreateTestMessageRepository()
	chatService := CreateChatService(chatRepository, messageRepository)
	chatID, err := chatRepository.InsertChat(types.DBChat{
		ClientID:    0,
		ModeratorID: 1,
		RepetitorID: 2,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting chat: %v", err)
	}
	resultChatID, err := chatService.GetChatIdByMIDAndRID(1, 2)
	if err != nil {
		t.Fatalf("Error getting chat id: %v", err)
	}
	if chatID != resultChatID {
		t.Fatalf("Chat id not updated: %v", chatID)
	}
}

func TestGetChatIdByMIDAndRIDCorrectClassic(t *testing.T) {
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	module, err := tu.SetupModule(db)
	if err != nil {
		t.Fatalf("Error setting up chat tables: %v", err)
	}
	chatRepository := module.ChatRepository
	messageRepository := module.MessageRepository
	chatService := CreateChatService(chatRepository, messageRepository)
	chatID, err := chatRepository.InsertChat(types.DBChat{
		ClientID:    0,
		ModeratorID: 1,
		RepetitorID: 2,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		t.Fatalf("Error inserting chat: %v", err)
	}
	resultChatID, err := chatService.GetChatIdByMIDAndRID(1, 2)
	if err != nil {
		t.Fatalf("Error getting chat id: %v", err)
	}
	if chatID != resultChatID {
		t.Fatalf("Chat id not updated: %v", chatID)
	}
}
