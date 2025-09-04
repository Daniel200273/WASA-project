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
                  <div class="conversation-header-right">
                    <span v-if="conversation.lastMessage" class="conversation-time">
                      {{ formatConversationTime(conversation.lastMessage.timestamp) }}
                    </span>
                  </div>
                </div>
                <p v-if="conversation.lastMessage" class="last-message">
                  <span v-if="isLastMessageFromCurrentUser(conversation)" class="you-prefix">You: </span>
                  <span v-else-if="conversation.type === 'group' && conversation.lastMessage.senderUsername" class="sender-prefix">{{ conversation.lastMessage.senderUsername }}: </span>
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
                  @forward="forwardMessage"
                  @refresh-message="handleRefreshMessage"
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
    <ForwardMessageModal
      v-if="showForwardModal"
      :conversations="conversations"
      @confirm="confirmForward"
      @cancel="cancelForward"
      @forward-to-user="forwardToNewUser"
    />
    <NotificationModal
      v-if="notificationModal.show"
      :type="notificationModal.type"
      :title="notificationModal.title"
      :message="notificationModal.message"
      @close="closeNotificationModal"
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
import ForwardMessageModal from '../components/modals/ForwardMessageModal.vue';
import NotificationModal from '../components/modals/NotificationModal.vue';
import { getConversationAvatar, getImageUrl } from '../utils/imageUtils.js';

export default {
  name: 'ChatView',
  components: {
    LoadingSpinner,
    MessageItem,
    MessageInput,
    UserSearchModal,
    ForwardMessageModal,
    NotificationModal
  },
  props: {
    conversationId: {
      type: String,
      default: null
    }
  },
  data() {
    return {
      // Component lifecycle
      isComponentDestroyed: false,
      
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

      // Forward modal state
      forwardingMessage: null,
      showForwardModal: false,
      
      // Notification modal state
      notificationModal: {
        show: false,
        type: 'info',
        title: '',
        message: ''
      },

      // Smart polling state
      pollingTimer: null,
      pollingInterval: 10000, // Default 10 seconds
      lastActivity: Date.now(),
      isUserActive: true,
      isTabVisible: true,
      lastConversationsUpdate: null,
      lastMessagesUpdate: null,
      
      // Activity tracking
      updateActivity: null,
      inactivityChecker: null,
      lastReadMarkTime: 0 // Track when we last marked conversation as read
    }
  },
  computed: {
    currentUserId() {
      return AuthService.getUserId();
    },
    currentUsername() {
      return AuthService.getUsername();
    },
    // Smart polling computed properties
    shouldPollFast() {
      // Poll fast when user is active and in a conversation
      return this.isUserActive && this.isTabVisible && this.selectedConversationId;
    },
    optimalPollingInterval() {
      if (!this.isTabVisible) {
        return 60000; // 1 minute when tab is hidden
      }
      if (!this.isUserActive) {
        return 30000; // 30 seconds when inactive
      }
      if (this.selectedConversationId) {
        return 3000; // 3 seconds when in active chat
      }
      return 10000; // 10 seconds default
    }
  },
  watch: {
    conversationId: {
      immediate: true,
      async handler(newId) {
        try {
          if (newId) {
            this.selectedConversationId = newId;
            // Load conversations and then the specific conversation
            await this.loadConversations();
            await this.loadConversationMessages(newId);
            
            // Reset polling timers when switching conversations
            this.lastMessagesUpdate = Date.now();
            if (this.restartPolling) {
              this.restartPolling();
            }
          } else {
            // Clear selection when no conversation ID
            this.selectedConversationId = null;
            this.selectedConversation = null;
            this.messages = [];
            this.replyingTo = null;
          }
        } catch (error) {
          console.error('Error in conversation watcher:', error);
        }
      }
    },
    
    // Watch for polling interval changes
    optimalPollingInterval: {
      handler(newInterval, oldInterval) {
        try {
          if (newInterval !== oldInterval && this.pollingTimer) {
            console.log(`Polling interval changed: ${oldInterval}ms → ${newInterval}ms`);
            this.restartPolling();
          }
        } catch (error) {
          console.error('Error in polling interval watcher:', error);
        }
      }
    }
  },
  async mounted() {
    try {
      // Request notification permission
      if ('Notification' in window && Notification.permission === 'default') {
        Notification.requestPermission();
      }
      
      // Simple: just load conversations, the watcher will handle the rest
      await this.loadConversations();
      // On desktop, auto-select first conversation if no specific one is requested
      if (!this.conversationId && this.conversations.length > 0 && window.innerWidth > 768) {
        this.$router.replace(`/chat/${this.conversations[0].id}`);
      }
      
      // Initialize smart polling with error handling
      try {
        this.initializeSmartPolling();
        this.startSmartPolling();
      } catch (pollingError) {
        console.error('Failed to initialize smart polling:', pollingError);
        // Continue without polling if it fails
      }
    } catch (error) {
      console.error('Initial load failed:', error);
      this.showNotification('error', 'Connection Error', 'Failed to load conversations. Please refresh the page or check your connection.');
    }
  },
  beforeUnmount() {
    this.isComponentDestroyed = true; // Mark component as destroyed
    this.stopSmartPolling();
    this.removeActivityListeners();
  },
  methods: {
    // === CONVERSATION MANAGEMENT ===
    async loadConversations() {
      if (this.isComponentDestroyed) return;
      
      try {
        this.conversationsLoading = true;
        
        const userId = AuthService.getUserId();
        const response = await axios.get(`/users/${userId}/conversations`);
        
        // Check if component is still mounted before updating reactive data
        if (this.isComponentDestroyed) return;
        
        const newConversations = response.data.conversations || [];
        
        // Always update conversations
        this.conversations = newConversations;
        
        // On desktop, select the first conversation if available and no specific one is selected
        // On mobile, keep conversations list visible by default
        if (!this.selectedConversationId && this.conversations.length > 0 && window.innerWidth > 768) {
          this.selectedConversationId = this.conversations[0].id;
        }
        
      } catch (error) {
        console.error('Error loading conversations:', error);
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
      if (!conversationId || this.isComponentDestroyed) return;
      
      try {
        this.messagesLoading = true;
        
        const userId = AuthService.getUserId();
        const response = await axios.get(`/users/${userId}/conversations/${conversationId}`);
        
        // Check if component is still mounted before updating reactive data
        if (this.isComponentDestroyed) return;
        
        this.selectedConversation = response.data;
        this.messages = response.data.messages || [];
        
        // Update the conversation in the sidebar list if needed
        const conversationIndex = this.conversations.findIndex(c => c.id === conversationId);
        if (conversationIndex !== -1) {
          // Update the conversation in place to maintain sidebar state
          this.conversations[conversationIndex] = {
            ...this.conversations[conversationIndex],
            ...response.data,
            // Keep the lastMessage from the list view but clear unreadCount since user is viewing it
            lastMessage: this.conversations[conversationIndex].lastMessage,
            unreadCount: 0 // Clear unread count when conversation is opened
          };
        }
        
        // Set the read timestamp to now (when the conversation is loaded)
        this.conversationReadAt = new Date().toISOString();
        
        this.$nextTick(() => {
          this.scrollToBottom();
        });
        
      } catch (error) {
        console.error('Error loading messages:', error);
        
        // Log more detailed error information for debugging
        if (error.response) {
          console.error('Response error:', {
            status: error.response.status,
            statusText: error.response.statusText,
            data: error.response.data,
            conversationId: conversationId
          });
        } else if (error.request) {
          console.error('Request error:', {
            message: error.message,
            code: error.code,
            conversationId: conversationId
          });
        } else {
          console.error('Unknown error:', error.message);
        }
        
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

    // === SMART POLLING SYSTEM ===
    initializeSmartPolling() {
      // Set up activity tracking
      this.setupActivityListeners();
      
      // Set up visibility change detection
      document.addEventListener('visibilitychange', this.handleVisibilityChange);
      
      // Set initial timestamps
      this.lastConversationsUpdate = Date.now();
      this.lastMessagesUpdate = Date.now();
    },

    setupActivityListeners() {
      if (this.isComponentDestroyed) return;
      
      const activityEvents = ['mousedown', 'mousemove', 'keypress', 'scroll', 'touchstart', 'click'];
      
      this.updateActivity = () => {
        if (this.isComponentDestroyed) return;
        this.lastActivity = Date.now();
        this.isUserActive = true;
      };
      
      // Add listeners for user activity
      activityEvents.forEach(event => {
        document.addEventListener(event, this.updateActivity, true);
      });
      
      // Check for inactivity every 30 seconds
      this.inactivityChecker = setInterval(() => {
        if (this.isComponentDestroyed) return;
        
        const now = Date.now();
        const timeSinceActivity = now - this.lastActivity;
        
        // Consider user inactive after 2 minutes of no activity
        this.isUserActive = timeSinceActivity < 120000;
      }, 30000);
    },

    removeActivityListeners() {
      const activityEvents = ['mousedown', 'mousemove', 'keypress', 'scroll', 'touchstart', 'click'];
      
      if (this.updateActivity) {
        activityEvents.forEach(event => {
          document.removeEventListener(event, this.updateActivity, true);
        });
      }
      
      if (this.inactivityChecker) {
        clearInterval(this.inactivityChecker);
        this.inactivityChecker = null;
      }
      
      document.removeEventListener('visibilitychange', this.handleVisibilityChange);
    },

    handleVisibilityChange() {
      if (this.isComponentDestroyed) return;
      
      this.isTabVisible = !document.hidden;
      
      if (this.isTabVisible) {
        // Tab became visible, restart polling immediately
        this.lastActivity = Date.now();
        this.isUserActive = true;
        this.restartPolling();
      }
    },

    startSmartPolling() {
      this.stopSmartPolling(); // Clear any existing timer
      this.scheduleNextPoll();
    },

    stopSmartPolling() {
      if (this.pollingTimer) {
        clearTimeout(this.pollingTimer);
        this.pollingTimer = null;
      }
    },

    restartPolling() {
      this.stopSmartPolling();
      this.scheduleNextPoll();
    },

    scheduleNextPoll() {
      if (this.isComponentDestroyed) return;
      
      const interval = this.optimalPollingInterval;
      
      this.pollingTimer = setTimeout(async () => {
        try {
          await this.performSmartPoll();
        } catch (error) {
          console.error('Polling error:', error);
        }
        
        // Schedule next poll
        this.scheduleNextPoll();
      }, interval);
    },

    async performSmartPoll() {
      if (this.isComponentDestroyed || this.conversationsLoading || this.messagesLoading) {
        return;
      }

      const now = Date.now();
      
      try {
        // Always check for conversation updates (but less frequently when inactive)
        const shouldUpdateConversations = 
          (now - this.lastConversationsUpdate) > (this.isTabVisible ? 15000 : 60000);
        
        if (shouldUpdateConversations) {
          await this.loadConversationsQuietly();
          this.lastConversationsUpdate = now;
        }
        
        // Check for new messages in active conversation
        if (this.selectedConversationId) {
          const shouldUpdateMessages = 
            (now - this.lastMessagesUpdate) > (this.shouldPollFast ? 3000 : 15000);
          
          if (shouldUpdateMessages) {
            await this.loadNewMessages();
            this.lastMessagesUpdate = now;
          }
          
          // Periodically ensure conversation is marked as read during active viewing
          // This helps with checkmark updates even when no new messages arrive
          if (this.isUserActive && this.isTabVisible) {
            await this.ensureConversationMarkedAsRead();
          }
        }
      } catch (error) {
        console.error('Smart polling error:', error);
      }
    },

    async loadConversationsQuietly() {
      // Load conversations without showing loading state
      try {
        const userId = AuthService.getUserId();
        const response = await axios.get(`/users/${userId}/conversations`);
        
        if (this.isComponentDestroyed) return;
        
        const newConversations = response.data.conversations || [];
        
        // Check if there are any meaningful changes
        if (this.hasConversationChanges(newConversations)) {
          this.conversations = newConversations;
        }
        
      } catch (error) {
        console.error('Error in quiet conversations refresh:', error);
      }
    },

    async loadNewMessages() {
      if (!this.selectedConversationId || this.isComponentDestroyed) return;
      
      try {
        // Store current message count and scroll position for comparison
        const currentMessageCount = this.messages.length;
        const wasAtBottom = this.isAtBottom();
        
        const userId = AuthService.getUserId();
        
        // Get the full conversation to check for new messages
        // This is more reliable than incremental loading
        const response = await axios.get(`/users/${userId}/conversations/${this.selectedConversationId}`);
        
        if (this.isComponentDestroyed) return;
        
        const freshMessages = response.data.messages || [];
        
        // Debug logging to understand what's happening
        console.log('Polling for new messages:', {
          conversationId: this.selectedConversationId,
          currentMessageCount: currentMessageCount,
          freshMessageCount: freshMessages.length,
          hasNewMessages: freshMessages.length > currentMessageCount
        });
        
        // Check if there are new messages OR if message status might have changed
        const hasNewMessages = freshMessages.length > currentMessageCount;
        const shouldUpdateStatus = this.shouldCheckForStatusUpdates(freshMessages);
        
        if (hasNewMessages || shouldUpdateStatus) {
          if (hasNewMessages) {
            console.log('New messages detected, updating chat view');
          } else {
            console.log('Potential message status changes detected, updating chat view');
          }
          
          // Update messages array with fresh data
          this.messages = freshMessages;
          
          // Ensure conversation is marked as read since user is actively viewing it
          await this.ensureConversationMarkedAsRead();
          
          // Clear unread count since user is viewing the conversation
          const conversation = this.conversations.find(c => c.id === this.selectedConversationId);
          if (conversation) {
            conversation.unreadCount = 0;
          }
          
          // Get new messages for notifications (only messages from others)
          const newMessages = freshMessages.slice(currentMessageCount);
          const messagesFromOthers = newMessages.filter(m => m.senderId !== this.currentUserId);
          
          // Show notification if tab is hidden and there are new messages from others
          if (!this.isTabVisible && messagesFromOthers.length > 0) {
            this.showNewMessageNotification(messagesFromOthers);
          }
          
          // Update conversation in sidebar with the latest message
          if (newMessages.length > 0) {
            this.updateConversationFromNewMessages(newMessages);
          }
          
          // Scroll to bottom if user was at bottom before the update
          this.$nextTick(() => {
            if (wasAtBottom) {
              this.scrollToBottom();
            }
          });
          
          console.log(`Added ${freshMessages.length - currentMessageCount} new messages to chat`);
        } else if (freshMessages.length < currentMessageCount) {
          // Messages were deleted, update the entire view
          console.log('Messages were deleted, refreshing chat view');
          this.messages = freshMessages;
        } else if (shouldUpdateStatus) {
          // No new messages, but status might have changed (e.g., read receipts)
          console.log('Message status updated, refreshing message display');
          this.messages = freshMessages;
          
          // Don't auto-scroll for status updates
          // Users shouldn't lose their scroll position just for checkmark updates
        }
        
      } catch (error) {
        console.error('Error loading new messages:', error);
        
        // If the conversation was deleted or access was revoked, redirect
        if (error.response?.status === 404 || error.response?.status === 403) {
          console.warn('Conversation no longer accessible, redirecting to chat list');
          this.$router.push('/chat');
        }
      }
    },

    hasConversationChanges(newConversations) {
      if (newConversations.length !== this.conversations.length) {
        return true;
      }
      
      // Check for changes in last message timestamps or unread counts
      for (let i = 0; i < newConversations.length; i++) {
        const oldConv = this.conversations[i];
        const newConv = newConversations[i];
        
        if (!oldConv || 
            oldConv.lastMessageAt !== newConv.lastMessageAt ||
            oldConv.unreadCount !== newConv.unreadCount) {
          return true;
        }
      }
      
      return false;
    },

    updateConversationFromNewMessages(newMessages) {
      if (newMessages.length === 0) return;
      
      const latestMessage = newMessages[newMessages.length - 1];
      const conversation = this.conversations.find(c => c.id === this.selectedConversationId);
      
      if (conversation) {
        conversation.lastMessage = {
          id: latestMessage.id,
          content: latestMessage.content,
          timestamp: latestMessage.timestamp,
          senderId: latestMessage.senderId,
          senderUsername: latestMessage.senderUsername
        };
        
        // Update unread count for messages from others
        const messagesFromOthers = newMessages.filter(m => m.senderId !== this.currentUserId);
        if (messagesFromOthers.length > 0 && !this.isTabVisible) {
          conversation.unreadCount = (conversation.unreadCount || 0) + messagesFromOthers.length;
        }
      }
    },

    showNewMessageNotification(newMessages) {
      if ('Notification' in window && Notification.permission === 'granted') {
        const messagesFromOthers = newMessages.filter(m => m.senderId !== this.currentUserId);
        
        if (messagesFromOthers.length > 0) {
          const latestMessage = messagesFromOthers[messagesFromOthers.length - 1];
          const title = this.selectedConversation?.name || 'New Message';
          const body = latestMessage.content || 'Photo';
          
          new Notification(title, {
            body: `${latestMessage.senderUsername}: ${body}`,
            icon: '/favicon.ico',
            tag: `message-${this.selectedConversationId}` // Prevent duplicate notifications
          });
        }
      }
    },

    isAtBottom() {
      if (!this.$refs || !this.$refs.messagesContainer) return true;
      
      const container = this.$refs.messagesContainer;
      const threshold = 100; // 100px threshold
      
      return (container.scrollHeight - container.scrollTop - container.clientHeight) < threshold;
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
          if (content && content.trim()) {
            formData.append('content', content.trim());
          }
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
        
        // Reset polling timers after sending a message
        this.lastMessagesUpdate = Date.now();
        this.lastConversationsUpdate = Date.now();
        this.lastActivity = Date.now();
        this.isUserActive = true;
        
        // Schedule a status update check shortly after sending
        // This helps with quick checkmark updates when the recipient reads the message
        setTimeout(() => {
          if (!this.isComponentDestroyed && this.selectedConversationId) {
            this.loadNewMessages();
          }
        }, 2000); // Check again in 2 seconds for status updates
        
        this.$nextTick(() => {
          this.scrollToBottom();
        });
        
      } catch (error) {
        console.error('Error sending message:', error);
        this.showNotification('error', 'Send Failed', 'Failed to send message. Please try again.');
      } finally {
        this.sendingMessage = false;
      }
    },

    async sendPhoto(photo, caption = null) {
      await this.sendMessage(caption, photo);
    },

    async deleteMessage(message) {
      if (!confirm('Are you sure you want to delete this message?')) return;
      
      try {
        const userId = AuthService.getUserId();
        await axios.delete(`/users/${userId}/messages/${message.id}`);
        
        // Only update state if component is still mounted
        if (!this.isComponentDestroyed) {
          // Remove from local state
          this.messages = this.messages.filter(m => m.id !== message.id);
        }
      } catch (error) {
        console.error('Error deleting message:', error);
        this.showNotification('error', 'Delete Failed', 'Failed to delete message. Please try again.');
      }
    },

    async reactToMessage(message, emoticon) {
      try {
        const userId = AuthService.getUserId();
        await axios.post(`/users/${userId}/messages/${message.id}/comments`, {
          emoticon
        });
        
        // Only refresh if component is still mounted
        if (!this.isComponentDestroyed && this.selectedConversationId) {
          // Refresh messages to show new reaction
          await this.loadConversationMessages(this.selectedConversationId);
        }
      } catch (error) {
        console.error('Error reacting to message:', error);
        this.showNotification('error', 'Reaction Failed', 'Failed to add reaction. Please try again.');
      }
    },

    async handleRefreshMessage(messageId) {
      // Refresh the conversation messages to update reactions
      if (!this.isComponentDestroyed && this.selectedConversationId) {
        try {
          await this.loadConversationMessages(this.selectedConversationId);
        } catch (error) {
          console.error('Error refreshing message:', error);
        }
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
        this.showNotification('error', 'Conversation Failed', 'Failed to start conversation. Please try again.');
      }
    },

    // === FORWARD MESSAGE ===
    async forwardMessage(message) {
      // Open modal to select target conversation
      this.forwardingMessage = message;
      this.showForwardModal = true;
    },

    async confirmForward(targetConversationId) {
      if (!this.forwardingMessage || !targetConversationId) return;
      try {
        const userId = AuthService.getUserId();
        await axios.post(
          `/users/${userId}/messages/${this.forwardingMessage.id}/forward`,
          { conversationId: targetConversationId }
        );
        this.showForwardModal = false;
        this.forwardingMessage = null;
        this.showNotification('success', 'Message Forwarded', 'Message forwarded successfully!');
      } catch (error) {
        console.error('Error forwarding message:', error);
        this.showNotification('error', 'Forward Failed', 'Failed to forward message.');
      }
    },

    async forwardToNewUser(user) {
      if (!this.forwardingMessage || !user) return;
      try {
        const userId = AuthService.getUserId();
        
        // First, create or get conversation with the user
        const conversationResponse = await axios.post(`/users/${userId}/conversations`, {
          userId: user.id
        });
        
        // Then forward the message to that conversation
        await axios.post(
          `/users/${userId}/messages/${this.forwardingMessage.id}/forward`,
          { conversationId: conversationResponse.data.id }
        );
        
        this.showForwardModal = false;
        this.forwardingMessage = null;
        this.showNotification('success', 'Message Forwarded', `Message forwarded to ${user.username}!`);
        
        // Optionally, refresh conversations to show the new chat
        await this.loadConversations();
        
      } catch (error) {
        console.error('Error forwarding message to new user:', error);
        this.showNotification('error', 'Forward Failed', 'Failed to forward message to new user.');
      }
    },

    // Helper method to ensure conversation is marked as read by making a strategic call to getConversation
    async ensureConversationMarkedAsRead() {
      if (!this.selectedConversationId || this.isComponentDestroyed) return;
      
      // Only mark as read if we haven't done so recently (avoid too frequent calls)
      const now = Date.now();
      if (this.lastReadMarkTime && (now - this.lastReadMarkTime) < 30000) {
        return; // Skip if we marked as read less than 30 seconds ago
      }
      
      try {
        const userId = AuthService.getUserId();
        // Make a call to getConversation which automatically marks it as read
        await axios.get(`/users/${userId}/conversations/${this.selectedConversationId}`);
        this.lastReadMarkTime = now;
      } catch (error) {
        // Don't show error to user, just log it - this is not critical functionality
        console.warn('Failed to mark conversation as read via getConversation:', error);
      }
    },

    // Check if message status might have changed (for checkmark updates)
    shouldCheckForStatusUpdates(freshMessages) {
      // If we don't have current messages, no comparison possible
      if (!this.messages || this.messages.length === 0) {
        return false;
      }
      
      // If message counts are different, new messages method will handle it
      if (freshMessages.length !== this.messages.length) {
        return false;
      }
      
      // Check if any message status has changed
      for (let i = 0; i < this.messages.length; i++) {
        const currentMsg = this.messages[i];
        const freshMsg = freshMessages[i];
        
        if (currentMsg.id === freshMsg.id && currentMsg.status !== freshMsg.status) {
          return true; // Status changed
        }
      }
      
      return false;
    },

    cancelForward() {
      this.showForwardModal = false;
      this.forwardingMessage = null;
    },

    // === NOTIFICATION MODAL ===
    showNotification(type, title, message) {
      this.notificationModal = {
        show: true,
        type,
        title,
        message
      };
    },

    closeNotificationModal() {
      this.notificationModal.show = false;
    }
  }
}
</script>

<style scoped>
.chat-view {
  --top-nav-height: 48px; /* Adjust this value to match your top navigation bar height */
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

.conversation-header-right {
  display: flex;
  align-items: center;
  gap: 0.25rem;
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

.sender-prefix {
  color: #6c757d;
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