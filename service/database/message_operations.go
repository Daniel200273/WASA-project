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
	// 2. Validate that either content or photoURL is provided (not both null)
	// 3. If replyToID is provided, validate that the message exists in the same conversation
	// 4. Generate message ID and insert into database
	// 5. Return created message with sender username

	messageID := uuid.Must(uuid.NewV4()).String()

	// TODO: Add your implementation here
	_ = messageID // Remove this line when implementing

	return nil, fmt.Errorf("CreateMessage not implemented")
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
	// TODO: Implement message forwarding
	// 1. Verify user has access to source message
	// 2. Verify user can send messages to target conversation
	// 3. Get original message content/photo
	// 4. Create new message in target conversation with forwarded flag
	// 5. Return the new forwarded message

	forwardedMessageID := uuid.Must(uuid.NewV4()).String()

	// TODO: Add your implementation here
	_ = forwardedMessageID // Remove this line when implementing

	return nil, fmt.Errorf("ForwardMessage not implemented")
}
