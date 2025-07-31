// Color utility for generating consistent colors for usernames
export function getUsernameColor(username) {
  if (!username) return "#6c757d";

  // Predefined set of colors that look good and are accessible
  const colors = [
    "#e74c3c", // Red
    "#3498db", // Blue
    "#2ecc71", // Green
    "#f39c12", // Orange
    "#9b59b6", // Purple
    "#1abc9c", // Turquoise
    "#e67e22", // Carrot
    "#34495e", // Wet Asphalt
    "#f1c40f", // Sun Flower
    "#16a085", // Green Sea
    "#27ae60", // Nephritis
    "#2980b9", // Belize Hole
    "#8e44ad", // Wisteria
    "#d35400", // Pumpkin
    "#c0392b", // Pomegranate
    "#7f8c8d", // Asbestos
  ];

  // Simple hash function to get consistent color for username
  let hash = 0;
  for (let i = 0; i < username.length; i++) {
    const char = username.charCodeAt(i);
    hash = (hash << 5) - hash + char;
    hash = hash & hash; // Convert to 32-bit integer
  }

  // Get color index from hash
  const colorIndex = Math.abs(hash) % colors.length;
  return colors[colorIndex];
}
