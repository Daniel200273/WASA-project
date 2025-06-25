package api

import (
	"net/http"
	"strings"

	"github.com/Daniel200273/WASA-project/service/database"
)

/*
isAuthorized checks if the user is authorized to perform the action, by checking the Authorization header.
The auth token must be in the format "Bearer <sessionToken>" where sessionToken is the string token
returned from the login endpoint (/session).
If the user is authorized, the function returns the userID, otherwise it returns an empty string.
*/
func isAuthorized(header http.Header, db database.AppDatabase) string {
	authHeader := header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Check if the header starts with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return ""
	}

	// Extract the token part after "Bearer "
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return ""
	}

	// Look up the user by session token
	user, err := db.GetUserByToken(token)
	if err != nil {
		// Token is invalid or expired
		return ""
	}

	return user.ID
}
