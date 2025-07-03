package api

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Daniel200273/WASA-project/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// setMyUserName handles updating a user's username
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Get user ID from URL parameter
	userID := ps.ByName("userId")
	if userID == "" {
		sendErrorResponse(w, http.StatusBadRequest, "User ID is required", ctx)
		return
	}

	// 2. Authorization check - user can only update their own username
	if userID != ctx.UserID {
		sendErrorResponse(w, http.StatusForbidden, "You can only update your own username", ctx)
		return
	}

	// 3. Parse request body to get new username
	var req UpdateUsernameRequest
	if err := parseJSONRequest(r, &req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request body", ctx)
		return
	}

	// 4. Validate username format (reuse validateUsername)
	if err := validateUsername(req.Name); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid username format", ctx)
		return
	}

	// 5. Check if new username is already taken
	_, err := rt.db.GetUserByUsername(req.Name)
	if err != nil && err.Error() != "user not found" {
		ctx.Logger.Error("Failed to check existing username", "error", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal server error", ctx)
		return
	}

	// 6. Update username in database
	if err := rt.db.UpdateUsername(ctx.UserID, req.Name); err != nil {
		ctx.Logger.Error("Failed to update username", "error", err, "ctxUserID", ctx.UserID)
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to update username", ctx)
		return
	}

	// 7. Return success response
	w.WriteHeader(http.StatusNoContent) // 204 No Content
	ctx.Logger.Info("Username updated successfully", "userID", ctx.UserID, "newUsername", req.Name)
}

// setMyPhoto handles updating a user's profile photo
func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Get user ID from URL parameter
	userID := ps.ByName("userId")
	if userID == "" {
		sendErrorResponse(w, http.StatusBadRequest, "User ID is required", ctx)
		return
	}

	// 2. Authorization check - user can only update their own photo
	if userID != ctx.UserID {
		sendErrorResponse(w, http.StatusForbidden, "You can only update your own photo", ctx)
		return
	}

	// 3. Get and validate uploaded file
	file, header, err := getUploadedFile(r, "photo")
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error(), ctx)
		return
	}
	defer file.Close()

	// 4. Generate filename (use userID to ensure uniqueness)
	fileExt := strings.ToLower(filepath.Ext(header.Filename))
	if fileExt == "" {
		fileExt = ".jpg" // default extension
	}
	filename := fmt.Sprintf("%s%s", ctx.UserID, fileExt)

	// 5. Save photo file to temporary storage
	photoURL, err := saveUploadedImage(file, "profiles", filename)
	if err != nil {
		ctx.Logger.Error("Failed to save profile photo", "error", err, "userID", ctx.UserID)
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to save photo", ctx)
		return
	}

	// 6. Update user's photo_url in database
	if err := rt.db.UpdateUserPhoto(ctx.UserID, photoURL); err != nil {
		ctx.Logger.Error("Failed to update user photo in database", "error", err, "userID", ctx.UserID)
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to update photo", ctx)
		return
	}

	// 7. Return success response
	w.WriteHeader(http.StatusNoContent)
	ctx.Logger.Info("Profile photo updated successfully", "userID", ctx.UserID, "photoURL", photoURL)
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

	// 3. Get current user from context/token (to exclude from results, already authenticated by wrapper)
	userID := ctx.UserID

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

// getUserProfile handles getting a specific user's profile information by user ID
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get user ID from URL parameter
	userID := ps.ByName("userId")
	if userID == "" {
		sendErrorResponse(w, http.StatusBadRequest, "User ID is required", ctx)
		return
	}

	// Validate user ID format (basic validation)
	if len(userID) < 1 || len(userID) > 64 {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid user ID format", ctx)
		return
	}

	// Get user by ID directly
	user, err := rt.db.GetUserByID(userID)
	if err != nil {
		if err.Error() == "user not found" {
			sendErrorResponse(w, http.StatusNotFound, "User not found", ctx)
			return
		}
		ctx.Logger.Error("Failed to get user profile", "error", err, "userID", userID)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal server error", ctx)
		return
	}

	// Convert to response format
	response := UserResponse{
		ID:       user.ID,
		Username: user.Username,
		PhotoURL: user.PhotoURL,
	}

	// Send response
	if err := sendJSONResponse(w, http.StatusOK, response); err != nil {
		ctx.Logger.WithError(err).Error("failed to send user profile response")
		return
	}

	ctx.Logger.Info("User profile retrieved successfully", "requestedUserID", userID, "requestingUserID", ctx.UserID)
}
