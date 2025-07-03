<script setup>
import { RouterLink, RouterView, useRoute, useRouter } from 'vue-router'
import { computed, ref } from 'vue'
import AuthService from './services/auth.js'
import UserSearchModal from './components/modals/UserSearchModal.vue'

const route = useRoute()
const router = useRouter()

// Check if current route is login page
const isLoginPage = computed(() => route.name === 'Login')

// Modal states
const showUserSearch = ref(false)
const showGroupCreate = ref(false)

// Modal handlers
const handleUserSelect = (user) => {
  // Handle user selection (e.g., start conversation)
  console.log('Selected user:', user)
  showUserSearch.value = false
}

const handleGroupCreate = (groupData) => {
  // Handle group creation - redirect to the new group's info page
  console.log('Created group:', groupData)
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
</script>
<script>
export default {}
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
                <RouterLink to="/chat" class="nav-link">
                  <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#message-circle" /></svg>
                  Messages
                </RouterLink>
              </li>
              <li class="nav-item">
                <RouterLink to="/profile?type=user&id=me" class="nav-link">
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
                <a href="#" class="nav-link" @click.prevent="showUserSearch = true">
                  <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#search" /></svg>
                  Find Users
                </a>
              </li>
              <li class="nav-item">
                <a href="#" class="nav-link" @click.prevent="showGroupCreate = true">
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
