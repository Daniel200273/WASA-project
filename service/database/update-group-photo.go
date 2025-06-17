package database

import "fmt"

// UpdateGroupPhoto updates the photo URL of a group
func (db *appdbimpl) UpdateGroupPhoto(groupID, photoURL string) error {
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

	// Update group photo
	_, err = db.c.Exec(`
		UPDATE conversations 
		SET photo_url = ? 
		WHERE id = ? AND type = 'group'
	`, photoURL, groupID)
	if err != nil {
		return fmt.Errorf("failed to update group photo: %w", err)
	}

	return nil
}
