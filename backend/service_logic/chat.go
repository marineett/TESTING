package service_logic

import (
	"data_base_project/data_base"
	"data_base_project/types"
	"time"
)

type IChatService interface {
	CreateCRChat(clientID int64, repetitorID int64) (int64, error)
	CreateRMChat(repetitorID int64, moderatorID int64) (int64, error)
	CreateCMChat(clientID int64, moderatorID int64) (int64, error)
	GetChatListByClientID(clientID int64, from int64, size int64) ([]types.ServiceChat, error)
	GetChatListByRepetitorID(repetitorID int64, from int64, size int64) ([]types.ServiceChat, error)
	GetChatListByModeratorID(moderatorID int64, from int64, size int64) ([]types.ServiceChat, error)
	GetChat(chatID int64) (*types.ServiceChat, error)
	SendMessage(chatID int64, senderID int64, message string) error
	GetMessages(chatID int64, from int64, size int64) ([]types.ServiceMessage, error)
	GetChatIdByCIDAndMID(clientID int64, moderatorID int64) (int64, error)
	GetChatIdByCIDAndRID(clientID int64, repetitorID int64) (int64, error)
	GetChatIdByMIDAndRID(moderatorID int64, repetitorID int64) (int64, error)
}

type ChatService struct {
	chatRepository    data_base.IChatRepository
	messageRepository data_base.IMessageRepository
}

func CreateChatService(chatRepository data_base.IChatRepository, messageRepository data_base.IMessageRepository) IChatService {
	return &ChatService{
		chatRepository:    chatRepository,
		messageRepository: messageRepository,
	}
}

func (s *ChatService) CreateCRChat(clientID int64, repetitorID int64) (int64, error) {
	chatID, err := s.chatRepository.GetChatIdByCIDAndRID(clientID, repetitorID)
	if err != nil {
		return 0, err
	}
	if chatID != 0 {
		return chatID, nil
	}
	return s.chatRepository.InsertChat(*types.MapperChatServiceToDB(&types.ServiceChat{
		ClientID:    clientID,
		RepetitorID: repetitorID,
		CreatedAt:   time.Now(),
	}))
}

func (s *ChatService) CreateRMChat(repetitorID int64, moderatorID int64) (int64, error) {
	chatID, err := s.chatRepository.GetChatIdByMIDAndRID(moderatorID, repetitorID)
	if err != nil {
		return 0, err
	}
	if chatID != 0 {
		return chatID, nil
	}
	chat := types.DBChat{
		RepetitorID: repetitorID,
		ModeratorID: moderatorID,
		CreatedAt:   time.Now(),
	}
	return s.chatRepository.InsertChat(chat)
}

func (s *ChatService) CreateCMChat(clientID int64, moderatorID int64) (int64, error) {
	chatID, err := s.chatRepository.GetChatIdByCIDAndMID(clientID, moderatorID)
	if err != nil {
		return 0, err
	}
	if chatID != 0 {
		return chatID, nil
	}
	chat := types.DBChat{
		ClientID:    clientID,
		ModeratorID: moderatorID,
		CreatedAt:   time.Now(),
	}
	return s.chatRepository.InsertChat(chat)
}

func (s *ChatService) GetChatListByClientID(clientID int64, from int64, size int64) ([]types.ServiceChat, error) {
	chats, err := s.chatRepository.GetChatListByClientID(clientID, from, size)
	if err != nil {
		return nil, err
	}
	serviceChats := make([]types.ServiceChat, 0)
	for _, chat := range chats {
		serviceChats = append(serviceChats, *types.MapperChatDBToService(&chat))
	}
	return serviceChats, nil
}

func (s *ChatService) GetChatListByRepetitorID(repetitorID int64, from int64, size int64) ([]types.ServiceChat, error) {
	chats, err := s.chatRepository.GetChatListByRepetitorID(repetitorID, from, size)
	if err != nil {
		return nil, err
	}
	serviceChats := make([]types.ServiceChat, 0)
	for _, chat := range chats {
		serviceChats = append(serviceChats, *types.MapperChatDBToService(&chat))
	}
	return serviceChats, nil
}

func (s *ChatService) GetChatListByModeratorID(moderatorID int64, from int64, size int64) ([]types.ServiceChat, error) {
	chats, err := s.chatRepository.GetChatListByModeratorID(moderatorID, from, size)
	if err != nil {
		return nil, err
	}
	serviceChats := make([]types.ServiceChat, 0)
	for _, chat := range chats {
		serviceChats = append(serviceChats, *types.MapperChatDBToService(&chat))
	}
	return serviceChats, nil
}

func (s *ChatService) GetChat(chatID int64) (*types.ServiceChat, error) {
	chat, err := s.chatRepository.GetChat(chatID)
	if err != nil {
		return nil, err
	}
	return types.MapperChatDBToService(chat), nil
}

func (s *ChatService) SendMessage(chatID int64, senderID int64, message string) error {
	_, err := s.messageRepository.InsertMessage(types.DBMessage{
		ChatID:    chatID,
		SenderID:  senderID,
		Content:   message,
		CreatedAt: time.Now(),
	})
	return err
}

func (s *ChatService) GetMessages(chatID int64, from int64, size int64) ([]types.ServiceMessage, error) {
	messages, err := s.messageRepository.GetMessages(chatID, from, size)
	if err != nil {
		return nil, err
	}
	serviceMessages := make([]types.ServiceMessage, 0)
	for _, message := range messages {
		serviceMessages = append(serviceMessages, *types.MapperMessageDBToService(&message))
	}
	return serviceMessages, nil
}

func (s *ChatService) GetChatIdByCIDAndMID(clientID int64, moderatorID int64) (int64, error) {
	return s.chatRepository.GetChatIdByCIDAndMID(clientID, moderatorID)
}

func (s *ChatService) GetChatIdByCIDAndRID(clientID int64, repetitorID int64) (int64, error) {
	return s.chatRepository.GetChatIdByCIDAndRID(clientID, repetitorID)
}

func (s *ChatService) GetChatIdByMIDAndRID(moderatorID int64, repetitorID int64) (int64, error) {
	return s.chatRepository.GetChatIdByMIDAndRID(moderatorID, repetitorID)
}
