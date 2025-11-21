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

type ServerChatV2 struct {
	ID          int64     `json:"id"`
	Type        string    `json:"type"`
	ClientID    int64     `json:"client_id"`
	RepetitorID int64     `json:"repetitor_id"`
	ModeratorID int64     `json:"moderator_id"`
	CreatedAt   time.Time `json:"created_at"`
	Status      string    `json:"status"`
}

type ServerChatCreateV2 struct {
	Type        string `json:"type"`
	ClientID    *int64 `json:"client_id,omitempty"`
	RepetitorID *int64 `json:"repetitor_id,omitempty"`
	ModeratorID *int64 `json:"moderator_id,omitempty"`
}

type ServerChatUpdateV2 struct {
	Type   string `json:"type"`
	Status string `json:"status"`
}

type ServerChatStatusPatchV2 struct {
	Status string `json:"status"`
}

type ServerMessageV2 struct {
	ID        int64     `json:"id"`
	SenderID  int64     `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type ServerMessageCreateV2 struct {
	Content  string `json:"content"`
	SenderID int64  `json:"senderId"`
}

type ServerMessageContentPatchV2 struct {
	Content string `json:"content"`
}
