package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
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
