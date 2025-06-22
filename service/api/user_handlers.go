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
	var req UpdateUsernameRequest
	if err := parseJSONRequest(r, &req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request body", ctx)
		return
	}

	// 2. Validate username format (reuse validateUsername)
	if err := validateUsername(req.Name); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid username format", ctx)
		return
	}

	// 3. Get current user from context/token
	userID := ctx.UserID
	if userID == "" {
		sendErrorResponse(w, http.StatusUnauthorized, "Authentication required", ctx)
		return
	}

	// 4. Check if new username is already taken
	existingUser, err := rt.db.GetUserByUsername(req.Name)
	if err != nil {
		ctx.Logger.Error("Failed to check existing username", "error", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal server error", ctx)
		return
	}

	if existingUser != nil {
		if existingUser.ID != userID {
			sendErrorResponse(w, http.StatusConflict, "Username already taken", ctx)
			return
		}
		// User is setting the same username they already have - allow it (idempotent operation)
		w.WriteHeader(http.StatusNoContent)
		ctx.Logger.Info("Username unchanged (same as current)", "userID", userID, "username", req.Name)
		return
	}

	// 5. Update username in database
	if err := rt.db.UpdateUsername(userID, req.Name); err != nil {
		ctx.Logger.Error("Failed to update username", "error", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to update username", ctx)
		return
	}
	// 6. Return success response
	w.WriteHeader(http.StatusNoContent) // 204 No Content
	ctx.Logger.Info("Username updated successfully", "userID", userID, "newUsername", req.Name)

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
	// 1. Get search query parameter 'q' from URL
	query := getQueryParam(r, "q")
	if query == "" {
		sendErrorResponse(w, http.StatusBadRequest, "Query parameter 'q' is required", ctx)
		return
	}

	// 2. Validate query parameter (required, max length, pattern)
	if len(query) < 1 || len(query) > 50 {
		sendErrorResponse(w, http.StatusBadRequest, "Query must be between 1 and 50 characters", ctx)
		return
	}

	// 3. Get current user from context/token (to exclude from results)
	userID := ctx.UserID
	if userID == "" {
		sendErrorResponse(w, http.StatusUnauthorized, "Authentication required", ctx)
		return
	}

	// 4. Search users in database using query (find users containing the search string)
	users, err := rt.db.SearchUsers(query, userID)
	if err != nil {
		ctx.Logger.Error("Failed to search users", "error", err, "query", query)
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to search users", ctx)
		return
	}

	// 5. Format response as JSON with list of matching users
	response := SearchUsersResponse{
		Users: make([]UserResponse, len(users)),
	}

	// Convert database users to response format
	for i, user := range users {
		response.Users[i] = UserResponse{
			ID:       user.ID,
			Username: user.Username,
			PhotoURL: user.PhotoURL,
		}
	}

	// 6. Return results (limit is handled in database layer)
	if err := sendJSONResponse(w, http.StatusOK, response); err != nil {
		ctx.Logger.WithError(err).Error("failed to send search response")
	}

	ctx.Logger.Info("User search completed", "query", query, "resultsCount", len(users))
}
