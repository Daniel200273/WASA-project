package database

import (
	"fmt"
	"time"
)

// AddUserToGroup adds a user to an existing group
func (db *appdbimpl) AddUserToGroup(groupID, userID string) error {
	// First check if it's actually a group
	var convType string
	err := db.c.QueryRow(`
		SELECT type FROM conversations WHERE id = ?
	`, groupID).Scan(&convType)
	if err != nil {
		return fmt.Errorf("failed to verify group exists: %w", err)
	}
	if convType != "group" {
		return fmt.Errorf("conversation %s is not a group", groupID)
	}

	// Check if user is already in the group
	isInGroup, err := db.IsUserInConversation(groupID, userID)
	if err != nil {
		return err
	}
	if isInGroup {
		return fmt.Errorf("user %s is already in group %s", userID, groupID)
	}

	// Add user to group
	_, err = db.c.Exec(`
		INSERT INTO conversation_participants (conversation_id, user_id, joined_at)
		VALUES (?, ?, ?)
	`, groupID, userID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to add user to group: %w", err)
	}

	return nil
}
