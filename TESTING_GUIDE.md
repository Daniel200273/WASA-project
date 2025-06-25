# WASA Project Testing Guide

This guide contains instructions and curl commands for testing the WASA messaging application API.

## Table of Contents

- [Getting Started](#getting-started)
- [Understanding curl Flags](#understanding-curl-flags)
- [Authentication](#authentication)
- [User Management](#user-management)
- [Conversations](#conversations)
- [Messages](#messages)
- [Groups](#groups)
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

### Starting the Server

Before running any API tests, you'll need to start the server:

```bash
# Build the application
go build ./cmd/webapi

# Start the server
./webapi &

# Store the process ID for later use
SERVER_PID=$!

# Verify the server is running
curl -s http://localhost:3000/liveness
```

### Checking if a Server is Running

Before starting a new server, you may want to check if one is already running:

```bash
# Check if the port is in use
lsof -i :3000

# Check for running webapi processes
ps aux | grep webapi

# Quick check if server is responding
curl -s http://localhost:3000/liveness && echo "✅ Server is running" || echo "❌ Server is not running"

# Kill any existing server if needed
pkill webapi
```

### Ending Your Test Session

When you're finished testing, make sure to properly shut down the server:

```bash
# Kill the server process
kill $SERVER_PID

# Or if you didn't save the PID:
pkill webapi

# Wait for the process to terminate
wait $SERVER_PID 2>/dev/null
```

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

## Authentication

```bash
# Create a new user session / login
curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{"name": "testuser"}' \
  http://localhost:3000/session

# Save the token for later use
TOKEN=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{"name": "testuser"}' \
  http://localhost:3000/session | \
  grep -o '"identifier":"[^"]*"' | cut -d'"' -f4)
```

## User Management

```bash
# Update your username
curl -s -X PUT \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name": "newusername"}' \
  http://localhost:3000/users/me/username

# Upload profile photo
curl -s -X PUT \
  -H "Authorization: Bearer $TOKEN" \
  -F "photo=@test_user.png" \
  http://localhost:3000/users/me/photo

# Get your user profile
curl -s -X GET \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/users/me

# Search for users
curl -s -X GET \
  -H "Authorization: Bearer $TOKEN" \
  "http://localhost:3000/users?query=username"
```

## Conversations

```bash
# Get your conversations
curl -s -X GET \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/conversations

# Get messages in a specific conversation
curl -s -X GET \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/conversations/{conversationId}

# Start a new direct conversation
curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"recipientId": "USER_ID"}' \
  http://localhost:3000/conversations
```

## Messages

```bash
# Send a text message
curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"content": "Hello there!"}' \
  http://localhost:3000/conversations/{conversationId}/messages

# Send a photo message
curl -s -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -F "photo=@test_user.png" \
  http://localhost:3000/conversations/{conversationId}/messages

# Add a reaction to a message
curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"emoji": "👍"}' \
  http://localhost:3000/conversations/{conversationId}/messages/{messageId}/reactions

# Delete a message
curl -s -X DELETE \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/conversations/{conversationId}/messages/{messageId}
```

## Groups

```bash
# Create a new group
curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name": "My Group", "members": ["USER_ID1", "USER_ID2"]}' \
  http://localhost:3000/groups

# Update a group's name
curl -s -X PUT \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name": "Updated Group Name"}' \
  http://localhost:3000/groups/{groupId}

# Upload group photo
curl -s -X PUT \
  -H "Authorization: Bearer $TOKEN" \
  -F "photo=@group_photo.png" \
  http://localhost:3000/groups/{groupId}/photo

# Add a member to a group
curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"userId": "USER_ID"}' \
  http://localhost:3000/groups/{groupId}/members

# Leave a group
curl -s -X DELETE \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/groups/{groupId}/members/me
```

## Testing Tips

### Understanding HTTP Status Codes

```
200 OK: Request succeeded
201 Created: Resource created
204 No Content: Request succeeded with no content returned
400 Bad Request: Invalid request format
401 Unauthorized: Authentication required
403 Forbidden: Not allowed to access the resource
404 Not Found: Resource not found
409 Conflict: Resource conflict (e.g., duplicate username)
500 Internal Server Error: Server-side error
```

### Viewing HTTP Status Codes

```bash
# Add -w "%{http_code}" to see the HTTP status code
curl -s -w "%{http_code}" \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/users/me \
  -o /dev/null
```

### Getting More Details

```bash
# Add -v to see detailed request and response information
curl -v -X GET \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/users/me
```

### Complete Test Flow Example

Here's a complete flow for testing basic functionality:

```bash
# 1. Start server
./webapi &
SERVER_PID=$!

# 2. Get authentication token
TOKEN=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{"name": "testuser"}' \
  http://localhost:3000/session | \
  grep -o '"identifier":"[^"]*"' | cut -d'"' -f4)
echo "Token: $TOKEN"

# 3. Update username
curl -s -X PUT \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name": "updatedname"}' \
  http://localhost:3000/users/me/username

# 4. Upload profile photo
curl -s -X PUT \
  -H "Authorization: Bearer $TOKEN" \
  -F "photo=@test_user.png" \
  http://localhost:3000/users/me/photo

# 5. Stop server when done
kill $SERVER_PID
wait $SERVER_PID 2>/dev/null
```

### Note About IDs

For endpoints that require IDs (conversationId, messageId, etc.), you'll need to:

1. First create the resource or retrieve it from a listing endpoint
2. Extract the ID from the response
3. Use the ID in subsequent requests

For example, to get a conversation ID:

```bash
CONVERSATION_ID=$(curl -s -X GET \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/conversations | jq -r '.[0].id')
```

This assumes you have `jq` installed for parsing JSON responses.

### Testing With Multiple Users

To test interaction between multiple users, you'll need to create and manage multiple sessions:

```bash
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

echo "User1 token: $TOKEN_USER1"
echo "User1 ID: $USER1_ID"
echo "User2 token: $TOKEN_USER2"
echo "User2 ID: $USER2_ID"
```

#### Example Multi-User Test Flow

Here's an example of how to test a conversation between two users:

```bash
# 1. User1 creates a conversation with User2
CONVERSATION_ID=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN_USER1" \
  -d "{\"recipientId\": \"$USER2_ID\"}" \
  http://localhost:3000/conversations | \
  grep -o '"id":"[^"]*"' | cut -d'"' -f4)
echo "Conversation ID: $CONVERSATION_ID"

# 2. User1 sends a message
curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN_USER1" \
  -d '{"content": "Hello from user1!"}' \
  http://localhost:3000/conversations/$CONVERSATION_ID/messages

# 3. User2 reads the conversation
curl -s -X GET \
  -H "Authorization: Bearer $TOKEN_USER2" \
  http://localhost:3000/conversations/$CONVERSATION_ID

# 4. User2 sends a response
curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN_USER2" \
  -d '{"content": "Hi user1, received your message!"}' \
  http://localhost:3000/conversations/$CONVERSATION_ID/messages

# 5. User1 reads the response
curl -s -X GET \
  -H "Authorization: Bearer $TOKEN_USER1" \
  http://localhost:3000/conversations/$CONVERSATION_ID
```

#### Testing Group Interactions

For testing group interactions with multiple users:

```bash
# 1. User1 creates a group with User2
GROUP_ID=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN_USER1" \
  -d "{\"name\": \"Test Group\", \"members\": [\"$USER2_ID\"]}" \
  http://localhost:3000/groups | \
  grep -o '"id":"[^"]*"' | cut -d'"' -f4)
echo "Group ID: $GROUP_ID"

# 2. User1 sends a message to the group
curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN_USER1" \
  -d '{"content": "Welcome to the group!"}' \
  http://localhost:3000/conversations/$GROUP_ID/messages

# 3. User2 reads the group conversation
curl -s -X GET \
  -H "Authorization: Bearer $TOKEN_USER2" \
  http://localhost:3000/conversations/$GROUP_ID

# 4. Create a third user to add to the group
TOKEN_USER3=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{"name": "user3"}' \
  http://localhost:3000/session | \
  grep -o '"identifier":"[^"]*"' | cut -d'"' -f4)
USER3_ID=$(curl -s -X GET \
  -H "Authorization: Bearer $TOKEN_USER3" \
  http://localhost:3000/users/me | \
  grep -o '"id":"[^"]*"' | cut -d'"' -f4)

# 5. User1 adds User3 to the group
curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN_USER1" \
  -d "{\"userId\": \"$USER3_ID\"}" \
  http://localhost:3000/groups/$GROUP_ID/members
```

Using this approach, you can simulate real-world interactions between users and test various scenarios involving multiple participants.
