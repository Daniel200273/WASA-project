package database

import (
	"fmt"

	"github.com/gofrs/uuid"
)

// === AUTHENTICATION OPERATIONS ===

// CreateUser creates a new user with the given username
func (db *appdbimpl) CreateUser(username string) (*User, error) {
	// TODO: Implement user creation
	// 1. Generate unique user ID
	// 2. Insert user into database
	// 3. Return created user

	userID := uuid.Must(uuid.NewV4()).String()

	// TODO: Add your implementation here
	_ = userID // Remove this line when implementing

	return nil, fmt.Errorf("CreateUser not implemented")
}

// GetUserByID retrieves a user by their ID
func (db *appdbimpl) GetUserByID(id string) (*User, error) {
	// TODO: Implement user retrieval by ID
	// 1. Query user from database by ID
	// 2. Handle user not found case
	// 3. Return user or error

	return nil, fmt.Errorf("GetUserByID not implemented")
}

// GetUserByUsername retrieves a user by their username
func (db *appdbimpl) GetUserByUsername(username string) (*User, error) {
	// TODO: Implement user retrieval by username
	// 1. Query user from database by username
	// 2. Handle user not found case
	// 3. Return user or error

	return nil, fmt.Errorf("GetUserByUsername not implemented")
}

// GetUserByToken retrieves a user by their session token
func (db *appdbimpl) GetUserByToken(token string) (*User, error) {
	// TODO: Implement user retrieval by token
	// 1. Join user_sessions with users table
	// 2. Query by token
	// 3. Handle token not found or expired case
	// 4. Return user or error

	return nil, fmt.Errorf("GetUserByToken not implemented")
}

// CreateUserSession creates a new session for a user
func (db *appdbimpl) CreateUserSession(userID string) (string, error) {
	// TODO: Implement session creation
	// 1. Generate unique session token
	// 2. Insert session into database
	// 3. Return token or error

	token := uuid.Must(uuid.NewV4()).String()

	// TODO: Add your implementation here
	_ = token // Remove this line when implementing

	return "", fmt.Errorf("CreateUserSession not implemented")
}

// DeleteUserSession deletes a user session
func (db *appdbimpl) DeleteUserSession(token string) error {
	// TODO: Implement session deletion
	// 1. Delete session from database by token
	// 2. Handle session not found case
	// 3. Return error if any

	return fmt.Errorf("DeleteUserSession not implemented")
}
