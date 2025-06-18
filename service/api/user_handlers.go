package api

import (
	"net/http"

	"github.com/Daniel200273/WASA-project/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// setMyUserName handles updating the current user's username
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// TODO: Implementation needed
	// 1. Parse request body to get new username
	// 2. Validate username format (reuse isValidUsername)
	// 3. Get current user from context/token
	// 4. Check if new username is already taken
	// 5. Update username in database
	// 6. Return success response

	ctx.Logger.Info("setMyUserName endpoint called - TODO: implement")
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// setMyPhoto handles updating the current user's profile photo
func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// TODO: Implementation needed
	// 1. Parse multipart form data to get photo file
	// 2. Validate file format (JPEG, PNG, etc.)
	// 3. Validate file size (max 10MB)
	// 4. Save photo file to storage
	// 5. Get current user from context/token
	// 6. Update user's photo_url in database
	// 7. Return success response

	ctx.Logger.Info("setMyPhoto endpoint called - TODO: implement")
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// searchUsers handles searching for users by username
func (rt *_router) searchUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// TODO: Implementation needed
	// 1. Get search query parameter 'q' from URL
	// 2. Validate query parameter (required, max length, pattern)
	// 3. Get current user from context/token (to exclude from results)
	// 4. Search users in database using query
	// 5. Return list of matching users (limit results)
	// 6. Format response as JSON

	ctx.Logger.Info("searchUsers endpoint called - TODO: implement")
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
