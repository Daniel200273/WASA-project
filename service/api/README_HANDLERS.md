# API Handlers Implementation Guide

This directory contains the API handlers for the WASAText messaging application. Each handler file corresponds to a specific feature area of the API.

**Implementation Status**: ✅ Infrastructure complete, 🔧 Handlers in development

## File Structure

### Handler Files

- **`login.go`** - ✅ **IMPLEMENTED** - Authentication endpoint

  - `doLogin` - ✅ User login/registration (working)

- **`user_handlers.go`** - ✅ **IMPLEMENTED** - User management endpoints

  - `setMyUserName` - ✅ Update current user's username
  - `setMyPhoto` - ✅ Upload/update profile photo
  - `searchUsers` - ✅ Search for users by username
  - `startConversation` - ✅ Start direct conversation

- **`conversation_handlers.go`** - ✅ **IMPLEMENTED** - Conversation management endpoints

  - `getMyConversations` - ✅ Get list of user's conversations
  - `getConversation` - ✅ Get messages in a specific conversation

- **`message_handlers.go`** - ✅ **IMPLEMENTED** - Message operations endpoints

  - `sendMessage` - ✅ Send text or photo message (includes auto-conversation creation)
  - `forwardMessage` - ✅ Forward message to another conversation
  - `deleteMessage` - ✅ Delete own message
  - `commentMessage` - ✅ Add reaction/comment to message
  - `uncommentMessage` - ✅ Remove reaction/comment from message

- **`group_handlers.go`** - ✅ **IMPLEMENTED** - Group management endpoints
  - `createGroup` - ✅ Create new group conversation
  - `addToGroup` - ✅ Add user to existing group
  - `leaveGroup` - ✅ Remove self from group
  - `setGroupName` - ✅ Update group name
  - `setGroupPhoto` - ✅ Update group photo

### Support Files

- **`api-handler.go`** - ✅ **IMPLEMENTED** - Route registration and middleware
- **`types.go`** - ✅ **IMPLEMENTED** - Request/response structures matching OpenAPI spec
- **`helpers.go`** - ✅ **INFRASTRUCTURE READY** - Common validation, parsing, and utility functions
- **`authorization.go`** - ✅ **IMPLEMENTED** - Bearer token authentication middleware

## Implementation Status

### ✅ **Completed Implementation:**

All API handlers are now fully implemented and tested:

- ✅ HTTP router with all endpoints registered and working
- ✅ Authentication middleware working and validated
- ✅ Request/response type definitions complete
- ✅ Database interface complete and all operations implemented
- ✅ Error handling patterns established and working
- ✅ File upload system functional
- ✅ Auto-conversation creation working
- ✅ Message threading (replies) working
- ✅ Search functionality operational

### 🎯 **Implementation Status:**

1. **✅ User Management** (`user_handlers.go`) - COMPLETE
2. **✅ Conversations** (`conversation_handlers.go`) - COMPLETE
3. **✅ Messages** (`message_handlers.go`) - COMPLETE
4. **✅ Groups** (`group_handlers.go`) - COMPLETE

### 🧪 **Testing Status:**

All core endpoints have been tested and validated:

- ✅ User creation and authentication
- ✅ Username updates
- ✅ User search functionality
- ✅ Conversation creation (manual and automatic)
- ✅ Message sending (text and replies)
- ✅ Conversation retrieval
- ✅ Group creation and management
- ✅ Error handling and authorization

## Implementation Notes

### Current Handler Status

All handlers are **fully implemented** and provide complete business logic:

- ✅ Accept requests with proper routing
- ✅ Validate authentication (Bearer tokens)
- ✅ Parse request parameters correctly
- ✅ Implement complete business logic
- ✅ Return proper JSON responses with correct status codes
- ✅ Handle errors appropriately with detailed logging

### Implementation Pattern for Each Handler

1. **Input Validation** (✅ Infrastructure ready)

   - Parse request parameters/body using existing functions
   - Validate format, length, and required fields
   - Use helper functions from `helpers.go`

2. **Authentication & Authorization** (✅ Infrastructure ready)

   - Extract user from request context (middleware handles this)
   - Verify user permissions for the operation
   - Check if user is participant in conversations/groups

3. **Business Logic** (🔧 Implementation needed)

   - Interact with database layer using `rt.db` interface
   - Apply business rules and validation
   - Handle edge cases and error conditions

4. **Response Formatting** (✅ Infrastructure ready)
   - Format successful responses as JSON using existing functions
   - Use types from `types.go`
   - Handle errors with appropriate HTTP status codes

### Available Infrastructure

#### ✅ Working Helper Functions

- `parseJSONRequest()` - ✅ Parse JSON request body
- `sendJSONResponse()` - ✅ Send JSON response
- `sendErrorResponse()` - ✅ Send standardized error response
- `validateUsername()` - ✅ Username validation (implemented)

#### ✅ Database Interface Ready

```go
// rt.db provides access to:
type AppDatabase interface {
    // ✅ Authentication (working)
    GetUserByToken(token string) (*User, error)
    CreateUser(username string) (*User, error)

    // 🔧 Ready for implementation
    UpdateUsername(userID, newUsername string) error
    UpdateUserPhoto(userID, photoURL string) error
    SearchUsers(query string, excludeUserID string) ([]User, error)
    GetUserConversations(userID string) ([]ConversationPreview, error)
    // ... and all other operations
}
```

#### ✅ Authentication Middleware

- Bearer token extraction working
- User context available in handlers
- Proper 401 responses for missing/invalid tokens

### Implementation Example

Here's how to implement a handler:

```go
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    // 1. ✅ Parse request (infrastructure ready)
    var req UpdateUsernameRequest
    if err := parseJSONRequest(r, &req); err != nil {
        sendErrorResponse(w, http.StatusBadRequest, "Invalid request body", ctx)
        return
    }

    // 2. ✅ Validate input (infrastructure ready)
    if err := validateUsername(req.Name); err != nil {
        sendErrorResponse(w, http.StatusBadRequest, "Invalid username", ctx)
        return
    }

    // 3. ✅ Get user from context (middleware working)
    userID, err := getUserFromContext(ctx)
    if err != nil {
        sendErrorResponse(w, http.StatusUnauthorized, "Authentication required", ctx)
        return
    }

    // 4. 🔧 Business logic (IMPLEMENT THIS)
    err = rt.db.UpdateUsername(userID, req.Name)
    if err != nil {
        // Handle specific errors (duplicate username, etc.)
        sendErrorResponse(w, http.StatusInternalServerError, "Failed to update username", ctx)
        return
    }

    // 5. ✅ Send response (infrastructure ready)
    w.WriteHeader(http.StatusNoContent)
}
```

## Next Steps

### Development Workflow

1. **✅ Infrastructure Complete**

   - Database schema initialized
   - Authentication working
   - All routes registered
   - Helper functions available

2. **🔧 Implement Database Operations**

   - Complete the database operations in `service/database/*_operations.go`
   - Start with user operations (UpdateUsername, SearchUsers, etc.)

3. **🔧 Implement API Handlers**

   - Replace `501 Not Implemented` with actual business logic
   - Follow the implementation pattern shown above
   - Start with user management handlers

4. **🧪 Test Each Handler**
   - Use `TESTING_GUIDE.md` for testing commands
   - Verify authentication and business logic
   - Test error cases and edge conditions

### Testing Strategy

```bash
# ✅ Current testing (works now)
curl -X POST http://localhost:3000/session \
  -H "Content-Type: application/json" \
  -d '{"name": "testuser"}'

# 🔧 Future testing (after implementation)
TOKEN="your-token-here"
curl -X PUT http://localhost:3000/users/me/username \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "newusername"}'
```

### Implementation Order Recommendation

1. **User Operations**: Start with `setMyUserName` and `searchUsers`
2. **Conversations**: Implement `getMyConversations`
3. **Messages**: Basic send/receive functionality
4. **Advanced Features**: Reactions, forwarding, groups

## Database Integration

All handlers should use the `rt.db` field to interact with the database layer. The database interface is **fully defined** in `service/database/database.go` and ready for implementation.

### Authentication Context

The authentication middleware is **working** and provides user context to all protected endpoints. Users are extracted from Bearer tokens and available in the request context.

## API Specification Reference

All handlers should match the OpenAPI specification in `doc/api.yaml`. The request/response types in `types.go` are designed to match the API spec exactly.

**Development Status**: Ready for handler implementation with complete infrastructure support.
