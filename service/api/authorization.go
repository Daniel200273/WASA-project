package api

import (
	"net/http"
	"strings"
)

/*
isAuthorized checks if the user is authorized to perform the action, by checking the Authorization header.
The auth token must be in the format "Bearer <userIdentifier>" where userIdentifier is the string token
returned from the login endpoint (/session).
If the user is authorized, the function returns the userIdentifier, otherwise it returns an empty string.
*/
func isAuthorized(header http.Header) string {
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

	return token
}
