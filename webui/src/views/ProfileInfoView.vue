<template>
  <div class="profile-info-view">
    <h2>{{ getTitle() }}</h2>
    
    <!-- Loading spinner -->
    <LoadingSpinner v-if="isLoading" />
    
    <!-- Error message -->
    <div v-else-if="error" class="error-message">
      <div class="error-icon">
        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#alert-circle" /></svg>
      </div>
      <h3>{{ error }}</h3>
      <p v-if="type === 'group' && error === 'Group not found'">
        The group you're looking for doesn't exist or you don't have access to it.
      </p>
      <p v-else-if="type === 'group' && error === 'Invalid group ID'">
        Invalid group ID. Please provide a valid group ID instead of "me".
      </p>
      <p v-else-if="type === 'user' && error === 'User not found'">
        The user you're looking for doesn't exist.
      </p>
      <button class="back-btn" @click="$router.go(-1)">Go Back</button>
    </div>
    
    <div v-else class="profile-sections">
      <!-- Personal Profile Section (for user type) -->
      <div v-if="type === 'user'" class="user-profile-section">
        <!-- Profile picture container + edit button (only for own profile) -->
        <div class="profile-picture-container">
          <img
            v-if="userData?.photoUrl" 
            :src="getImageUrl(userData.photoUrl)" 
            alt="Profile Picture" 
            class="profile-picture"
          >
          <img
            v-else src="/default-avatar.svg" 
            alt="Default Avatar" 
            class="profile-picture"
          >
                
          <button v-if="isOwnProfile" class="edit-profile-picture-btn" @click="changeProfilePicture">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#camera" /></svg>
          </button>
        </div>

        <!-- Username + Edit Button (only for own profile) -->
        <div class="username-section">
          <h4 v-if="isOwnProfile">Hi {{ userData?.username || currentUsername }}!</h4>
          <h4 v-else>{{ userData?.username || 'Loading...' }}</h4>
          <button v-if="isOwnProfile" class="edit-username-btn" @click="changeUsername">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#edit-2" /></svg>
          </button>
        </div>
            
        <!-- User ID display -->
        <div v-if="userData" class="user-id-section">
          <p v-if="isOwnProfile"><strong>Your User ID:</strong> {{ userData.id }}</p>
          <p v-else><strong>User ID:</strong> {{ userData.id }}</p>
        </div>
        
        <button v-if="isOwnProfile" class="logout-btn" @click="logout">Logout</button>
      </div>
      
      <!-- Group Info Section (only for group type) -->
      <div v-if="type === 'group'" class="group-info-section">
        <div class="group-header">
          <h4>{{ getGroupTitle() }}</h4>
          <span v-if="groupData" class="group-type-badge">Group</span>
        </div>
        
        <!-- Group picture + edit button -->
        <div v-if="groupData" class="group-picture-container">
          <img
            v-if="groupData?.photoUrl" 
            :src="getImageUrl(groupData.photoUrl)" 
            alt="Group Picture" 
            class="group-picture"
          >
          <img
            v-else src="/default-group.svg" 
            alt="Default Group" 
            class="group-picture"
          >
          
          <button v-if="groupData && groupData.type === 'group'" class="edit-group-picture-btn" @click="changeGroupPicture">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#camera" /></svg>
          </button>
        </div>
        
        <!-- Group name + edit button -->
        <div v-if="groupData" class="group-name-section">
          <h5>{{ groupData?.name || 'Unnamed Group' }}</h5>
          <button v-if="groupData.type === 'group'" class="edit-group-name-btn" @click="changeGroupName">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#edit-2" /></svg>
          </button>
        </div>
        
        <!-- Group info details -->
        <div v-if="groupData" class="group-details">
          <div class="group-info-item">
            <strong>{{ groupData.type === 'group' ? 'Group' : 'Chat' }} ID:</strong> {{ groupData.id }}
          </div>
          <div v-if="groupData.createdAt" class="group-info-item">
            <strong>Created:</strong> {{ formatDate(groupData.createdAt) }}
          </div>
          <div v-if="groupData.lastMessageAt" class="group-info-item">
            <strong>Last activity:</strong> {{ formatDate(groupData.lastMessageAt) }}
          </div>
        </div>
        
        <!-- Group participants -->
        <div v-if="groupData?.participants || groupData?.members" class="participants-section">
          <h5>Participants ({{ (groupData.participants || groupData.members || []).length }})</h5>
          <div class="participants-list">
            <div v-for="participant in (groupData.participants || groupData.members || [])" :key="participant.id" class="participant-item">
              <img
                v-if="participant.photoUrl" 
                :src="getImageUrl(participant.photoUrl)" 
                alt="Profile" 
                class="participant-avatar"
              >
              <img
                v-else src="/default-avatar.svg" 
                alt="Default Avatar" 
                class="participant-avatar"
              >
              <div class="participant-info">
                <span class="participant-name">{{ participant.username }}</span>
                <span v-if="participant.id === currentUserId" class="you-badge">You</span>
              </div>
              <button 
                v-if="participant.id !== currentUserId" 
                class="remove-member-btn" 
                :title="`Remove ${participant.username} from group`"
                @click="removeMember(participant)"
              >
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x" /></svg>
              </button>
            </div>
          </div>
        </div>
        
        <!-- Group actions -->
        <div class="group-actions">
          <button v-if="groupData?.type === 'group'" class="action-btn" @click="addMember">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#user-plus" /></svg>
            Add Member
          </button>
          <button class="action-btn secondary" @click="goToChat">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#message-square" /></svg>
            Go to Chat
          </button>
          <button class="action-btn danger" @click="leaveGroup">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#log-out" /></svg>
            Leave {{ groupData?.type === 'group' ? 'Group' : 'Chat' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import AuthService from '../services/auth.js';
import axios from '../services/axios.js';
import LoadingSpinner from '../components/LoadingSpinner.vue';

export default {
  name: 'ProfileInfoView',
  components: {
    LoadingSpinner
  },
  props: {
    type: {
      type: String,
      default: 'user' // 'user' or 'group'
    },
    id: {
      type: String,
      default: 'me'
    }
  },
  data() {
    return {
      userData: null,
      groupData: null,
      isLoading: false,
      error: null // Add error state
    }
  },
  computed: {
    currentUser() {
      return AuthService.getCurrentUser();
    },
    currentUsername() {
      return AuthService.getUsername();
    },
    currentUserId() {
      return AuthService.getUserId();
    },
    isOwnProfile() {
      return this.id === 'me';
    }
  },
  watch: {
    // Watch for changes to the id prop to reload data when navigating between profiles
    id: {
      immediate: false,
      handler() {
        this.reloadData();
      }
    }
  },
  async mounted() {
    // Load appropriate data based on type
    await this.reloadData();
  },
  methods: {
    // New method to reload data based on type
    async reloadData() {
      if (this.type === 'user') {
        await this.loadUser();
      } else if (this.type === 'group') {
        await this.loadGroup();
      }
    },
    
    getTitle() {
      if (this.type === 'group') {
        return 'Group Info';
      } else if (this.isOwnProfile) {
        return 'Your Profile';
      } else {
        return this.userData?.username ? `${this.userData.username}'s Profile` : 'User Profile';
      }
    },
    
    getGroupTitle() {
      if (!this.groupData) return 'Loading...';
      return this.groupData.name || 'Unnamed Group';
    },
    
    formatDate(dateString) {
      if (!dateString) return 'Unknown';
      try {
        const date = new Date(dateString);
        // Check if date is valid
        if (isNaN(date.getTime())) return 'Unknown';
        return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
      } catch (error) {
        console.error('Error formatting date:', error);
        return 'Unknown';
      }
    },
    
    canRemoveMember(participant) {
      // Can remove any member except yourself
      return participant.id !== AuthService.getUserId();
    },
    
    getImageUrl(photoUrl) {
      if (!photoUrl) return '/default-avatar.svg';
      
      // photoUrl comes as "/uploads/profiles/filename.jpg" from backend
      // We need to prepend the API base URL and add cache busting
      const baseURL = axios.defaults.baseURL || 'http://localhost:3000';
      const timestamp = Date.now(); // Cache busting parameter
      return `${baseURL}${photoUrl}?t=${timestamp}`;
    },
    async loadUser() {
        try {
            this.isLoading = true;
            this.error = null; // Clear any previous error
            let userIdParam;
            
            if (this.id === 'me') {
                // For own profile, use the actual user ID from auth service
                userIdParam = AuthService.getUserId();
                if (!userIdParam) {
                    console.error('Could not get current user ID');
                    this.error = 'Could not get current user ID';
                    return;
                }
            } else {
                // For other users, use the provided ID directly
                userIdParam = this.id;
            }
            
            const response = await axios.get(`/users/${userIdParam}`);
            this.userData = response.data;
            
            // Force Vue to reactively update the DOM
            this.$nextTick(() => {
                this.$forceUpdate();
            });
        }
        catch (error) {
            console.error('Error loading user:', error);
            if (error.response?.status === 404) {
                this.error = 'User not found';
            } else {
                this.error = 'Failed to load user information';
            }
        } finally {
            this.isLoading = false;
        }
    },

    async loadGroup() {
        try {
            this.isLoading = true;
            this.error = null; // Clear any previous error
            
            // Handle special case: type=group&id=me (invalid combination)
            if (this.id === 'me') {
                this.error = 'Invalid group ID';
                return;
            }
            
            // Use the user ID for the user-centric API pattern
            const userId = AuthService.getUserId();
            if (!userId) {
                console.error('Could not get current user ID');
                this.error = 'Could not get current user ID';
                return;
            }
            const response = await axios.get(`/users/${userId}/conversations/${this.id}`);
            
            // The backend returns ConversationDetailResponse with correct field names
            this.groupData = response.data;
            
            // Ensure participants array is available (use members as participants)
            if (this.groupData.members && !this.groupData.participants) {
                this.groupData.participants = this.groupData.members;
            }
            
            console.log('Group data loaded:', this.groupData); // Debug log
        }
        catch (error) {
            console.error('Error loading group:', error);
            if (error.response?.status === 404) {
                this.error = 'Group not found';
            } else {
                this.error = 'Failed to load group information';
            }
        } finally {
            this.isLoading = false;
        }
    },

    async changeUsername() {
      // Open modal to change username
      const newUsername = prompt('Enter new username:', this.currentUsername);
      if (newUsername && newUsername.trim()) {
        try {
          // Make API call to update username using user ID
          const userId = AuthService.getUserId();
          if (!userId) {
            console.error('Could not get current user ID');
            return;
          }
          
          await axios.put(`/users/${userId}/username`, { 
            name: newUsername 
          });
          
          // Update auth service with new username
          AuthService.setAuthData(AuthService.getAuthToken(), newUsername, userId);
          
          // Reload user data to get the updated username
          await this.loadUser();
        } catch (error) {
          console.error('Error updating username:', error);
          console.error('Full error response:', error.response);
          
          // Show more detailed error message
          const errorMessage = error.response?.data?.message || error.message || 'Unknown error';
          alert(`Failed to update username: ${errorMessage}`);
        }
      }
    },
    async changeProfilePicture() {
      // Create file input
      const input = document.createElement('input');
      input.type = 'file';
      input.accept = 'image/*';
      
      input.onchange = async (e) => {
        const file = e.target.files[0];
        if (!file) return;
        
        try {
          const formData = new FormData();
          formData.append('photo', file);
          
          // Use user ID for the endpoint
          const userId = AuthService.getUserId();
          if (!userId) {
            console.error('Could not get current user ID for photo upload');
            return;
          }
          
          await axios.put(`/users/${userId}/photo`, formData);
          
          // Reload user data to get the new photo URL
          await this.loadUser();
        } catch (error) {
          console.error('Error updating photo:', error);
        }
      };
      
      input.click();
    },
    logout() {
      // Clear auth data from sessionStorage
      AuthService.clearAuthData();
      
      // Redirect to login
      this.$router.push('/login');
    },

    // Group management methods
    async changeGroupName() {
      if (!this.groupData?.name) return;
      
      const newName = prompt('Enter new group name:', this.groupData.name);
      if (newName && newName.trim()) {
        try {
          const userId = AuthService.getUserId();
          if (!userId) {
            console.error('Could not get current user ID');
            return;
          }
          
          await axios.put(`/users/${userId}/groups/${this.id}/name`, { 
            name: newName.trim() 
          });
          
          // Reload group data to show new name
          await this.loadGroup();
        } catch (error) {
          console.error('Error updating group name:', error);
          const errorMessage = error.response?.data?.message || error.message || 'Unknown error';
          alert(`Failed to update group name: ${errorMessage}`);
        }
      }
    },

    async changeGroupPicture() {
      // Create file input
      const input = document.createElement('input');
      input.type = 'file';
      input.accept = 'image/*';
      
      input.onchange = async (e) => {
        const file = e.target.files[0];
        if (!file) return;
        
        try {
          const formData = new FormData();
          formData.append('photo', file);
          
          const userId = AuthService.getUserId();
          if (!userId) {
            console.error('Could not get current user ID');
            return;
          }
          
          await axios.put(`/users/${userId}/groups/${this.id}/photo`, formData);
          
          // Reload group data to show new picture
          await this.loadGroup();
        } catch (error) {
          console.error('Error updating group picture:', error);
          const errorMessage = error.response?.data?.message || error.message || 'Unknown error';
          alert(`Failed to update group picture: ${errorMessage}`);
        }
      };
      
      input.click();
    },

    async addMember() {
      const username = prompt('Enter username to add:');
      if (username && username.trim()) {
        try {
          const userId = AuthService.getUserId();
          if (!userId) {
            console.error('Could not get current user ID');
            return;
          }
          
          // First, search for the user to get their ID
          const searchResponse = await axios.get(`/users?q=${encodeURIComponent(username.trim())}`);
          const users = searchResponse.data.users || [];
          
          // Find exact username match
          const targetUser = users.find(user => user.username === username.trim());
          if (!targetUser) {
            alert('User not found');
            return;
          }
          
          // Now add the user using their ID
          await axios.post(`/users/${userId}/groups/${this.id}/members`, { 
            userId: targetUser.id 
          });
          
          // Reload group data to show new member
          await this.loadGroup();
        } catch (error) {
          console.error('Error adding member:', error);
          const errorMessage = error.response?.data?.message || error.message || 'Unknown error';
          alert(`Failed to add member: ${errorMessage}`);
        }
      }
    },

    async removeMember(participant) {
      if (confirm(`Are you sure you want to remove ${participant.username} from this group?`)) {
        try {
          const userId = AuthService.getUserId();
          if (!userId) {
            console.error('Could not get current user ID');
            return;
          }
          
          await axios.delete(`/users/${userId}/groups/${this.id}/members/${participant.id}`);
          
          // Reload group data to show updated member list
          await this.loadGroup();
        } catch (error) {
          console.error('Error removing member:', error);
          const errorMessage = error.response?.data?.message || error.message || 'Unknown error';
          alert(`Failed to remove member: ${errorMessage}`);
        }
      }
    },

    async leaveGroup() {
      if (confirm(`Are you sure you want to leave this group?`)) {
        try {
          const userId = AuthService.getUserId();
          if (!userId) {
            console.error('Could not get current user ID');
            return;
          }
          
          await axios.delete(`/users/${userId}/groups/${this.id}/members`);
          
          // Redirect back to chat view or home
          this.$router.push('/');
        } catch (error) {
          console.error('Error leaving group:', error);
          const errorMessage = error.response?.data?.message || error.message || 'Unknown error';
          alert(`Failed to leave group: ${errorMessage}`);
        }
      }
    },

    goToChat() {
      // Navigate to the chat view for this conversation
      this.$router.push(`/chat/${this.id}`);
    }
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
.group-info-section {
  padding: 0;
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

.username-section {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 1rem;
}

.username-section h4 {
    margin: 0; /* Remove default margin from h4 */
}

.user-id-section {
    margin-bottom: 1rem;
}

.user-id-section p {
    margin: 0;
    color: #6c757d;
    font-size: 14px;
}

.edit-username-btn {
    color: rgb(0, 0, 0);
    border: none;
    padding: 6px 8px;
    border-radius: 5px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    scale: 0.8
}

.edit-username-btn:hover {
    background: #218838;
    color: white;
    transition-duration: 0.15s;
}

.edit-profile-picture-btn {
    color: rgb(0, 0, 0);
    border: none;
    margin-left:15px;
    padding: 6px 8px;
    border-radius: 5px;
    cursor: pointer;
    scale:0.8
}
.edit-profile-picture-btn:hover {
    background: #218838;
    color:white;
    transition-duration:0.15s;
}

.profile-picture-container {
    position: relative;
    display: inline-block;
    margin-bottom: 1rem;
}

.profile-picture {
    width: 100px;
    height: 100px;
    border-radius: 50%;
    object-fit: cover;
    border: 2px solid #ddd;
}

.group-picture-container {
    position: relative;
    display: inline-block;
    margin-bottom: 1rem;
}

.group-picture {
    width: 100px;
    height: 100px;
    border-radius: 50%;
    object-fit: cover;
    border: 2px solid #ddd;
}

.edit-group-picture-btn {
    position: absolute;
    bottom: 0;
    right: 0;
    color: rgb(0, 0, 0);
    border: none;
    padding: 6px 8px;
    border-radius: 50%;
    cursor: pointer;
    background: white;
    box-shadow: 0 2px 4px rgba(0,0,0,0.2);
    scale: 0.8;
}

.edit-group-picture-btn:hover {
    background: #218838;
    color: white;
    transition-duration: 0.15s;
}

.group-header {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 1rem;
}

.group-header h4 {
    margin: 0;
}

.group-type-badge {
    background-color: #007bff;
    color: white;
    padding: 2px 8px;
    border-radius: 12px;
    font-size: 12px;
    font-weight: bold;
}

.group-name-section {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 1rem;
}

.group-name-section h5 {
    margin: 0;
    font-size: 16px;
}

.edit-group-name-btn {
    color: rgb(0, 0, 0);
    border: none;
    padding: 4px 6px;
    border-radius: 4px;
    cursor: pointer;
    scale: 0.8;
}

.edit-group-name-btn:hover {
    background: #218838;
    color: white;
    transition-duration: 0.15s;
}

.group-details {
    margin-bottom: 1rem;
    padding: 1rem;
    background-color: #f8f9fa;
    border-radius: 8px;
}

.group-info-item {
    margin-bottom: 0.5rem;
    font-size: 14px;
}

.group-info-item:last-child {
    margin-bottom: 0;
}

.group-info-item strong {
    color: #495057;
}

.participants-section {
    margin: 1rem 0;
}

.participants-section h5 {
    margin-bottom: 0.5rem;
    color: #333;
}

.participants-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    max-height: 200px;
    overflow-y: auto;
}

.participant-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.75rem;
    background-color: #f8f9fa;
    border-radius: 8px;
    position: relative;
}

.participant-item:hover {
    background-color: #e9ecef;
}

.participant-avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    object-fit: cover;
    border: 1px solid #ddd;
}

.participant-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 2px;
}

.participant-name {
    font-weight: 500;
    font-size: 14px;
}

.you-badge {
    background-color: #007bff;
    color: white;
    padding: 1px 6px;
    border-radius: 10px;
    font-size: 10px;
    font-weight: bold;
    width: fit-content;
}

.remove-member-btn {
    color: #dc3545;
    background: none;
    border: none;
    padding: 4px;
    border-radius: 4px;
    cursor: pointer;
    opacity: 0.7;
    scale: 0.8;
}

.remove-member-btn:hover {
    background-color: #dc3545;
    color: white;
    opacity: 1;
    transition-duration: 0.15s;
}

.group-actions {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    margin-top: 1rem;
}

.action-btn {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 8px 16px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    background-color: #007bff;
    color: white;
    font-size: 14px;
    transition: background-color 0.15s;
}

.action-btn:hover {
    background-color: #0056b3;
}

.action-btn.secondary {
    background-color: #6c757d;
}

.action-btn.secondary:hover {
    background-color: #545b62;
}

.action-btn.danger {
    background-color: #dc3545;
}

.action-btn.danger:hover {
    background-color: #c82333;
}

.loading {
    width: 100px;
    height: 100px;
    border-radius: 50%;
    border: 2px solid #ddd;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: #f8f9fa;
    color: #6c757d;
    font-size: 12px;
}

.error-message {
    text-align: center;
    padding: 3rem 2rem;
    background-color: #f8f9fa;
    border-radius: 8px;
    margin: 2rem 0;
}

.error-icon {
    margin-bottom: 1rem;
}

.error-icon .feather {
    width: 48px;
    height: 48px;
    color: #dc3545;
}

.error-message h3 {
    color: #dc3545;
    margin-bottom: 1rem;
    font-size: 1.5rem;
}

.error-message p {
    color: #6c757d;
    margin-bottom: 1.5rem;
    font-size: 1rem;
}

.back-btn {
    background-color: #007bff;
    color: white;
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.9rem;
    transition: background-color 0.15s;
}

.back-btn:hover {
    background-color: #0056b3;
}
</style>
