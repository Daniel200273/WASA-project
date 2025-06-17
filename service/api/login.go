package api

import (
	"encoding/json"
	"net/http"

	"github.com/Daniel200273/WASA-project/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// LoginRequest represents the login request body
type LoginRequest struct {
	Name string `json:"name"`
}

// LoginResponse represents the login response body
type LoginResponse struct {
	Identifier string `json:"identifier"`
}

// doLogin handles user login/registration
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Parse request body
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ctx.Logger.WithError(err).Error("invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate username
	if req.Name == "" || len(req.Name) < 3 || len(req.Name) > 16 {
		ctx.Logger.Error("invalid username length")
		http.Error(w, "Username must be between 3 and 16 characters", http.StatusBadRequest)
		return
	}

	// Try to get existing user by username
	user, err := rt.db.GetUserByUsername(req.Name)
	if err != nil {
		// User doesn't exist, create new user
		ctx.Logger.WithField("username", req.Name).Info("creating new user")
		user, err = rt.db.CreateUser(req.Name)
		if err != nil {
			ctx.Logger.WithError(err).Error("failed to create user")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	// Create user session (token)
	token, err := rt.db.CreateUserSession(user.ID)
	if err != nil {
		ctx.Logger.WithError(err).Error("failed to create user session")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := LoginResponse{
		Identifier: token,
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		ctx.Logger.WithError(err).Error("failed to encode response")
	}

	ctx.Logger.WithField("user_id", user.ID).Info("user login successful")
}
