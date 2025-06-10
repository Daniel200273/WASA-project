# WASAText API Summary

## 🎯 **Overview**

**WASAText** is a desktop messaging application API that enables users to connect and communicate through both direct messages and group chats. The API provides a comprehensive set of endpoints for user management, conversations, messaging, and group functionality.

**Base URL:** `http://localhost:3000`

---

## 🏗️ **System Architecture**

### **Core Components:**

```
┌─────────────────────────────────────────────────────────────┐
│                    WASAText Messaging API                   │
│                  (http://localhost:3000)                   │
└─────────────────────────────────────────────────────────────┘
                                │
                                ▼
        ┌─────────────────────────────────────────────────┐
        │              Authentication                     │
        │  • Simple login with username                   │
        │  • Auto-creates users if not exist             │
        │  • Bearer token authentication                  │
        └─────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────┬─────────────────┬─────────────────┬─────────────────┐
│   User Mgmt     │  Conversations  │    Messages     │     Groups      │
│                 │                 │                 │                 │
│ • Profile setup │ • List chats    │ • Text msgs     │ • Create groups │
│ • Username      │ • Direct msgs   │ • Photo msgs    │ • Add members   │
│ • Profile pics  │ • Group chats   │ • Forward       │ • Leave groups  │
│ • Search users  │ • Unread counts │ • Delete        │ • Update info   │
│                 │                 │ • Reactions     │                 │
└─────────────────┴─────────────────┴─────────────────┴─────────────────┘
```

---

## 👤 **User Journey Flow**

### **1. Authentication & Setup**
```
User opens app → Enter username → System creates/finds user → Returns auth token
                                                                      │
                                                                      ▼
                                                            User can now access app
```

### **2. Profile Management**
```
Update username ← → Upload profile photo ← → Search for other users
```

### **3. Messaging Flow**
```
View conversations list → Select conversation → View messages → Send message/photo
         │                        │                   │              │
         │                        │                   │              ▼
         │                        │                   │        React with emoji
         │                        │                   │              │
         │                        │                   ▼              ▼
         │                        │            Forward message → Delete message
         │                        │
         ▼                        ▼
Create new group → Add members → Start messaging
```

---

## 💬 **Conversation Types**

### **Direct Messages**
- **1:1 conversations** between two users
- Display contact's name and profile photo
- Private messaging

### **Group Chats**
- **Multi-user conversations** (1-100 members)
- Custom group names and photos
- Member management (add/remove)
- Group administration

---

## 📱 **Key Features**

### **Message Types:**
```
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   Text Message  │  │  Photo Message  │  │ Forwarded Msg   │
│                 │  │                 │  │                 │
│ "Hello there!"  │  │ [📷 image.jpg]  │  │ ↪️ "Fwd: Hello" │
│                 │  │                 │  │                 │
│ 👍 ❤️ 😂       │  │ 👍 ❤️ 😂       │  │ 👍 ❤️ 😂       │
└─────────────────┘  └─────────────────┘  └─────────────────┘
```

### **Message Features:**
- ✅ **Delivery Status**: sent → delivered → read
- 💬 **Reactions**: Emoji comments on messages  
- ↪️ **Replies**: Reference previous messages
- ⏭️ **Forwarding**: Share messages across conversations
- 🗑️ **Deletion**: Remove own messages

---

## 🎨 **User Interface Concept**

### **Main Screen Layout:**
```
┌─────────────────────────────────────────────────────────────┐
│  WASAText                                   [User Profile]   │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Conversations List          │        Active Chat           │
│  ┌─────────────────────┐    │  ┌─────────────────────────┐  │
│  │ 👤 Maria            │    │  │ 👥 Project Team        │  │
│  │ "Hello there!"      │ 2  │  │                         │  │
│  │ 2 min ago          │    │  │ Maria: Hello everyone!   │  │
│  └─────────────────────┘    │  │ John: Hi Maria! 👋      │  │
│  ┌─────────────────────┐    │  │ You: How's the project? │  │
│  │ 👥 Project Team     │    │  │                         │  │
│  │ "How's the project?"│    │  │ [Type message here...] │  │
│  │ 5 min ago          │    │  └─────────────────────────┘  │
│  └─────────────────────┘    │                             │
│                             │                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 🔐 **Security Model**

```
Bearer Token Authentication
         │
         ▼
┌─────────────────────────────────────┐
│        Access Control               │
│                                     │
│ ✅ Users can only:                  │
│   • Access their own conversations  │
│   • Delete their own messages       │
│   • Remove their own reactions      │
│   • Leave groups they're in         │
│                                     │
│ ❌ Users cannot:                    │
│   • Access others' private data     │
│   • Delete others' messages         │
│   • Add to groups they're not in    │
└─────────────────────────────────────┘
```

---

## 📊 **Data Structures**

### **Core Entities:**
```
User                    Conversation              Message
├── id                 ├── id                   ├── id
├── username           ├── type (direct/group)  ├── content/photoUrl
├── photoUrl           ├── name                 ├── senderId
                       ├── photoUrl             ├── timestamp
Group                  ├── members              ├── status
├── id                 ├── lastMessage          ├── replyTo
├── name               └── unreadCount          ├── forwarded
├── photoUrl                                    └── comments[]
├── members[]          Comment
├── createdBy          ├── id
└── createdAt          ├── userId
                       ├── emoticon
                       └── timestamp
```

---

## 🚀 **API Endpoints Reference**

### **Authentication**
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/session` | Login/Register user |

### **User Management**
| Method | Endpoint | Description |
|--------|----------|-------------|
| `PUT` | `/users/me/username` | Update current user's username |
| `PUT` | `/users/me/photo` | Upload profile photo |
| `GET` | `/users` | Search for users by username |

### **Conversations**
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/conversations` | Get user's conversations list |
| `GET` | `/conversations/{conversationId}` | Get conversation details with messages |

### **Messages**
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/conversations/{conversationId}/messages` | Send text or photo message |
| `DELETE` | `/messages/{messageId}` | Delete own message |
| `POST` | `/messages/{messageId}/forward` | Forward message to another conversation |

### **Message Reactions**
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/messages/{messageId}/comments` | Add emoji reaction to message |
| `DELETE` | `/messages/{messageId}/comments/{commentId}` | Remove own reaction |

### **Groups**
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/groups` | Create new group |
| `POST` | `/groups/{groupId}/members` | Add user to group |
| `DELETE` | `/groups/{groupId}/members/me` | Leave group |
| `PUT` | `/groups/{groupId}/name` | Update group name |
| `PUT` | `/groups/{groupId}/photo` | Update group photo |

---

## 📋 **Key Constraints & Limits**

### **Data Limits:**
- **Usernames**: 3-16 characters, alphanumeric + underscore/dash
- **Messages**: Max 1000 characters
- **Photos**: Max 10MB file size
- **Group size**: 1-100 members
- **Conversations per user**: Max 500
- **Messages per conversation**: Max 500 returned
- **Reactions per message**: Max 50

### **ID Patterns:**
- All IDs follow pattern: `^[a-zA-Z0-9_-]+$`
- Max 36 characters length

---

## 🎯 **Business Logic**

### **User Creation:**
- Users are auto-created on first login
- Username uniqueness enforced
- Bearer token returned for authentication

### **Conversation Management:**
- Direct conversations created implicitly when messaging
- Group conversations require explicit creation
- Conversations sorted by latest message timestamp

### **Message Delivery:**
- Status tracking: sent → delivered → read
- Users can only delete their own messages
- Photo messages stored with URL references

### **Group Administration:**
- Any member can add new users
- Users can leave groups voluntarily
- Group creators have no special privileges (democratic model)

---

## 🔧 **Technical Specifications**

### **Authentication:**
- **Type**: Bearer Token
- **Header**: `Authorization: Bearer <token>`
- **Token source**: Login endpoint response

### **Content Types:**
- **JSON**: `application/json` (most endpoints)
- **Multipart**: `multipart/form-data` (photo uploads)

### **HTTP Status Codes:**
- **200**: Success (updates)
- **201**: Success (creation)
- **204**: Success (deletion)
- **400**: Bad request
- **401**: Unauthorized
- **403**: Forbidden
- **404**: Not found
- **409**: Conflict (duplicate username)

---

## 🎨 **Response Examples**

### **Login Response:**
```json
{
  "identifier": "abcdef012345"
}
```

### **Conversations List:**
```json
{
  "conversations": [
    {
      "id": "conv123",
      "type": "direct",
      "name": "Maria",
      "photoUrl": "/photos/user123.jpg",
      "lastMessage": {
        "id": "msg789",
        "content": "Hello there!",
        "timestamp": "2023-06-15T14:30:00Z",
        "senderUsername": "Maria",
        "hasPhoto": false
      },
      "unreadCount": 2
    }
  ]
}
```

### **Message with Reactions:**
```json
{
  "id": "msg123",
  "senderId": "user123",
  "senderUsername": "Maria",
  "content": "Hello everyone!",
  "timestamp": "2023-06-15T14:30:00Z",
  "status": "read",
  "comments": [
    {
      "id": "comment123",
      "userId": "user456",
      "username": "John",
      "emoticon": "👍",
      "timestamp": "2023-06-15T17:30:00Z"
    }
  ]
}
```

---

This API provides a complete foundation for building a modern desktop messaging application with all the essential features users expect from contemporary chat applications.
