  <!-- LoginView.vue -->
<template>
    <div class="login">
      <h2>Enter Your Name</h2>
      <input v-model="username" placeholder="Username" />
      <button @click="doLogin">Log in</button>
    </div>
</template>
  
<script>
    export default {
    data() {
        return { username: '' }
    },
    methods: {
        async doLogin() {
            // Check if username is provided
            if (!this.username) {
                alert('Please enter a username');
                return;
            }
            // Make a POST request to the server to create a session
            try {
                let response = await this.$axios.post('/session', {
                    name: this.username  // API expects 'name' not 'username'
                });
                
                // Store the token in sessionStorage (separate per tab for testing)
                sessionStorage.setItem('authToken', response.data.identifier);
                sessionStorage.setItem('username', this.username);
                
                console.log('Login successful!', response.data);
                
                // Redirect to home page
                this.$router.push('/');
                
            } catch (error) {
                console.error('Login failed:', error);
                alert('Login failed. Please try again.');
            }
        }
    }
}
</script>