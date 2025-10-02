package types

func MapperChatDBToService(chat *DBChat) *ServiceChat {
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
