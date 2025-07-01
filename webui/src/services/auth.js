// Authentication service for managing user authentication state
// Uses sessionStorage for temporary authentication (cleared when browser closes)

export class AuthService {
  // Keys for sessionStorage
  static AUTH_TOKEN_KEY = "authToken";
  static USERNAME_KEY = "username";
  static USER_ID_KEY = "userId";

  /**
   * Check if user is currently authenticated
   * @returns {boolean}
   */
  static isAuthenticated() {
    return !!sessionStorage.getItem(AuthService.AUTH_TOKEN_KEY);
  }

  /**
   * Get the current authentication token
   * @returns {string|null}
   */
  static getAuthToken() {
    return sessionStorage.getItem(AuthService.AUTH_TOKEN_KEY);
  }

  /**
   * Get the current username
   * @returns {string|null}
   */
  static getUsername() {
    return sessionStorage.getItem(AuthService.USERNAME_KEY);
  }

  /**
   * Get the current user ID
   * @returns {string|null}
   */
  static getUserId() {
    return sessionStorage.getItem(AuthService.USER_ID_KEY);
  }

  /**
   * Store authentication data
   * @param {string} token - Authentication token
   * @param {string} username - Username
   * @param {string} userId - User ID (optional)
   */
  static setAuthData(token, username, userId = null) {
    sessionStorage.setItem(AuthService.AUTH_TOKEN_KEY, token);
    sessionStorage.setItem(AuthService.USERNAME_KEY, username);
    if (userId) {
      sessionStorage.setItem(AuthService.USER_ID_KEY, userId);
    }
  }

  /**
   * Clear all authentication data (logout)
   */
  static clearAuthData() {
    sessionStorage.removeItem(AuthService.AUTH_TOKEN_KEY);
    sessionStorage.removeItem(AuthService.USERNAME_KEY);
    sessionStorage.removeItem(AuthService.USER_ID_KEY);
  }

  /**
   * Get current user info object
   * @returns {object|null}
   */
  static getCurrentUser() {
    if (!AuthService.isAuthenticated()) {
      return null;
    }

    return {
      token: AuthService.getAuthToken(),
      username: AuthService.getUsername(),
      userId: AuthService.getUserId(),
    };
  }
}

export default AuthService;
