package types

import "time"

type ServerChat struct {
	ID          int64     `json:"id"`
	ClientID    int64     `json:"client_id"`
	RepetitorID int64     `json:"repetitor_id"`
	ModeratorID int64     `json:"moderator_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type ServerMessage struct {
	SenderID  int64     `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
