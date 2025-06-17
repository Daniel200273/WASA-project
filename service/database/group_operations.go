package database

import (
	"fmt"

	"github.com/gofrs/uuid"
)

// === GROUP OPERATIONS ===

// CreateGroup creates a new group conversation
func (db *appdbimpl) CreateGroup(name, createdBy string, memberIDs []string) (*Conversation, error) {
	// TODO: Implement group creation
	// 1. Generate group conversation ID
	// 2. Insert conversation with type='group'
	// 3. Add creator to the group as participant
	// 4. Add all specified members as participants
	// 5. Return created conversation with member details

	groupID := uuid.Must(uuid.NewV4()).String()

	// TODO: Add your implementation here
	_ = groupID // Remove this line when implementing

	return nil, fmt.Errorf("CreateGroup not implemented")
}

// AddUserToGroup adds a user to an existing group
func (db *appdbimpl) AddUserToGroup(groupID, userID string) error {
	// TODO: Implement user addition to group
	// 1. Verify group exists and is of type 'group'
	// 2. Verify user is not already a member
	// 3. Add user to conversation_participants

	return fmt.Errorf("AddUserToGroup not implemented")
}

// RemoveUserFromGroup removes a user from a group
func (db *appdbimpl) RemoveUserFromGroup(groupID, userID string) error {
	// TODO: Implement user removal from group
	// 1. Verify group exists and is of type 'group'
	// 2. Verify user is currently a member
	// 3. Remove user from conversation_participants
	// 4. Consider what happens if removing the group creator

	return fmt.Errorf("RemoveUserFromGroup not implemented")
}

// UpdateGroupName updates a group's name
func (db *appdbimpl) UpdateGroupName(groupID, name string) error {
	// TODO: Implement group name update
	// 1. Verify group exists and is of type 'group'
	// 2. Update conversation name

	return fmt.Errorf("UpdateGroupName not implemented")
}

// UpdateGroupPhoto updates a group's photo URL
func (db *appdbimpl) UpdateGroupPhoto(groupID, photoURL string) error {
	// TODO: Implement group photo update
	// 1. Verify group exists and is of type 'group'
	// 2. Update conversation photo_url

	return fmt.Errorf("UpdateGroupPhoto not implemented")
}
