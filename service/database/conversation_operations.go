package database

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

// === CONVERSATION OPERATIONS ===

// GetUserConversations retrieves all conversations for a user
func (db *appdbimpl) GetUserConversations(userID string) ([]ConversationPreview, error) {
	// 1. Get all conversations where the user is a participant
	query := `
		SELECT c.id, c.type, c.name, c.photo_url, c.last_message_at
		FROM conversations c
		JOIN conversation_participants cp ON c.id = cp.conversation_id
		WHERE cp.user_id = ?
		ORDER BY c.last_message_at DESC
	`
	rows, err := db.c.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user conversations: %w", err)
	}
	defer rows.Close()

	// Prepare result array
	var conversations []ConversationPreview

	// Process each conversation
	for rows.Next() {
		var conv ConversationPreview
		err := rows.Scan(
			&conv.ID,
			&conv.Type,
			&conv.Name,
			&conv.PhotoURL,
			&conv.LastMessageAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning conversation: %w", err)
		}

		// For direct conversations, get the other participant's info
		if conv.Type == "direct" {
			otherUserQuery := `
				SELECT u.id, u.username, u.photo_url
				FROM conversation_participants cp
				JOIN users u ON cp.user_id = u.id
				WHERE cp.conversation_id = ? AND cp.user_id != ?
				LIMIT 1
			`
			var otherID string
			var otherUsername string
			var otherPhotoURL *string

			err = db.c.QueryRow(otherUserQuery, conv.ID, userID).Scan(
				&otherID,
				&otherUsername,
				&otherPhotoURL,
			)
			if err != nil && !isNotFoundError(err) {
				return nil, fmt.Errorf("error retrieving other participant: %w", err)
			}

			if err == nil {
				// Set other participant info
				conv.OtherParticipant = &struct {
					ID       string  `json:"id"`
					Username string  `json:"username"`
					PhotoURL *string `json:"photoUrl,omitempty"`
				}{
					ID:       otherID,
					Username: otherUsername,
					PhotoURL: otherPhotoURL,
				}

				// For direct chats, use the other participant's name
				if conv.Name == nil || *conv.Name == "" {
					conv.Name = &otherUsername
				}

				// For direct chats without a photo, use the other participant's photo
				if conv.PhotoURL == nil {
					conv.PhotoURL = otherPhotoURL
				}
			}
		}

		// Get last message for the conversation
		lastMessageQuery := `
			SELECT m.id, m.content, m.photo_url, m.created_at, u.username
			FROM messages m
			JOIN users u ON m.sender_id = u.id
			WHERE m.conversation_id = ?
			ORDER BY m.created_at DESC
			LIMIT 1
		`
		var msgID string
		var content *string
		var photoURL *string
		var timestamp time.Time
		var senderUsername string

		err = db.c.QueryRow(lastMessageQuery, conv.ID).Scan(
			&msgID,
			&content,
			&photoURL,
			&timestamp,
			&senderUsername,
		)

		// If we found a message, add it to the conversation
		if err == nil {
			hasPhoto := photoURL != nil
			conv.LastMessage = &MessagePreview{
				ID:             msgID,
				Content:        content,
				Timestamp:      timestamp,
				SenderUsername: senderUsername,
				HasPhoto:       hasPhoto,
			}
		} else if !isNotFoundError(err) {
			return nil, fmt.Errorf("error retrieving last message: %w", err)
		}

		// Get unread count for the conversation
		// Count messages newer than the user's last_read_at timestamp
		unreadQuery := `
			SELECT COUNT(m.id)
			FROM messages m
			JOIN conversation_participants cp ON m.conversation_id = cp.conversation_id
			WHERE m.conversation_id = ? 
			AND cp.user_id = ?
			AND m.sender_id != ?
			AND m.created_at > cp.last_read_at
		`
		var unreadCount int
		err = db.c.QueryRow(unreadQuery, conv.ID, userID, userID).Scan(&unreadCount)
		if err != nil {
			// If there's an error, default to 0
			unreadCount = 0
		}
		conv.UnreadCount = unreadCount

		conversations = append(conversations, conv)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over conversations: %w", err)
	}

	return conversations, nil
}

// GetConversation retrieves a specific conversation for a user
func (db *appdbimpl) GetConversation(conversationID, userID string) (*Conversation, error) {

	// 1. Get conversation details
	query := `
		SELECT c.id, c.type, c.name, c.photo_url, c.created_by, c.created_at, c.last_message_at
		FROM conversations c
		WHERE c.id = ?
	`

	row := db.c.QueryRow(query, conversationID)
	conv, err := scanConversation(row)
	if err != nil {
		return nil, fmt.Errorf("error retrieving conversation: %w", err)
	}

	// 2. Get all participants
	participantsQuery := `
		SELECT u.id, u.username, u.photo_url, u.created_at
		FROM users u
		JOIN conversation_participants cp ON u.id = cp.user_id
		WHERE cp.conversation_id = ?
	`
	rows, err := db.c.Query(participantsQuery, conversationID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving participants: %w", err)
	}

	participants, err := scanUsers(rows)
	if err != nil {
		return nil, fmt.Errorf("error scanning participants: %w", err)
	}

	// 4. Populate conversation details
	conv.Participants = participants

	// For direct conversations, find the other participant
	if conv.Type == "direct" {
		for _, p := range participants {
			if p.ID != userID {
				conv.OtherParticipant = &p
				break
			}
		}
	}

	// Calculate unread count for this conversation
	unreadQuery := `
		SELECT COUNT(m.id)
		FROM messages m
		JOIN conversation_participants cp ON m.conversation_id = cp.conversation_id
		WHERE m.conversation_id = ? 
		AND cp.user_id = ?
		AND m.sender_id != ?
		AND m.created_at > cp.last_read_at
	`
	var unreadCount int
	err = db.c.QueryRow(unreadQuery, conv.ID, userID, userID).Scan(&unreadCount)
	if err != nil {
		// If there's an error, default to 0
		unreadCount = 0
	}
	conv.UnreadCount = unreadCount

	return conv, nil
}

// GetOrCreateDirectConversation gets or creates a direct conversation between two users
func (db *appdbimpl) GetOrCreateDirectConversation(user1ID, user2ID string) (*Conversation, error) {
	// 1. Check if direct conversation already exists between the two users
	query := `
		SELECT c.id, c.type, c.name, c.photo_url, c.created_by, c.created_at, c.last_message_at
		FROM conversations c
		JOIN conversation_participants cp1 ON c.id = cp1.conversation_id
		JOIN conversation_participants cp2 ON c.id = cp2.conversation_id
		WHERE c.type = 'direct' AND cp1.user_id = ? AND cp2.user_id = ?
		LIMIT 1
	`

	row := db.c.QueryRow(query, user1ID, user2ID)
	existingConversation, err := scanConversation(row)

	// If conversation exists, return it
	if err == nil {
		return existingConversation, nil
	}

	// If error is not "no rows", return the error
	if !isNotFoundError(err) {
		return nil, fmt.Errorf("error checking existing conversation: %w", err)
	}

	// 3. Create new direct conversation
	tx, err := db.c.Begin()
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", err)
	}

	// Generate a new conversation ID
	conversationID := uuid.Must(uuid.NewV4()).String()

	// Insert the new conversation
	createQuery := `
		INSERT INTO conversations (id, type, created_by, created_at) 
		VALUES (?, 'direct', ?, CURRENT_TIMESTAMP)
	`
	_, err = tx.Exec(createQuery, conversationID, user1ID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("error creating conversation: %w (rollback failed: %w)", err, rollbackErr)
		}
		return nil, fmt.Errorf("error creating conversation: %w", err)
	}

	// 4. Add both users as participants
	addParticipantQuery := `
		INSERT INTO conversation_participants (conversation_id, user_id, joined_at, last_read_at)
		VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	_, err = tx.Exec(addParticipantQuery, conversationID, user1ID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("error adding user1 to conversation: %w (rollback failed: %w)", err, rollbackErr)
		}
		return nil, fmt.Errorf("error adding user1 to conversation: %w", err)
	}

	_, err = tx.Exec(addParticipantQuery, conversationID, user2ID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("error adding user2 to conversation: %w (rollback failed: %w)", err, rollbackErr)
		}
		return nil, fmt.Errorf("error adding user2 to conversation: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	// 5. Return the newly created conversation
	return db.GetConversation(conversationID, user1ID)
}

// IsUserInConversation checks if a user is a participant in a conversation
func (db *appdbimpl) IsUserInConversation(conversationID, userID string) (bool, error) {
	query := `SELECT COUNT(*) FROM conversation_participants 
			  WHERE conversation_id = ? AND user_id = ?`

	var count int
	err := db.c.QueryRow(query, conversationID, userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking conversation participation: %w", err)
	}

	return count > 0, nil
}
