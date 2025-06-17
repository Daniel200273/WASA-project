package database

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

// CreateMessage creates a new message in a conversation
func (db *appdbimpl) CreateMessage(conversationID, senderID string, content *string, photoURL *string, replyToID *string) (*Message, error) {
	// Verify sender is participant in conversation
	isParticipant, err := db.IsUserInConversation(conversationID, senderID)
	if err != nil {
		return nil, err
	}
	if !isParticipant {
		return nil, fmt.Errorf("user %s is not a participant in conversation %s", senderID, conversationID)
	}

	// Generate message ID
	messageID, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("failed to generate message ID: %w", err)
	}

	now := time.Now()

	// Insert message
	_, err = db.c.Exec(`
		INSERT INTO messages (id, conversation_id, sender_id, content, photo_url, reply_to_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, messageID.String(), conversationID, senderID, content, photoURL, replyToID, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	// Get sender username for response
	var senderUsername string
	err = db.c.QueryRow(`
		SELECT username FROM users WHERE id = ?
	`, senderID).Scan(&senderUsername)
	if err != nil {
		return nil, fmt.Errorf("failed to get sender username: %w", err)
	}

	return &Message{
		ID:             messageID.String(),
		ConversationID: conversationID,
		SenderID:       senderID,
		SenderUsername: senderUsername,
		Content:        content,
		PhotoURL:       photoURL,
		ReplyToID:      replyToID,
		Forwarded:      false,
		Status:         "sent",
		CreatedAt:      now,
	}, nil
}
