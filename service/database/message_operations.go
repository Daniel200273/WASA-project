package database

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

// Constants for message status
const (
	MessageStatusSent = "sent"
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
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("error creating message: %w (rollback failed: %w)", err, rollbackErr)
		}
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
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("error updating conversation last_message_at: %w (rollback failed: %w)", err, rollbackErr)
		}
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
	msg.Status = MessageStatusSent

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
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("error creating forwarded message: %w (rollback failed: %w)", err, rollbackErr)
		}
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
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("error updating conversation last_message_at: %w (rollback failed: %w)", err, rollbackErr)
		}
		return nil, fmt.Errorf("error updating conversation last_message_at: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	// 5. Return the new forwarded message
	return db.GetMessage(forwardedMessageID)
}

// GetConversationMessages retrieves all messages in a conversation with read status
func (db *appdbimpl) GetConversationMessages(conversationID, currentUserID string) ([]Message, error) {
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

		// Determine message status based on read receipts
		msg.Status = db.calculateMessageStatus(msg, conversationID, currentUserID)

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

// calculateMessageStatus determines if a message has been read by checking last_read_at timestamps
func (db *appdbimpl) calculateMessageStatus(msg Message, conversationID, currentUserID string) string {
	// If the message is not from the current user, status is always "sent" from their perspective
	if msg.SenderID != currentUserID {
		return MessageStatusSent
	}

	// For messages sent by the current user, check if all other participants have read them
	// First, get the conversation type to handle group vs direct differently
	conversationType, err := db.getConversationType(conversationID)
	if err != nil {
		return MessageStatusSent
	}

	if conversationType == "group" {
		// For group chats, check if ALL members (except sender) have read the message
		// First, let's get the total count of active participants (excluding sender)
		totalQuery := `
			SELECT COUNT(*)
			FROM conversation_participants cp
			WHERE cp.conversation_id = ? AND cp.user_id != ?
		`

		var totalOthers int
		err := db.c.QueryRow(totalQuery, conversationID, currentUserID).Scan(&totalOthers)
		if err != nil || totalOthers == 0 {
			return MessageStatusSent
		}

		// Now check how many have read the message after it was created
		readQuery := `
			SELECT COUNT(*)
			FROM conversation_participants cp
			WHERE cp.conversation_id = ? 
			  AND cp.user_id != ? 
			  AND cp.last_read_at IS NOT NULL 
			  AND cp.last_read_at > ?
		`

		var othersRead int
		err = db.c.QueryRow(readQuery, conversationID, currentUserID, msg.CreatedAt).Scan(&othersRead)
		if err != nil {
			return MessageStatusSent
		}

		// If all other participants have read the message, show double check
		if othersRead == totalOthers {
			return "read"
		}
	} else {
		// For direct messages, check if the other person has read it
		query := `
			SELECT last_read_at
			FROM conversation_participants 
			WHERE conversation_id = ? AND user_id != ?
		`

		var readAt *time.Time
		err := db.c.QueryRow(query, conversationID, currentUserID).Scan(&readAt)
		if err != nil {
			return MessageStatusSent
		}

		// If last_read_at is NULL, the other user hasn't read anything yet
		if readAt == nil {
			return MessageStatusSent
		}

		// If the other participant's read time is after this message was created, it's been read
		// We use strictly After() because Equal() would mean they read exactly when message was created,
		// which is unlikely and shouldn't count as "read"
		if readAt.After(msg.CreatedAt) {
			return "read"
		}
	}

	// Otherwise, it's just been sent
	return MessageStatusSent
}

// GetConversationMessagesAfter retrieves messages in a conversation created after a specific timestamp
func (db *appdbimpl) GetConversationMessagesAfter(conversationID, currentUserID, afterTimestamp string, limit int) ([]Message, error) {
	// Parse the afterTimestamp to validate it
	_, err := time.Parse(time.RFC3339, afterTimestamp)
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp format: %w", err)
	}

	// Query to get messages created after the specified timestamp
	query := `
		SELECT m.id, m.conversation_id, m.sender_id, u.username, m.content, 
			   m.photo_url, m.reply_to_id, m.forwarded, m.created_at
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.conversation_id = ? AND m.created_at > ?
		ORDER BY m.created_at ASC
		LIMIT ?
	`

	rows, err := db.c.Query(query, conversationID, afterTimestamp, limit)
	if err != nil {
		return nil, fmt.Errorf("error querying conversation messages after timestamp: %w", err)
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

		// Determine message status based on read receipts
		msg.Status = db.calculateMessageStatus(msg, conversationID, currentUserID)

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

// IsReadByAllGroupMembers checks if all members of a group conversation have read the conversation after its last message
func (db *appdbimpl) IsReadByAllGroupMembers(conversationID string, lastMessageAt time.Time) (bool, error) {
	// First, check if this is a group conversation
	query := `SELECT type FROM conversations WHERE id = ?`
	var conversationType string
	err := db.c.QueryRow(query, conversationID).Scan(&conversationType)
	if err != nil {
		return false, fmt.Errorf("error checking conversation type: %w", err)
	}

	// If it's not a group conversation, return false
	if conversationType != "group" {
		return false, nil
	}

	// Check if all group members have a last_read_at timestamp after the last message timestamp
	query = `
		SELECT COUNT(*) as total_members,
		       COUNT(CASE WHEN cp.last_read_at >= ? THEN 1 END) as members_read
		FROM conversation_participants cp
		WHERE cp.conversation_id = ?
	`
	var totalMembers, membersRead int
	err = db.c.QueryRow(query, lastMessageAt, conversationID).Scan(&totalMembers, &membersRead)
	if err != nil {
		return false, fmt.Errorf("error checking group read status: %w", err)
	}

	// All members have read if members_read equals total_members
	return totalMembers > 0 && membersRead == totalMembers, nil
}

// getConversationType is a helper method to get the type of a conversation
func (db *appdbimpl) getConversationType(conversationID string) (string, error) {
	query := `SELECT type FROM conversations WHERE id = ?`
	var conversationType string
	err := db.c.QueryRow(query, conversationID).Scan(&conversationType)
	if err != nil {
		return "", fmt.Errorf("error getting conversation type: %w", err)
	}
	return conversationType, nil
}
