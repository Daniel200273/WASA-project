<template>
  <div class="profile-info-view">
    <h2>{{ getTitle() }}</h2>
    
    <div class="profile-sections">
        <!-- Personal Profile Section (for user type) -->
        <div v-if="type === 'user'" class="user-profile-section">
        <!-- Profile picture container + edit button (only for own profile) -->
            <div class="profile-picture-container">
                <img v-if="userData?.photoUrl" 
                    :src="getImageUrl(userData.photoUrl)" 
                    alt="Profile Picture" 
                    class="profile-picture" />
                <img v-else src="/default-avatar.svg" 
                    alt="Default Avatar" 
                    class="profile-picture" />
                
                <button v-if="isOwnProfile" @click="changeProfilePicture" class="edit-profile-picture-btn">
                    <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#camera"/></svg>
                </button>
            </div>

        <!-- Username + Edit Button (only for own profile) -->
            <div class="username-section">
                <h4 v-if="isOwnProfile">Hi {{ userData?.username || currentUsername }}!</h4>
                <h4 v-else>{{ userData?.username || 'Loading...' }}</h4>
                <button v-if="isOwnProfile" @click="changeUsername" class="edit-username-btn">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#edit-2"/></svg>
                </button>
            </div>
            
            <!-- User ID display -->
            <div v-if="userData" class="user-id-section">
                <p v-if="isOwnProfile"><strong>Your User ID:</strong> {{ userData.id }}</p>
                <p v-else><strong>User ID:</strong> {{ userData.id }}</p>
            </div>
        
            <button v-if="isOwnProfile" @click="logout" class="logout-btn">Logout</button>
        </div>
      
      <!-- Group/Conversation Info Section (only for conversation type) -->
      <div v-if="type === 'conversation'" class="conversation-info-section">
        <h4>Group: {{ groupData?.name || 'Loading...' }}</h4>
        
        <!-- Group picture -->
        <div v-if="groupData" class="group-picture-container">
          <img v-if="groupData?.photoUrl" 
              :src="getImageUrl(groupData.photoUrl)" 
              alt="Group Picture" 
              class="group-picture" />
          <img v-else src="/default-group.svg" 
              alt="Default Group" 
              class="group-picture" />
        </div>
        
        <!-- Group participants -->
        <div v-if="groupData?.participants" class="participants-section">
          <h5>Participants ({{ groupData.participants.length }})</h5>
          <div class="participants-list">
            <div v-for="participant in groupData.participants" :key="participant.id" class="participant-item">
              <img v-if="participant.photoUrl" 
                  :src="getImageUrl(participant.photoUrl)" 
                  alt="Profile" 
                  class="participant-avatar" />
              <img v-else src="/default-avatar.svg" 
                  alt="Default Avatar" 
                  class="participant-avatar" />
              <span>{{ participant.username }}</span>
            </div>
          </div>
        </div>
        
        <!-- Group actions -->
        <div class="group-actions">
          <button @click="addMember" class="action-btn">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#user-plus"/></svg>
            Add Member
          </button>
          <button @click="leaveGroup" class="action-btn danger">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#log-out"/></svg>
            Leave Group
          </button>
        </div>
      </div>
    </div>
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
      userData: null,
      groupData: null
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
  async mounted() {
    // Load appropriate data based on type
    if (this.type === 'user') {
      await this.loadUser();
    } else if (this.type === 'conversation') {
      await this.loadGroup();
    }
  },
  methods: {
    getTitle() {
      if (this.type === 'conversation') {
        return 'Group Info';
      } else if (this.isOwnProfile) {
        return 'Your Profile';
      } else {
        return this.userData?.username ? `${this.userData.username}'s Profile` : 'User Profile';
      }
    },
    
    getImageUrl(photoUrl) {
      // photoUrl comes as "/uploads/profiles/filename.jpg" from backend
      // We need to prepend the API base URL and add cache busting
      const baseURL = this.$axios.defaults.baseURL || 'http://localhost:3000';
      const timestamp = Date.now(); // Cache busting parameter
      return `${baseURL}${photoUrl}?t=${timestamp}`;
    },
    async loadUser() {
        try {
            const userId = this.id === 'me' ? this.currentUser.token : this.id;
            const response = await this.$axios.get(`/users/${userId}`);
            this.userData = response.data;
        }
        catch (error) {
            console.error('Error loading user:', error);
        }
    },

    async loadGroup() {
        try {
            const response = await this.$axios.get(`/groups/${this.id}`);
            this.groupData = response.data;
        }
        catch (error) {
            console.error('Error loading group:', error);
        }
    },

    async changeUsername() {
      // Open modal to change username
      const newUsername = prompt('Enter new username:', this.currentUsername);
      if (newUsername && newUsername.trim()) {
        try {
          // Make API call to update username using token as userID
          await this.$axios.put(`/users/${this.currentUser.token}/username`, { 
            name: newUsername 
          });
          
          // Update auth service with new username
          AuthService.setAuthData(this.currentUser.token, newUsername, this.currentUser.userId);
          
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
          
          await this.$axios.put(`/users/${this.currentUser.token}/photo`, formData);
          
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
    async addMember() {
      // TODO: Open modal to add member
      const username = prompt('Enter username to add:');
      if (username && username.trim()) {
        try {
          await this.$axios.post(`/groups/${this.id}/members`, { 
            username: username.trim() 
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

    async leaveGroup() {
      if (confirm('Are you sure you want to leave this group?')) {
        try {
          await this.$axios.delete(`/groups/${this.id}/members/me`);
          
          // Redirect back to chat view or home
          this.$router.push('/');
        } catch (error) {
          console.error('Error leaving group:', error);
          const errorMessage = error.response?.data?.message || error.message || 'Unknown error';
          alert(`Failed to leave group: ${errorMessage}`);
        }
      }
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
.conversation-info-section {
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
    margin-bottom: 1rem;
}

.group-picture {
    width: 80px;
    height: 80px;
    border-radius: 50%;
    object-fit: cover;
    border: 2px solid #ddd;
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
    padding: 0.5rem;
    background-color: #f8f9fa;
    border-radius: 8px;
}

.participant-avatar {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    object-fit: cover;
    border: 1px solid #ddd;
}

.group-actions {
    display: flex;
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
}

.action-btn:hover {
    background-color: #0056b3;
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
</style>
