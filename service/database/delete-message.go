package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// DeleteMessage deletes a message if the user is the sender
func (db *appdbimpl) DeleteMessage(messageID, userID string) error {
	// First check if message exists and user is the sender
	var senderID string
	err := db.c.QueryRow(`
		SELECT sender_id FROM messages WHERE id = ?
	`, messageID).Scan(&senderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("message not found")
		}
		return fmt.Errorf("failed to check message ownership: %w", err)
	}

	// Check if user is the sender
	if senderID != userID {
		return errors.New("user is not authorized to delete this message")
	}

	// Delete the message
	result, err := db.c.Exec(`
		DELETE FROM messages WHERE id = ? AND sender_id = ?
	`, messageID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check deletion result: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("message not found or user not authorized")
	}

	return nil
}
