<template>
  <div class="modal-backdrop" style="position:fixed;top:0;left:0;width:100vw;height:100vh;background:rgba(0,0,0,0.3);z-index:2000;display:flex;align-items:center;justify-content:center;">
    <div class="modal-dialog" style="background:white;border-radius:8px;box-shadow:0 2px 16px rgba(0,0,0,0.15);max-width:500px;width:100%;">
      <div class="modal-content" style="padding:1.5rem;">
        <div class="modal-header" style="display:flex;align-items:center;justify-content:space-between;">
          <h5 class="modal-title">Forward Message</h5>
          <button type="button" class="close" style="background:none;border:none;font-size:1.5rem;" @click="$emit('cancel')">&times;</button>
        </div>
        <div class="modal-body">
          <!-- Tab Navigation -->
          <div class="tab-nav" style="display:flex;margin-bottom:1rem;border-bottom:1px solid #dee2e6;">
            <button 
              class="tab-btn" 
              :class="{ active: activeTab === 'conversations' }"
              style="flex:1;padding:0.5rem;border:none;background:none;border-bottom:2px solid transparent;"
              @click="activeTab = 'conversations'"
            >
              Existing Chats
            </button>
            <button 
              class="tab-btn" 
              :class="{ active: activeTab === 'search' }"
              style="flex:1;padding:0.5rem;border:none;background:none;border-bottom:2px solid transparent;"
              @click="activeTab = 'search'"
            >
              New Chat
            </button>
          </div>

          <!-- Existing Conversations Tab -->
          <div v-if="activeTab === 'conversations'">
            <p>Select a conversation to forward this message to:</p>
            <div v-if="conversations.length === 0" style="text-align:center;padding:2rem;color:#6c757d;">
              <p>No existing conversations</p>
              <p><small>Use the "New Chat" tab to forward to a new user</small></p>
            </div>
            <ul v-else class="list-group" style="padding:0;list-style:none;">
              <li v-for="conv in conversations" :key="conv.id" class="list-group-item" style="margin-bottom:0.5rem;display:flex;align-items:center;justify-content:space-between;padding:0.5rem;border:1px solid #dee2e6;border-radius:4px;">
                <span style="font-weight:500;">{{ conv.name }}</span>
                <button class="btn btn-sm btn-primary" style="padding:0.25rem 0.75rem;" @click="$emit('confirm', conv.id)">
                  Forward
                </button>
              </li>
            </ul>
          </div>

          <!-- User Search Tab -->
          <div v-if="activeTab === 'search'">
            <p>Search for a user to start a new conversation:</p>
            
            <!-- Search Input -->
            <div class="search-input-container" style="position:relative;margin-bottom:1rem;">
              <input 
                v-model="searchQuery" 
                type="text"
                class="form-control" 
                placeholder="Search for users..." 
                style="padding-right:2.5rem;"
                @input="onSearch"
              >
              <svg class="search-icon feather" style="position:absolute;right:12px;top:50%;transform:translateY(-50%);color:#6c757d;width:16px;height:16px;">
                <use href="/feather-sprite-v4.29.0.svg#search" />
              </svg>
            </div>

            <!-- Loading State -->
            <div v-if="searching" class="text-center py-3">
              <div class="spinner-border spinner-border-sm" role="status">
                <span class="visually-hidden">Loading...</span>
              </div>
            </div>

            <!-- Search Results -->
            <div v-else-if="searchResults.length" class="search-results" style="max-height:200px;overflow-y:auto;">
              <div
                v-for="user in searchResults" 
                :key="user.id" 
                class="user-result"
                style="display:flex;align-items:center;padding:12px;border-radius:8px;transition:background-color 0.2s;margin-bottom:4px;border:1px solid #e9ecef;"
              >
                <img :src="getImageUrl(user.photoUrl)" :alt="user.username" class="user-avatar" style="width:40px;height:40px;border-radius:50%;object-fit:cover;margin-right:12px;">
                <div class="user-info" style="flex:1;">
                  <div class="user-name" style="font-weight:500;margin-bottom:2px;">{{ user.username }}</div>
                </div>
                <button 
                  type="button" 
                  class="btn btn-primary btn-sm"
                  @click="forwardToNewUser(user)"
                >
                  <svg class="feather" style="width:16px;height:16px;margin-right:4px;"><use href="/feather-sprite-v4.29.0.svg#message-circle" /></svg>
                  Forward
                </button>
              </div>
            </div>

            <!-- Empty State -->
            <div v-else-if="searchQuery && !searching" class="empty-results" style="text-align:center;padding:2rem;color:#6c757d;">
              <svg class="feather empty-icon" style="width:48px;height:48px;margin-bottom:1rem;color:#dee2e6;"><use href="/feather-sprite-v4.29.0.svg#users" /></svg>
              <p>No users found</p>
              <small class="text-muted">Try a different search term</small>
            </div>
            
            <!-- Initial State -->
            <div v-else class="initial-state" style="text-align:center;padding:2rem;color:#6c757d;">
              <svg class="feather search-large" style="width:48px;height:48px;margin-bottom:1rem;color:#dee2e6;"><use href="/feather-sprite-v4.29.0.svg#search" /></svg>
              <p>Search for users to forward this message to</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from '../../services/axios.js';
import { getImageUrl } from '../../utils/imageUtils.js';

export default {
  name: 'ForwardMessageModal',
  props: {
    conversations: {
      type: Array,
      required: true
    }
  },
  emits: ['cancel', 'confirm', 'forward-to-user'],
  data() {
    return {
      activeTab: 'conversations',
      searchQuery: '',
      searchResults: [],
      searching: false,
      searchTimeout: null
    }
  },
  beforeUnmount() {
    if (this.searchTimeout) {
      clearTimeout(this.searchTimeout);
    }
  },
  methods: {
    getImageUrl,
    
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
      
      this.searching = true;
      
      try {
        const response = await axios.get('/users', {
          params: {
            q: this.searchQuery.trim()
          }
        });
        
        if (response.data && response.data.users) {
          this.searchResults = response.data.users;
        } else {
          this.searchResults = [];
        }
          
      } catch (error) {
        console.error('Search error:', error);
        this.searchResults = [];
      } finally {
        this.searching = false;
      }
    },
    
    forwardToNewUser(user) {
      this.$emit('forward-to-user', user);
    }
  }
}
</script>

<style scoped>
.tab-btn {
  cursor: pointer;
  transition: all 0.2s;
}

.tab-btn:hover {
  background-color: #f8f9fa;
}

.tab-btn.active {
  border-bottom-color: #007bff !important;
  color: #007bff;
  font-weight: 500;
}

.user-result:hover {
  background-color: #f8f9fa;
}

.search-input-container .form-control {
  border: 1px solid #ced4da;
  border-radius: 0.375rem;
  padding: 0.375rem 0.75rem;
  font-size: 0.875rem;
}

.search-input-container .form-control:focus {
  border-color: #86b7fe;
  outline: 0;
  box-shadow: 0 0 0 0.25rem rgba(13, 110, 253, 0.25);
}

.spinner-border-sm {
  width: 1rem;
  height: 1rem;
  border-width: 0.125em;
}

.spinner-border {
  display: inline-block;
  vertical-align: text-bottom;
  border: 0.25em solid currentColor;
  border-right-color: transparent;
  border-radius: 50%;
  animation: spinner-border 0.75s linear infinite;
}

@keyframes spinner-border {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.btn {
  display: inline-block;
  font-weight: 400;
  line-height: 1.5;
  color: #212529;
  text-align: center;
  text-decoration: none;
  vertical-align: middle;
  cursor: pointer;
  user-select: none;
  background-color: transparent;
  border: 1px solid transparent;
  padding: 0.375rem 0.75rem;
  font-size: 0.875rem;
  border-radius: 0.375rem;
  transition: color 0.15s ease-in-out, background-color 0.15s ease-in-out, border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out;
}

.btn-primary {
  color: #fff;
  background-color: #0d6efd;
  border-color: #0d6efd;
}

.btn-primary:hover {
  color: #fff;
  background-color: #0b5ed7;
  border-color: #0a58ca;
}

.btn-sm {
  padding: 0.25rem 0.5rem;
  font-size: 0.8125rem;
  border-radius: 0.25rem;
}
</style>
