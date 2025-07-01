import axios from "axios";
import AuthService from "./auth.js";

const instance = axios.create({
  baseURL: __API_URL__,
  timeout: 1000 * 5,
});

// Add authentication header to all requests
instance.interceptors.request.use((config) => {
  // Get token from AuthService
  const token = AuthService.getAuthToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle authentication errors
instance.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Token expired or invalid, clear session
      AuthService.clearAuthData();
      window.location.href = "/login";
    }
    return Promise.reject(error);
  }
);

export default instance;
