package database

import (
	"database/sql"
	"errors"
)

// === DATABASE UTILITIES ===
//
// These helper functions provide common database operations for converting
// database rows into Go structs and handling common database tasks.

// === USER SCANNING FUNCTIONS ===

// scanUser converts a single database row into a User struct.
// Used after QueryRow() calls to map database columns to User fields.
// Expected column order: id, username, photo_url, created_at
func scanUser(row *sql.Row) (*User, error) {
	var user User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.PhotoURL,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// scanUsers converts multiple database rows into a slice of User structs.
// Used after Query() calls when retrieving multiple users.
// Expected column order: id, username, photo_url, created_at
// Automatically handles closing rows and iterating through all records.
func scanUsers(rows *sql.Rows) ([]User, error) {
	var users []User
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.PhotoURL,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// === CONVERSATION SCANNING FUNCTIONS ===

// scanConversation converts a single database row into a Conversation struct.
// Used after QueryRow() calls to map database columns to Conversation fields.
// Expected column order: id, type, name, photo_url, created_by, created_at, last_message_at
func scanConversation(row *sql.Row) (*Conversation, error) {
	var conv Conversation
	err := row.Scan(
		&conv.ID,
		&conv.Type,
		&conv.Name,
		&conv.PhotoURL,
		&conv.CreatedBy,
		&conv.CreatedAt,
		&conv.LastMessageAt,
	)
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

// === UTILITY FUNCTIONS ===

// isNotFoundError checks if an error is sql.ErrNoRows (record not found).
// Used to distinguish between "not found" vs actual database errors.
// Returns true if the error indicates no rows were found in the query result.
func isNotFoundError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
