package database

import (
	"fmt"
)

// === USER MANAGEMENT OPERATIONS ===

// UpdateUsername updates a user's username
func (db *appdbimpl) UpdateUsername(userID, newUsername string) error {
	// 1. Check if the new username is already taken by another user
	checkQuery := `SELECT id FROM users WHERE username = ? AND id != ?`
	var existingUserID string
	err := db.c.QueryRow(checkQuery, newUsername, userID).Scan(&existingUserID)
	if err == nil {
		return fmt.Errorf("username already taken")
	} else if !isNotFoundError(err) {
		return fmt.Errorf("error checking username availability: %w", err)
	}

	// 2. Update user's username in database
	query := `	
		UPDATE users
		SET username = ?
 		WHERE id = ?
	`
	result, err := db.c.Exec(query, newUsername, userID)
	if err != nil {
		return fmt.Errorf("error updating username: %w", err)
	}

	// 3. Check if user was found and updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking update result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// UpdateUserPhoto updates a user's profile photo URL
func (db *appdbimpl) UpdateUserPhoto(userID, photoURL string) error {
	// 1. Update user's photo_url in database
	query := `
	 	UPDATE users
	 	SET photo_url = ?
	   	WHERE id = ?
	 `
	result, err := db.c.Exec(query, photoURL, userID)
	if err != nil {
		return fmt.Errorf("error updating user photo: %w", err)
	}

	// 2. Check if user was found and updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking update result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// SearchUsers searches for users by query string, excluding the specified user
func (db *appdbimpl) SearchUsers(query string, excludeUserID string) ([]User, error) {
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

	return users, nil
}

// GetUser retrieves a user by their ID
func (db *appdbimpl) GetUser(userID string) (*User, error) {
	query := `
		SELECT id, username, photo_url, created_at
		FROM users 
		WHERE id = ?
	`

	row := db.c.QueryRow(query, userID)

	var user User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.PhotoURL,
		&user.CreatedAt,
	)

	if err != nil {
		if isNotFoundError(err) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	return &user, nil
}
