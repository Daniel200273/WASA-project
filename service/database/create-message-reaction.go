package database

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

// CreateMessageReaction creates a new reaction to a message
func (db *appdbimpl) CreateMessageReaction(messageID, userID, emoticon string) (*MessageReaction, error) {
	// First verify the message exists and get conversation ID
	var conversationID string
	err := db.c.QueryRow(`
		SELECT conversation_id FROM messages WHERE id = ?
	`, messageID).Scan(&conversationID)
	if err != nil {
		return nil, fmt.Errorf("message not found: %w", err)
	}

	// Verify user is participant in conversation
	isParticipant, err := db.IsUserInConversation(conversationID, userID)
	if err != nil {
		return nil, err
	}
	if !isParticipant {
		return nil, fmt.Errorf("user %s is not a participant in the conversation", userID)
	}

	// Check if user already reacted to this message (should update existing reaction)
	var existingReactionID string
	err = db.c.QueryRow(`
		SELECT id FROM message_reactions 
		WHERE message_id = ? AND user_id = ?
	`, messageID, userID).Scan(&existingReactionID)

	now := time.Now()

	if err == nil {
		// Update existing reaction
		_, err = db.c.Exec(`
			UPDATE message_reactions 
			SET emoticon = ?, created_at = ?
			WHERE id = ?
		`, emoticon, now, existingReactionID)
		if err != nil {
			return nil, fmt.Errorf("failed to update reaction: %w", err)
		}

		// Get username for response
		var username string
		err = db.c.QueryRow(`SELECT username FROM users WHERE id = ?`, userID).Scan(&username)
		if err != nil {
			return nil, fmt.Errorf("failed to get username: %w", err)
		}

		return &MessageReaction{
			ID:        existingReactionID,
			MessageID: messageID,
			UserID:    userID,
			Username:  username,
			Emoticon:  emoticon,
			CreatedAt: now,
		}, nil
	}

	// Create new reaction
	reactionID, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("failed to generate reaction ID: %w", err)
	}

	_, err = db.c.Exec(`
		INSERT INTO message_reactions (id, message_id, user_id, emoticon, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, reactionID.String(), messageID, userID, emoticon, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create reaction: %w", err)
	}

	// Get username for response
	var username string
	err = db.c.QueryRow(`SELECT username FROM users WHERE id = ?`, userID).Scan(&username)
	if err != nil {
		return nil, fmt.Errorf("failed to get username: %w", err)
	}

	return &MessageReaction{
		ID:        reactionID.String(),
		MessageID: messageID,
		UserID:    userID,
		Username:  username,
		Emoticon:  emoticon,
		CreatedAt: now,
	}, nil
}
