package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Daniel200273/WASA-project/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

// createGroup handles creating a new group conversation
func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Parse request body to get group name and member IDs
	var req CreateGroupRequest
	if err := parseJSONRequest(r, &req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request body", ctx)
		return
	}

	// 2. Validate group name (required, length, pattern)
	if err := validateGroupName(req.Name); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error(), ctx)
		return
	}

	// 3. Validate member IDs array (required, length, format)
	if len(req.Members) == 0 {
		sendErrorResponse(w, http.StatusBadRequest, "At least one member is required", ctx)
		return
	}

	if len(req.Members) > 99 { // Max 99 members + creator = 100 total
		sendErrorResponse(w, http.StatusBadRequest, "Maximum 99 members allowed", ctx)
		return
	}

	for _, memberID := range req.Members {
		if err := validateID(memberID, "member ID"); err != nil {
			sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Invalid member ID: %s", err.Error()), ctx)
			return
		}
	}

	// 4. Get current user from context/token
	creatorID := ctx.UserID

	// 5. Validate all member IDs exist in database and are not duplicates
	memberIDsSet := make(map[string]bool)
	for _, memberID := range req.Members {
		if memberIDsSet[memberID] {
			sendErrorResponse(w, http.StatusBadRequest, "Duplicate member IDs are not allowed", ctx)
			return
		}
		memberIDsSet[memberID] = true

		// Check if member exists
		_, err := rt.db.GetUserByID(memberID)
		if err != nil {
			ctx.Logger.WithError(err).Error("Member not found", "memberID", memberID)
			sendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("User with ID %s not found", memberID), ctx)
			return
		}
	}

	// Check if creator is trying to add themselves (not allowed, they're added automatically)
	if memberIDsSet[creatorID] {
		sendErrorResponse(w, http.StatusBadRequest, "Cannot add yourself as a member - you are automatically the group creator", ctx)
		return
	}

	// 6. Create group conversation in database
	group, err := rt.db.CreateGroup(req.Name, creatorID, req.Members)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to create group")
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to create group", ctx)
		return
	}

	// 7. Convert to response format
	members := make([]UserResponse, len(group.Participants))
	for i, participant := range group.Participants {
		members[i] = UserResponse{
			ID:       participant.ID,
			Username: participant.Username,
			PhotoURL: participant.PhotoURL,
		}
	}

	response := GroupResponse{
		ID:        group.ID,
		Name:      *group.Name, // Groups always have names
		PhotoURL:  group.PhotoURL,
		Members:   members,
		CreatedBy: *group.CreatedBy,
		CreatedAt: group.CreatedAt,
	}

	// 8. Return created group details as JSON response
	if err := sendJSONResponse(w, http.StatusCreated, response); err != nil {
		ctx.Logger.WithError(err).Error("failed to send group creation response")
	}

	ctx.Logger.Info("Group created successfully", "groupID", group.ID, "name", req.Name, "memberCount", len(req.Members)+1)
}

// addToGroup handles adding a user to an existing group
func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Get groupId from URL path parameters
	groupID := ps.ByName("groupId")

	// 2. Validate groupId format
	if err := validateID(groupID, "groupId"); err != nil {
		ctx.Logger.Error("Invalid group ID", "error", err)
		sendErrorResponse(w, http.StatusBadRequest, err.Error(), ctx)
		return
	}

	// 3. Parse request body to get userId to add
	var req AddToGroupRequest
	if err := parseJSONRequest(r, &req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request body", ctx)
		return
	}

	// 4. Validate userId format
	if err := validateID(req.UserID, "userId"); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error(), ctx)
		return
	}

	// 5. Get current user from context/token
	currentUserID := ctx.UserID

	// 6. Check if current user is member of the group
	isCurrentUserMember, err := rt.db.IsUserInConversation(groupID, currentUserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to check user membership")
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to verify group membership", ctx)
		return
	}

	if !isCurrentUserMember {
		sendErrorResponse(w, http.StatusForbidden, "Only group members can add new users", ctx)
		return
	}

	// 7. Check if target user exists and is not already in group
	_, err = rt.db.GetUserByID(req.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Target user not found", "userID", req.UserID)
		sendErrorResponse(w, http.StatusBadRequest, "User not found", ctx)
		return
	}

	isTargetUserMember, err := rt.db.IsUserInConversation(groupID, req.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to check target user membership")
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to verify user membership", ctx)
		return
	}

	if isTargetUserMember {
		sendErrorResponse(w, http.StatusConflict, "User is already a member of this group", ctx)
		return
	}

	// 8. Add user to group participants
	err = rt.db.AddUserToGroup(groupID, req.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to add user to group")
		if strings.Contains(err.Error(), "not found") {
			sendErrorResponse(w, http.StatusNotFound, "Group not found", ctx)
		} else {
			sendErrorResponse(w, http.StatusInternalServerError, "Failed to add user to group", ctx)
		}
		return
	}

	// 9. Return success response
	w.WriteHeader(http.StatusNoContent)
	ctx.Logger.Info("User added to group successfully", "groupID", groupID, "userID", req.UserID, "addedBy", currentUserID)
}

// leaveGroup handles removing current user from a group
func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Get groupId from URL path parameters
	groupID := ps.ByName("groupId")

	// 2. Validate groupId format
	if err := validateID(groupID, "groupId"); err != nil {
		ctx.Logger.Error("Invalid group ID", "error", err)
		sendErrorResponse(w, http.StatusBadRequest, err.Error(), ctx)
		return
	}

	// 3. Get current user from context/token
	userID := ctx.UserID

	// 4. Check if user is member of the group
	isMember, err := rt.db.IsUserInConversation(groupID, userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to check user membership")
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to verify group membership", ctx)
		return
	}

	if !isMember {
		sendErrorResponse(w, http.StatusForbidden, "You are not a member of this group", ctx)
		return
	}

	// 5. Remove user from group participants
	err = rt.db.RemoveUserFromGroup(groupID, userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to remove user from group")
		if strings.Contains(err.Error(), "not found") {
			sendErrorResponse(w, http.StatusNotFound, "Group not found", ctx)
		} else {
			sendErrorResponse(w, http.StatusInternalServerError, "Failed to leave group", ctx)
		}
		return
	}

	// 6. Return 204 No Content response
	w.WriteHeader(http.StatusNoContent)
	ctx.Logger.Info("User left group successfully", "groupID", groupID, "userID", userID)
}

// setGroupName handles updating a group's name
func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Get groupId from URL path parameters
	groupID := ps.ByName("groupId")

	// 2. Validate groupId format
	if err := validateID(groupID, "groupId"); err != nil {
		ctx.Logger.Error("Invalid group ID", "error", err)
		sendErrorResponse(w, http.StatusBadRequest, err.Error(), ctx)
		return
	}

	// 3. Parse request body to get new group name
	var req UpdateGroupNameRequest
	if err := parseJSONRequest(r, &req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request body", ctx)
		return
	}

	// 4. Validate group name (required, length, pattern)
	if err := validateGroupName(req.Name); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error(), ctx)
		return
	}

	// 5. Get current user from context/token
	userID := ctx.UserID

	// 6. Check if user is member of the group
	isMember, err := rt.db.IsUserInConversation(groupID, userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to check user membership")
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to verify group membership", ctx)
		return
	}

	if !isMember {
		sendErrorResponse(w, http.StatusForbidden, "Only group members can update the group name", ctx)
		return
	}

	// 7. Update group name in database
	err = rt.db.UpdateGroupName(groupID, req.Name)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to update group name")
		if strings.Contains(err.Error(), "not found") {
			sendErrorResponse(w, http.StatusNotFound, "Group not found", ctx)
		} else {
			sendErrorResponse(w, http.StatusInternalServerError, "Failed to update group name", ctx)
		}
		return
	}

	// 8. Return success response
	w.WriteHeader(http.StatusNoContent)
	ctx.Logger.Info("Group name updated successfully", "groupID", groupID, "newName", req.Name, "updatedBy", userID)
}

// setGroupPhoto handles updating a group's photo
func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Get groupId from URL path parameters
	groupID := ps.ByName("groupId")

	// 2. Validate groupId format
	if err := validateID(groupID, "groupId"); err != nil {
		ctx.Logger.Error("Invalid group ID", "error", err)
		sendErrorResponse(w, http.StatusBadRequest, err.Error(), ctx)
		return
	}

	// 3. Validate image file format and size
	if err := validateImageFile(r, "photo"); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error(), ctx)
		return
	}

	// 4. Parse multipart form data to get photo file
	file, header, err := getUploadedFile(r, "photo")
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error(), ctx)
		return
	}
	defer file.Close()

	// 5. Get current user from context/token
	userID := ctx.UserID

	// 6. Check if user is member of the group
	isMember, err := rt.db.IsUserInConversation(groupID, userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to check user membership")
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to verify group membership", ctx)
		return
	}

	if !isMember {
		sendErrorResponse(w, http.StatusForbidden, "Only group members can update the group photo", ctx)
		return
	}

	// 7. Generate filename for photo
	filename := fmt.Sprintf("%s_%s", uuid.Must(uuid.NewV4()).String(), header.Filename)

	// 8. Save photo file to storage
	photoURL, err := saveUploadedImage(file, "groups", filename)
	if err != nil {
		ctx.Logger.Error("Failed to save group photo", "error", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to save photo", ctx)
		return
	}

	// 9. Update group photo_url in database
	err = rt.db.UpdateGroupPhoto(groupID, photoURL)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to update group photo in database")
		if strings.Contains(err.Error(), "not found") {
			sendErrorResponse(w, http.StatusNotFound, "Group not found", ctx)
		} else {
			sendErrorResponse(w, http.StatusInternalServerError, "Failed to update group photo", ctx)
		}
		return
	}

	// 10. Return success response
	w.WriteHeader(http.StatusNoContent)
	ctx.Logger.Info("Group photo updated successfully", "groupID", groupID, "photoURL", photoURL, "updatedBy", userID)
}
