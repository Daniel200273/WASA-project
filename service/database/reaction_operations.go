package database

import (
	"fmt"

	"github.com/gofrs/uuid"
)

// === REACTION OPERATIONS ===

// CreateMessageReaction creates a reaction to a message
func (db *appdbimpl) CreateMessageReaction(messageID, userID, emoticon string) (*MessageReaction, error) {
	// TODO: Implement message reaction creation
	// 1. Verify user has access to the message (is in the conversation)
	// 2. Check if user already reacted to this message (update existing or create new)
	// 3. Insert or update reaction in database
	// 4. Return created/updated reaction with username

	reactionID := uuid.Must(uuid.NewV4()).String()

	// TODO: Add your implementation here
	_ = reactionID // Remove this line when implementing

	return nil, fmt.Errorf("CreateMessageReaction not implemented")
}

// DeleteMessageReaction deletes a reaction (only by the user who created it)
func (db *appdbimpl) DeleteMessageReaction(reactionID, userID string) error {
	// TODO: Implement reaction deletion
	// 1. Verify that the user owns the reaction
	// 2. Delete the reaction from database
	// 3. Handle reaction not found or unauthorized cases

	return fmt.Errorf("DeleteMessageReaction not implemented")
}
