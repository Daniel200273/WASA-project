package database

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

// CreateGroup creates a new group conversation with specified members
func (db *appdbimpl) CreateGroup(name, createdBy string, memberIDs []string) (*Conversation, error) {
	// Generate group ID
	groupID, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("failed to generate group ID: %w", err)
	}

	tx, err := db.c.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create group conversation
	now := time.Now()
	_, err = tx.Exec(`
		INSERT INTO conversations (id, type, name, created_by, created_at)
		VALUES (?, 'group', ?, ?, ?)
	`, groupID.String(), name, createdBy, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	// Add creator to group
	_, err = tx.Exec(`
		INSERT INTO conversation_participants (conversation_id, user_id, joined_at)
		VALUES (?, ?, ?)
	`, groupID.String(), createdBy, now)
	if err != nil {
		return nil, fmt.Errorf("failed to add creator to group: %w", err)
	}

	// Add other members to group
	for _, memberID := range memberIDs {
		if memberID != createdBy { // Don't add creator twice
			_, err = tx.Exec(`
				INSERT INTO conversation_participants (conversation_id, user_id, joined_at)
				VALUES (?, ?, ?)
			`, groupID.String(), memberID, now)
			if err != nil {
				return nil, fmt.Errorf("failed to add member %s to group: %w", memberID, err)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &Conversation{
		ID:        groupID.String(),
		Type:      "group",
		Name:      &name,
		CreatedBy: &createdBy,
		CreatedAt: now,
	}, nil
}
