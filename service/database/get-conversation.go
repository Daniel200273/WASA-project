package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// GetConversation returns a specific conversation if the user is a participant
func (db *appdbimpl) GetConversation(conversationID, userID string) (*Conversation, error) {
	// First check if user is participant
	isParticipant, err := db.IsUserInConversation(conversationID, userID)
	if err != nil {
		return nil, err
	}
	if !isParticipant {
		return nil, errors.New("user is not a participant in this conversation")
	}

	query := `
		SELECT id, type, name, photo_url, created_by, created_at
		FROM conversations 
		WHERE id = ?
	`

	var conv Conversation
	err = db.c.QueryRow(query, conversationID).Scan(
		&conv.ID, &conv.Type, &conv.Name, &conv.PhotoURL, &conv.CreatedBy, &conv.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("conversation not found")
		}
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	// Get participants
	participantsQuery := `
		SELECT u.id, u.username, u.photo_url
		FROM users u
		JOIN conversation_participants cp ON u.id = cp.user_id
		WHERE cp.conversation_id = ?
	`
	rows, err := db.c.Query(participantsQuery, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var participant User
		err := rows.Scan(&participant.ID, &participant.Username, &participant.PhotoURL)
		if err != nil {
			return nil, fmt.Errorf("failed to scan participant: %w", err)
		}
		conv.Participants = append(conv.Participants, participant)
	}

	return &conv, rows.Err()
}
