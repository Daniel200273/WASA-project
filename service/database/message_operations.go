package database

import (
	"fmt"

	"github.com/gofrs/uuid"
)

// === MESSAGE OPERATIONS ===

// CreateMessage creates a new message in a conversation
func (db *appdbimpl) CreateMessage(conversationID, senderID string, content *string, photoURL *string, replyToID *string) (*Message, error) {
	// TODO: Implement message creation
	// 1. Validate that sender is a participant in the conversation
	isParticipant, err := db.IsUserInConversation(conversationID, senderID)
	if err != nil {
		return nil, fmt.Errorf("error checking conversation participation: %w", err)
	}
	if !isParticipant {
		return nil, fmt.Errorf("user is not a participant in this conversation")
	}

	// 2. Validate that either content or photoURL is provided (not both null)
	if (content == nil && photoURL == nil) || (content != nil && photoURL != nil) {
		return nil, fmt.Errorf("must provide either content or photo, not both or neither")
	}

	// 3. If replyToID is provided, validate that the message exists in the same conversation
	if replyToID != nil && *replyToID != "" {
		// Validation logic for reply would go here
	}

	// 4. Generate message ID and insert into database
	messageID := uuid.Must(uuid.NewV4()).String()

	// Begin transaction to ensure both operations complete together
	tx, err := db.c.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}

	// Insert the message
	messageQuery := `
		INSERT INTO messages (id, conversation_id, sender_id, content, photo_url, reply_to_id, forwarded, created_at)
		VALUES (?, ?, ?, ?, ?, ?, FALSE, CURRENT_TIMESTAMP)
	`
	_, err = tx.Exec(messageQuery, messageID, conversationID, senderID, content, photoURL, replyToID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating message: %w", err)
	}

	// Update the conversation's last_message_at field
	updateQuery := `
		UPDATE conversations 
		SET last_message_at = CURRENT_TIMESTAMP 
		WHERE id = ?
	`
	_, err = tx.Exec(updateQuery, conversationID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error updating conversation last_message_at: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	// 5. Return created message with sender username
	return db.GetMessage(messageID)
}

// GetMessage retrieves a message by its ID
func (db *appdbimpl) GetMessage(messageID string) (*Message, error) {
	// TODO: Implement message retrieval
	// 1. Query message from database by ID
	// 2. Join with users table to get sender username
	// 3. Handle message not found case
	// 4. Return message or error

	return nil, fmt.Errorf("GetMessage not implemented")
}

// DeleteMessage deletes a message (only by the sender)
func (db *appdbimpl) DeleteMessage(messageID, userID string) error {
	// TODO: Implement message deletion
	// 1. Verify that the user is the sender of the message
	// 2. Delete the message from database
	// 3. Handle message not found or unauthorized cases
	// 4. Note: Consider cascade effects on reactions and replies

	return fmt.Errorf("DeleteMessage not implemented")
}

// ForwardMessage forwards a message to another conversation
func (db *appdbimpl) ForwardMessage(messageID, targetConversationID, userID string) (*Message, error) {
	// 1. Verify user has access to source message
	// This would require checking if the user can access the conversation containing the message

	// 2. Verify user can send messages to target conversation
	isParticipant, err := db.IsUserInConversation(targetConversationID, userID)
	if err != nil {
		return nil, fmt.Errorf("error checking conversation participation: %w", err)
	}
	if !isParticipant {
		return nil, fmt.Errorf("user is not a participant in the target conversation")
	}

	// 3. Get original message content/photo
	originalMessage, err := db.GetMessage(messageID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving original message: %w", err)
	}

	// 4. Create new message in target conversation with forwarded flag
	forwardedMessageID := uuid.Must(uuid.NewV4()).String()

	// Begin transaction
	tx, err := db.c.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}

	// Insert the forwarded message
	insertQuery := `
		INSERT INTO messages (id, conversation_id, sender_id, content, photo_url, reply_to_id, forwarded, created_at)
		VALUES (?, ?, ?, ?, ?, NULL, TRUE, CURRENT_TIMESTAMP)
	`
	_, err = tx.Exec(insertQuery, forwardedMessageID, targetConversationID, userID,
		originalMessage.Content, originalMessage.PhotoURL)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating forwarded message: %w", err)
	}

	// Update the conversation's last_message_at field
	updateQuery := `
		UPDATE conversations 
		SET last_message_at = CURRENT_TIMESTAMP 
		WHERE id = ?
	`
	_, err = tx.Exec(updateQuery, targetConversationID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error updating conversation last_message_at: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	// 5. Return the new forwarded message
	return db.GetMessage(forwardedMessageID)
}
