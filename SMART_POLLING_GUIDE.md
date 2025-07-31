# Smart Polling Implementation Guide

## Overview

I've implemented a smart polling system for your WASA messaging application that provides near real-time messaging without WebSockets. The system adapts its polling frequency based on user activity and visibility.

## Features Implemented

### Frontend (ChatView.vue)

1. **Adaptive Polling Frequency**

   - Active chat: 3 seconds
   - Background tab: 60 seconds
   - Inactive user: 30 seconds
   - Default: 10 seconds

2. **Activity Tracking**

   - Monitors mouse, keyboard, scroll, and touch events
   - Automatically adjusts polling based on user interaction
   - Detects tab visibility changes

3. **Incremental Loading**

   - Only fetches new messages since last update
   - Prevents duplicate message loading
   - Efficient bandwidth usage

4. **Browser Notifications**

   - Shows notifications for new messages when tab is hidden
   - Requests permission on component mount
   - Prevents duplicate notifications

5. **Smart State Management**
   - Tracks last update timestamps
   - Optimistic UI updates
   - Graceful error handling

### Backend Optimizations

1. **New API Endpoint Parameters**

   - `GET /users/{userId}/conversations/{conversationId}?after={timestamp}&limit={number}`
   - Supports incremental message loading
   - Configurable limits (max 500 messages per request)

2. **Database Method**
   - `GetConversationMessagesAfter()` - fetches messages after a specific timestamp
   - Proper indexing for performance
   - Time-based filtering

## How It Works

### Polling Logic

```javascript
// Determines optimal polling interval based on user state
optimalPollingInterval() {
  if (!this.isTabVisible) return 60000;     // 1 minute when hidden
  if (!this.isUserActive) return 30000;     // 30s when inactive
  if (this.selectedConversationId) return 3000; // 3s in active chat
  return 10000; // 10s default
}
```

### Activity Detection

The system tracks user activity through:

- Mouse movements and clicks
- Keyboard input
- Scroll events
- Touch interactions
- Tab visibility changes

### Incremental Loading

Instead of loading all messages every time:

```javascript
// Get timestamp of last message
const lastMessageTime =
  this.messages.length > 0
    ? this.messages[this.messages.length - 1].timestamp
    : this.conversationReadAt;

// Request only newer messages
const response = await axios.get(
  `/users/${userId}/conversations/${conversationId}`,
  { params: { after: lastMessageTime, limit: 50 } }
);
```

## Performance Benefits

1. **Reduced Database Load**

   - Time-based queries instead of full table scans
   - Indexed queries for better performance
   - Configurable limits prevent large responses

2. **Bandwidth Efficiency**

   - Only transmits new data
   - Smaller payload sizes
   - Adaptive frequency reduces unnecessary requests

3. **Better User Experience**
   - Near real-time updates (3s in active chats)
   - Background notifications
   - Smooth scrolling and UI updates
   - No page refresh needed

## Testing the Implementation

### 1. Basic Functionality Test

```bash
# Start the backend
cd /path/to/WASA-project
go run ./cmd/webapi

# Start the frontend
cd webui
npm run dev
```

### 2. Test Scenarios

1. **Active Chat Testing**

   - Open two browser windows with different users
   - Send messages and verify 3-second updates
   - Check browser console for polling logs

2. **Background Tab Testing**

   - Send message from one window
   - Switch to another tab
   - Verify notification appears after ~60 seconds

3. **Activity Detection Testing**
   - Stop interacting with the page
   - Wait 2 minutes and check console logs
   - Verify polling slows to 30 seconds

### 3. Monitor Performance

Check browser console for:

```
Polling interval changed: 10000ms → 3000ms
```

Monitor network tab for:

- Reduced request sizes when using `?after=` parameter
- Adaptive request frequency

## Database Indexes

Add these indexes for optimal performance:

```sql
-- Critical for message timestamp queries
CREATE INDEX idx_messages_conversation_timestamp ON messages(conversation_id, created_at);

-- For conversation participant queries
CREATE INDEX idx_conversation_participants_user ON conversation_participants(user_id);

-- For session validation
CREATE INDEX idx_user_sessions_token ON user_sessions(token);
```

## Configuration Options

You can adjust polling behavior by modifying these values in ChatView.vue:

```javascript
// Polling intervals (milliseconds)
const ACTIVE_CHAT_INTERVAL = 3000;
const INACTIVE_TAB_INTERVAL = 60000;
const INACTIVE_USER_INTERVAL = 30000;
const DEFAULT_INTERVAL = 10000;

// Activity timeout (milliseconds)
const INACTIVITY_TIMEOUT = 120000; // 2 minutes

// Message loading limit
const MESSAGE_LIMIT = 50;
```

## Future Enhancements

1. **Connection Quality Adaptation**

   - Slower polling on poor connections
   - Exponential backoff on errors

2. **Typing Indicators**

   - Send typing events via separate endpoint
   - Show "user is typing" notifications

3. **Read Receipts**

   - Real-time read status updates
   - Mark conversations as read automatically

4. **Push Notifications** (Progressive Web App)
   - Service worker integration
   - Offline message queuing

## Troubleshooting

### Common Issues

1. **High CPU Usage**

   - Check polling intervals in console
   - Verify activity detection is working
   - Consider increasing intervals

2. **Messages Not Updating**

   - Check network connectivity
   - Verify API responses in Network tab
   - Check for JavaScript errors

3. **Duplicate Messages**
   - Verify timestamp formatting
   - Check message deduplication logic
   - Monitor `existingIds` Set in console

### Debug Logging

Enable debug logging by adding to ChatView.vue:

```javascript
console.log(
  `Polling: active=${this.isUserActive}, visible=${this.isTabVisible}, interval=${this.optimalPollingInterval}ms`
);
```

## Conclusion

This smart polling implementation provides an excellent balance between real-time responsiveness and system efficiency. It automatically adapts to user behavior while maintaining good performance characteristics for both frontend and backend systems.
