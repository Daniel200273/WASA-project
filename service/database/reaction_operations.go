package database

import (
	"fmt"

	"github.com/gofrs/uuid"
)

// === REACTION OPERATIONS ===

// CreateMessageReaction creates a reaction to a message
func (db *appdbimpl) CreateMessageReaction(messageID, userID, emoticon string) (*MessageReaction, error) {
	// 1. Verify message exists and user has access to it
	message, err := db.GetMessage(messageID)
	if err != nil {
		return nil, fmt.Errorf("message not found: %w", err)
	}

	// Check if user is in the conversation containing this message
	isParticipant, err := db.IsUserInConversation(message.ConversationID, userID)
	if err != nil {
		return nil, fmt.Errorf("error checking conversation participation: %w", err)
	}
	if !isParticipant {
		return nil, fmt.Errorf("user not authorized to react to this message")
	}

	// 2. Check if user already reacted to this message (UPSERT behavior)
	checkQuery := `SELECT id FROM message_reactions WHERE message_id = ? AND user_id = ?`
	var existingID string
	err = db.c.QueryRow(checkQuery, messageID, userID).Scan(&existingID)

	if err == nil {
		// Update existing reaction
		updateQuery := `UPDATE message_reactions SET emoticon = ? WHERE id = ?`
		_, err = db.c.Exec(updateQuery, emoticon, existingID)
		if err != nil {
			return nil, fmt.Errorf("error updating reaction: %w", err)
		}

		// Return updated reaction
		return db.getMessageReactionByID(existingID)
	} else if !isNotFoundError(err) {
		return nil, fmt.Errorf("error checking existing reaction: %w", err)
	}

	// 3. Create new reaction
	reactionID := uuid.Must(uuid.NewV4()).String()
	insertQuery := `
		INSERT INTO message_reactions (id, message_id, user_id, emoticon, created_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
	`
	_, err = db.c.Exec(insertQuery, reactionID, messageID, userID, emoticon)
	if err != nil {
		return nil, fmt.Errorf("error creating reaction: %w", err)
	}

	// 4. Return created reaction with username
	return db.getMessageReactionByID(reactionID)
}

// getMessageReactionByID retrieves a single reaction by ID with username
func (db *appdbimpl) getMessageReactionByID(reactionID string) (*MessageReaction, error) {
	query := `
		SELECT mr.id, mr.message_id, mr.user_id, u.username, mr.emoticon, mr.created_at
		FROM message_reactions mr
		JOIN users u ON mr.user_id = u.id
		WHERE mr.id = ?
	`

	row := db.c.QueryRow(query, reactionID)
	var reaction MessageReaction
	err := row.Scan(
		&reaction.ID,
		&reaction.MessageID,
		&reaction.UserID,
		&reaction.Username,
		&reaction.Emoticon,
		&reaction.CreatedAt,
	)

	if err != nil {
		if isNotFoundError(err) {
			return nil, fmt.Errorf("reaction not found")
		}
		return nil, fmt.Errorf("error retrieving reaction: %w", err)
	}

	return &reaction, nil
}

// DeleteMessageReaction deletes a reaction (only by the user who created it)
func (db *appdbimpl) DeleteMessageReaction(reactionID, userID string) error {
	// 1. Verify that the user owns the reaction
	checkQuery := `SELECT user_id FROM message_reactions WHERE id = ?`
	var ownerID string
	err := db.c.QueryRow(checkQuery, reactionID).Scan(&ownerID)
	if err != nil {
		if isNotFoundError(err) {
			return fmt.Errorf("reaction not found")
		}
		return fmt.Errorf("error checking reaction ownership: %w", err)
	}

	if ownerID != userID {
		return fmt.Errorf("unauthorized: user can only delete their own reactions")
	}

	// 2. Delete the reaction from database
	deleteQuery := `DELETE FROM message_reactions WHERE id = ?`
	result, err := db.c.Exec(deleteQuery, reactionID)
	if err != nil {
		return fmt.Errorf("error deleting reaction: %w", err)
	}

	// 3. Verify deletion was successful
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking deletion result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("reaction not found or already deleted")
	}

	return nil
}
