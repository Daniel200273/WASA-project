<template>
  <div class="home-view">
    <h2>Welcome to WASAText</h2>
    <p>Your messaging application is ready!</p>
    <div class="features">
      <h3>Features:</h3>
      <ul>
        <li>User authentication ✅</li>
        <li>Real-time messaging (coming soon)</li>
        <li>Group chats (coming soon)</li>
        <li>File sharing (coming soon)</li>
      </ul>
    </div>
  </div>
</template>

<script>
export default {
  name: 'HomeView',
  data() {
    return {
      // Add data properties here when needed
    }
  },
  methods: {
    // Add methods here when needed
  }
}
</script>

<style scoped>
.home-view {
  padding: 20px;
}

.features {
  margin-top: 30px;
}

.features ul {
  text-align: left;
  display: inline-block;
}
</style>
