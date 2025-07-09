import axios from "../services/axios.js";

/**
 * Converts a relative photo URL to an absolute URL with cache busting
 * @param {string} photoUrl - The relative photo URL from the backend (e.g., "/uploads/profiles/user123.jpg")
 * @param {string} defaultImage - The default image to use if photoUrl is null/undefined
 * @returns {string} The absolute URL with cache busting parameter
 */
export function getImageUrl(photoUrl, defaultImage = "/default-avatar.svg") {
  // Return default image if no photoUrl provided
  if (!photoUrl) return defaultImage;

  // If it's already a full URL (starts with http/https), return as-is
  if (photoUrl.startsWith("http")) return photoUrl;

  // Build full URL with cache busting
  const baseURL = axios.defaults.baseURL || "http://localhost:3000";
  const timestamp = Date.now();
  return `${baseURL}${photoUrl}?t=${timestamp}`;
}

/**
 * Gets the appropriate avatar image URL for a conversation
 * @param {object} conversation - The conversation object
 * @returns {string} The avatar URL
 */
export function getConversationAvatar(conversation) {
  if (conversation?.photoUrl) {
    return getImageUrl(conversation.photoUrl);
  }
  return conversation?.type === "group"
    ? "/default-group.svg"
    : "/default-avatar.svg";
}

/**
 * Gets the appropriate avatar image URL for a user
 * @param {object} user - The user object
 * @returns {string} The avatar URL
 */
export function getUserAvatar(user) {
  return getImageUrl(user?.photoUrl, "/default-avatar.svg");
}

/**
 * Gets the appropriate avatar image URL for a group
 * @param {object} group - The group object
 * @returns {string} The avatar URL
 */
export function getGroupAvatar(group) {
  return getImageUrl(group?.photoUrl, "/default-group.svg");
}
