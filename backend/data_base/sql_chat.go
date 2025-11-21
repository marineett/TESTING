package data_base

import (
	"data_base_project/types"
	"database/sql"
	"fmt"
)

type IChatRepository interface {
	InsertChat(chat types.DBChat) (int64, error)
	GetChat(id int64) (*types.DBChat, error)
	GetChatListByClientID(clientID int64, from int64, size int64) ([]types.DBChat, error)
	GetChatListByRepetitorID(repetitorID int64, from int64, size int64) ([]types.DBChat, error)
	GetChatListByModeratorID(moderatorID int64, from int64, size int64) ([]types.DBChat, error)
	GetChatIdByCIDAndMID(clientID int64, moderatorID int64) (int64, error)
	GetChatIdByCIDAndRID(clientID int64, repetitorID int64) (int64, error)
	GetChatIdByMIDAndRID(moderatorID int64, repetitorID int64) (int64, error)
	DeleteChat(id int64) error
	UpdateChat(chatID int64, chatStatus string) error
}

func CreateSqlChatTable(db *sql.DB, chatTableName string, userTableName string) error {
	query := `
	CREATE TABLE IF NOT EXISTS ` + chatTableName + ` (
		id INTEGER PRIMARY KEY,
		client_id INTEGER NOT NULL,
		repetitor_id INTEGER NOT NULL,
		moderator_id INTEGER NOT NULL,
		created_at TIMESTAMP NOT NULL,
		type VARCHAR(255) NOT NULL,
		status VARCHAR(255) NOT NULL
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", chatTableName, err)
	}
	return nil
}

type SqlChatRepository struct {
	db           *sql.DB
	chatTable    string
	sequenceName string
}

func CreateSqlChatRepository(db *sql.DB, chatTable string, sequenceName string) *SqlChatRepository {
	return &SqlChatRepository{
		db:           db,
		chatTable:    chatTable,
		sequenceName: sequenceName,
	}
}

func (r *SqlChatRepository) InsertChat(chat types.DBChat) (int64, error) {
	var id int64
	err := r.db.QueryRow("SELECT nextval('" + r.sequenceName + "')").Scan(&id)
	if err != nil {
		return 0, err
	}
	query := `
	INSERT INTO ` + r.chatTable + ` (id, client_id, repetitor_id, moderator_id, created_at, type, status) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = r.db.Exec(query, id, chat.ClientID, chat.RepetitorID, chat.ModeratorID, chat.CreatedAt, chat.Type, chat.Status)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *SqlChatRepository) GetChat(id int64) (*types.DBChat, error) {
	query := `
	SELECT * FROM ` + r.chatTable + ` WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	var chat types.DBChat
	err := row.Scan(&chat.ID, &chat.ClientID, &chat.RepetitorID, &chat.ModeratorID, &chat.CreatedAt, &chat.Type, &chat.Status)
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

func (r *SqlChatRepository) GetChatListByClientID(clientID int64, from int64, size int64) ([]types.DBChat, error) {
	query := `
	SELECT * FROM ` + r.chatTable + ` WHERE client_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, clientID, size, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	chats := make([]types.DBChat, 0)
	for rows.Next() {
		var chat types.DBChat
		err := rows.Scan(&chat.ID, &chat.ClientID, &chat.RepetitorID, &chat.ModeratorID, &chat.CreatedAt, &chat.Type, &chat.Status)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, nil
}

func (r *SqlChatRepository) GetChatListByRepetitorID(repetitorID int64, from int64, size int64) ([]types.DBChat, error) {
	query := `
	SELECT * FROM ` + r.chatTable + ` WHERE repetitor_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, repetitorID, size, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	chats := make([]types.DBChat, 0)
	for rows.Next() {
		var chat types.DBChat
		err := rows.Scan(&chat.ID, &chat.ClientID, &chat.RepetitorID, &chat.ModeratorID, &chat.CreatedAt, &chat.Type, &chat.Status)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, nil
}

func (r *SqlChatRepository) GetChatListByModeratorID(moderatorID int64, from int64, size int64) ([]types.DBChat, error) {
	query := `
	SELECT * FROM ` + r.chatTable + ` WHERE moderator_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, moderatorID, size, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	chats := make([]types.DBChat, 0)
	for rows.Next() {
		var chat types.DBChat
		err := rows.Scan(&chat.ID, &chat.ClientID, &chat.RepetitorID, &chat.ModeratorID, &chat.CreatedAt, &chat.Type, &chat.Status)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, nil
}

func (r *SqlChatRepository) GetChatIdByCIDAndMID(clientID int64, moderatorID int64) (int64, error) {
	query := `
	SELECT id FROM ` + r.chatTable + ` WHERE client_id = $1 AND moderator_id = $2
	`
	row := r.db.QueryRow(query, clientID, moderatorID)
	var chatID int64
	err := row.Scan(&chatID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return chatID, nil
}

func (r *SqlChatRepository) GetChatIdByCIDAndRID(clientID int64, repetitorID int64) (int64, error) {
	query := `
	SELECT id FROM ` + r.chatTable + ` WHERE client_id = $1 AND repetitor_id = $2
	`
	row := r.db.QueryRow(query, clientID, repetitorID)
	var chatID int64
	err := row.Scan(&chatID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return chatID, nil
}

func (r *SqlChatRepository) GetChatIdByMIDAndRID(moderatorID int64, repetitorID int64) (int64, error) {
	query := `
	SELECT id FROM ` + r.chatTable + ` WHERE moderator_id = $1 AND repetitor_id = $2
	`
	row := r.db.QueryRow(query, moderatorID, repetitorID)
	var chatID int64
	err := row.Scan(&chatID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return chatID, nil
}

func (r *SqlChatRepository) DeleteChat(id int64) error {
	query := `
	DELETE FROM ` + r.chatTable + ` WHERE id = $1
	`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *SqlChatRepository) UpdateChat(chatID int64, chatStatus string) error {
	query := `
	UPDATE ` + r.chatTable + ` SET status = $1 WHERE id = $2
	`
	_, err := r.db.Exec(query, chatStatus, chatID)
	if err != nil {
		return err
	}
	return nil
}
