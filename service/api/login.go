package api

import (
	"net/http"

	"github.com/Daniel200273/WASA-project/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// doLogin handles user login/registration
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Parse request body
	var req LoginRequest
	if err := parseJSONRequest(r, &req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request body", ctx)
		return
	}

	// Validate username
	if err := validateUsername(req.Name); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid username", ctx)
		return
	}

	// Try to get existing user by username
	user, err := rt.db.GetUserByUsername(req.Name)
	if err != nil {
		// User doesn't exist, create new user
		ctx.Logger.WithField("username", req.Name).Info("creating new user")
		user, err = rt.db.CreateUser(req.Name)
		if err != nil {
			sendErrorResponse(w, http.StatusInternalServerError, "Internal server error", ctx)
			return
		}
	}

	// Create user session (token)
	token, err := rt.db.CreateUserSession(user.ID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal server error", ctx)
		return
	}

	// Prepare response
	response := LoginResponse{
		Identifier: token,
		UserID:     user.ID,
	}

	// Send response
	if err := sendJSONResponse(w, http.StatusCreated, response); err != nil {
		ctx.Logger.WithError(err).Error("failed to encode response")
	}

	ctx.Logger.WithField("user_id", user.ID).Info("user login successful")
}
