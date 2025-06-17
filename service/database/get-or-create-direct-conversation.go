package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

// GetOrCreateDirectConversation gets existing direct conversation or creates new one
func (db *appdbimpl) GetOrCreateDirectConversation(user1ID, user2ID string) (*Conversation, error) {
	// Check if direct conversation already exists
	query := `
		SELECT c.id, c.type, c.name, c.photo_url, c.created_by, c.created_at
		FROM conversations c
		JOIN conversation_participants cp1 ON c.id = cp1.conversation_id
		JOIN conversation_participants cp2 ON c.id = cp2.conversation_id
		WHERE c.type = 'direct' 
		AND cp1.user_id = ? AND cp2.user_id = ?
		LIMIT 1
	`

	var conv Conversation
	err := db.c.QueryRow(query, user1ID, user2ID).Scan(
		&conv.ID, &conv.Type, &conv.Name, &conv.PhotoURL, &conv.CreatedBy, &conv.CreatedAt,
	)

	if err == nil {
		// Conversation exists, return it
		return &conv, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to check existing conversation: %w", err)
	}

	// Create new direct conversation
	conversationID, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("failed to generate conversation ID: %w", err)
	}

	tx, err := db.c.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create conversation
	_, err = tx.Exec(`
		INSERT INTO conversations (id, type, created_by, created_at)
		VALUES (?, 'direct', ?, ?)
	`, conversationID.String(), user1ID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	// Add participants
	_, err = tx.Exec(`
		INSERT INTO conversation_participants (conversation_id, user_id, joined_at)
		VALUES (?, ?, ?), (?, ?, ?)
	`, conversationID.String(), user1ID, time.Now(), conversationID.String(), user2ID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to add participants: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return the new conversation
	return &Conversation{
		ID:        conversationID.String(),
		Type:      "direct",
		CreatedBy: &user1ID,
		CreatedAt: time.Now(),
	}, nil
}
