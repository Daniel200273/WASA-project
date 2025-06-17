# 🗄️ WASAText Database Design

## 📋 **Panoramica**

Database **SQLite3** per l'applicazione di messaggistica WASAText, progettato per essere semplice, efficiente e adatto al progetto universitario WASA.

---

## 🏗️ **Schema Database Semplificato**

### **1. Tabella USERS**

```sql
-- Utenti dell'applicazione
CREATE TABLE users (
    id TEXT PRIMARY KEY,                    -- UUID come string
    username TEXT UNIQUE NOT NULL,          -- Username univoco
    photo_url TEXT,                        -- Path foto profilo (opzionale)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Indice per ricerca veloce
CREATE INDEX idx_users_username ON users(username);
```

### **2. Tabella USER_SESSIONS**

```sql
-- Sessioni/token di autenticazione
CREATE TABLE user_sessions (
    token TEXT PRIMARY KEY,                -- Bearer token (UUID)
    user_id TEXT NOT NULL,                -- Riferimento a users.id
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
```

### **3. Tabella CONVERSATIONS**

```sql
-- Conversazioni (dirette e di gruppo)
CREATE TABLE conversations (
    id TEXT PRIMARY KEY,                   -- UUID
    type TEXT NOT NULL CHECK (type IN ('direct', 'group')),
    name TEXT,                            -- Nome gruppo (NULL per direct)
    photo_url TEXT,                       -- Foto gruppo (opzionale)
    created_by TEXT,                      -- Chi ha creato (solo gruppi)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE INDEX idx_conversations_type ON conversations(type);
```

### **4. Tabella CONVERSATION_PARTICIPANTS**

```sql
-- Chi partecipa a quali conversazioni
CREATE TABLE conversation_participants (
    conversation_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (conversation_id, user_id),
    FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_participants_user_id ON conversation_participants(user_id);
```

### **5. Tabella MESSAGES**

```sql
-- Messaggi nelle conversazioni
CREATE TABLE messages (
    id TEXT PRIMARY KEY,                   -- UUID
    conversation_id TEXT NOT NULL,
    sender_id TEXT NOT NULL,
    content TEXT,                         -- Testo messaggio (NULL se foto)
    photo_url TEXT,                       -- Foto messaggio (NULL se testo)
    reply_to_id TEXT,                     -- Risposta a messaggio
    forwarded BOOLEAN DEFAULT FALSE,       -- Messaggio inoltrato
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (reply_to_id) REFERENCES messages(id) ON DELETE SET NULL,

    -- Deve avere contenuto O foto, non entrambi NULL
    CHECK ((content IS NOT NULL AND photo_url IS NULL) OR
           (content IS NULL AND photo_url IS NOT NULL))
);

-- Indici per performance
CREATE INDEX idx_messages_conversation_id ON messages(conversation_id);
CREATE INDEX idx_messages_sender_id ON messages(sender_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);
```

### **6. Tabella MESSAGE_REACTIONS**

```sql
-- Reazioni emoji ai messaggi
CREATE TABLE message_reactions (
    id TEXT PRIMARY KEY,                   -- UUID
    message_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    emoticon TEXT NOT NULL,               -- Emoji
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

    -- Un utente può reagire una sola volta per messaggio
    UNIQUE (message_id, user_id)
);

CREATE INDEX idx_reactions_message_id ON message_reactions(message_id);
```

---

## 📊 **Modelli Go**

### **service/database/models.go**

```go
package database

import "time"

// User rappresenta un utente
type User struct {
    ID       string    `json:"id" db:"id"`
    Username string    `json:"username" db:"username"`
    PhotoURL *string   `json:"photoUrl,omitempty" db:"photo_url"`
    CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

// UserSession rappresenta una sessione utente
type UserSession struct {
    Token     string    `db:"token"`
    UserID    string    `db:"user_id"`
    CreatedAt time.Time `db:"created_at"`
}

// Conversation rappresenta una conversazione
type Conversation struct {
    ID        string    `json:"id" db:"id"`
    Type      string    `json:"type" db:"type"` // "direct" | "group"
    Name      *string   `json:"name,omitempty" db:"name"`
    PhotoURL  *string   `json:"photoUrl,omitempty" db:"photo_url"`
    CreatedBy *string   `json:"createdBy,omitempty" db:"created_by"`
    CreatedAt time.Time `json:"createdAt" db:"created_at"`

    // Campi calcolati (non in DB)
    Members      []User           `json:"members,omitempty"`
    LastMessage  *MessagePreview  `json:"lastMessage,omitempty"`
    UnreadCount  int             `json:"unreadCount"`
}

// Message rappresenta un messaggio
type Message struct {
    ID             string    `json:"id" db:"id"`
    ConversationID string    `json:"-" db:"conversation_id"`
    SenderID       string    `json:"senderId" db:"sender_id"`
    SenderUsername string    `json:"senderUsername"` // Campo joined
    Content        *string   `json:"content,omitempty" db:"content"`
    PhotoURL       *string   `json:"photoUrl,omitempty" db:"photo_url"`
    ReplyToID      *string   `json:"replyTo,omitempty" db:"reply_to_id"`
    Forwarded      bool      `json:"forwarded" db:"forwarded"`
    Status         string    `json:"status"` // "sent", "delivered", "read"
    CreatedAt      time.Time `json:"timestamp" db:"created_at"`

    Comments []MessageReaction `json:"comments,omitempty"`
}

// MessagePreview per la lista conversazioni
type MessagePreview struct {
    ID             string    `json:"id"`
    Content        *string   `json:"content,omitempty"`
    Timestamp      time.Time `json:"timestamp"`
    SenderUsername string    `json:"senderUsername"`
    HasPhoto       bool      `json:"hasPhoto"`
}

// MessageReaction rappresenta una reazione
type MessageReaction struct {
    ID        string    `json:"id" db:"id"`
    MessageID string    `json:"-" db:"message_id"`
    UserID    string    `json:"userId" db:"user_id"`
    Username  string    `json:"username"` // Campo joined
    Emoticon  string    `json:"emoticon" db:"emoticon"`
    CreatedAt time.Time `json:"timestamp" db:"created_at"`
}
```

---

## 🔧 **Interfaccia Database**

### **Aggiorna service/database/database.go**

```go
package database

import (
    "database/sql"
    "errors"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
)

// AppDatabase interfaccia principale
type AppDatabase interface {
    // Esistenti (mantieni per compatibilità)
    GetName() (string, error)
    SetName(name string) error
    Ping() error

    // === AUTHENTICATION ===
    CreateUser(username string) (*User, error)
    GetUserByID(id string) (*User, error)
    GetUserByToken(token string) (*User, error)
    CreateUserSession(userID string) (string, error)
    DeleteUserSession(token string) error

    // === USER MANAGEMENT ===
    UpdateUsername(userID, newUsername string) error
    UpdateUserPhoto(userID, photoURL string) error
    SearchUsers(query string, excludeUserID string) ([]User, error)

    // === CONVERSATIONS ===
    GetUserConversations(userID string) ([]Conversation, error)
    GetConversation(conversationID, userID string) (*Conversation, error)
    GetOrCreateDirectConversation(user1ID, user2ID string) (*Conversation, error)

    // === MESSAGES ===
    CreateMessage(conversationID, senderID string, content *string, photoURL *string, replyToID *string) (*Message, error)
    GetMessage(messageID string) (*Message, error)
    DeleteMessage(messageID, userID string) error
    ForwardMessage(messageID, targetConversationID, userID string) (*Message, error)

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

// New crea una nuova istanza del database
func New(db *sql.DB) (AppDatabase, error) {
    if db == nil {
        return nil, errors.New("database is required when building a AppDatabase")
    }

    appDB := &appdbimpl{c: db}

    // Inizializza lo schema se necessario
    if err := appDB.initializeSchema(); err != nil {
        return nil, fmt.Errorf("error initializing database schema: %w", err)
    }

    return appDB, nil
}

// initializeSchema crea le tabelle se non esistono
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
        FOREIGN KEY (created_by) REFERENCES users(id)
    );

    -- Conversation participants table
    CREATE TABLE IF NOT EXISTS conversation_participants (
        conversation_id TEXT NOT NULL,
        user_id TEXT NOT NULL,
        joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
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
    return err
}

func (db *appdbimpl) Ping() error {
    return db.c.Ping()
}
```

---

## 📝 **Query di Esempio**

### **Lista Conversazioni per Utente**

```sql
SELECT DISTINCT
    c.id,
    c.type,
    CASE
        WHEN c.type = 'direct' THEN (
            SELECT u2.username
            FROM conversation_participants cp2
            JOIN users u2 ON cp2.user_id = u2.id
            WHERE cp2.conversation_id = c.id AND cp2.user_id != ?
            LIMIT 1
        )
        ELSE c.name
    END as display_name,
    CASE
        WHEN c.type = 'direct' THEN (
            SELECT u2.photo_url
            FROM conversation_participants cp2
            JOIN users u2 ON cp2.user_id = u2.id
            WHERE cp2.conversation_id = c.id AND cp2.user_id != ?
            LIMIT 1
        )
        ELSE c.photo_url
    END as display_photo
FROM conversations c
JOIN conversation_participants cp ON c.id = cp.conversation_id
WHERE cp.user_id = ?
ORDER BY c.created_at DESC;
```

### **Messaggi di una Conversazione**

```sql
SELECT
    m.id,
    m.sender_id,
    u.username as sender_username,
    m.content,
    m.photo_url,
    m.reply_to_id,
    m.forwarded,
    m.created_at
FROM messages m
JOIN users u ON m.sender_id = u.id
WHERE m.conversation_id = ?
ORDER BY m.created_at DESC
LIMIT 100;
```

---

## 🚀 **Prossimi Passi**

1. **Crea il file models.go** con le strutture dati
2. **Aggiorna database.go** con il nuovo schema
3. **Implementa i metodi del database** uno alla volta
4. **Testa ogni metodo** con unit test semplici
5. **Integra con l'API** esistente

## 💡 **Vantaggi di questo Schema**

- ✅ **Semplice**: Poche tabelle, relazioni chiare
- ✅ **Efficace**: Supporta tutti i requisiti dell'API
- ✅ **Performante**: Indici ottimizzati
- ✅ **Testabile**: Struttura lineare e prevedibile
- ✅ **Estendibile**: Facile aggiungere nuove feature

Vuoi che iniziamo implementando questo schema nel tuo progetto?
