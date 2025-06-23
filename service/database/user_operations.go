package database

import (
	"fmt"
)

// === USER MANAGEMENT OPERATIONS ===

// UpdateUsername updates a user's username
func (db *appdbimpl) UpdateUsername(userID, newUsername string) error {
	// TODO: Implement username update
	query := `	
		UPDATE users
		SET username = ?
 		WHERE id = ?
	`
	// 1. Update user's username in database
	_, err := db.c.Exec(query, newUsername, userID)

	// 2. Handle user not found case
	if err != nil && isNotFoundError(err) {
		return fmt.Errorf("user not found")
	}
	if err != nil {
		return fmt.Errorf("error updating username: %w", err)
	}

	// 3. return nil if successful
	return nil
}

// UpdateUserPhoto updates a user's profile photo URL
func (db *appdbimpl) UpdateUserPhoto(userID, photoURL string) error {
	// TODO: Implement user photo update
	// 1. Update user's photo_url in database
	query := `
	 	UPDATE users
	 	SET photo_url = ?
	   	WHERE id = ?
	 `
	_, err := db.c.Exec(query, photoURL, userID)

	// 2. Handle user not found case
	if err != nil && isNotFoundError(err) {
		return fmt.Errorf("user not found")
	}

	// 3. Handle other errors
	if err != nil {
		return fmt.Errorf("error updating user photo: %w", err)
	}

	// 3. Return error nil if successful
	return nil
}

// SearchUsers searches for users by query string, excluding the specified user
func (db *appdbimpl) SearchUsers(query string, excludeUserID string) ([]User, error) {
	// TODO: Implement user search
	// 1. Search users by username (case-insensitive LIKE query)
	sqlQuery := `
		SELECT id, username, photo_url, created_at
		FROM users
		WHERE username LIKE ? COLLATE NOCASE
		AND id != ?
		ORDER BY username
		LIMIT 20
	`
	// 2. Prepare search pattern (% for substring matching)
	searchPattern := "%" + query + "%"

	// 3. Execute query
	rows, err := db.c.Query(sqlQuery, searchPattern, excludeUserID)
	if err != nil {
		return nil, fmt.Errorf("error searching users: %w", err)
	}
	defer rows.Close()
	// 4. Scan results using existing helper function
	users, err := scanUsers(rows)
	if err != nil {
		return nil, fmt.Errorf("error scanning user results: %w", err)
	}
	// 4. Return list of matching users

	return users, nil
}
