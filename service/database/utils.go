package database

import (
	"database/sql"
	"time"
)

// === DATABASE UTILITIES ===

// Helper function to scan user from database row
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

// Helper function to scan users from database rows
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

// Helper function to scan conversation from database row
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

// Helper function to scan conversations from database rows
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

// Helper function to scan message from database row
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

// Helper function to scan messages from database rows
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

// Helper function to scan message reaction from database row
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

// Helper function to scan message reactions from database rows
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

// Helper function to check if error is "not found"
func isNotFoundError(err error) bool {
	return err == sql.ErrNoRows
}

// Helper function to format timestamp for SQLite
func formatTimestamp(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// Helper function to parse timestamp from SQLite
func parseTimestamp(s string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", s)
}
