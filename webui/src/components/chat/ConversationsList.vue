<template>
  <div class="conversations-list">
    <div class="conversations-header">
      <h5>Conversations</h5>
      <button class="btn btn-primary btn-sm" @click="$emit('new-chat')">
        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#plus"/></svg>
      </button>
    </div>
    
    <div class="conversations-scroll">
      <div v-for="conversation in conversations" 
           :key="conversation.id" 
           class="conversation-item" 
           :class="{ active: selectedConversationId === conversation.id }"
           @click="$emit('select-conversation', conversation)">
        
        <div class="conversation-avatar">
          <img :src="conversation.avatar || '/default-avatar.svg'" :alt="conversation.name">
        </div>
        
        <div class="conversation-info">
          <div class="conversation-name">{{ conversation.name }}</div>
          <div class="conversation-preview">{{ conversation.lastMessage || 'No messages yet' }}</div>
        </div>
        
        <div class="conversation-meta">
          <div class="conversation-time">{{ formatTime(conversation.lastMessageTime) }}</div>
          <div v-if="conversation.unreadCount" class="unread-badge">{{ conversation.unreadCount }}</div>
        </div>
      </div>
      
      <!-- Empty State -->
      <div v-if="!conversations.length" class="empty-conversations">
        <p>No conversations yet</p>
        <button class="btn btn-outline-primary btn-sm" @click="$emit('new-chat')">Start a conversation</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ConversationsList',
  props: {
    conversations: {
      type: Array,
      default: () => []
    },
    selectedConversationId: {
      type: String,
      default: null
    }
  },
  emits: ['select-conversation', 'new-chat'],
  methods: {
    formatTime(timestamp) {
      if (!timestamp) return '';
      
      const date = new Date(timestamp);
      const now = new Date();
      const diffMs = now - date;
      const diffMins = Math.floor(diffMs / 60000);
      const diffHours = Math.floor(diffMs / 3600000);
      const diffDays = Math.floor(diffMs / 86400000);
      
      if (diffMins < 1) return 'now';
      if (diffMins < 60) return `${diffMins}m`;
      if (diffHours < 24) return `${diffHours}h`;
      if (diffDays < 7) return `${diffDays}d`;
      
      return date.toLocaleDateString();
    }
  }
}
</script>

<style scoped>
.conversations-list {
  width: 320px;
  height: 100%;
  border-right: 1px solid #dee2e6;
  background: #f8f9fa;
  display: flex;
  flex-direction: column;
}

.conversations-header {
  padding: 1rem;
  border-bottom: 1px solid #dee2e6;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: white;
}

.conversations-header h5 {
  margin: 0;
}

.conversations-scroll {
  flex: 1;
  overflow-y: auto;
}

.conversation-item {
  padding: 1rem;
  border-bottom: 1px solid #eee;
  cursor: pointer;
  display: flex;
  align-items: center;
  transition: background-color 0.2s;
  background: white;
}

.conversation-item:hover {
  background-color: #f8f9fa;
}

.conversation-item.active {
  background-color: #007bff;
  color: white;
}

.conversation-avatar img {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
  margin-right: 0.75rem;
}

.conversation-info {
  flex: 1;
  min-width: 0;
}

.conversation-name {
  font-weight: 500;
  margin-bottom: 0.25rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.conversation-preview {
  font-size: 0.875rem;
  opacity: 0.8;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.conversation-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 0.25rem;
}

.conversation-time {
  font-size: 0.75rem;
  opacity: 0.7;
}

.unread-badge {
  background: #dc3545;
  color: white;
  border-radius: 10px;
  padding: 0.125rem 0.5rem;
  font-size: 0.75rem;
  font-weight: 500;
  min-width: 20px;
  text-align: center;
}

.conversation-item.active .unread-badge {
  background: white;
  color: #007bff;
}

.empty-conversations {
  padding: 2rem;
  text-align: center;
  color: #6c757d;
}

.feather {
  width: 16px;
  height: 16px;
}

/* Mobile responsive */
@media (max-width: 768px) {
  .conversations-list {
    width: 100%;
  }
}
</style>
