package database

import "fmt"

// RemoveUserFromGroup removes a user from a group
func (db *appdbimpl) RemoveUserFromGroup(groupID, userID string) error {
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

	// Check if user is in the group
	isInGroup, err := db.IsUserInConversation(groupID, userID)
	if err != nil {
		return err
	}
	if !isInGroup {
		return fmt.Errorf("user %s is not in group %s", userID, groupID)
	}

	// Remove user from group
	_, err = db.c.Exec(`
		DELETE FROM conversation_participants 
		WHERE conversation_id = ? AND user_id = ?
	`, groupID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove user from group: %w", err)
	}

	return nil
}
