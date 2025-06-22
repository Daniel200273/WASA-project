package database

import (
	"database/sql"
	"errors"
	"time"
)

// === DATABASE UTILITIES ===
//
// These helper functions provide common database operations for converting
// database rows into Go structs and handling common database tasks.

// === USER SCANNING FUNCTIONS ===

// scanUser converts a single database row into a User struct.
// Used after QueryRow() calls to map database columns to User fields.
// Expected column order: id, username, photo_url, created_at
func scanUser(row *sql.Row) (*User, error) {
	var user User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.PhotoURL,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// scanUsers converts multiple database rows into a slice of User structs.
// Used after Query() calls when retrieving multiple users.
// Expected column order: id, username, photo_url, created_at
// Automatically handles closing rows and iterating through all records.
func scanUsers(rows *sql.Rows) ([]User, error) {
	var users []User
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.PhotoURL,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// === CONVERSATION SCANNING FUNCTIONS ===

// scanConversation converts a single database row into a Conversation struct.
// Used after QueryRow() calls to map database columns to Conversation fields.
// Expected column order: id, type, name, photo_url, created_by, created_at
func scanConversation(row *sql.Row) (*Conversation, error) {
	var conv Conversation
	err := row.Scan(
		&conv.ID,
		&conv.Type,
		&conv.Name,
		&conv.PhotoURL,
		&conv.CreatedBy,
		&conv.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

// scanConversations converts multiple database rows into a slice of Conversation structs.
// Used after Query() calls when retrieving multiple conversations.
// Expected column order: id, type, name, photo_url, created_by, created_at
// Automatically handles closing rows and iterating through all records.
func scanConversations(rows *sql.Rows) ([]Conversation, error) {
	var conversations []Conversation
	defer rows.Close()

	for rows.Next() {
		var conv Conversation
		err := rows.Scan(
			&conv.ID,
			&conv.Type,
			&conv.Name,
			&conv.PhotoURL,
			&conv.CreatedBy,
			&conv.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, conv)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return conversations, nil
}

// === MESSAGE SCANNING FUNCTIONS ===

// scanMessage converts a single database row into a Message struct.
// Used after QueryRow() calls to map database columns to Message fields.
// Expected column order: id, conversation_id, sender_id, sender_username, content, photo_url, reply_to_id, forwarded, created_at
// Also extracts sender username and sets default status to "sent"
func scanMessage(row *sql.Row) (*Message, error) {
	var msg Message
	var senderUsername string
	err := row.Scan(
		&msg.ID,
		&msg.ConversationID,
		&msg.SenderID,
		&senderUsername,
		&msg.Content,
		&msg.PhotoURL,
		&msg.ReplyToID,
		&msg.Forwarded,
		&msg.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	msg.SenderUsername = senderUsername
	msg.Status = "sent" // Default status
	return &msg, nil
}

// scanMessages converts multiple database rows into a slice of Message structs.
// Used after Query() calls when retrieving multiple messages.
// Expected column order: id, conversation_id, sender_id, sender_username, content, photo_url, reply_to_id, forwarded, created_at
// Also extracts sender username and sets default status to "sent"
// Automatically handles closing rows and iterating through all records.
func scanMessages(rows *sql.Rows) ([]Message, error) {
	var messages []Message
	defer rows.Close()

	for rows.Next() {
		var msg Message
		var senderUsername string
		err := rows.Scan(
			&msg.ID,
			&msg.ConversationID,
			&msg.SenderID,
			&senderUsername,
			&msg.Content,
			&msg.PhotoURL,
			&msg.ReplyToID,
			&msg.Forwarded,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		msg.SenderUsername = senderUsername
		msg.Status = "sent" // Default status
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// === MESSAGE REACTION SCANNING FUNCTIONS ===

// scanMessageReaction converts a single database row into a MessageReaction struct.
// Used after QueryRow() calls to map database columns to MessageReaction fields.
// Expected column order: id, message_id, user_id, username, emoticon, created_at
// Also extracts username of the user who reacted
func scanMessageReaction(row *sql.Row) (*MessageReaction, error) {
	var reaction MessageReaction
	var username string
	err := row.Scan(
		&reaction.ID,
		&reaction.MessageID,
		&reaction.UserID,
		&username,
		&reaction.Emoticon,
		&reaction.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	reaction.Username = username
	return &reaction, nil
}

// scanMessageReactions converts multiple database rows into a slice of MessageReaction structs.
// Used after Query() calls when retrieving multiple message reactions.
// Expected column order: id, message_id, user_id, username, emoticon, created_at
// Also extracts username of the user who reacted
// Automatically handles closing rows and iterating through all records.
func scanMessageReactions(rows *sql.Rows) ([]MessageReaction, error) {
	var reactions []MessageReaction
	defer rows.Close()

	for rows.Next() {
		var reaction MessageReaction
		var username string
		err := rows.Scan(
			&reaction.ID,
			&reaction.MessageID,
			&reaction.UserID,
			&username,
			&reaction.Emoticon,
			&reaction.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		reaction.Username = username
		reactions = append(reactions, reaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reactions, nil
}

// === UTILITY FUNCTIONS ===

// isNotFoundError checks if an error is sql.ErrNoRows (record not found).
// Used to distinguish between "not found" vs actual database errors.
// Returns true if the error indicates no rows were found in the query result.
func isNotFoundError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

// formatTimestamp converts Go time.Time to SQLite format string.
// SQLite uses the format "2006-01-02 15:04:05" for datetime storage.
// This function helps when you need to format timestamps for SQL queries.
func formatTimestamp(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// parseTimestamp converts SQLite timestamp string back to Go time.Time.
// This function helps when you need to parse timestamps from SQL query results.
// Counterpart to formatTimestamp for bidirectional timestamp conversion.
func parseTimestamp(s string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", s)
}
