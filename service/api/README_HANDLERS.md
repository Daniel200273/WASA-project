# API Handlers Implementation Guide

This directory contains skeleton API handlers for the WASAText messaging application. Each handler file corresponds to a specific feature area of the API.

## File Structure

### Handler Files

- **`user_handlers.go`** - User management endpoints

  - `setMyUserName` - Update current user's username
  - `setMyPhoto` - Upload/update profile photo
  - `searchUsers` - Search for users by username

- **`conversation_handlers.go`** - Conversation management endpoints

  - `getMyConversations` - Get list of user's conversations
  - `getConversation` - Get messages in a specific conversation

- **`message_handlers.go`** - Message operations endpoints

  - `sendMessage` - Send text or photo message
  - `forwardMessage` - Forward message to another conversation
  - `deleteMessage` - Delete own message
  - `commentMessage` - Add reaction/comment to message
  - `uncommentMessage` - Remove reaction/comment from message

- **`group_handlers.go`** - Group management endpoints
  - `createGroup` - Create new group conversation
  - `addToGroup` - Add user to existing group
  - `leaveGroup` - Remove self from group
  - `setGroupName` - Update group name
  - `setGroupPhoto` - Update group photo

### Support Files

- **`types.go`** - Request/response structures matching OpenAPI spec
- **`helpers.go`** - Common validation, parsing, and utility functions
- **`validateUsername.go`** - Username validation (already implemented)
- **`login.go`** - Authentication login handler (already implemented)

## Implementation Notes

### Current Status

All handlers are **skeleton implementations** that return `501 Not Implemented`. Each handler contains detailed TODO comments outlining the implementation steps.

### Key Implementation Steps for Each Handler

1. **Input Validation**

   - Parse request parameters/body
   - Validate format, length, and required fields
   - Use helper functions from `helpers.go`

2. **Authentication & Authorization**

   - Extract user from request context
   - Verify user permissions for the operation
   - Check if user is participant in conversations/groups

3. **Business Logic**

   - Interact with database layer
   - Apply business rules
   - Handle edge cases and error conditions

4. **Response Formatting**
   - Format successful responses as JSON
   - Use types from `types.go`
   - Handle errors with appropriate HTTP status codes

### Helper Functions Available

#### Validation

- `validateID()` - Validate UUID-format IDs
- `validateUsername()` - Validate username format (already implemented)
- `validateGroupName()` - Validate group name
- `validateMessageContent()` - Validate message text
- `validateEmoticon()` - Validate reaction emoticons
- `validateSearchQuery()` - Validate search queries
- `validateImageFile()` - Validate uploaded images

#### HTTP Utilities

- `parseJSONRequest()` - Parse JSON request body
- `sendJSONResponse()` - Send JSON response
- `sendErrorResponse()` - Send standardized error response
- `getPathParam()` - Extract URL path parameters
- `getQueryParam()` - Extract URL query parameters

#### File Handling

- `validateImageFile()` - Validate image uploads for photos

### Error Handling Pattern

```go
// Validate input
if err := validateID(messageId, "messageId"); err != nil {
    sendErrorResponse(w, http.StatusBadRequest, err.Error(), ctx)
    return
}

// Get user context
userID, err := getUserFromContext(ctx)
if err != nil {
    sendErrorResponse(w, http.StatusUnauthorized, "Authentication required", ctx)
    return
}

// Business logic with database
result, err := rt.db.SomeOperation(userID, messageId)
if err != nil {
    ctx.Logger.WithError(err).Error("database operation failed")
    sendErrorResponse(w, http.StatusInternalServerError, "Internal server error", ctx)
    return
}

// Success response
sendJSONResponse(w, http.StatusOK, result)
```

### Database Integration

All handlers should use the `rt.db` field to interact with the database layer. The database interface is defined in `service/database/database.go`.

### Authentication Context

The `getUserFromContext()` function in `helpers.go` is a placeholder. You'll need to implement this based on your authentication middleware that sets user information in the request context.

## Next Steps

1. **Implement Authentication Context** - Complete `getUserFromContext()` function
2. **Choose Handler Priority** - Start with core functionality (login, conversations, messages)
3. **Test Each Handler** - Use tools like Postman or curl to test endpoints
4. **Add Business Logic** - Implement the TODO items in each handler
5. **Error Handling** - Ensure proper error responses for all edge cases
6. **File Upload** - Implement photo/image handling for profile and group photos

## API Specification Reference

All handlers should match the OpenAPI specification in `doc/api.yaml`. The request/response types in `types.go` are designed to match the API spec exactly.
