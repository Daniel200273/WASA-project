package database

import "time"

// User rappresenta un utente dell'applicazione
type User struct {
	ID        string    `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	PhotoURL  *string   `json:"photoUrl,omitempty" db:"photo_url"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

// UserSession rappresenta una sessione di autenticazione
type UserSession struct {
	Token     string    `db:"token"`
	UserID    string    `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}

// Conversation rappresenta una conversazione (diretta o di gruppo)
type Conversation struct {
	ID        string    `json:"id" db:"id"`
	Type      string    `json:"type" db:"type"` // "direct" | "group"
	Name      *string   `json:"name,omitempty" db:"name"`
	PhotoURL  *string   `json:"photoUrl,omitempty" db:"photo_url"`
	CreatedBy *string   `json:"createdBy,omitempty" db:"created_by"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`

	// Campi calcolati per l'API (non salvati nel DB)
	Members          []User    `json:"members,omitempty"`
	Participants     []User    `json:"participants,omitempty"`
	OtherParticipant *User     `json:"otherParticipant,omitempty"` // Per conversazioni dirette
	LastMessage      *Message  `json:"lastMessage,omitempty"`
	UnreadCount      int       `json:"unreadCount"`
	Messages         []Message `json:"messages,omitempty"`
}

// Message rappresenta un messaggio in una conversazione
type Message struct {
	ID             string    `json:"id" db:"id"`
	ConversationID string    `json:"-" db:"conversation_id"`
	SenderID       string    `json:"senderId" db:"sender_id"`
	SenderUsername string    `json:"senderUsername"` // Campo joined dalle query
	Content        *string   `json:"content,omitempty" db:"content"`
	PhotoURL       *string   `json:"photoUrl,omitempty" db:"photo_url"`
	ReplyToID      *string   `json:"replyTo,omitempty" db:"reply_to_id"`
	Forwarded      bool      `json:"forwarded" db:"forwarded"`
	Status         string    `json:"status"` // "sent", "delivered", "read"
	CreatedAt      time.Time `json:"timestamp" db:"created_at"`

	Comments []MessageReaction `json:"comments,omitempty"`
}

// MessagePreview rappresenta un'anteprima di messaggio per la lista conversazioni
type MessagePreview struct {
	ID             string    `json:"id"`
	Content        *string   `json:"content,omitempty"`
	Timestamp      time.Time `json:"timestamp"`
	SenderUsername string    `json:"senderUsername"`
	HasPhoto       bool      `json:"hasPhoto"`
}

// MessageReaction rappresenta una reazione emoji a un messaggio
type MessageReaction struct {
	ID        string    `json:"id" db:"id"`
	MessageID string    `json:"-" db:"message_id"`
	UserID    string    `json:"userId" db:"user_id"`
	Username  string    `json:"username"` // Campo joined dalle query
	Emoticon  string    `json:"emoticon" db:"emoticon"`
	CreatedAt time.Time `json:"timestamp" db:"created_at"`
}

// ConversationParticipant rappresenta la partecipazione di un utente a una conversazione
type ConversationParticipant struct {
	ConversationID string    `db:"conversation_id"`
	UserID         string    `db:"user_id"`
	JoinedAt       time.Time `db:"joined_at"`
}

// Group rappresenta le informazioni specifiche di un gruppo (estende Conversation)
type Group struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	PhotoURL  *string   `json:"photoUrl,omitempty"`
	Members   []User    `json:"members"`
	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
}
