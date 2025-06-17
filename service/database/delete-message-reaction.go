package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// DeleteMessageReaction deletes a message reaction if the user is the owner
func (db *appdbimpl) DeleteMessageReaction(reactionID, userID string) error {
	// First check if reaction exists and user is the owner
	var ownerID string
	err := db.c.QueryRow(`
		SELECT user_id FROM message_reactions WHERE id = ?
	`, reactionID).Scan(&ownerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("reaction not found")
		}
		return fmt.Errorf("failed to check reaction ownership: %w", err)
	}

	// Check if user is the owner of the reaction
	if ownerID != userID {
		return errors.New("user is not authorized to delete this reaction")
	}

	// Delete the reaction
	result, err := db.c.Exec(`
		DELETE FROM message_reactions WHERE id = ? AND user_id = ?
	`, reactionID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete reaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check deletion result: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("reaction not found or user not authorized")
	}

	return nil
}
