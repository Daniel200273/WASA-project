/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	// Health check
	Ping() error

	// === AUTHENTICATION ===
	CreateUser(username string) (*User, error)
	GetUserByID(id string) (*User, error)
	GetUserByUsername(username string) (*User, error)
	GetUserByToken(token string) (*User, error)
	CreateUserSession(userID string) (string, error)
	DeleteUserSession(token string) error

	// === USER MANAGEMENT ===
	UpdateUsername(userID, newUsername string) error
	UpdateUserPhoto(userID, photoURL string) error
	SearchUsers(query string, excludeUserID string) ([]User, error)

	// === CONVERSATIONS ===
	GetUserConversations(userID string) ([]ConversationPreview, error)
	GetConversation(conversationID, userID string) (*Conversation, error)
	GetOrCreateDirectConversation(user1ID, user2ID string) (*Conversation, error)

	// === MESSAGES ===
	CreateMessage(conversationID, senderID string, content *string, photoURL *string, replyToID *string) (*Message, error)
	GetMessage(messageID string) (*Message, error)
	GetConversationMessages(conversationID string) ([]Message, error)
	DeleteMessage(messageID, userID string) error
	ForwardMessage(messageID, targetConversationID, userID string) (*Message, error)
	MarkConversationAsRead(conversationID, userID string) error

	// === REACTIONS ===
	CreateMessageReaction(messageID, userID, emoticon string) (*MessageReaction, error)
	DeleteMessageReaction(reactionID, userID string) error

	// === GROUPS ===
	CreateGroup(name, createdBy string, memberIDs []string) (*Conversation, error)
	AddUserToGroup(groupID, userID string) error
	RemoveUserFromGroup(groupID, userID string) error
	UpdateGroupName(groupID, name string) error
	UpdateGroupPhoto(groupID, photoURL string) error
	IsUserInConversation(conversationID, userID string) (bool, error)
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	appDB := &appdbimpl{c: db}

	// Inizializza lo schema WASAText se necessario
	if err := appDB.initializeSchema(); err != nil {
		return nil, fmt.Errorf("error initializing database schema: %w", err)
	}

	return appDB, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

// initializeSchema crea le tabelle WASAText se non esistono
func (db *appdbimpl) initializeSchema() error {
	schema := `
	-- Users table
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		photo_url TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	
	-- User sessions table
	CREATE TABLE IF NOT EXISTS user_sessions (
		token TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	
	-- Conversations table
	CREATE TABLE IF NOT EXISTS conversations (
		id TEXT PRIMARY KEY,
		type TEXT NOT NULL CHECK (type IN ('direct', 'group')),
		name TEXT,
		photo_url TEXT,
		created_by TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_message_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (created_by) REFERENCES users(id)
	);
	
	-- Conversation participants table
	CREATE TABLE IF NOT EXISTS conversation_participants (
		conversation_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_read_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (conversation_id, user_id),
		FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	
	-- Messages table
	CREATE TABLE IF NOT EXISTS messages (
		id TEXT PRIMARY KEY,
		conversation_id TEXT NOT NULL,
		sender_id TEXT NOT NULL,
		content TEXT,
		photo_url TEXT,
		reply_to_id TEXT,
		forwarded BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
		FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (reply_to_id) REFERENCES messages(id) ON DELETE SET NULL,
		CHECK ((content IS NOT NULL AND photo_url IS NULL) OR 
			   (content IS NULL AND photo_url IS NOT NULL))
	);
	
	-- Message reactions table
	CREATE TABLE IF NOT EXISTS message_reactions (
		id TEXT PRIMARY KEY,
		message_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		emoticon TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE (message_id, user_id)
	);
	
	-- Indices for performance
	CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
	CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions(user_id);
	CREATE INDEX IF NOT EXISTS idx_conversations_type ON conversations(type);
	CREATE INDEX IF NOT EXISTS idx_participants_user_id ON conversation_participants(user_id);
	CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON messages(conversation_id);
	CREATE INDEX IF NOT EXISTS idx_messages_sender_id ON messages(sender_id);
	CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);
	CREATE INDEX IF NOT EXISTS idx_reactions_message_id ON message_reactions(message_id);
	`

	_, err := db.c.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to initialize database schema: %w", err)
	}

	return nil
}
