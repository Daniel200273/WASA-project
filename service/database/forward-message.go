package database

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

// ForwardMessage forwards an existing message to another conversation
func (db *appdbimpl) ForwardMessage(messageID, targetConversationID, userID string) (*Message, error) {
	// Get original message
	originalMsg, err := db.GetMessage(messageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get original message: %w", err)
	}

	// Verify user can access original message (is participant in original conversation)
	isParticipant, err := db.IsUserInConversation(originalMsg.ConversationID, userID)
	if err != nil {
		return nil, err
	}
	if !isParticipant {
		return nil, fmt.Errorf("user %s cannot access original message", userID)
	}

	// Verify user is participant in target conversation
	isTargetParticipant, err := db.IsUserInConversation(targetConversationID, userID)
	if err != nil {
		return nil, err
	}
	if !isTargetParticipant {
		return nil, fmt.Errorf("user %s is not a participant in target conversation", userID)
	}

	// Generate new message ID for forwarded message
	forwardedMessageID, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("failed to generate forwarded message ID: %w", err)
	}

	now := time.Now()

	// Insert forwarded message
	_, err = db.c.Exec(`
		INSERT INTO messages (id, conversation_id, sender_id, content, photo_url, forwarded, created_at)
		VALUES (?, ?, ?, ?, ?, true, ?)
	`, forwardedMessageID.String(), targetConversationID, userID, originalMsg.Content, originalMsg.PhotoURL, now)
	if err != nil {
		return nil, fmt.Errorf("failed to forward message: %w", err)
	}

	// Get sender username for response
	var senderUsername string
	err = db.c.QueryRow(`
		SELECT username FROM users WHERE id = ?
	`, userID).Scan(&senderUsername)
	if err != nil {
		return nil, fmt.Errorf("failed to get sender username: %w", err)
	}

	return &Message{
		ID:             forwardedMessageID.String(),
		ConversationID: targetConversationID,
		SenderID:       userID,
		SenderUsername: senderUsername,
		Content:        originalMsg.Content,
		PhotoURL:       originalMsg.PhotoURL,
		ReplyToID:      nil,
		Forwarded:      true,
		Status:         "sent",
		CreatedAt:      now,
	}, nil
}
