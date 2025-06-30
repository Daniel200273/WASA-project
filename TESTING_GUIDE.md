# WASAText API Testing Guide

This guide contains instructions and curl commands for testing the WASAText messaging application API.

**Current Implementation Status**: ✅ Authentication working, 🔧 Other endpoints in development

## Table of Contents

- [Getting Started](#getting-started)
- [Understanding curl Flags](#understanding-curl-flags)
- [Authentication (✅ Working)](#authentication)
- [User Management (🔧 In Development)](#user-management)
- [Conversations (🔧 In Development)](#conversations)
- [Messages (🔧 In Development)](#messages)
- [Groups (🔧 In Development)](#groups)
- [Testing Tips](#testing-tips)
  - [Understanding HTTP Status Codes](#understanding-http-status-codes)
  - [Viewing HTTP Status Codes](#viewing-http-status-codes)
  - [Getting More Details](#getting-more-details)
  - [Complete Test Flow Example](#complete-test-flow-example)
  - [Note About IDs](#note-about-ids)
  - [Testing With Multiple Users](#testing-with-multiple-users)
    - [Example Multi-User Test Flow](#example-multi-user-test-flow)
    - [Testing Group Interactions](#testing-group-interactions)

## Getting Started

### Summary

Backend (Go API Server)

```bash
# Build and start the API server
cd /Users/daniel/Desktop/WASA-project
go build ./cmd/webapi
./webapi &
SERVER_PID=$!

# Quick health check
curl -s http://localhost:3000/liveness
```

Frontend (Vue.js Developement)

```bash
# Start development server with hot reload
cd /Users/daniel/Desktop/WASA-project/webui
yarn dev
# Runs on http://localhost:5173 (or similar)
```

Building Frontend (do before commit)

```bash
cd /Users/daniel/Desktop/WASA-project/webui

# Development build
yarn build-dev

# Production build
yarn build-prod

# Preview production build
yarn preview
```

Stopping Server

```bash
# Stop backend server (triggers cleanup)
kill $SERVER_PID

# Or if you don't have the PID
pkill -SIGTERM webapi

# Stop frontend dev server
# Press Ctrl+C in the terminal running yarn dev
```

Quick Server Checks

```bash
# Check if servers are running
lsof -i :3000  # Backend
lsof -i :5173  # Frontend (or check yarn dev output for actual port)

# Health check backend
curl -s http://localhost:3000/liveness && echo "✅ Backend running" || echo "❌ Backend down"
```

### Starting the Server

Before running any API tests, you'll need to start the WASAText server:

```bash
# Build the WASAText application
go build ./cmd/webapi

# Start the server (runs on port 3000)
./webapi &

# Store the process ID for later use
SERVER_PID=$!

# ✅ Verify the server is running (this should work)
curl -s http://localhost:3000/liveness
# Expected response: Liveness check passed
```

### Checking if a Server is Running

Before starting a new server, you may want to check if one is already running:

```bash
# Check if port 3000 is in use
lsof -i :3000

# Check for running webapi processes
ps aux | grep webapi

# ✅ Quick health check (this endpoint is working)
curl -s http://localhost:3000/liveness && echo "✅ Server is running" || echo "❌ Server is not running"

# Kill any existing server if needed (use graceful shutdown)
pkill -SIGTERM webapi
```

### Database Initialization

The WASAText server automatically initializes the SQLite database on startup:

```bash
# The database file will be created automatically
# Check if database was created (after starting server)
ls -la *.db

# You can also check server logs for database initialization messages
```

### Ending Your Test Session

When you're finished testing, make sure to properly shut down the server to trigger cleanup:

```bash
# ✅ Graceful shutdown (recommended - triggers tmp cleanup)
kill $SERVER_PID

# Or if you didn't save the PID, use SIGTERM for graceful shutdown:
pkill -SIGTERM webapi

# Wait for the process to terminate
wait $SERVER_PID 2>/dev/null

# ❌ Avoid force kill (prevents cleanup):
# kill -9 $SERVER_PID  # This prevents tmp directory cleanup
# pkill -9 webapi      # This also prevents cleanup
```

**Important**: Use graceful shutdown (`kill` or `pkill -SIGTERM`) to ensure the `tmp/uploads` directory is properly cleaned up. Force killing the process (`kill -9` or `pkill -9`) will leave temporary files behind.

## Understanding curl Flags

The curl commands in this guide use several flags. Here's what each one does:

| Flag                   | Description                                                           |
| ---------------------- | --------------------------------------------------------------------- |
| `-s` or `--silent`     | Silent mode that suppresses the progress meter but still shows errors |
| `-X [METHOD]`          | Specifies the HTTP method (GET, POST, PUT, DELETE, etc.)              |
| `-H "Header: Value"`   | Sets an HTTP header for the request                                   |
| `-d 'data'`            | Sends data in the request body (for POST, PUT)                        |
| `-F "field=@filename"` | Uploads a file using multipart/form-data encoding                     |
| `-w "format"`          | Displays information after a completed transfer (e.g., status code)   |
| `-o /dev/null`         | Redirects the response output to /dev/null (discards it)              |
| `-v`                   | Verbose mode that shows detailed request and response information     |

## Authentication (✅ Working)

The authentication system is **fully implemented and working**:

```bash
# ✅ Create a new user session / login (THIS WORKS)
curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{"name": "testuser"}' \
  http://localhost:3000/session

# Expected response:
# {"identifier":"some-uuid-token-here"}

# ✅ Save the token for later use (THIS WORKS)
TOKEN=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{"name": "testuser"}' \
  http://localhost:3000/session | \
  grep -o '"identifier":"[^"]*"' | cut -d'"' -f4)

echo "Your token: $TOKEN"

# ✅ Test with a different user (THIS WORKS)
TOKEN2=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{"name": "anotheruser"}' \
  http://localhost:3000/session | \
  grep -o '"identifier":"[^"]*"' | cut -d'"' -f4)

echo "Second user token: $TOKEN2"

# ✅ Verify authentication works by testing any protected endpoint
# (This will return 501 Not Implemented, but not 401 Unauthorized)
curl -s -X GET \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/conversations

# Expected: 501 Not Implemented (handler not implemented yet)
# If you get 401, there's an auth problem
```

### Authentication Details

- **User Creation**: Users are automatically created on first login
- **Token Format**: UUID v4 strings stored in database
- **Token Usage**: Include in `Authorization: Bearer <token>` header
- **Session Storage**: Tokens are stored in SQLite `user_sessions` table

## User Management (🔧 In Development)

**Status**: Database operations ready, API handlers being implemented

```bash
# 🔧 Update your username (Handler in development)
curl -s -X PUT \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name": "newusername"}' \
  http://localhost:3000/users/me/username
# Expected: 501 Not Implemented (for now)

# 🔧 Upload profile photo (Handler in development)
curl -s -X PUT \
  -H "Authorization: Bearer $TOKEN" \
  -F "photo=@test_user.png" \
  http://localhost:3000/users/me/photo
# Expected: 501 Not Implemented (for now)

# 🔧 Search for users (Handler in development)
curl -s -X GET \
  -H "Authorization: Bearer $TOKEN" \
  "http://localhost:3000/users?query=testuser" \
# Expected: 501 Not Implemented (for now)

# 🔧 Start a conversation with another user (Handler in development)
curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{}' \
  http://localhost:3000/users/USER_ID/conversations
# Expected: 501 Not Implemented (for now)
```

**Note**: All user management endpoints are registered and will authenticate properly, but return `501 Not Implemented` until handlers are completed.

## Conversations (🔧 In Development)

**Status**: Database schema ready, API handlers being implemented

```bash
# 🔧 Get your conversations list (Handler in development)
curl -s -X GET \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/conversations
# Expected: 501 Not Implemented (for now)

# 🔧 Get messages in a specific conversation (Handler in development)
curl -s -X GET \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/conversations/CONVERSATION_ID
# Expected: 501 Not Implemented (for now)
```

## Messages (🔧 In Development)

**Status**: Database schema ready, API handlers being implemented

```bash
# 🔧 Send a text message (Handler in development)
curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"content": "Hello there!"}' \
  http://localhost:3000/conversations/CONVERSATION_ID/messages
# Expected: 501 Not Implemented (for now)

# 🔧 Send a photo message (Handler in development)
curl -s -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -F "photo=@test_user.png" \
  http://localhost:3000/conversations/CONVERSATION_ID/messages
# Expected: 501 Not Implemented (for now)

# 🔧 Reply to a message (Handler in development)
curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"content": "Great point!", "replyTo": "MESSAGE_ID"}' \
  http://localhost:3000/conversations/CONVERSATION_ID/messages
# Expected: 501 Not Implemented (for now)

# 🔧 Add a reaction to a message (Handler in development)
curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"emoticon": "👍"}' \
  http://localhost:3000/messages/MESSAGE_ID/comments
# Expected: 501 Not Implemented (for now)

# 🔧 Forward a message (Handler in development)
curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"conversationId": "TARGET_CONVERSATION_ID"}' \
  http://localhost:3000/messages/MESSAGE_ID/forward
# Expected: 501 Not Implemented (for now)

# 🔧 Delete a message (Handler in development)
curl -s -X DELETE \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/messages/MESSAGE_ID
# Expected: 501 Not Implemented (for now)
```

**Database Ready**: Messages, reactions, and forwarding are fully designed in the SQLite schema.

## Groups (🔧 In Development)

**Status**: Database schema ready, API handlers being implemented

```bash
# 🔧 Create a new group (Handler in development)
curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name": "Project Team", "members": ["USER_ID_1", "USER_ID_2"]}' \
  http://localhost:3000/groups
# Expected: 501 Not Implemented (for now)

# 🔧 Add user to group (Handler in development)
curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"userId": "USER_ID"}' \
  http://localhost:3000/groups/GROUP_ID/members
# Expected: 501 Not Implemented (for now)

# 🔧 Leave a group (Handler in development)
curl -s -X DELETE \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/groups/GROUP_ID/members/me
# Expected: 501 Not Implemented (for now)

# 🔧 Update group name (Handler in development)
curl -s -X PUT \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name": "New Group Name"}' \
  http://localhost:3000/groups/GROUP_ID/name
# Expected: 501 Not Implemented (for now)

# 🔧 Update group photo (Handler in development)
curl -s -X PUT \
  -H "Authorization: Bearer $TOKEN" \
  -F "photo=@group_photo.png" \
  http://localhost:3000/groups/GROUP_ID/photo
# Expected: 501 Not Implemented (for now)
```

**Database Ready**: Group conversations, membership, and management are fully implemented in the database schema.

## Testing Tips

### Understanding HTTP Status Codes

**Current Implementation Status:**

```
✅ Working Endpoints:
200 OK: Request succeeded (login works)
201 Created: Resource created (user creation works)

🔧 Development Endpoints:
501 Not Implemented: Handler not yet implemented (most endpoints)

❌ Error Cases:
400 Bad Request: Invalid request format
401 Unauthorized: Missing/invalid Bearer token
403 Forbidden: Not allowed to access the resource
404 Not Found: Resource not found
500 Internal Server Error: Server-side error
```

````

### Viewing HTTP Status Codes

```bash
# Add -w "%{http_code}" to see the HTTP status code
curl -s -w "%{http_code}" \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/users/me \
  -o /dev/null
````

### Getting More Details

```bash
# Add -v to see detailed request and response information
curl -v -X GET \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/users/me
```

### Complete Test Flow Example

Here's a complete flow for testing **currently working** functionality:

```bash
# 1. Start the WASAText server
go build ./cmd/webapi
./webapi &
SERVER_PID=$!

# 2. ✅ Health check (should work)
curl -s http://localhost:3000/liveness
echo "✅ Health check passed"

# 3. ✅ Create first user (should work)
TOKEN1=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{"name": "alice"}' \
  http://localhost:3000/session | \
  grep -o '"identifier":"[^"]*"' | cut -d'"' -f4)
echo "✅ Alice token: $TOKEN1"

# 4. ✅ Create second user (should work)
TOKEN2=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{"name": "bob"}' \
  http://localhost:3000/session | \
  grep -o '"identifier":"[^"]*"' | cut -d'"' -f4)
echo "✅ Bob token: $TOKEN2"

# 5. 🔧 Test authentication on protected endpoint (should return 501, not 401)
echo "Testing protected endpoint (should return 501):"
curl -s -w "Status: %{http_code}\n" -o /dev/null \
  -H "Authorization: Bearer $TOKEN1" \
  http://localhost:3000/conversations

# 6. 🔧 Test missing auth (should return 401)
echo "Testing missing auth (should return 401):"
curl -s -w "Status: %{http_code}\n" -o /dev/null \
  http://localhost:3000/conversations

# 7. Stop server when done
kill $SERVER_PID
wait $SERVER_PID 2>/dev/null
echo "✅ Server stopped"
```

### Note About IDs

For endpoints that require IDs (conversationId, messageId, etc.), **these are not yet implemented**:

- 🔧 Conversation IDs: Will be available once conversation handlers are implemented
- 🔧 Message IDs: Will be available once message handlers are implemented
- 🔧 User IDs: Database has them, but search endpoint needs to be implemented

**Current ID Format**: All IDs are UUID v4 strings (e.g., `"550e8400-e29b-41d4-a716-446655440000"`)

Once endpoints are implemented, you can extract IDs like this:

```bash
# This will work once conversation endpoint is implemented:
# CONVERSATION_ID=$(curl -s -X GET \
#   -H "Authorization: Bearer $TOKEN" \
#   http://localhost:3000/conversations | jq -r '.[0].id')
```

This assumes you have `jq` installed for parsing JSON responses.

### Testing With Multiple Users

To test interaction between multiple users, you'll need to create and manage multiple sessions:

````bash
# Create the first user and store their token
TOKEN_USER1=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{"name": "user1"}' \
  http://localhost:3000/session | \
  grep -o '"identifier":"[^"]*"' | cut -d'"' -f4)

# Create a second user and store their token
TOKEN_USER2=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{"name": "user2"}' \
  http://localhost:3000/session | \
  grep -o '"identifier":"[^"]*"' | cut -d'"' -f4)

# Get user1's ID
USER1_ID=$(curl -s -X GET \
  -H "Authorization: Bearer $TOKEN_USER1" \
  http://localhost:3000/users/me | \
  grep -o '"id":"[^"]*"' | cut -d'"' -f4)

# Get user2's ID
USER2_ID=$(curl -s -X GET \
  -H "Authorization: Bearer $TOKEN_USER2" \
  http://localhost:3000/users/me | \
  grep -o '"id":"[^"]*"' | cut -d'"' -f4)

### Testing With Multiple Users

**Note**: Multi-user testing will work once the API handlers are implemented. For now, you can test authentication with multiple users:

```bash
# ✅ Create multiple users (authentication works)
TOKEN_USER1=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{"name": "alice"}' \
  http://localhost:3000/session | \
  grep -o '"identifier":"[^"]*"' | cut -d'"' -f4)

TOKEN_USER2=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{"name": "bob"}' \
  http://localhost:3000/session | \
  grep -o '"identifier":"[^"]*"' | cut -d'"' -f4)

TOKEN_USER3=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{"name": "charlie"}' \
  http://localhost:3000/session | \
  grep -o '"identifier":"[^"]*"' | cut -d'"' -f4)

echo "✅ User tokens created:"
echo "Alice: $TOKEN_USER1"
echo "Bob: $TOKEN_USER2"
echo "Charlie: $TOKEN_USER3"

# ✅ Verify each user can authenticate
for TOKEN in $TOKEN_USER1 $TOKEN_USER2 $TOKEN_USER3; do
  STATUS=$(curl -s -w "%{http_code}" -o /dev/null \
    -H "Authorization: Bearer $TOKEN" \
    http://localhost:3000/conversations)
  echo "Auth test with token: $STATUS (should be 501, not 401)"
done
````

#### Example Multi-User Test Flow (🔧 Future)

Once the API handlers are implemented, this will work:

```bash
# 🔧 Future: User1 creates a conversation with User2
# CONVERSATION_ID=$(curl -s -X POST \
#   -H "Content-Type: application/json" \
#   -H "Authorization: Bearer $TOKEN_USER1" \
#   -d "{}" \
#   http://localhost:3000/users/$USER2_ID/conversations | \
#   jq -r '.id')

# 🔧 Future: User1 sends a message
# curl -s -X POST \
#   -H "Content-Type: application/json" \
#   -H "Authorization: Bearer $TOKEN_USER1" \
#   -d '{"content": "Hello from Alice!"}' \
#   http://localhost:3000/conversations/$CONVERSATION_ID/messages

# 🔧 Future: User2 reads the conversation
# curl -s -X GET \
#   -H "Authorization: Bearer $TOKEN_USER2" \
#   http://localhost:3000/conversations/$CONVERSATION_ID
```

#### Testing Group Interactions (🔧 Future)

Once group handlers are implemented:

```bash
# 🔧 Future: Create a group with multiple users
# GROUP_ID=$(curl -s -X POST \
#   -H "Content-Type: application/json" \
#   -H "Authorization: Bearer $TOKEN_USER1" \
#   -d "{\"name\": \"Test Group\", \"members\": [\"$USER2_ID\", \"$USER3_ID\"]}" \
#   http://localhost:3000/groups | \
#   jq -r '.id')

# 🔧 Future: Send group messages
# curl -s -X POST \
#   -H "Content-Type: application/json" \
#   -H "Authorization: Bearer $TOKEN_USER1" \
#   -d '{"content": "Welcome to the group!"}' \
#   http://localhost:3000/conversations/$GROUP_ID/messages
```

---

## 🚀 **Development Status Summary**

### ✅ **Currently Working:**

- Server startup and health check
- User registration/login
- Authentication with Bearer tokens
- Database initialization

### 🔧 **In Development:**

- All API handlers (return 501 Not Implemented)
- User management operations
- Conversation and message handling
- Group management

### 📋 **Testing Strategy:**

1. **Phase 1**: Test authentication (working now)
2. **Phase 2**: Test individual API handlers as they're implemented
3. **Phase 3**: Test end-to-end workflows with multiple users

**Ready for**: Authentication testing and database verification  
**Next milestone**: Complete user management handlers
USER3_ID=$(curl -s -X GET \
 -H "Authorization: Bearer $TOKEN_USER3" \
 http://localhost:3000/users/me | \
 grep -o '"id":"[^"]\*"' | cut -d'"' -f4)

# 5. User1 adds User3 to the group

curl -s -X POST \
 -H "Content-Type: application/json" \
 -H "Authorization: Bearer $TOKEN_USER1" \
  -d "{\"userId\": \"$USER3_ID\"}" \
 http://localhost:3000/groups/$GROUP_ID/members

```

Using this approach, you can simulate real-world interactions between users and test various scenarios involving multiple participants.
```
