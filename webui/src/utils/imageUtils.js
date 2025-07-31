import axios from "../services/axios.js";

// Cache for image URLs to prevent constant reloading
const imageUrlCache = new Map();

/**
 * Converts a relative photo URL to an absolute URL with cache busting
 * @param {string} photoUrl - The relative photo URL from the backend (e.g., "/uploads/profiles/user123.jpg")
 * @param {string} defaultImage - The default image to use if photoUrl is null/undefined
 * @param {boolean} forceRefresh - Force refresh the cache (e.g., after uploading new image)
 * @returns {string} The absolute URL with cache busting parameter
 */
export function getImageUrl(
  photoUrl,
  defaultImage = "/default-avatar.svg",
  forceRefresh = false
) {
  // Return default image if no photoUrl provided
  if (!photoUrl) return defaultImage;

  // If it's already a full URL (starts with http/https), return as-is
  if (photoUrl.startsWith("http")) return photoUrl;

  // Check cache first unless force refresh is requested
  if (!forceRefresh && imageUrlCache.has(photoUrl)) {
    return imageUrlCache.get(photoUrl);
  }

  // Build full URL with cache busting only when needed
  const baseURL = axios.defaults.baseURL || window.location.origin;
  const timestamp = forceRefresh
    ? Date.now()
    : Math.floor(Date.now() / (5 * 60 * 1000)); // 5-minute cache window
  const cachedUrl = `${baseURL}${photoUrl}?t=${timestamp}`;

  // Store in cache
  imageUrlCache.set(photoUrl, cachedUrl);

  return cachedUrl;
}

/**
 * Clears the image cache for a specific photo URL (useful after uploading new image)
 * @param {string} photoUrl - The photo URL to clear from cache
 */
export function clearImageCache(photoUrl) {
  if (photoUrl && imageUrlCache.has(photoUrl)) {
    imageUrlCache.delete(photoUrl);
  }
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
