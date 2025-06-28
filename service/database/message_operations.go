package database

import (
	"fmt"

	"github.com/gofrs/uuid"
)

// === MESSAGE OPERATIONS ===

// CreateMessage creates a new message in a conversation
func (db *appdbimpl) CreateMessage(conversationID, senderID string, content *string, photoURL *string, replyToID *string) (*Message, error) {
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
		replyMessage, err := db.GetMessage(*replyToID)
		if err != nil {
			return nil, fmt.Errorf("reply target message not found: %w", err)
		}
		if replyMessage.ConversationID != conversationID {
			return nil, fmt.Errorf("cannot reply to message from different conversation")
		}
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
	// Query message from database by ID with sender username
	query := `
		SELECT m.id, m.conversation_id, m.sender_id, u.username, m.content, 
			   m.photo_url, m.reply_to_id, m.forwarded, m.created_at
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.id = ?
	`

	row := db.c.QueryRow(query, messageID)

	var msg Message
	err := row.Scan(
		&msg.ID,
		&msg.ConversationID,
		&msg.SenderID,
		&msg.SenderUsername,
		&msg.Content,
		&msg.PhotoURL,
		&msg.ReplyToID,
		&msg.Forwarded,
		&msg.CreatedAt,
	)

	if err != nil {
		if isNotFoundError(err) {
			return nil, fmt.Errorf("message not found")
		}
		return nil, fmt.Errorf("error retrieving message: %w", err)
	}

	// Set message status (for now, just set as "sent")
	msg.Status = "sent"

	// Get reactions/comments for this message
	msg.Comments, err = db.getMessageReactions(msg.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting message reactions: %w", err)
	}

	return &msg, nil
}

// DeleteMessage deletes a message (only by the sender)
func (db *appdbimpl) DeleteMessage(messageID, userID string) error {
	// 1. Verify that the user is the sender of the message
	query := `SELECT sender_id FROM messages WHERE id = ?`
	var senderID string
	err := db.c.QueryRow(query, messageID).Scan(&senderID)
	if err != nil {
		if isNotFoundError(err) {
			return fmt.Errorf("message not found")
		}
		return fmt.Errorf("error checking message ownership: %w", err)
	}

	if senderID != userID {
		return fmt.Errorf("unauthorized: user can only delete their own messages")
	}

	// 2. Delete the message from database
	// Note: Reactions and replies will be handled by cascade DELETE constraints
	deleteQuery := `DELETE FROM messages WHERE id = ?`
	result, err := db.c.Exec(deleteQuery, messageID)
	if err != nil {
		return fmt.Errorf("error deleting message: %w", err)
	}

	// 3. Verify deletion was successful
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking deletion result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("message not found or already deleted")
	}

	return nil
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

// GetConversationMessages retrieves all messages in a conversation
func (db *appdbimpl) GetConversationMessages(conversationID string) ([]Message, error) {
	// Query to get all messages in the conversation with sender usernames
	query := `
		SELECT m.id, m.conversation_id, m.sender_id, u.username, m.content, 
			   m.photo_url, m.reply_to_id, m.forwarded, m.created_at
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.conversation_id = ?
		ORDER BY m.created_at ASC
	`

	rows, err := db.c.Query(query, conversationID)
	if err != nil {
		return nil, fmt.Errorf("error querying conversation messages: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(
			&msg.ID,
			&msg.ConversationID,
			&msg.SenderID,
			&msg.SenderUsername,
			&msg.Content,
			&msg.PhotoURL,
			&msg.ReplyToID,
			&msg.Forwarded,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning message: %w", err)
		}

		// Set message status (for now, just set as "sent" - this would be enhanced with read receipts)
		msg.Status = "sent"

		// Get reactions/comments for this message
		msg.Comments, err = db.getMessageReactions(msg.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting message reactions: %w", err)
		}

		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over messages: %w", err)
	}

	return messages, nil
}

// getMessageReactions retrieves all reactions for a specific message
func (db *appdbimpl) getMessageReactions(messageID string) ([]MessageReaction, error) {
	query := `
		SELECT mr.id, mr.message_id, mr.user_id, u.username, mr.emoticon, mr.created_at
		FROM message_reactions mr
		JOIN users u ON mr.user_id = u.id
		WHERE mr.message_id = ?
		ORDER BY mr.created_at ASC
	`

	rows, err := db.c.Query(query, messageID)
	if err != nil {
		return nil, fmt.Errorf("error querying message reactions: %w", err)
	}
	defer rows.Close()

	var reactions []MessageReaction
	for rows.Next() {
		var reaction MessageReaction
		err := rows.Scan(
			&reaction.ID,
			&reaction.MessageID,
			&reaction.UserID,
			&reaction.Username,
			&reaction.Emoticon,
			&reaction.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning reaction: %w", err)
		}
		reactions = append(reactions, reaction)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over reactions: %w", err)
	}

	return reactions, nil
}

// === SIMPLE READ STATUS OPERATIONS ===

// MarkConversationAsRead updates the last_read_at timestamp for a user in a conversation
func (db *appdbimpl) MarkConversationAsRead(conversationID, userID string) error {
	// Update the last_read_at timestamp for this user in this conversation
	query := `
		UPDATE conversation_participants 
		SET last_read_at = CURRENT_TIMESTAMP 
		WHERE conversation_id = ? AND user_id = ?
	`
	result, err := db.c.Exec(query, conversationID, userID)
	if err != nil {
		return fmt.Errorf("error marking conversation as read: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking update result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user is not a participant in this conversation")
	}

	return nil
}
