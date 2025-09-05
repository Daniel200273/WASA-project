<script setup>
import { RouterLink, RouterView, useRoute, useRouter } from 'vue-router'
import { computed, ref } from 'vue'
import AuthService from './services/auth.js'
import UserSearchModal from './components/modals/UserSearchModal.vue'
import axios from './services/axios.js'

const route = useRoute()
const router = useRouter()

// Check if current route is login page
const isLoginPage = computed(() => route.name === 'Login')

// Modal states
const showUserSearch = ref(false)
const showGroupCreate = ref(false)

// Modal handlers
const handleUserSelect = async (user) => {
  try {
    // Close the modal first
    showUserSearch.value = false
    
    const userId = AuthService.getUserId();
    
    // Create or get conversation with the selected user
    const response = await axios.post(`/users/${userId}/conversations`, {
      userId: user.id
    });
    
    // Navigate to the conversation in chat view
    router.push(`/chat/${response.data.id}`)
    
  } catch (error) {
    console.error('Error starting conversation:', error)
    showUserSearch.value = false
  }
}

const handleGroupCreate = (groupData) => {
  // Handle group creation - redirect to the new group's info page
  showGroupCreate.value = false
  
  // Redirect to group info view
  if (groupData && groupData.id) {
    router.push(`/profile?type=group&id=${groupData.id}`)
  }
}

const closeUserSearch = () => {
  showUserSearch.value = false
}

const closeGroupCreate = () => {
  showGroupCreate.value = false
}

// Function to close mobile sidebar menu
const closeMobileSidebar = () => {
  // Check if we're on mobile and the sidebar is open
  const sidebarElement = document.getElementById('sidebarMenu')
  if (sidebarElement && sidebarElement.classList.contains('show')) {
    // Remove the 'show' class to close the sidebar
    sidebarElement.classList.remove('show')
    
    // Also update the navbar toggler button state
    const togglerButton = document.querySelector('.navbar-toggler')
    if (togglerButton) {
      togglerButton.classList.add('collapsed')
      togglerButton.setAttribute('aria-expanded', 'false')
    }
  }
}
</script>

<template>
  <!-- Show full layout only when not on login page -->
  <div v-if="!isLoginPage">
    <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
      <a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WASAText</a>
      <button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
        <span class="navbar-toggler-icon" />
      </button>
    </header>

    <div class="container-fluid">
      <div class="row">
        <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
          <div class="position-sticky pt-3 sidebar-sticky">
            <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
              <span>General</span>
            </h6>
            <ul class="nav flex-column">
              <li class="nav-item">
                <RouterLink to="/chat" class="nav-link" @click="closeMobileSidebar">
                  <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#message-circle" /></svg>
                  Messages
                </RouterLink>
              </li>
              <li class="nav-item">
                <RouterLink to="/profile?type=user&id=me" class="nav-link" @click="closeMobileSidebar">
                  <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#user" /></svg>
                  My Profile
                </RouterLink>
              </li>
            </ul>

            <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
              <span>Quick Actions</span>
            </h6>
            <ul class="nav flex-column">
              <li class="nav-item">
                <a href="#" class="nav-link" @click.prevent="showUserSearch = true; closeMobileSidebar()">
                  <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#search" /></svg>
                  Find Users
                </a>
              </li>
              <li class="nav-item">
                <a href="#" class="nav-link" @click.prevent="showGroupCreate = true; closeMobileSidebar()">
                  <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#users" /></svg>
                  Create Group
                </a>
              </li>
            </ul>
          </div>
        </nav>

        <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
          <RouterView />
        </main>
      </div>
    </div>
  </div>

  <!-- Show only RouterView for login page (clean, no sidebar/header) -->
  <div v-else class="login-container">
    <RouterView />
  </div>

  <!-- Modals -->
  <UserSearchModal 
    v-if="showUserSearch" 
    mode="user-search"
    @close="closeUserSearch"
    @select-user="handleUserSelect"
  />
  
  <UserSearchModal 
    v-if="showGroupCreate" 
    mode="group-create"
    @close="closeGroupCreate"
    @group-created="handleGroupCreate"
  />
</template>

<style>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f8f9fa;
}

/* Modal backdrop */
.user-search-modal {
  background-color: rgba(0, 0, 0, 0.5);
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1050;
}
</style>
