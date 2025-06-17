package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// GetMessage returns a specific message by ID
func (db *appdbimpl) GetMessage(messageID string) (*Message, error) {
	query := `
		SELECT m.id, m.conversation_id, m.sender_id, m.content, m.photo_url, 
		       m.reply_to_id, m.forwarded, m.created_at, u.username
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.id = ?
	`

	var msg Message
	var content, photoURL, replyToID sql.NullString
	err := db.c.QueryRow(query, messageID).Scan(
		&msg.ID, &msg.ConversationID, &msg.SenderID, &content, &photoURL,
		&replyToID, &msg.Forwarded, &msg.CreatedAt, &msg.SenderUsername,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("message not found")
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	// Handle nullable fields
	if content.Valid {
		msg.Content = &content.String
	}
	if photoURL.Valid {
		msg.PhotoURL = &photoURL.String
	}
	if replyToID.Valid {
		msg.ReplyToID = &replyToID.String
	}

	msg.Status = "sent" // Default status

	return &msg, nil
}
