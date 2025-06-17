package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
)

// === AUTHENTICATION METHODS ===

// CreateUser crea un nuovo utente con username univoco
func (db *appdbimpl) CreateUser(username string) (*User, error) {
	// Genera UUID per il nuovo utente
	userID := uuid.Must(uuid.NewV4()).String()

	query := `INSERT INTO users (id, username) VALUES (?, ?)`
	_, err := db.c.Exec(query, userID, username)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return db.GetUserByID(userID)
}

// GetUserByID recupera un utente per ID
func (db *appdbimpl) GetUserByID(id string) (*User, error) {
	var user User
	query := `SELECT id, username, photo_url, created_at FROM users WHERE id = ?`

	err := db.c.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.PhotoURL, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetUserByToken recupera un utente dal token di sessione
func (db *appdbimpl) GetUserByToken(token string) (*User, error) {
	var userID string
	query := `SELECT user_id FROM user_sessions WHERE token = ?`

	err := db.c.QueryRow(query, token).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("invalid token")
		}
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	return db.GetUserByID(userID)
}

// CreateUserSession crea una nuova sessione per l'utente
func (db *appdbimpl) CreateUserSession(userID string) (string, error) {
	// Genera token UUID
	token := uuid.Must(uuid.NewV4()).String()

	query := `INSERT INTO user_sessions (token, user_id) VALUES (?, ?)`
	_, err := db.c.Exec(query, token, userID)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	return token, nil
}

// DeleteUserSession elimina una sessione utente
func (db *appdbimpl) DeleteUserSession(token string) error {
	query := `DELETE FROM user_sessions WHERE token = ?`
	_, err := db.c.Exec(query, token)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

// === USER MANAGEMENT METHODS ===

// UpdateUsername aggiorna il username di un utente
func (db *appdbimpl) UpdateUsername(userID, newUsername string) error {
	query := `UPDATE users SET username = ? WHERE id = ?`
	result, err := db.c.Exec(query, newUsername, userID)
	if err != nil {
		return fmt.Errorf("failed to update username: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check update result: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// UpdateUserPhoto aggiorna la foto profilo di un utente
func (db *appdbimpl) UpdateUserPhoto(userID, photoURL string) error {
	query := `UPDATE users SET photo_url = ? WHERE id = ?`
	result, err := db.c.Exec(query, photoURL, userID)
	if err != nil {
		return fmt.Errorf("failed to update user photo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check update result: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// SearchUsers cerca utenti per username (escludendo l'utente corrente)
func (db *appdbimpl) SearchUsers(query string, excludeUserID string) ([]User, error) {
	sql := `SELECT id, username, photo_url, created_at 
			FROM users 
			WHERE username LIKE ? AND id != ? 
			ORDER BY username 
			LIMIT 100`

	rows, err := db.c.Query(sql, query+"%", excludeUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.PhotoURL, &user.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}
