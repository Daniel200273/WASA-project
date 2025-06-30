package api

import "time"

// Common request/response structures for API handlers
// These match the OpenAPI specification in doc/api.yaml

// === REQUEST STRUCTURES ===

// LoginRequest represents the login request body
type LoginRequest struct {
	Name string `json:"name"`
}

// UpdateUsernameRequest represents username update request
type UpdateUsernameRequest struct {
	Name string `json:"name"`
}

// SendMessageRequest represents text message sending request
type SendMessageRequest struct {
	Content string  `json:"content"`
	ReplyTo *string `json:"replyTo,omitempty"`
}

// ForwardMessageRequest represents message forwarding request
type ForwardMessageRequest struct {
	ConversationID string `json:"conversationId"`
}

// CommentMessageRequest represents message reaction request
type CommentMessageRequest struct {
	Emoticon string `json:"emoticon"`
}

// CreateGroupRequest represents group creation request
type CreateGroupRequest struct {
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

// AddToGroupRequest represents adding user to group request
type AddToGroupRequest struct {
	UserID string `json:"userId"`
}

// UpdateGroupNameRequest represents group name update request
type UpdateGroupNameRequest struct {
	Name string `json:"name"`
}

// StartConversationRequest represents starting a conversation request
type StartConversationRequest struct {
	UserID string `json:"userId"`
}

// === RESPONSE STRUCTURES ===

// LoginResponse represents the login response body
type LoginResponse struct {
	Identifier string `json:"identifier"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Message string `json:"message"`
}

// UserResponse represents user information
type UserResponse struct {
	ID       string  `json:"id"`
	Username string  `json:"username"`
	PhotoURL *string `json:"photoUrl,omitempty"`
}

// SearchUsersResponse represents user search results
type SearchUsersResponse struct {
	Users []UserResponse `json:"users"`
}

// MessagePreview represents a message preview in conversation list
type MessagePreview struct {
	ID             string    `json:"id"`
	Content        *string   `json:"content,omitempty"`
	Timestamp      time.Time `json:"timestamp"`
	SenderUsername string    `json:"senderUsername"`
	HasPhoto       bool      `json:"hasPhoto"`
}

// ConversationResponse represents a conversation in the list
type ConversationResponse struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"` // "direct" or "group"
	Name        string          `json:"name"`
	PhotoURL    *string         `json:"photoUrl,omitempty"`
	LastMessage *MessagePreview `json:"lastMessage,omitempty"`
	UnreadCount int             `json:"unreadCount"`
}

// ConversationsResponse represents the list of user's conversations
type ConversationsResponse struct {
	Conversations []ConversationResponse `json:"conversations"`
}

// CommentResponse represents a message reaction/comment
type CommentResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Username  string    `json:"username"`
	Emoticon  string    `json:"emoticon"`
	Timestamp time.Time `json:"timestamp"`
}

// MessageResponse represents a message with all details
type MessageResponse struct {
	ID             string            `json:"id"`
	SenderID       string            `json:"senderId"`
	SenderUsername string            `json:"senderUsername"`
	Content        *string           `json:"content,omitempty"`
	PhotoURL       *string           `json:"photoUrl,omitempty"`
	ReplyToID      *string           `json:"replyToId,omitempty"`
	Forwarded      bool              `json:"forwarded,omitempty"`
	Timestamp      time.Time         `json:"timestamp"`
	Status         string            `json:"status"` // "sent", "delivered", "read"
	Comments       []CommentResponse `json:"comments"`
}

// ConversationDetailResponse represents conversation details with messages
type ConversationDetailResponse struct {
	ID       string            `json:"id"`
	Type     string            `json:"type"`
	Name     string            `json:"name"`
	PhotoURL *string           `json:"photoUrl,omitempty"`
	Members  []UserResponse    `json:"members"`
	Messages []MessageResponse `json:"messages"`
}

// GroupResponse represents a group with all details
type GroupResponse struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	PhotoURL  *string        `json:"photoUrl,omitempty"`
	Members   []UserResponse `json:"members"`
	CreatedBy string         `json:"createdBy"`
	CreatedAt time.Time      `json:"createdAt"`
}

// EmptyResponse represents an empty success response
type EmptyResponse struct{}
