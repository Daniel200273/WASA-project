<template>
  <div class="chat-view">
    <!-- Two-panel layout: Conversations sidebar + Chat area -->
    <div class="chat-layout">
      <!-- Conversations sidebar -->
      <div class="conversations-sidebar" :class="{ 'mobile-hidden': selectedConversation }">
        <div class="conversations-header">
          <h3>Conversations</h3>
          <div class="header-actions">
            <button 
              class="btn btn-sm btn-outline-secondary refresh-btn" 
              :disabled="isRefreshing"
              :title="isRefreshing ? 'Refreshing...' : 'Refresh all'"
              @click="refreshAll"
            >
              <svg class="feather" :class="{ 'spinning': isRefreshing }">
                <use href="/feather-sprite-v4.29.0.svg#refresh-cw" />
              </svg>
            </button>
            <button 
              class="btn btn-sm btn-primary" 
              title="Start new chat"
              @click="showUserSearch = true"
            >
              <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#plus" /></svg>
            </button>
          </div>
        </div>
        
        <div class="conversations-content">
          <LoadingSpinner v-if="conversationsLoading" :loading="true" />
          
          <div v-else-if="conversations.length === 0" class="no-conversations">
            <svg class="feather empty-icon"><use href="/feather-sprite-v4.29.0.svg#message-circle" /></svg>
            <p>No conversations yet</p>
            <button class="btn btn-sm btn-outline-primary" @click="showUserSearch = true">
              Start chatting
            </button>
          </div>
          
          <div v-else class="conversations-list">
            <div 
              v-for="conversation in conversations" 
              :key="conversation.id"
              class="conversation-item"
              :class="{ active: selectedConversationId === conversation.id }"
              @click="selectConversation(conversation)"
            >
              <div class="conversation-avatar-container">
                <img 
                  :src="getConversationAvatar(conversation)" 
                  :alt="conversation.name"
                  class="conversation-avatar"
                >
                <div v-if="conversation.unreadCount > 0" class="unread-badge">
                  {{ conversation.unreadCount > 99 ? '99+' : conversation.unreadCount }}
                </div>
              </div>
              
              <div class="conversation-info">
                <div class="conversation-header">
                  <h4 class="conversation-name">{{ conversation.name }}</h4>
                  <span v-if="conversation.lastMessage" class="conversation-time">
                    {{ formatConversationTime(conversation.lastMessage.timestamp) }}
                  </span>
                </div>
                <p v-if="conversation.lastMessage" class="last-message">
                  <span v-if="isLastMessageFromCurrentUser(conversation)" class="you-prefix">You: </span>
                  {{ conversation.lastMessage.content || 'Photo' }}
                </p>
                <p v-else class="no-messages">No messages yet</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Chat area -->
      <div class="chat-main" :class="{ 'mobile-hidden': !selectedConversation }">
        <div v-if="!selectedConversation" class="no-conversation-selected">
          <div class="welcome-message">
            <svg class="feather welcome-icon"><use href="/feather-sprite-v4.29.0.svg#message-circle" /></svg>
            <h3>Select a conversation</h3>
            <p>Choose a conversation from the sidebar to start messaging</p>
          </div>
        </div>

        <!-- Active conversation -->
        <div v-else class="active-chat">
          <!-- Chat header -->
          <div class="chat-header">
            <div class="chat-header-info">
              <!-- Back button for mobile -->
              <button class="back-btn" title="Back to conversations" @click="goBackToConversations">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#arrow-left" /></svg>
              </button>
              
              <img 
                :src="getConversationAvatar(selectedConversation)" 
                :alt="selectedConversation.name"
                class="conversation-avatar"
              >
              <div class="conversation-details">
                <h4 class="conversation-name">{{ selectedConversation.name }}</h4>
                <span v-if="selectedConversation.type === 'group'" class="participant-count">
                  {{ selectedConversation.members?.length || 0 }} members
                </span>
              </div>
            </div>
            <div class="chat-header-actions">
              <button 
                class="btn btn-sm btn-outline-secondary refresh-btn" 
                :disabled="isRefreshing"
                :title="isRefreshing ? 'Refreshing...' : 'Refresh all'"
                @click="refreshAll"
              >
                <svg class="feather" :class="{ 'spinning': isRefreshing }">
                  <use href="/feather-sprite-v4.29.0.svg#refresh-cw" />
                </svg>
              </button>
              <button 
                class="btn btn-sm btn-outline-secondary" 
                :title="selectedConversation.type === 'group' ? 'Group Info' : 'User Info'"
                @click="goToConversationInfo"
              >
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#info" /></svg>
              </button>
            </div>
          </div>

          <!-- Messages area -->
          <div ref="messagesContainer" class="messages-container">
            <!-- Only show loading spinner on initial load, not on refresh -->
            <LoadingSpinner v-if="messagesLoading && messages.length === 0" :loading="true" />
            
            <div v-else class="messages-list">
              <!-- Individual messages -->
              <div v-for="message in messages" :key="message.id" class="message-wrapper">
                <MessageItem 
                  :message="message"
                  :is-own="message.senderId === currentUserId"
                  :show-sender="shouldShowSender(message)"
                  :is-group-chat="selectedConversation?.type === 'group'"
                  :conversation-read-at="conversationReadAt"
                  :all-messages="messages"
                  @reply="replyToMessage"
                  @react="reactToMessage"
                  @delete="deleteMessage"
                />
              </div>

              <!-- Empty state -->
              <div v-if="messages.length === 0" class="empty-messages">
                <svg class="feather empty-icon"><use href="/feather-sprite-v4.29.0.svg#message-square" /></svg>
                <p>No messages yet. Start the conversation!</p>
              </div>
            </div>
          </div>

          <!-- Message input area -->
          <div class="message-input-container">
            <!-- Reply preview -->
            <div v-if="replyingTo" class="reply-preview">
              <div class="reply-content">
                <span class="reply-label">Replying to {{ replyingTo.senderUsername }}</span>
                <p class="reply-text">{{ replyingTo.content || 'Photo' }}</p>
              </div>
              <button class="cancel-reply" @click="cancelReply">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x" /></svg>
              </button>
            </div>

            <!-- Message input -->
            <MessageInput 
              :placeholder="getInputPlaceholder()"
              :disabled="sendingMessage"
              :replying-to="replyingTo"
              @send-message="sendMessage"
              @send-photo="sendPhoto"
              @cancel-reply="cancelReply"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Modals -->
    <UserSearchModal 
      v-if="showUserSearch"
      @close="showUserSearch = false"
      @select-user="startConversationWithUser"
    />
  </div>
</template>

<script>
import AuthService from '../services/auth.js';
import axios from '../services/axios.js';
import LoadingSpinner from '../components/LoadingSpinner.vue';
import MessageItem from '../components/chat/MessageItem.vue';
import MessageInput from '../components/chat/MessageInput.vue';
import UserSearchModal from '../components/modals/UserSearchModal.vue';
import { getConversationAvatar, getImageUrl } from '../utils/imageUtils.js';

export default {
  name: 'ChatView',
  components: {
    LoadingSpinner,
    MessageItem,
    MessageInput,
    UserSearchModal
  },
  props: {
    conversationId: {
      type: String,
      default: null
    }
  },
  data() {
    return {
      // Conversations
      conversations: [],
      conversationsLoading: false,
      selectedConversationId: null,
      selectedConversation: null,
      conversationReadAt: null,
      
      // Messages
      messages: [],
      messagesLoading: false,
      sendingMessage: false,
      
      // UI State
      showUserSearch: false,
      replyingTo: null,
      
      // Refresh state
      lastRefreshTime: null,
      refreshDebounceMs: 1000 // Prevent multiple refreshes within 1 second
    }
  },
  computed: {
    currentUserId() {
      return AuthService.getUserId();
    },
    currentUsername() {
      return AuthService.getUsername();
    },
    isRefreshing() {
      return this.conversationsLoading || this.messagesLoading;
    },

  },
  watch: {
    conversationId: {
      immediate: true,
      async handler(newId) {
        if (newId) {
          this.selectedConversationId = newId;
          // Simple: just refresh everything to ensure consistency
          await this.refreshAll();
        } else {
          // Clear selection when no conversation ID
          this.selectedConversationId = null;
          this.selectedConversation = null;
          this.messages = [];
          this.replyingTo = null;
        }
      }
    }
  },
  async mounted() {
    try {
      // Simple: just load conversations, the watcher will handle the rest
      await this.loadConversations();
      
      // On desktop, auto-select first conversation if no specific one is requested
      if (!this.conversationId && this.conversations.length > 0 && window.innerWidth > 768) {
        this.$router.replace(`/chat/${this.conversations[0].id}`);
      }
    } catch (error) {
      console.error('Initial load failed:', error);
      // Show a more user-friendly message for initial load failures
      if (error.code !== 'ECONNABORTED') {
        alert('Failed to load conversations. Please refresh the page or check your connection.');
      }
    }
  },
  methods: {
    // === CONVERSATION MANAGEMENT ===
    async loadConversations() {
      try {
        this.conversationsLoading = true;
        
        const result = await this.retryApiCall(async () => {
          const userId = AuthService.getUserId();
          return await axios.get(`/users/${userId}/conversations`, {
            timeout: 5000 // 5 second timeout
          });
        });
        
        const newConversations = result.data.conversations || [];
        
        // Always update conversations
        this.conversations = newConversations;
        
        // On desktop, select the first conversation if available and no specific one is selected
        // On mobile, keep conversations list visible by default
        if (!this.selectedConversationId && this.conversations.length > 0 && window.innerWidth > 768) {
          this.selectedConversationId = this.conversations[0].id;
        }
        
      } catch (error) {
        console.error('Error loading conversations:', error);
        
        // Re-throw the error so the caller can handle it
        throw error;
      } finally {
        this.conversationsLoading = false;
      }
    },

    async selectConversation(conversation) {
      // Simple: just navigate to the conversation, the watcher will handle the rest
      this.$router.replace(`/chat/${conversation.id}`);
    },

    async loadConversationMessages(conversationId) {
      if (!conversationId) return;
      
      try {
        this.messagesLoading = true;
        
        const result = await this.retryApiCall(async () => {
          const userId = AuthService.getUserId();
          return await axios.get(`/users/${userId}/conversations/${conversationId}`, {
            timeout: 5000 // 5 second timeout
          });
        });
        
        this.selectedConversation = result.data;
        this.messages = result.data.messages || [];
        
        // Update the conversation in the sidebar list if needed
        const conversationIndex = this.conversations.findIndex(c => c.id === conversationId);
        if (conversationIndex !== -1) {
          // Update the conversation in place to maintain sidebar state
          this.conversations[conversationIndex] = {
            ...this.conversations[conversationIndex],
            ...result.data,
            // Keep the lastMessage and unreadCount from the list view
            lastMessage: this.conversations[conversationIndex].lastMessage,
            unreadCount: this.conversations[conversationIndex].unreadCount
          };
        }
        
        // Set the read timestamp to now (when the conversation is loaded)
        this.conversationReadAt = new Date().toISOString();
        
        this.$nextTick(() => {
          this.scrollToBottom();
        });
        
      } catch (error) {
        console.error('Error loading messages:', error);
        
        if (error.response?.status === 404) {
          // Conversation not found, redirect back to chat list
          this.$router.push('/chat');
        }
        
        // Re-throw the error so the caller can handle it
        throw error;
      } finally {
        this.messagesLoading = false;
      }
    },

    // === MESSAGE HANDLING ===
    async sendMessage(content, photo = null) {
      if (!this.selectedConversationId || this.sendingMessage) return;
      
      try {
        this.sendingMessage = true;
        const userId = AuthService.getUserId();
        
        let response;
        if (photo) {
          const formData = new FormData();
          formData.append('photo', photo);
          if (this.replyingTo) {
            formData.append('replyTo', this.replyingTo.id);
          }
          
          response = await axios.post(
            `/users/${userId}/conversations/${this.selectedConversationId}/messages`,
            formData,
            { headers: { 'Content-Type': 'multipart/form-data' } }
          );
        } else {
          const messageData = { content };
          if (this.replyingTo) {
            messageData.replyTo = this.replyingTo.id;
          }
          
          response = await axios.post(
            `/users/${userId}/conversations/${this.selectedConversationId}/messages`,
            messageData
          );
        }
        
        // Add the new message to the messages array
        this.messages.push(response.data);
        this.replyingTo = null;
        
        // Update the conversation's last message in local state
        this.updateConversationLastMessage(response.data);
        
        // Update read timestamp since sender is actively in the conversation
        this.conversationReadAt = new Date().toISOString();
        
        this.$nextTick(() => {
          this.scrollToBottom();
        });
        
      } catch (error) {
        console.error('Error sending message:', error);
        alert('Failed to send message. Please try again.');
      } finally {
        this.sendingMessage = false;
      }
    },

    async sendPhoto(photo) {
      await this.sendMessage(null, photo);
    },

    async deleteMessage(message) {
      if (!confirm('Are you sure you want to delete this message?')) return;
      
      try {
        const userId = AuthService.getUserId();
        await axios.delete(`/users/${userId}/messages/${message.id}`);
        
        // Remove from local state
        this.messages = this.messages.filter(m => m.id !== message.id);
      } catch (error) {
        console.error('Error deleting message:', error);
        alert('Failed to delete message');
      }
    },

    async reactToMessage(message, emoticon) {
      try {
        const userId = AuthService.getUserId();
        await axios.post(`/users/${userId}/messages/${message.id}/comments`, {
          emoticon
        });
        
        // Refresh messages to show new reaction
        await this.loadConversationMessages(this.selectedConversationId);
      } catch (error) {
        console.error('Error reacting to message:', error);
      }
    },

    replyToMessage(message) {
      this.replyingTo = message;
    },

    cancelReply() {
      this.replyingTo = null;
    },

    // === LOCAL STATE UPDATES ===
    updateConversationLastMessage(message) {
      const conversation = this.conversations.find(c => c.id === this.selectedConversationId);
      if (conversation) {
        conversation.lastMessage = {
          id: message.id,
          content: message.content,
          timestamp: message.timestamp,
          senderId: message.senderId,
          senderUsername: message.senderUsername
        };
        // Reset unread count to 0 when sending a message
        conversation.unreadCount = 0;
        // Move conversation to top of list
        const index = this.conversations.indexOf(conversation);
        if (index > 0) {
          this.conversations.splice(index, 1);
          this.conversations.unshift(conversation);
        }
      }
    },

    // === UI HELPERS ===
    shouldShowSender(message) {
      if (this.selectedConversation?.type !== 'group') {
        // In direct messages, never show sender names
        return false;
      }
      
      // In group chats, show sender name if it's the first message from this sender
      const messageIndex = this.messages.findIndex(m => m.id === message.id);
      const prevMessage = messageIndex > 0 ? this.messages[messageIndex - 1] : null;
      
      // Show sender name if it's the first message from this sender
      return !prevMessage || prevMessage.senderId !== message.senderId;
    },

    isLastMessageFromCurrentUser(conversation) {
      if (!conversation.lastMessage || !conversation.lastMessage.senderId) return false;
      // Determine if the message is from the current user by comparing senderId
      return conversation.lastMessage.senderId === this.currentUserId;
    },

    getConversationAvatar,

    getImageUrl,

    getInputPlaceholder() {
      if (!this.selectedConversation) return 'Type a message...';
      return this.selectedConversation.type === 'group' 
        ? `Message ${this.selectedConversation.name}...`
        : `Message ${this.selectedConversation.name}...`;
    },

    formatConversationTime(timestamp) {
      if (!timestamp) return '';
      const date = new Date(timestamp);
      const now = new Date();
      const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
      const messageDate = new Date(date.getFullYear(), date.getMonth(), date.getDate());
      
      if (messageDate.getTime() === today.getTime()) {
        // Today - show time
        return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
      } else if (messageDate.getTime() === today.getTime() - 86400000) {
        // Yesterday
        return 'Yesterday';
      } else if (now.getTime() - messageDate.getTime() < 7 * 86400000) {
        // This week - show day name
        return date.toLocaleDateString([], { weekday: 'short' });
      } else {
        // Older - show date
        return date.toLocaleDateString([], { month: 'short', day: 'numeric' });
      }
    },

    scrollToBottom() {
      if (this.$refs.messagesContainer) {
        this.$refs.messagesContainer.scrollTop = this.$refs.messagesContainer.scrollHeight;
      }
    },

    // === NAVIGATION ===
    goBackToConversations() {
      this.selectedConversationId = null;
      this.selectedConversation = null;
      this.messages = [];
      this.replyingTo = null;
      this.$router.push('/chat');
    },

    goToConversationInfo() {
      if (!this.selectedConversation) return;
      
      if (this.selectedConversation.type === 'group') {
        this.$router.push(`/profile?type=group&id=${this.selectedConversation.id}`);
      } else {
        // For direct messages, find the other participant
        const otherParticipant = this.selectedConversation.members?.find(
          member => member.id !== this.currentUserId
        );
        if (otherParticipant) {
          this.$router.push(`/profile?type=user&id=${otherParticipant.id}`);
        }
      }
    },

    async startConversationWithUser(user) {
      try {
        const userId = AuthService.getUserId();
        
        const response = await axios.post(`/users/${userId}/conversations`, {
          userId: user.id
        }, {
          timeout: 10000 // 10 second timeout for conversation creation
        });
        
        // Create conversation object for the list
        const conversationForList = {
          id: response.data.id,
          type: 'direct',
          name: user.username,
          photoUrl: user.photoUrl,
          lastMessage: null,
          unreadCount: 0
        };
        
        // Add to conversations list if not already there
        const existingConv = this.conversations.find(c => c.id === response.data.id);
        if (!existingConv) {
          this.conversations.unshift(conversationForList);
        }
        
        // Close the modal first
        this.showUserSearch = false;
        
        // Set up the conversation directly using the API response
        this.selectedConversation = response.data;
        this.selectedConversationId = response.data.id;
        this.messages = response.data.messages || [];
        this.conversationReadAt = new Date().toISOString();
        
        // Navigate to the conversation
        this.$router.replace(`/chat/${response.data.id}`);
        
        // Scroll to bottom if there are messages
        this.$nextTick(() => {
          this.scrollToBottom();
        });
        
      } catch (error) {
        console.error('Error starting conversation:', error);
        alert('Failed to start conversation. Please try again.');
      }
    },

    // === REFRESH FUNCTIONALITY ===
    async refreshAll() {
      // Debounce rapid refresh calls
      const now = Date.now();
      if (this.lastRefreshTime && (now - this.lastRefreshTime) < this.refreshDebounceMs) {
        return;
      }
      
      // Prevent multiple simultaneous refreshes
      if (this.conversationsLoading || this.messagesLoading) {
        return;
      }
      
      this.lastRefreshTime = now;
      
      try {
        // Always refresh conversations
        await this.loadConversations();
        
        // If we have a selected conversation ID, update the UI and load messages
        if (this.selectedConversationId) {
          const conversation = this.conversations.find(c => c.id === this.selectedConversationId);
          if (conversation) {
            this.selectedConversation = conversation;
            this.conversationReadAt = new Date().toISOString();
            this.replyingTo = null; // Clear any reply state
            if (conversation.unreadCount > 0) {
              conversation.unreadCount = 0;
            }
          }
          await this.loadConversationMessages(this.selectedConversationId);
        }
      } catch (error) {
        console.error('Refresh failed:', error);
        // Don't show alert for timeout/network errors, just log them
        if (error.code !== 'ECONNABORTED' && error.code !== 'NETWORK_ERROR') {
          alert('Failed to refresh. Please check your connection and try again.');
        }
      }
    },

    // === ERROR HANDLING ===
    async retryApiCall(apiCall, maxRetries = 2, delay = 500) {
      for (let i = 0; i < maxRetries; i++) {
        try {
          return await apiCall();
        } catch (error) {
          // If it's a database lock error, wait a bit but not too long
          if (error.response?.status === 500 && error.response?.data?.includes?.('database is locked')) {
            await new Promise(resolve => setTimeout(resolve, delay * (i + 1)));
          } else if (i === maxRetries - 1) {
            throw error; // Re-throw if it's the last attempt
          } else {
            await new Promise(resolve => setTimeout(resolve, delay));
          }
        }
      }
    },
  }
}
</script>

<style scoped>
.chat-view {
  --top-nav-height: 48px;
  height: calc(100vh - var(--top-nav-height));
  width: 100%;
  overflow: hidden;
  display: grid;
  grid-template-rows: 1fr;
  grid-template-columns: 320px 1fr;
  position: relative;
}

.chat-layout {
  display: contents; /* Let children participate in the parent grid */
}

/* Conversations sidebar */
.conversations-sidebar {
  background-color: white;
  border-right: 1px solid #e9ecef;
  overflow: hidden;
  display: grid;
  grid-template-rows: 70px 1fr;
}

.conversations-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1.5rem;
  border-bottom: 1px solid #e9ecef;
  background-color: white;
  height: 70px;
  box-sizing: border-box;
}

.conversations-header h3 {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
  color: #212529;
}

.header-actions {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.refresh-btn {
  min-width: 32px;
  justify-content: center;
}

.refresh-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.conversations-content {
  overflow: hidden;
  display: grid;
  grid-template-rows: 1fr;
}

.no-conversations {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  padding: 2rem;
  text-align: center;
  color: #6c757d;
}

.conversations-list {
  overflow-y: auto;
}

.conversation-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1.25rem;
  cursor: pointer;
  border-bottom: 1px solid #f1f3f4;
  transition: background-color 0.15s ease;
}

.conversation-item:hover {
  background-color: #f8f9fa;
}

.conversation-item.active {
  background-color: #e3f2fd;
  border-right: 3px solid #007bff;
}

.conversation-avatar-container {
  position: relative;
  flex-shrink: 0;
}

.conversation-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #e9ecef;
}

.unread-badge {
  position: absolute;
  top: -4px;
  right: -4px;
  background-color: #dc3545;
  color: white;
  border-radius: 10px;
  font-size: 0.625rem;
  font-weight: bold;
  min-width: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 4px;
}

.conversation-info {
  flex: 1;
  min-width: 0;
}

.conversation-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 0.25rem;
}

.conversation-name {
  margin: 0;
  font-size: 0.9rem;
  font-weight: 600;
  color: #212529;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
}

.conversation-time {
  font-size: 0.75rem;
  color: #6c757d;
  flex-shrink: 0;
  margin-left: 0.5rem;
}

.last-message,
.no-messages {
  margin: 0;
  font-size: 0.8rem;
  color: #6c757d;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.you-prefix {
  color: #007bff;
  font-weight: 500;
}

/* Chat main area */
.chat-main {
  background-color: #f8f9fa;
  overflow: hidden;
  display: grid;
  grid-template-rows: 1fr;
}

/* No conversation selected state */
.no-conversation-selected {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.welcome-message {
  text-align: center;
  max-width: 400px;
  padding: 2rem;
}

.welcome-icon {
  width: 64px;
  height: 64px;
  color: #6c757d;
  margin-bottom: 1rem;
}

.welcome-message h3 {
  color: #495057;
  margin-bottom: 0.5rem;
}

.welcome-message p {
  color: #6c757d;
  margin-bottom: 1.5rem;
}

/* Active chat area */
.active-chat {
  background-color: white;
  overflow: hidden;
  display: grid;
  grid-template-rows: 70px 1fr auto;
}

/* Chat header */
.chat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1.5rem;
  border-bottom: 1px solid #e9ecef;
  background-color: white;
}

.chat-header-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.back-btn {
  background: none;
  border: none;
  padding: 0.5rem;
  border-radius: 50%;
  cursor: pointer;
  color: #6c757d;
  display: none; /* Hidden by default */
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
}

.back-btn:hover {
  background-color: #e9ecef;
  color: #495057;
}

.back-btn .feather {
  width: 18px;
  height: 18px;
}

.chat-header .conversation-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #e9ecef;
}

.conversation-details {
  display: flex;
  flex-direction: column;
}

.chat-header .conversation-name {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
  color: #212529;
}

.participant-count,
.online-status {
  font-size: 0.875rem;
  color: #6c757d;
}

.chat-header-actions {
  display: flex;
  gap: 0.5rem;
}

/* Messages container */
.messages-container {
  overflow-y: auto;
  overflow-x: hidden;
  padding: 1rem;
  background-color: #f8f9fa;
}

.messages-list {
  display: flex;
  flex-direction: column;
}

/* Message wrapper */
.message-wrapper {
  margin-bottom: 0;
}

/* Empty messages state */
.empty-messages {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 200px;
  color: #6c757d;
}

.empty-icon {
  width: 48px;
  height: 48px;
  margin-bottom: 1rem;
  opacity: 0.5;
}

/* Message input area */
.message-input-container {
  border-top: 1px solid #e9ecef;
  background-color: white;
}

/* Reply preview */
.reply-preview {
  display: flex;
  align-items: center;
  padding: 0.75rem 1rem;
  background-color: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
}

.reply-content {
  flex: 1;
}

.reply-label {
  font-size: 0.75rem;
  color: #007bff;
  font-weight: 600;
  display: block;
  margin-bottom: 0.25rem;
}

.reply-text {
  font-size: 0.875rem;
  color: #6c757d;
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 300px;
}

.cancel-reply {
  background: none;
  border: none;
  color: #6c757d;
  padding: 0.25rem;
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.cancel-reply:hover {
  background-color: #e9ecef;
}

.cancel-reply .feather {
  width: 16px;
  height: 16px;
}

/* Responsive design */
@media (max-width: 768px) {
  .chat-view {
    grid-template-columns: 1fr;
    position: relative;
  }
  
  .conversations-sidebar {
    position: fixed;
    top: var(--top-nav-height);
    left: 0;
    width: 100%;
    height: calc(100vh - var(--top-nav-height));
    z-index: 1040;
    transform: translateX(0);
    transition: transform 0.3s ease;
    box-shadow: none;
  }

  .conversations-sidebar.mobile-hidden {
    transform: translateX(-100%);
  }

  .chat-main {
    position: fixed;
    top: var(--top-nav-height);
    left: 0;
    width: 100%;
    height: calc(100vh - var(--top-nav-height));
    z-index: 1030;
    transform: translateX(100%);
    transition: transform 0.3s ease;
  }

  .chat-main.mobile-hidden {
    transform: translateX(100%);
  }

  .chat-main:not(.mobile-hidden) {
    transform: translateX(0);
  }

  .back-btn {
    display: flex !important; /* Show back button on mobile */
  }

  .chat-header {
    padding: 0.75rem 1rem;
  }

  .chat-header .conversation-avatar {
    width: 32px;
    height: 32px;
  }

  .chat-header .conversation-name {
    font-size: 1rem;
  }

  .messages-container {
    padding: 0.75rem;
  }
}

/* Utility classes */
.btn {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  border: 1px solid transparent;
  border-radius: 4px;
  font-size: 0.875rem;
  font-weight: 500;
  text-decoration: none;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-primary {
  background-color: #007bff;
  border-color: #007bff;
  color: white;
}

.btn-primary:hover {
  background-color: #0056b3;
  border-color: #0056b3;
}

.btn-outline-primary {
  color: #007bff;
  border-color: #007bff;
  background-color: transparent;
}

.btn-outline-primary:hover {
  color: white;
  background-color: #007bff;
  border-color: #007bff;
}

.btn-sm {
  padding: 0.25rem 0.5rem;
  font-size: 0.75rem;
}

.btn-outline-secondary {
  color: #6c757d;
  border-color: #6c757d;
  background-color: transparent;
}

.btn-outline-secondary:hover {
  color: white;
  background-color: #6c757d;
  border-color: #6c757d;
}

.feather {
  width: 16px;
  height: 16px;
}
</style>