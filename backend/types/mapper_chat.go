package types

import "time"

func MapperChatDBToService(chat *DBChat) *ServiceChat {
	if chat == nil {
		return nil
	}
	return &ServiceChat{
		ID:          chat.ID,
		Type:        chat.Type,
		Status:      chat.Status,
		ClientID:    chat.ClientID,
		RepetitorID: chat.RepetitorID,
		ModeratorID: chat.ModeratorID,
		CreatedAt:   chat.CreatedAt,
	}
}

func MapperChatServiceToDB(chat *ServiceChat) *DBChat {
	if chat == nil {
		return nil
	}
	return &DBChat{
		ID:          chat.ID,
		ClientID:    chat.ClientID,
		RepetitorID: chat.RepetitorID,
		ModeratorID: chat.ModeratorID,
		CreatedAt:   chat.CreatedAt,
		Type:        chat.Type,
		Status:      chat.Status,
	}
}

func MapperMessageDBToService(message *DBMessage) *ServiceMessage {
	if message == nil {
		return nil
	}
	return &ServiceMessage{
		ID:        message.ID,
		ChatID:    message.ChatID,
		SenderID:  message.SenderID,
		Content:   message.Content,
		CreatedAt: message.CreatedAt,
	}
}

func MapperMessageServiceToDB(message *ServiceMessage) *DBMessage {
	if message == nil {
		return nil
	}
	return &DBMessage{
		ID:        message.ID,
		ChatID:    message.ChatID,
		SenderID:  message.SenderID,
		Content:   message.Content,
		CreatedAt: message.CreatedAt,
	}
}

func MapperChatServiceToServer(chat *ServiceChat) *ServerChat {
	if chat == nil {
		return nil
	}
	return &ServerChat{
		ID:          chat.ID,
		ClientID:    chat.ClientID,
		RepetitorID: chat.RepetitorID,
		ModeratorID: chat.ModeratorID,
		CreatedAt:   chat.CreatedAt,
	}
}

func MapperChatServerToService(chat *ServerChat) *ServiceChat {
	if chat == nil {
		return nil
	}
	return &ServiceChat{
		ID:          chat.ID,
		ClientID:    chat.ClientID,
		RepetitorID: chat.RepetitorID,
		ModeratorID: chat.ModeratorID,
		CreatedAt:   chat.CreatedAt,
	}
}

func MapperMessageServiceToServer(message *ServiceMessage) *ServerMessage {
	if message == nil {
		return nil
	}
	return &ServerMessage{
		SenderID:  message.SenderID,
		Content:   message.Content,
		CreatedAt: message.CreatedAt,
	}
}

func MapperMessageServerToService(message *ServerMessage) *ServiceMessage {
	if message == nil {
		return nil
	}
	return &ServiceMessage{
		SenderID:  message.SenderID,
		Content:   message.Content,
		CreatedAt: message.CreatedAt,
	}
}

func MapperChatServiceToServerV2(chat *ServiceChat) *ServerChatV2 {
	if chat == nil {
		return nil
	}
	return &ServerChatV2{
		ID:          chat.ID,
		Type:        chat.Type,
		ClientID:    chat.ClientID,
		RepetitorID: chat.RepetitorID,
		ModeratorID: chat.ModeratorID,
		CreatedAt:   chat.CreatedAt,
		Status:      chat.Status,
	}
}

func MapperChatServerV2ToService(chat *ServerChatV2) *ServiceChat {
	if chat == nil {
		return nil
	}
	return &ServiceChat{
		ID:          chat.ID,
		Type:        chat.Type,
		ClientID:    chat.ClientID,
		RepetitorID: chat.RepetitorID,
		ModeratorID: chat.ModeratorID,
		CreatedAt:   chat.CreatedAt,
		Status:      chat.Status,
	}
}

// --- V2 Message mappers ---
func MapperMessageServiceToServerV2(msg *ServiceMessage) *ServerMessageV2 {
	if msg == nil {
		return nil
	}
	return &ServerMessageV2{
		ID:        msg.ID,
		SenderID:  msg.SenderID,
		Content:   msg.Content,
		CreatedAt: msg.CreatedAt,
	}
}

func MapperMessageCreateV2ServerToService(chatID int64, senderID int64, req *ServerMessageCreateV2) *ServiceMessage {
	if req == nil {
		return nil
	}
	return &ServiceMessage{
		ChatID:    chatID,
		SenderID:  senderID,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}
}
