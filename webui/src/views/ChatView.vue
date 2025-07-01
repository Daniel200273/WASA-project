<template>
  <div class="chat-view">
    <h2>Chat View - Skeleton</h2>
    <div class="layout-preview">
      <div class="conversations-sidebar">
        <h5>Conversations List</h5>
        <p>TODO: ConversationsList component</p>
      </div>
      <div class="chat-area">
        <h5>Active Chat</h5>
        <p>TODO: Messages + Input area</p>
      </div>
    </div>
    
    <!-- TODO: Modals -->
    <p>TODO: UserSearchModal, GroupCreateModal</p>
  </div>
</template>

<script>
import AuthService from '../services/auth.js';

export default {
  name: 'ChatView',
  data() {
    return {
      // TODO: Add state for conversations, messages, etc.
    }
  },
  computed: {
    currentUser() {
      return AuthService.getCurrentUser();
    },
    currentUsername() {
      return AuthService.getUsername();
    }
  },
  methods: {
    // TODO: loadConversations()
    // TODO: selectConversation() 
    // TODO: sendMessage()
    // TODO: openUserSearch()
  }
}
</script>

<style scoped>
.chat-view {
  padding: 20px;
}

.layout-preview {
  display: flex;
  gap: 20px;
  margin-top: 20px;
}

.conversations-sidebar {
  width: 300px;
  border: 2px dashed #ccc;
  padding: 20px;
}

.chat-area {
  flex: 1;
  border: 2px dashed #ccc;
  padding: 20px;
}
</style>
