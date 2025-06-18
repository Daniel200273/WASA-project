package database

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

// === AUTHENTICATION OPERATIONS ===

// CreateUser creates a new user with the given username
func (db *appdbimpl) CreateUser(username string) (*User, error) {
	// 1. Generate unique user ID
	// 1.1 Validate username
	if !isValidUsername(username) {
		return nil, fmt.Errorf("username contains invalid characters")
	}

	// 1.2 Check if username already exists
	existingUser, err := db.GetUserByUsername(username)
	if err != nil && !isNotFoundError(err) {
		// Database error occurred, don't proceed
		return nil, fmt.Errorf("error checking username availability: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("username already exists")
	}

	// 1.3 Generate user ID
	userID := uuid.Must(uuid.NewV4()).String()

	// 2. Insert user into database
	query := `
		INSERT INTO users (id, username, photo_url, created_at)
		VALUES (?, ?, NULL, ?)
	`
	createdAt := time.Now().UTC()
	_, err = db.c.Exec(query, userID, username, createdAt)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	// 3. Return created user
	return db.GetUserByID(userID)
}

// GetUserByID retrieves a user by their ID
func (db *appdbimpl) GetUserByID(id string) (*User, error) {
	// 1. Query user from database by ID
	query := `
		SELECT id, username, photo_url, created_at
		FROM users
		WHERE id = ?
	`
	row := db.c.QueryRow(query, id)
	user, err := scanUser(row)

	// 2. Handle user not found case
	if err != nil {
		if isNotFoundError(err) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}
	// 3. Return user or error
	return user, nil
}

// GetUserByUsername retrieves a user by their username
func (db *appdbimpl) GetUserByUsername(username string) (*User, error) {
	// 1. Query user from database by username
	query := `
		SELECT id, username, photo_url, created_at
		FROM users
		WHERE username = ?
	`
	row := db.c.QueryRow(query, username)
	user, err := scanUser(row)
	// 2. Handle user not found case
	if err != nil {
		if isNotFoundError(err) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}
	// 3. Return user or error

	return user, nil
}

// GetUserByToken retrieves a user by their session token.
// This query returns user details (id, username, photo_url, created_at)
// for all users who have active sessions in the user_sessions table.
func (db *appdbimpl) GetUserByToken(token string) (*User, error) {
	// 1. Join user_sessions with users table
	query := `		
		SELECT u.id, u.username, u.photo_url, u.created_at
		FROM user_sessions us
		JOIN users u ON us.user_id = u.id
		WHERE us.token = ? 
	`
	// 2. Query by token
	row := db.c.QueryRow(query, token)
	user, err := scanUser(row)

	// 3. Handle token not found or expired case
	if err != nil {
		if isNotFoundError(err) {
			return nil, fmt.Errorf("session not found or expired")
		}
		return nil, fmt.Errorf("error retrieving user by token: %w", err)
	}
	// 4. Return user or error
	return user, nil
}

// CreateUserSession creates a new session for a user
func (db *appdbimpl) CreateUserSession(userID string) (string, error) {
	// 1. Generate unique session token
	token := uuid.Must(uuid.NewV4()).String()

	// 2. Insert session into database
	query := `		
		INSERT INTO user_sessions (token, user_id, created_at)
		VALUES (?, ?, ?)
	`
	// 3. Return token or error
	createdAt := time.Now().UTC()
	_, err := db.c.Exec(query, token, userID, createdAt)
	if err != nil {
		return "", fmt.Errorf("error creating user session: %w", err)
	}

	// 4. Return the generated token
	return token, nil

}

// DeleteUserSession deletes a user session
func (db *appdbimpl) DeleteUserSession(token string) error {
	// 1. Delete session from database by token
	query := `		
		DELETE FROM user_sessions
		WHERE token = ?
	`
	result, err := db.c.Exec(query, token)
	if err != nil {
		return fmt.Errorf("error deleting user session: %w", err)
	}

	// 2. Handle session not found case
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking deletion outcome: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("session not found or already deleted")
	}

	// 3. Return nil if deletion was successful
	return nil
}
