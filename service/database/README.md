# Database Operations Implementation Guide

This directory contains the skeleton implementation for all database operations defined in the `AppDatabase` interface. Each operation has been organized into separate files for better maintainability.

## File Structure

- `database.go` - Main database interface and initialization
- `models.go` - Data models and structures
- `auth_operations.go` - Authentication related operations
- `user_operations.go` - User management operations
- `conversation_operations.go` - Conversation management operations
- `message_operations.go` - Message operations
- `reaction_operations.go` - Message reaction operations
- `group_operations.go` - Group management operations
- `utils.go` - Helper functions and utilities

## Implementation Guidelines

### 1. Error Handling

- Use `sql.ErrNoRows` for "not found" cases
- Return meaningful error messages
- Use `fmt.Errorf()` to wrap errors with context

### 2. SQL Queries

- Use prepared statements with placeholders (`?`) for security
- Join tables when you need related data (e.g., usernames)
- Use indexes defined in the schema for performance

### 3. Transactions

- Use transactions for operations that modify multiple tables
- Always rollback on errors and commit on success
- Example pattern:

```go
tx, err := db.c.Begin()
if err != nil {
    return err
}
defer tx.Rollback()

// ... perform operations

return tx.Commit()
```

### 4. Data Validation

- Validate input parameters before database operations
- Check foreign key constraints manually where needed
- Verify user permissions for operations

### 5. Common Patterns

#### User Authorization

Many operations require checking if a user has permission:

```go
// Check if user is in conversation
func (db *appdbimpl) IsUserInConversation(conversationID, userID string) (bool, error)
```

#### Pagination

For operations that return lists, consider adding pagination:

```go
// Add LIMIT and OFFSET to queries
query := "SELECT ... LIMIT ? OFFSET ?"
```

#### Joins

When you need related data, use JOINs instead of separate queries:

```go
query := `
    SELECT m.*, u.username
    FROM messages m
    JOIN users u ON m.sender_id = u.id
    WHERE m.conversation_id = ?
`
```

## Implementation Priority

1. **Authentication Operations** (`auth_operations.go`)

   - `CreateUser` - Basic user creation
   - `GetUserByUsername` - For login
   - `CreateUserSession` - For session management

2. **User Operations** (`user_operations.go`)

   - `GetUserByID` - Basic user retrieval
   - `SearchUsers` - For finding users to chat with

3. **Conversation Operations** (`conversation_operations.go`)

   - `GetOrCreateDirectConversation` - For starting chats
   - `GetUserConversations` - For listing user's chats

4. **Message Operations** (`message_operations.go`)

   - `CreateMessage` - For sending messages
   - `GetMessage` - For retrieving messages

5. **Group Operations** (`group_operations.go`)

   - `CreateGroup` - For group creation
   - `AddUserToGroup` - For group management

6. **Reaction Operations** (`reaction_operations.go`)
   - `CreateMessageReaction` - For emoji reactions

## Database Schema Reference

The schema is defined in `database.go`. Key tables:

- `users` - User information
- `user_sessions` - Authentication sessions
- `conversations` - Chat conversations (direct/group)
- `conversation_participants` - Who's in each conversation
- `messages` - Chat messages
- `message_reactions` - Emoji reactions to messages

## Testing

After implementing each operation:

1. Test with valid data
2. Test error cases (user not found, permission denied, etc.)
3. Test edge cases (empty strings, null values, etc.)
4. Test concurrent operations if applicable

## Example Implementation

Here's an example of how to implement `GetUserByID`:

```go
func (db *appdbimpl) GetUserByID(id string) (*User, error) {
    query := `
        SELECT id, username, photo_url, created_at
        FROM users
        WHERE id = ?
    `

    row := db.c.QueryRow(query, id)
    user, err := scanUser(row)
    if err != nil {
        if isNotFoundError(err) {
            return nil, fmt.Errorf("user not found")
        }
        return nil, fmt.Errorf("error getting user: %w", err)
    }

    return user, nil
}
```

Remember to remove the placeholder return statements and the `_ = variable` lines when implementing each function.
