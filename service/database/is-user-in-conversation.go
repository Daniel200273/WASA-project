package database

import "fmt"

// IsUserInConversation checks if a user is a participant in a conversation
func (db *appdbimpl) IsUserInConversation(conversationID, userID string) (bool, error) {
	var count int
	err := db.c.QueryRow(`
		SELECT COUNT(*) FROM conversation_participants 
		WHERE conversation_id = ? AND user_id = ?
	`, conversationID, userID).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("failed to check conversation participation: %w", err)
	}

	return count > 0, nil
}
