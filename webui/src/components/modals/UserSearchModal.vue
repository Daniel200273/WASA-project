<template>
  <div class="user-search-modal modal d-block" tabindex="-1">
    <div class="modal-dialog modal-dialog-centered">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">Find Users</h5>
          <button type="button" class="btn-close" @click="$emit('close')"></button>
        </div>
        
        <div class="modal-body">
          <div class="search-input-container">
            <input 
              v-model="searchQuery" 
              @input="onSearch"
              type="text" 
              class="form-control" 
              placeholder="Search for users..."
              autofocus>
            <svg class="search-icon feather"><use href="/feather-sprite-v4.29.0.svg#search"/></svg>
          </div>
          
          <!-- Loading State -->
          <div v-if="loading" class="text-center py-3">
            <LoadingSpinner :loading="true" />
          </div>
          
          <!-- Search Results -->
          <div v-else-if="searchResults.length" class="search-results">
            <div v-for="user in searchResults" 
                 :key="user.id" 
                 class="user-result" 
                 @click="selectUser(user)">
              <img :src="user.photo || '/default-avatar.svg'" :alt="user.name" class="user-avatar">
              <div class="user-info">
                <div class="user-name">{{ user.name }}</div>
              </div>
              <svg class="feather select-arrow"><use href="/feather-sprite-v4.29.0.svg#chevron-right"/></svg>
            </div>
          </div>
          
          <!-- Empty State -->
          <div v-else-if="searchQuery && !loading" class="empty-results">
            <svg class="feather empty-icon"><use href="/feather-sprite-v4.29.0.svg#users"/></svg>
            <p>No users found</p>
            <small class="text-muted">Try a different search term</small>
          </div>
          
          <!-- Initial State -->
          <div v-else class="initial-state">
            <svg class="feather search-large"><use href="/feather-sprite-v4.29.0.svg#search"/></svg>
            <p>Search for users to start a conversation</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import LoadingSpinner from '@/components/LoadingSpinner.vue'

export default {
  name: 'UserSearchModal',
  components: {
    LoadingSpinner
  },
  emits: ['close', 'select-user'],
  data() {
    return {
      searchQuery: '',
      searchResults: [],
      loading: false,
      searchTimeout: null
    }
  },
  methods: {
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
        // TODO: Replace with actual API call
        // const response = await this.$axios.get(`/users?query=${encodeURIComponent(this.searchQuery)}`);
        // this.searchResults = response.data;
        
        // Mock search results for now
        await new Promise(resolve => setTimeout(resolve, 500)); // Simulate API delay
        this.searchResults = [
          { id: '1', name: 'Alice Johnson', photo: null },
          { id: '2', name: 'Bob Smith', photo: null },
          { id: '3', name: 'Charlie Brown', photo: null }
        ].filter(user => user.name.toLowerCase().includes(this.searchQuery.toLowerCase()));
        
      } catch (error) {
        console.error('Search error:', error);
        this.searchResults = [];
      } finally {
        this.loading = false;
      }
    },
    
    selectUser(user) {
      this.$emit('select-user', user);
    }
  },
  
  beforeUnmount() {
    if (this.searchTimeout) {
      clearTimeout(this.searchTimeout);
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
  cursor: pointer;
  border-radius: 8px;
  transition: background-color 0.2s;
  margin-bottom: 4px;
}

.user-result:hover {
  background-color: #f8f9fa;
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
