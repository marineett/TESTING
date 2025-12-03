package data_base

import (
	"data_base_project/types"
	"database/sql"
	"fmt"
)

type IMessageRepository interface {
	InsertMessage(message types.DBMessage) (int64, error)
	GetMessages(chatID int64, from int64, size int64) ([]types.DBMessage, error)
	DeleteMessages(chatID int64) error
	DeleteMessage(messageID int64) error
	UpdateMessageContent(messageID int64, content string) error
	GetMessage(messageID int64) (*types.DBMessage, error)
}

func CreateSqlMessageTable(db *sql.DB, messageTableName string, chatTableName string, userTableName string) error {
	query := `
		CREATE TABLE IF NOT EXISTS ` + messageTableName + ` (
		id INTEGER PRIMARY KEY,
		chat_id INTEGER NOT NULL,
		sender_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL,
		FOREIGN KEY (chat_id) REFERENCES ` + chatTableName + `(id),
		FOREIGN KEY (sender_id) REFERENCES ` + userTableName + `(id)
	)`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table %s: %v", messageTableName, err)
	}
	return nil
}

type SqlMessageRepository struct {
	db           *sql.DB
	messageTable string
	sequenceName string
}

func CreateSqlMessageRepository(db *sql.DB, messageTable string, sequenceName string) *SqlMessageRepository {
	return &SqlMessageRepository{
		db:           db,
		messageTable: messageTable,
		sequenceName: sequenceName,
	}
}

func (r *SqlMessageRepository) InsertMessage(message types.DBMessage) (int64, error) {
	var id int64
	err := r.db.QueryRow("SELECT nextval('" + r.sequenceName + "')").Scan(&id)
	if err != nil {
		return 0, err
	}
	query := `
	INSERT INTO ` + r.messageTable + ` (id, chat_id, sender_id, content, created_at)
	VALUES ($1, $2, $3, $4, $5)
	`
	_, err = r.db.Exec(query, id, message.ChatID, message.SenderID, message.Content, message.CreatedAt)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *SqlMessageRepository) GetMessages(chatID int64, from int64, size int64) ([]types.DBMessage, error) {
	query := `
	SELECT * FROM ` + r.messageTable + ` WHERE chat_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, chatID, size, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []types.DBMessage{}
	for rows.Next() {
		var message types.DBMessage
		err := rows.Scan(&message.ID, &message.ChatID, &message.SenderID, &message.Content, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (r *SqlMessageRepository) DeleteMessages(chatID int64) error {
	query := `
	DELETE FROM ` + r.messageTable + ` WHERE chat_id = $1
	`
	_, err := r.db.Exec(query, chatID)
	if err != nil {
		return err
	}
	return nil
}

func (r *SqlMessageRepository) DeleteMessage(messageID int64) error {
	query := `
	DELETE FROM ` + r.messageTable + ` WHERE id = $1
	`
	_, err := r.db.Exec(query, messageID)
	if err != nil {
		return err
	}
	return nil
}

func (r *SqlMessageRepository) UpdateMessageContent(messageID int64, content string) error {
	query := `
	UPDATE ` + r.messageTable + ` SET content = $1 WHERE id = $2
	`
	_, err := r.db.Exec(query, content, messageID)
	if err != nil {
		return err
	}
	return nil
}

func (r *SqlMessageRepository) GetMessage(messageID int64) (*types.DBMessage, error) {
	query := `
	SELECT * FROM ` + r.messageTable + ` WHERE id = $1
	`
	var message types.DBMessage
	err := r.db.QueryRow(query, messageID).Scan(&message.ID, &message.ChatID, &message.SenderID, &message.Content, &message.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &message, nil
}
