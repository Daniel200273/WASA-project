<template>
  <div class="message-item" :class="{ 'own-message': isOwn }">
    <!-- Message content -->
    <div class="message-content">
      <!-- Sender name (only for group chats and non-own messages) -->
      <div v-if="!isOwn && showSender && isGroupChat" class="message-sender">
        {{ message.senderUsername }}
      </div>

      <!-- Reply context -->
      <div v-if="message.replyToId" class="reply-context">
        <div class="reply-indicator"></div>
        <div class="reply-info">
          <span class="reply-to">Reply to message</span>
        </div>
      </div>

      <!-- Message bubble -->
      <div class="message-bubble" :class="{ 'own-bubble': isOwn, 'optimistic-message': message.isOptimistic }">
        <!-- Photo message -->
        <div v-if="message.photoUrl" class="message-photo">
          <img :src="getImageUrl(message.photoUrl)" :alt="message.content || 'Photo'" @click="openPhotoModal" />
        </div>

        <!-- Text content -->
        <div v-if="message.content" class="message-text">
          {{ message.content }}
        </div>

        <!-- Forwarded indicator -->
        <div v-if="message.forwarded" class="forwarded-indicator">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#corner-up-right" /></svg>
          <span>Forwarded</span>
        </div>

        <!-- Message footer -->
        <div class="message-footer">
          <span class="message-time">{{ formatTime(message.timestamp) }}</span>
          
          <!-- Message status for own messages -->
          <div v-if="isOwn" class="message-status">
            <svg v-if="computedMessageStatus === 'sent'" class="feather status-icon">
              <use href="/feather-sprite-v4.29.0.svg#check" />
            </svg>
            <svg v-else-if="computedMessageStatus === 'read'" class="feather status-icon status-read">
              <use href="/feather-sprite-v4.29.0.svg#check-check" />
            </svg>
          </div>
        </div>
      </div>

      <!-- Reactions -->
      <div v-if="message.comments && message.comments.length > 0" class="message-reactions">
        <div
          v-for="reaction in groupedReactions"
          :key="reaction.emoticon"
          class="reaction-pill"
          :class="{ 'own-reaction': reaction.hasOwnReaction }"
          @click="toggleReaction(reaction.emoticon)"
        >
          <span class="reaction-emoji">{{ reaction.emoticon }}</span>
          <span class="reaction-count">{{ reaction.count }}</span>
        </div>
      </div>
    </div>

    <!-- Message actions -->
    <div class="message-actions" v-show="showActions">
      <button class="action-btn" @click="$emit('reply', message)" title="Reply">
        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#corner-up-left" /></svg>
      </button>
      
      <button class="action-btn" @click="showReactionPicker = !showReactionPicker" title="React">
        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#smile" /></svg>
      </button>
      
      <button v-if="isOwn" class="action-btn danger" @click="$emit('delete', message)" title="Delete">
        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#trash-2" /></svg>
      </button>

      <!-- Quick reaction picker -->
      <div v-if="showReactionPicker" class="reaction-picker" @click.stop>
        <button
          v-for="emoji in quickReactions"
          :key="emoji"
          class="quick-reaction"
          @click="addReaction(emoji)"
        >
          {{ emoji }}
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'MessageItem',
  props: {
    message: {
      type: Object,
      required: true
    },
    isOwn: {
      type: Boolean,
      default: false
    },
    showAvatar: {
      type: Boolean,
      default: true
    },
    showSender: {
      type: Boolean,
      default: true
    },
    isGroupChat: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      showActions: false,
      showReactionPicker: false,
      quickReactions: ['👍', '❤️', '😂', '😮', '😢', '😡']
    }
  },
  computed: {
    groupedReactions() {
      if (!this.message.comments) return [];
      
      const groups = {};
      this.message.comments.forEach(comment => {
        if (!groups[comment.emoticon]) {
          groups[comment.emoticon] = {
            emoticon: comment.emoticon,
            count: 0,
            hasOwnReaction: false,
            users: []
          };
        }
        groups[comment.emoticon].count++;
        groups[comment.emoticon].users.push(comment.username);
        
        // Check if current user has this reaction (you'd need to pass current user ID)
        // groups[comment.emoticon].hasOwnReaction = comment.userId === this.currentUserId;
      });
      
      return Object.values(groups);
    },
    
    computedMessageStatus() {
      // Simple logic: If the message is older than 5 minutes, consider it "read"
      // In a real app, this would be based on actual read receipts from the backend
      const messageTime = new Date(this.message.timestamp);
      const now = new Date();
      const fiveMinutesAgo = new Date(now.getTime() - 5 * 60 * 1000);
      
      // If message is older than 5 minutes, show as read
      if (messageTime < fiveMinutesAgo) {
        return 'read';
      }
      
      // Otherwise, show as sent
      return 'sent';
    }
  },
  methods: {
    getImageUrl(photoUrl) {
      if (!photoUrl) return null;
      // Use the same logic as ChatView for consistent URL handling
      const baseURL = 'http://localhost:3000'; // You might want to import this from a config
      return photoUrl.startsWith('http') ? photoUrl : `${baseURL}${photoUrl}?t=${Date.now()}`;
    },
    
    formatTime(timestamp) {
      const date = new Date(timestamp);
      return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    },
    
    addReaction(emoji) {
      this.$emit('react', this.message, emoji);
      this.showReactionPicker = false;
    },
    
    toggleReaction(emoji) {
      this.$emit('react', this.message, emoji);
    },
    
    openPhotoModal() {
      // TODO: Implement photo modal
      console.log('Open photo modal for:', this.message.photoUrl);
    }
  }
}
</script>

<style scoped>
.message-item {
  margin-bottom: 0.5rem;
  position: relative;
  padding: 0.25rem;
  border-radius: 8px;
  transition: background-color 0.15s ease;
  display: flex;
  flex-direction: column;
}

.message-item:hover .message-actions {
  opacity: 1;
}

/* Own messages aligned to right */
.message-item.own-message {
  align-items: flex-end;
}

.message-item.own-message .message-content {
  align-items: flex-end;
}

/* Message content */
.message-content {
  display: flex;
  flex-direction: column;
  max-width: 70%;
  position: relative;
  width: fit-content;
}

.message-sender {
  font-size: 0.75rem;
  color: #6c757d;
  font-weight: 600;
  margin-bottom: 0.25rem;
  padding-left: 0.75rem;
}

/* Reply context */
.reply-context {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.25rem;
  padding-left: 0.75rem;
}

.reply-indicator {
  width: 3px;
  height: 20px;
  background-color: #007bff;
  border-radius: 2px;
}

.reply-to {
  font-size: 0.75rem;
  color: #007bff;
  font-weight: 500;
}

/* Message bubble */
.message-bubble {
  background-color: white;
  border-radius: 18px;
  padding: 0.75rem 1rem;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  position: relative;
  word-wrap: break-word;
  width: fit-content;
  min-width: 60px;
  max-width: 100%;
}

.message-bubble.own-bubble {
  background-color: #007bff;
  color: white;
}

/* Photo messages */
.message-photo {
  margin: -0.25rem -0.5rem 0.5rem -0.5rem;
  border-radius: 14px 14px 8px 8px;
  overflow: hidden;
  cursor: pointer;
}

.message-photo img {
  width: 100%;
  max-width: 300px;
  height: auto;
  display: block;
}

/* Text content */
.message-text {
  line-height: 1.4;
  font-size: 0.9rem;
}

/* Forwarded indicator */
.forwarded-indicator {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  font-size: 0.75rem;
  color: #6c757d;
  margin-bottom: 0.25rem;
  font-style: italic;
}

.own-bubble .forwarded-indicator {
  color: rgba(255, 255, 255, 0.8);
}

/* Message footer */
.message-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 0.25rem;
  margin-top: 0.25rem;
}

.message-time {
  font-size: 0.7rem;
  color: #6c757d;
}

.own-bubble .message-time {
  color: rgba(255, 255, 255, 0.8);
}

/* Message status */
.message-status {
  display: flex;
  align-items: center;
}

.status-icon {
  width: 12px;
  height: 12px;
  color: rgba(255, 255, 255, 0.8);
}

.status-icon.status-read {
  color: #28a745;
}

/* Reactions */
.message-reactions {
  display: flex;
  flex-wrap: wrap;
  gap: 0.25rem;
  margin-top: 0.25rem;
  padding-left: 0.75rem;
}

.reaction-pill {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  background-color: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 12px;
  padding: 0.125rem 0.5rem;
  font-size: 0.75rem;
  cursor: pointer;
  transition: all 0.15s ease;
}

.reaction-pill:hover {
  background-color: #e9ecef;
}

.reaction-pill.own-reaction {
  background-color: #e3f2fd;
  border-color: #007bff;
  color: #007bff;
}

.reaction-emoji {
  font-size: 0.875rem;
}

.reaction-count {
  font-weight: 500;
  min-width: 1ch;
}

/* Message actions */
.message-actions {
  position: absolute;
  top: -10px;
  right: 10px;
  display: flex;
  gap: 0.125rem;
  background-color: white;
  border-radius: 6px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  padding: 0.25rem;
  opacity: 0;
  transition: opacity 0.15s ease;
  z-index: 10;
}

.own-message .message-actions {
  right: auto;
  left: 10px;
}

.action-btn {
  background: none;
  border: none;
  padding: 0.25rem;
  border-radius: 4px;
  cursor: pointer;
  color: #6c757d;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
}

.action-btn:hover {
  background-color: #f8f9fa;
  color: #212529;
}

.action-btn.danger:hover {
  background-color: #dc3545;
  color: white;
}

.action-btn .feather {
  width: 14px;
  height: 14px;
}

/* Reaction picker */
.reaction-picker {
  position: absolute;
  top: 100%;
  right: 0;
  background-color: white;
  border-radius: 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
  padding: 0.5rem;
  display: flex;
  gap: 0.25rem;
  z-index: 20;
  margin-top: 0.25rem;
}

.quick-reaction {
  background: none;
  border: none;
  padding: 0.25rem;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.15s ease;
}

.quick-reaction:hover {
  background-color: #f8f9fa;
}

/* Optimistic message styling */
.optimistic-message {
  opacity: 0.7;
  position: relative;
}

.optimistic-message::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(45deg, transparent 30%, rgba(255, 255, 255, 0.3) 50%, transparent 70%);
  animation: shimmer 1.5s infinite;
  border-radius: inherit;
  pointer-events: none;
}

@keyframes shimmer {
  0% { transform: translateX(-100%); }
  100% { transform: translateX(100%); }
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .message-content {
    max-width: 85%;
  }
  
  .message-bubble {
    padding: 0.5rem 0.75rem;
  }
  
  .message-photo {
    margin: -0.125rem -0.25rem 0.25rem -0.25rem;
  }
}
</style>
