<template>
  <div class="profile-info-view">
    <h2>Profile Info View - Skeleton</h2>
    
    <div class="profile-sections">
      <div class="user-profile-section">
        <h4>Hi {{ currentUsername }}!</h4>
        <p>TODO: Show user photo, name, edit buttons</p>
        <button v-if="id === 'me'" @click="logout" class="logout-btn">Logout</button>
      </div>
      
      <div class="conversation-info-section">
        <h4>Conversation/Group Info Section</h4>
        <p>TODO: Show participants list</p>
        <p>TODO: Add member button, leave group button</p>
      </div>
    </div>
    
    <!-- TODO: Modals -->
    <p>TODO: AddMemberModal, EditNameModal</p>
  </div>
</template>

<script>
import AuthService from '../services/auth.js';

export default {
  name: 'ProfileInfoView',
  props: {
    type: {
      type: String,
      default: 'user' // 'user' or 'conversation'
    },
    id: {
      type: String,
      default: 'me'
    }
  },
  data() {
    return {
      // TODO: Add state for user, conversation, participants
    }
  },
  computed: {
    currentUser() {
      return AuthService.getCurrentUser();
    },
    currentUsername() {
      return AuthService.getUsername();
    },
    isOwnProfile() {
      return this.id === 'me';
    }
  },
  methods: {
    // TODO: loadUser()
    // TODO: loadConversation()
    // TODO: updateProfile()
    logout() {
      // Clear auth data from sessionStorage
      AuthService.clearAuthData();
      
      // Redirect to login
      this.$router.push('/login');
    }
    // TODO: addMember()
    // TODO: leaveGroup()
  }
}
</script>

<style scoped>
.profile-info-view {
  padding: 20px;
}

.profile-sections {
  margin-top: 20px;
}

.user-profile-section,
.conversation-info-section {
  border: 2px dashed #ccc;
  padding: 20px;
  margin-bottom: 20px;
}

.logout-btn {
  background-color: #dc3545;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  margin-top: 10px;
}

.logout-btn:hover {
  background-color: #c82333;
}
</style>
