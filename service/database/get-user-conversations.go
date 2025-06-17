package database

import (
	"database/sql"
	"fmt"
)

// GetUserConversations returns all conversations for a user
func (db *appdbimpl) GetUserConversations(userID string) ([]Conversation, error) {
	query := `
		SELECT DISTINCT 
			c.id, c.type, c.name, c.photo_url, c.created_by, c.created_at,
			m.id as last_message_id, m.content, m.photo_url as msg_photo, m.sender_id, m.created_at as msg_created_at,
			u.username as sender_username
		FROM conversations c
		JOIN conversation_participants cp ON c.id = cp.conversation_id
		LEFT JOIN messages m ON m.id = (
			SELECT m2.id FROM messages m2 
			WHERE m2.conversation_id = c.id 
			ORDER BY m2.created_at DESC 
			LIMIT 1
		)
		LEFT JOIN users u ON m.sender_id = u.id
		WHERE cp.user_id = ?
		ORDER BY COALESCE(m.created_at, c.created_at) DESC
	`

	rows, err := db.c.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query conversations: %w", err)
	}
	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		var conv Conversation
		var lastMsg *Message
		var msgID, content, photoURL, senderID, senderUsername sql.NullString
		var msgCreatedAt sql.NullTime

		err := rows.Scan(
			&conv.ID, &conv.Type, &conv.Name, &conv.PhotoURL, &conv.CreatedBy, &conv.CreatedAt,
			&msgID, &content, &photoURL, &senderID, &msgCreatedAt, &senderUsername,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan conversation: %w", err)
		}

		if msgID.Valid {
			lastMsg = &Message{
				ID:             msgID.String,
				ConversationID: conv.ID,
				SenderID:       senderID.String,
				Content:        &content.String,
				CreatedAt:      msgCreatedAt.Time,
			}
			if photoURL.Valid {
				lastMsg.PhotoURL = &photoURL.String
			}
		}

		conv.LastMessage = lastMsg

		// For direct conversations, get the other participant's info
		if conv.Type == "direct" {
			participantQuery := `
				SELECT u.id, u.username, u.photo_url
				FROM users u
				JOIN conversation_participants cp ON u.id = cp.user_id
				WHERE cp.conversation_id = ? AND u.id != ?
				LIMIT 1
			`
			var participant User
			err = db.c.QueryRow(participantQuery, conv.ID, userID).Scan(
				&participant.ID, &participant.Username, &participant.PhotoURL,
			)
			if err == nil {
				conv.OtherParticipant = &participant
			}
		}

		conversations = append(conversations, conv)
	}

	return conversations, rows.Err()
}
