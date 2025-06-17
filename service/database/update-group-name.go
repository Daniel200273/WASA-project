package database

import "fmt"

// UpdateGroupName updates the name of a group
func (db *appdbimpl) UpdateGroupName(groupID, name string) error {
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

	// Update group name
	_, err = db.c.Exec(`
		UPDATE conversations 
		SET name = ? 
		WHERE id = ? AND type = 'group'
	`, name, groupID)
	if err != nil {
		return fmt.Errorf("failed to update group name: %w", err)
	}

	return nil
}
