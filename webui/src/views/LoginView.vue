<!-- LoginView.vue -->
<template>
  <div class="login">
    <h2>Enter Your Name</h2>
    <input v-model="username" placeholder="Username">
    <button @click="doLogin">Log in</button>
  </div>
</template>
  
<script>
import AuthService from '../services/auth.js';

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
                
                // Store the token in sessionStorage
                AuthService.setAuthData(response.data.identifier, this.username);
                
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

<style scoped>
.login {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 400px;
  text-align: center;
}

.login h2 {
  margin-bottom: 1.5rem;
  color: #333;
  font-weight: 300;
}

.login input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  margin-bottom: 1rem;
  box-sizing: border-box;
}

.login input:focus {
  outline: none;
  border-color: #007bff;
  box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
}

.login button {
  width: 100%;
  padding: 0.75rem;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.login button:hover {
  background-color: #0056b3;
}

.login button:disabled {
  background-color: #6c757d;
  cursor: not-allowed;
}
</style>