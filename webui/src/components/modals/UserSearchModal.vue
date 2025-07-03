<template>
  <div class="user-search-modal modal d-block" tabindex="-1">
    <div class="modal-dialog modal-dialog-centered">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">{{ modalTitle }}</h5>
          <button type="button" class="btn-close" @click="$emit('close')" />
        </div>
        
        <div class="modal-body">
          <!-- Group Name Input (only in group-create mode) -->
          <div v-if="mode === 'group-create'" class="mb-3">
            <label for="groupName" class="form-label">Group Name</label>
            <input 
              id="groupName"
              v-model="groupName" 
              type="text"
              class="form-control" 
              placeholder="Enter group name..." 
              required
            >
          </div>
          
          <!-- Selected Users (only in group-create mode) -->
          <div v-if="mode === 'group-create' && selectedUsers.length" class="selected-users mb-3">
            <label class="form-label">Selected Members ({{ selectedUsers.length }})</label>
            <div class="selected-users-list">
              <div
                v-for="user in selectedUsers" 
                :key="user.id"
                class="selected-user-chip"
              >
                <img :src="getImageUrl(user.photoUrl)" :alt="user.username" class="chip-avatar">
                <span>{{ user.username }}</span>
                <button type="button" class="btn-remove" @click="removeUser(user)">×</button>
              </div>
            </div>
          </div>
          
          <div class="search-input-container">
            <input 
              v-model="searchQuery" 
              type="text"
              class="form-control" 
              :placeholder="searchPlaceholder" 
              autofocus
              @input="onSearch"
            >
            <svg class="search-icon feather"><use href="/feather-sprite-v4.29.0.svg#search" /></svg>
          </div>
          
          <!-- Loading State -->
          <div v-if="loading" class="text-center py-3">
            <LoadingSpinner :loading="true" />
          </div>
          
          <!-- Search Results -->
          <div v-else-if="searchResults.length" class="search-results">
            <div
              v-for="user in searchResults" 
              :key="user.id" 
              class="user-result"
              :class="{ 'selected': isUserSelected(user) }"
            >
              <img :src="getImageUrl(user.photoUrl)" :alt="user.username" class="user-avatar">
              <div class="user-info">
                <div class="user-name">{{ user.username }}</div>
              </div>
              
              <!-- User Search Mode Actions -->
              <div v-if="mode === 'user-search'" class="user-actions">
                <button 
                  type="button" 
                  class="btn btn-outline-primary btn-sm me-2"
                  @click="viewProfile(user)"
                >
                  <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#user" /></svg>
                  Profile
                </button>
                <button 
                  type="button" 
                  class="btn btn-primary btn-sm"
                  @click="sendMessage(user)"
                >
                  <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#message-circle" /></svg>
                  Message
                </button>
              </div>
              
              <!-- Group Create Mode Actions -->
              <div v-else-if="mode === 'group-create'" class="group-actions">
                <button 
                  type="button" 
                  class="btn btn-outline-primary btn-sm"
                  @click="toggleUserSelection(user)"
                >
                  <svg v-if="isUserSelected(user)" class="feather"><use href="/feather-sprite-v4.29.0.svg#check" /></svg>
                  <svg v-else class="feather"><use href="/feather-sprite-v4.29.0.svg#plus" /></svg>
                  {{ isUserSelected(user) ? 'Selected' : 'Add' }}
                </button>
              </div>
            </div>
          </div>
          
          <!-- Empty State -->
          <div v-else-if="searchQuery && !loading" class="empty-results">
            <svg class="feather empty-icon"><use href="/feather-sprite-v4.29.0.svg#users" /></svg>
            <p>No users found</p>
            <small class="text-muted">Try a different search term</small>
          </div>
          
          <!-- Initial State -->
          <div v-else class="initial-state">
            <svg class="feather search-large"><use href="/feather-sprite-v4.29.0.svg#search" /></svg>
            <p>{{ initialStateText }}</p>
          </div>
          
          <!-- Group Create Actions -->
          <div v-if="mode === 'group-create' && selectedUsers.length" class="group-create-actions mt-3">
            <button 
              type="button" 
              class="btn btn-success w-100"
              :disabled="!groupName.trim() || selectedUsers.length === 0"
              @click="createGroup"
            >
              <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#users" /></svg>
              Create Group ({{ selectedUsers.length }} members)
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import axios from '@/services/axios.js'
import AuthService from '@/services/auth.js'
import { useRouter } from 'vue-router'

export default {
  name: 'UserSearchModal',
  components: {
    LoadingSpinner
  },
  props: {
    mode: {
      type: String,
      default: 'user-search',
      validator: (value) => ['user-search', 'group-create'].includes(value)
    }
  },
  emits: ['close', 'select-user', 'group-created'],
  setup() {
    const router = useRouter()
    return { router }
  },
  data() {
    return {
      searchQuery: '',
      searchResults: [],
      loading: false,
      searchTimeout: null,
      selectedUsers: [],
      groupName: ''
    }
  },
  
  computed: {
    modalTitle() {
      return this.mode === 'group-create' ? 'Create Group' : 'Find Users'
    },
    searchPlaceholder() {
      return this.mode === 'group-create' 
        ? 'Search for users to add to group...' 
        : 'Search for users...'
    },
    initialStateText() {
      return this.mode === 'group-create'
        ? 'Search for users to add to your group'
        : 'Search for users to message or view profiles'
    }
  },
  
  beforeUnmount() {
    if (this.searchTimeout) {
      clearTimeout(this.searchTimeout);
    }
  },
  methods: {
    getImageUrl(photoUrl) {
      if (!photoUrl) return '/default-avatar.svg';
      
      // photoUrl comes as "/uploads/profiles/filename.jpg" from backend
      // We need to prepend the API base URL and add cache busting
      const baseURL = axios.defaults.baseURL || '';
      const timestamp = Date.now(); // Cache busting parameter
      return `${baseURL}${photoUrl}?t=${timestamp}`;
    },
    
    onSearch() {
      // Clear previous timeout
      if (this.searchTimeout) {
        clearTimeout(this.searchTimeout);
      }
      
      // Debounce search
      this.searchTimeout = setTimeout(() => {
        this.performSearch();
      }, 300);
    },
     async performSearch() {
      if (!this.searchQuery.trim()) {
        this.searchResults = [];
        return;
      }
      
      this.loading = true;
      
      try {
        const response = await axios.get('/users', {
          params: {
            q: this.searchQuery.trim()
          }
        });
        
        if (response.data && response.data.users) {
          // Filter out already selected users in group create mode
          this.searchResults = this.mode === 'group-create' 
            ? response.data.users.filter(user => !this.isUserSelected(user))
            : response.data.users;
        } else {
          this.searchResults = [];
        }
          
      } catch (error) {
        console.error('Search error:', error);
        this.searchResults = [];
      } finally {
        this.loading = false;
      }
    },
    
    // User Search Mode Methods
    sendMessage(user) {
      this.$emit('select-user', user);
    },
    
    viewProfile(user) {
      if (!AuthService.isAuthenticated()) {
        this.router.push('/login');
        return;
      }
      
      if (!user.id) {
        alert('Cannot view profile: User ID not found');
        return;
      }
      
      // Navigate to user profile
      this.router.push(`/profile?type=user&id=${user.id}`);
      this.$emit('close');
    },
    
    // Group Create Mode Methods
    toggleUserSelection(user) {
      if (this.isUserSelected(user)) {
        this.removeUser(user);
      } else {
        this.selectedUsers.push(user);
        // Remove from search results to avoid confusion
        this.searchResults = this.searchResults.filter(u => u.id !== user.id);
      }
    },
    
    removeUser(user) {
      this.selectedUsers = this.selectedUsers.filter(u => u.id !== user.id);
      // Add back to search results if they match current query
      if (this.searchQuery && user.username.toLowerCase().includes(this.searchQuery.toLowerCase())) {
        this.searchResults.push(user);
      }
    },
    
    isUserSelected(user) {
      return this.selectedUsers.some(u => u.id === user.id);
    },
    
    async createGroup() {
      if (!this.groupName.trim() || this.selectedUsers.length === 0) {
        return;
      }
      
      try {
        const currentUserId = AuthService.getUserId();
        
        if (!currentUserId) {
          console.error('No authenticated user found');
          return;
        }
        
        const groupData = {
          name: this.groupName.trim(),
          members: this.selectedUsers.map(user => user.id)
        };
        
        const response = await axios.post(`/users/${currentUserId}/groups`, groupData);
        
        this.$emit('group-created', response.data);
      } catch (error) {
        console.error('Group creation error:', error);
        // You might want to show an error message to the user
      }
    },
    
    selectUser(user) {
      this.$emit('select-user', user);
    }
  }
}
</script>

<style scoped>
.user-search-modal {
  background-color: rgba(0, 0, 0, 0.5);
}

.search-input-container {
  position: relative;
  margin-bottom: 1rem;
}

.search-icon {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: #6c757d;
  width: 16px;
  height: 16px;
  pointer-events: none;
}

.search-results {
  max-height: 300px;
  overflow-y: auto;
}

.user-result {
  display: flex;
  align-items: center;
  padding: 12px;
  border-radius: 8px;
  transition: background-color 0.2s;
  margin-bottom: 4px;
}

.user-result:hover {
  background-color: #f8f9fa;
}

.user-result.selected {
  background-color: #e7f3ff;
  border: 1px solid #0d6efd;
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  margin-right: 12px;
}

.user-info {
  flex: 1;
}

.user-name {
  font-weight: 500;
  margin-bottom: 2px;
}

.user-actions, .group-actions {
  display: flex;
  gap: 8px;
}

.user-actions .btn, .group-actions .btn {
  display: flex;
  align-items: center;
  gap: 4px;
}

.selected-users {
  border: 1px solid #dee2e6;
  border-radius: 8px;
  padding: 12px;
  background-color: #f8f9fa;
}

.selected-users-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 8px;
}

.selected-user-chip {
  display: flex;
  align-items: center;
  background-color: #0d6efd;
  color: white;
  border-radius: 20px;
  padding: 4px 8px 4px 4px;
  font-size: 0.875rem;
  gap: 6px;
}

.chip-avatar {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  object-fit: cover;
}

.btn-remove {
  background: none;
  border: none;
  color: white;
  font-size: 16px;
  line-height: 1;
  padding: 0;
  margin-left: 4px;
  cursor: pointer;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: background-color 0.2s;
}

.btn-remove:hover {
  background-color: rgba(255, 255, 255, 0.2);
}

.group-create-actions {
  border-top: 1px solid #dee2e6;
  padding-top: 1rem;
}

.select-arrow {
  width: 16px;
  height: 16px;
  color: #6c757d;
}

.empty-results, .initial-state {
  text-align: center;
  padding: 2rem;
  color: #6c757d;
}

.empty-icon, .search-large {
  width: 48px;
  height: 48px;
  margin-bottom: 1rem;
  color: #dee2e6;
}

.empty-results p, .initial-state p {
  margin-bottom: 0.5rem;
  font-weight: 500;
}

.feather {
  width: 16px;
  height: 16px;
}
</style>
