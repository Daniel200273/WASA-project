package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

// === GROUP OPERATIONS ===

// CreateGroup creates a new group conversation
func (db *appdbimpl) CreateGroup(name, createdBy string, memberIDs []string) (*Conversation, error) {
	// Start transaction
	tx, err := db.c.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Generate group conversation ID
	groupID := uuid.Must(uuid.NewV4()).String()
	now := time.Now()

	// 2. Insert conversation with type='group'
	_, err = tx.Exec(`
		INSERT INTO conversations (id, type, name, created_by, created_at, last_message_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		groupID, "group", name, createdBy, now, now)
	if err != nil {
		return nil, fmt.Errorf("error creating group conversation: %w", err)
	}

	// 3. Add creator to the group as participant
	_, err = tx.Exec(`
		INSERT INTO conversation_participants (conversation_id, user_id, joined_at, last_read_at)
		VALUES (?, ?, ?, ?)`,
		groupID, createdBy, now, now)
	if err != nil {
		return nil, fmt.Errorf("error adding creator to group: %w", err)
	}

	// 4. Add all specified members as participants
	for _, memberID := range memberIDs {
		_, err = tx.Exec(`
			INSERT INTO conversation_participants (conversation_id, user_id, joined_at, last_read_at)
			VALUES (?, ?, ?, ?)`,
			groupID, memberID, now, now)
		if err != nil {
			return nil, fmt.Errorf("error adding member %s to group: %w", memberID, err)
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing group creation transaction: %w", err)
	}

	// 5. Return created conversation with member details
	group, err := db.getConversationWithParticipants(groupID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving created group: %w", err)
	}

	return group, nil
}

// AddUserToGroup adds a user to an existing group
func (db *appdbimpl) AddUserToGroup(groupID, userID string) error {
	// 1. Verify group exists and is of type 'group'
	var conversationType string
	err := db.c.QueryRow(`
		SELECT type FROM conversations WHERE id = ?`, groupID).Scan(&conversationType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("group not found")
		}
		return fmt.Errorf("error checking group existence: %w", err)
	}

	if conversationType != "group" {
		return fmt.Errorf("conversation is not a group")
	}

	// 2. Verify user is not already a member
	isAlreadyMember, err := db.IsUserInConversation(groupID, userID)
	if err != nil {
		return fmt.Errorf("error checking user membership: %w", err)
	}
	if isAlreadyMember {
		return fmt.Errorf("user is already a member of this group")
	}

	// 3. Add user to conversation_participants
	now := time.Now()
	_, err = db.c.Exec(`
		INSERT INTO conversation_participants (conversation_id, user_id, joined_at, last_read_at)
		VALUES (?, ?, ?, ?)`,
		groupID, userID, now, now)
	if err != nil {
		return fmt.Errorf("error adding user to group: %w", err)
	}

	return nil
}

// RemoveUserFromGroup removes a user from a group
func (db *appdbimpl) RemoveUserFromGroup(groupID, userID string) error {
	// 1. Verify group exists and is of type 'group'
	var conversationType string
	err := db.c.QueryRow(`
		SELECT type FROM conversations WHERE id = ?`, groupID).Scan(&conversationType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("group not found")
		}
		return fmt.Errorf("error checking group existence: %w", err)
	}

	if conversationType != "group" {
		return fmt.Errorf("conversation is not a group")
	}

	// 2. Verify user is currently a member
	isMember, err := db.IsUserInConversation(groupID, userID)
	if err != nil {
		return fmt.Errorf("error checking user membership: %w", err)
	}
	if !isMember {
		return fmt.Errorf("user is not a member of this group")
	}

	// 3. Remove user from conversation_participants
	result, err := db.c.Exec(`
		DELETE FROM conversation_participants 
		WHERE conversation_id = ? AND user_id = ?`,
		groupID, userID)
	if err != nil {
		return fmt.Errorf("error removing user from group: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user was not found in group")
	}

	return nil
}

// UpdateGroupName updates a group's name
func (db *appdbimpl) UpdateGroupName(groupID, name string) error {
	// 1. Verify group exists and is of type 'group'
	var conversationType string
	err := db.c.QueryRow(`
		SELECT type FROM conversations WHERE id = ?`, groupID).Scan(&conversationType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("group not found")
		}
		return fmt.Errorf("error checking group existence: %w", err)
	}

	if conversationType != "group" {
		return fmt.Errorf("conversation is not a group")
	}

	// 2. Update conversation name
	result, err := db.c.Exec(`
		UPDATE conversations SET name = ? WHERE id = ? AND type = 'group'`,
		name, groupID)
	if err != nil {
		return fmt.Errorf("error updating group name: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("group not found or not updated")
	}

	return nil
}

// UpdateGroupPhoto updates a group's photo URL
func (db *appdbimpl) UpdateGroupPhoto(groupID, photoURL string) error {
	// 1. Verify group exists and is of type 'group'
	var conversationType string
	err := db.c.QueryRow(`
		SELECT type FROM conversations WHERE id = ?`, groupID).Scan(&conversationType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("group not found")
		}
		return fmt.Errorf("error checking group existence: %w", err)
	}

	if conversationType != "group" {
		return fmt.Errorf("conversation is not a group")
	}

	// 2. Update conversation photo_url
	result, err := db.c.Exec(`
		UPDATE conversations SET photo_url = ? WHERE id = ? AND type = 'group'`,
		photoURL, groupID)
	if err != nil {
		return fmt.Errorf("error updating group photo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("group not found or not updated")
	}

	return nil
}

// getConversationWithParticipants retrieves a conversation with all its participants
func (db *appdbimpl) getConversationWithParticipants(conversationID string) (*Conversation, error) {
	// Get conversation details
	conversation := &Conversation{}
	err := db.c.QueryRow(`
		SELECT id, type, name, photo_url, created_by, created_at, last_message_at
		FROM conversations WHERE id = ?`, conversationID).Scan(
		&conversation.ID,
		&conversation.Type,
		&conversation.Name,
		&conversation.PhotoURL,
		&conversation.CreatedBy,
		&conversation.CreatedAt,
		&conversation.LastMessageAt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving conversation: %w", err)
	}

	// Get participants
	rows, err := db.c.Query(`
		SELECT u.id, u.username, u.photo_url
		FROM users u
		JOIN conversation_participants cp ON u.id = cp.user_id
		WHERE cp.conversation_id = ?
		ORDER BY u.username`, conversationID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving conversation participants: %w", err)
	}
	defer rows.Close()

	var participants []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.PhotoURL)
		if err != nil {
			return nil, fmt.Errorf("error scanning participant: %w", err)
		}
		participants = append(participants, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating participants: %w", err)
	}

	conversation.Participants = participants
	return conversation, nil
}
