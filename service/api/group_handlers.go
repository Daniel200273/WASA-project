package api

import (
	"net/http"

	"github.com/Daniel200273/WASA-project/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// createGroup handles creating a new group conversation
func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// TODO: Implementation needed
	// 1. Parse request body to get group name and member IDs
	// 2. Validate group name (required, length, pattern)
	// 3. Validate member IDs array (required, length, format)
	// 4. Get current user from context/token
	// 5. Validate all member IDs exist in database
	// 6. Create group conversation in database
	// 7. Add all members (including creator) as participants
	// 8. Return created group details as JSON response

	ctx.Logger.Info("createGroup endpoint called - TODO: implement")
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// addToGroup handles adding a user to an existing group
func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// TODO: Implementation needed
	// 1. Get groupId from URL path parameters
	// 2. Validate groupId format
	// 3. Parse request body to get userId to add
	// 4. Validate userId format
	// 5. Get current user from context/token
	// 6. Check if current user is member of the group
	// 7. Check if target user exists and is not already in group
	// 8. Add user to group participants
	// 9. Return success response

	ctx.Logger.Info("addToGroup endpoint called - TODO: implement")
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// leaveGroup handles removing current user from a group
func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// TODO: Implementation needed
	// 1. Get groupId from URL path parameters
	// 2. Validate groupId format
	// 3. Get current user from context/token
	// 4. Check if user is member of the group
	// 5. Remove user from group participants
	// 6. Handle special case if user is the last member (delete group?)
	// 7. Return 204 No Content response

	ctx.Logger.Info("leaveGroup endpoint called - TODO: implement")
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// setGroupName handles updating a group's name
func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// TODO: Implementation needed
	// 1. Get groupId from URL path parameters
	// 2. Validate groupId format
	// 3. Parse request body to get new group name
	// 4. Validate group name (required, length, pattern)
	// 5. Get current user from context/token
	// 6. Check if user is member of the group
	// 7. Update group name in database
	// 8. Return success response

	ctx.Logger.Info("setGroupName endpoint called - TODO: implement")
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// setGroupPhoto handles updating a group's photo
func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// TODO: Implementation needed
	// 1. Get groupId from URL path parameters
	// 2. Validate groupId format
	// 3. Parse multipart form data to get photo file
	// 4. Validate file format (JPEG, PNG, etc.)
	// 5. Validate file size (max 10MB)
	// 6. Get current user from context/token
	// 7. Check if user is member of the group
	// 8. Save photo file to storage
	// 9. Update group photo_url in database
	// 10. Return success response

	ctx.Logger.Info("setGroupPhoto endpoint called - TODO: implement")
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
