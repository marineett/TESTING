package types

import "time"

type DBChat struct {
	ID          int64     `json:"id"`
	ClientID    int64     `json:"client_id"`
	RepetitorID int64     `json:"repetitor_id"`
	ModeratorID int64     `json:"moderator_id"`
	CreatedAt   time.Time `json:"created_at"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
}

type DBMessage struct {
	ID        int64     `json:"id"`
	ChatID    int64     `json:"chat_id"`
	SenderID  int64     `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
