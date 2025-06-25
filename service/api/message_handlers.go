package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Daniel200273/WASA-project/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

// sendMessage handles sending a new message to a conversation
func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Get conversationId from URL path parameters
	conversationID := ps.ByName("conversationId")

	// 2. Validate conversationId format
	if err := validateID(conversationID, "conversationId"); err != nil {
		ctx.Logger.Error("Invalid conversation ID", "error", err)
		sendErrorResponse(w, http.StatusBadRequest, err.Error(), ctx)
		return
	}

	// 3. Get current user from context/token
	userID := ctx.UserID

	// 4. Check if conversation exists and if user is participant
	isParticipant, err := rt.db.IsUserInConversation(conversationID, userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to check user participation")
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to check conversation participation", ctx)
		return
	}

	if !isParticipant {
		// If user is not a participant, check if conversationID is actually a userID
		// In this case, we'll try to create or get a direct conversation
		targetUser, userErr := rt.db.GetUserByID(conversationID)
		if userErr != nil {
			ctx.Logger.Error("User not authorized to send message", "userID", userID, "conversationID", conversationID)
			sendErrorResponse(w, http.StatusForbidden, "Unauthorized access to conversation", ctx)
			return
		}

		// Validate that user is not trying to message themselves
		if userID == targetUser.ID {
			sendErrorResponse(w, http.StatusBadRequest, "Cannot send message to yourself", ctx)
			return
		}

		// Create or get direct conversation between current user and target user
		conversation, createErr := rt.db.GetOrCreateDirectConversation(userID, targetUser.ID)
		if createErr != nil {
			ctx.Logger.WithError(createErr).Error("Failed to create direct conversation")
			sendErrorResponse(w, http.StatusInternalServerError, "Failed to create conversation", ctx)
			return
		}

		// Update conversationID to the actual conversation ID
		conversationID = conversation.ID
		ctx.Logger.Info("Auto-created direct conversation", "conversationID", conversationID, "targetUserID", targetUser.ID)
	}

	// 5. Parse request body for text message or multipart for photo
	contentType := r.Header.Get("Content-Type")
	var content *string
	var photoURL *string
	var replyTo *string

	if strings.Contains(contentType, "application/json") {
		// Text message
		var req SendMessageRequest
		if err := parseJSONRequest(r, &req); err != nil {
			sendErrorResponse(w, http.StatusBadRequest, "Invalid request body", ctx)
			return
		}

		// 6. Validate message content
		if err := validateMessageContent(req.Content); err != nil {
			sendErrorResponse(w, http.StatusBadRequest, err.Error(), ctx)
			return
		}

		content = &req.Content
		replyTo = req.ReplyTo
	} else if strings.Contains(contentType, "multipart/form-data") {
		// Photo message
		file, header, err := getUploadedFile(r, "photo")
		if err != nil {
			sendErrorResponse(w, http.StatusBadRequest, err.Error(), ctx)
			return
		}
		defer file.Close()

		// Generate filename for photo
		filename := fmt.Sprintf("%s_%s", uuid.Must(uuid.NewV4()).String(), header.Filename)

		// Save photo file
		savedPhotoURL, err := saveUploadedImage(file, "messages", filename)
		if err != nil {
			ctx.Logger.Error("Failed to save message photo", "error", err)
			sendErrorResponse(w, http.StatusInternalServerError, "Failed to save photo", ctx)
			return
		}

		photoURL = &savedPhotoURL

		// Check for optional replyTo in form data
		if replyToStr := r.FormValue("replyTo"); replyToStr != "" {
			replyTo = &replyToStr
		}
	} else {
		sendErrorResponse(w, http.StatusBadRequest, "Content-Type must be application/json or multipart/form-data", ctx)
		return
	}

	// 7. Create message in database
	message, err := rt.db.CreateMessage(conversationID, userID, content, photoURL, replyTo)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to create message")
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to create message", ctx)
		return
	}

	// 8. Convert to response format
	comments := make([]CommentResponse, len(message.Comments))
	for i, comment := range message.Comments {
		comments[i] = CommentResponse{
			ID:        comment.ID,
			UserID:    comment.UserID,
			Username:  comment.Username,
			Emoticon:  comment.Emoticon,
			Timestamp: comment.CreatedAt,
		}
	}

	response := MessageResponse{
		ID:             message.ID,
		SenderID:       message.SenderID,
		SenderUsername: message.SenderUsername,
		Content:        message.Content,
		PhotoURL:       message.PhotoURL,
		ReplyToID:      message.ReplyToID,
		Forwarded:      message.Forwarded,
		Timestamp:      message.CreatedAt,
		Status:         message.Status,
		Comments:       comments,
	}

	// 9. Return created message as JSON response
	if err := sendJSONResponse(w, http.StatusCreated, response); err != nil {
		ctx.Logger.WithError(err).Error("failed to send message response")
	}

	ctx.Logger.Info("Message sent successfully", "messageID", message.ID, "conversationID", conversationID)
}

// forwardMessage handles forwarding an existing message to another conversation
func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Get messageId from URL path parameters
	messageID := ps.ByName("messageId")

	// 2. Validate messageId format
	if err := validateID(messageID, "messageId"); err != nil {
		ctx.Logger.Error("Invalid message ID", "error", err)
		sendErrorResponse(w, http.StatusBadRequest, err.Error(), ctx)
		return
	}

	// 3. Parse request body to get target conversationId
	var req ForwardMessageRequest
	if err := parseJSONRequest(r, &req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request body", ctx)
		return
	}

	// 4. Validate conversationId format
	if err := validateID(req.ConversationID, "conversationId"); err != nil {
		ctx.Logger.Error("Invalid target conversation ID", "error", err)
		sendErrorResponse(w, http.StatusBadRequest, err.Error(), ctx)
		return
	}

	// 5. Get current user from context/token
	userID := ctx.UserID

	// 6. Forward message using database operation (it handles all validation)
	forwardedMessage, err := rt.db.ForwardMessage(messageID, req.ConversationID, userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to forward message")
		if strings.Contains(err.Error(), "not found") {
			sendErrorResponse(w, http.StatusNotFound, "Message not found", ctx)
		} else if strings.Contains(err.Error(), "not authorized") || strings.Contains(err.Error(), "not a participant") {
			sendErrorResponse(w, http.StatusForbidden, "Unauthorized to forward message", ctx)
		} else {
			sendErrorResponse(w, http.StatusInternalServerError, "Failed to forward message", ctx)
		}
		return
	}

	// 7. Convert to response format
	comments := make([]CommentResponse, len(forwardedMessage.Comments))
	for i, comment := range forwardedMessage.Comments {
		comments[i] = CommentResponse{
			ID:        comment.ID,
			UserID:    comment.UserID,
			Username:  comment.Username,
			Emoticon:  comment.Emoticon,
			Timestamp: comment.CreatedAt,
		}
	}

	response := MessageResponse{
		ID:             forwardedMessage.ID,
		SenderID:       forwardedMessage.SenderID,
		SenderUsername: forwardedMessage.SenderUsername,
		Content:        forwardedMessage.Content,
		PhotoURL:       forwardedMessage.PhotoURL,
		ReplyToID:      forwardedMessage.ReplyToID,
		Forwarded:      forwardedMessage.Forwarded,
		Timestamp:      forwardedMessage.CreatedAt,
		Status:         forwardedMessage.Status,
		Comments:       comments,
	}

	// 8. Return created message as JSON response
	if err := sendJSONResponse(w, http.StatusCreated, response); err != nil {
		ctx.Logger.WithError(err).Error("failed to send forward message response")
	}

	ctx.Logger.Info("Message forwarded successfully", "originalMessageID", messageID, "forwardedMessageID", forwardedMessage.ID, "targetConversationID", req.ConversationID)
}

// deleteMessage handles deleting a message sent by the current user
func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Get messageId from URL path parameters
	messageID := ps.ByName("messageId")

	// 2. Validate messageId format
	if err := validateID(messageID, "messageId"); err != nil {
		ctx.Logger.Error("Invalid message ID", "error", err)
		sendErrorResponse(w, http.StatusBadRequest, err.Error(), ctx)
		return
	}

	// 3. Get current user from context/token
	userID := ctx.UserID

	// 4. Delete message using database operation (it handles ownership validation)
	err := rt.db.DeleteMessage(messageID, userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to delete message")
		if strings.Contains(err.Error(), "not found") {
			sendErrorResponse(w, http.StatusNotFound, "Message not found", ctx)
		} else if strings.Contains(err.Error(), "unauthorized") {
			sendErrorResponse(w, http.StatusForbidden, "Unauthorized to delete message", ctx)
		} else {
			sendErrorResponse(w, http.StatusInternalServerError, "Failed to delete message", ctx)
		}
		return
	}

	// 5. Return 204 No Content response
	w.WriteHeader(http.StatusNoContent)
	ctx.Logger.Info("Message deleted successfully", "messageID", messageID, "userID", userID)
}

// commentMessage handles adding a reaction/comment to a message
func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Get messageId from URL path parameters
	messageID := ps.ByName("messageId")

	// 2. Validate messageId format
	if err := validateID(messageID, "messageId"); err != nil {
		ctx.Logger.Error("Invalid message ID", "error", err)
		sendErrorResponse(w, http.StatusBadRequest, err.Error(), ctx)
		return
	}

	// 3. Parse request body to get emoticon
	var req CommentMessageRequest
	if err := parseJSONRequest(r, &req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request body", ctx)
		return
	}

	// 4. Validate emoticon format and length
	if err := validateEmoticon(req.Emoticon); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error(), ctx)
		return
	}

	// 5. Get current user from context/token
	userID := ctx.UserID

	// 6. Create or update reaction using database operation (it handles access validation)
	reaction, err := rt.db.CreateMessageReaction(messageID, userID, req.Emoticon)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to create message reaction")
		if strings.Contains(err.Error(), "not found") {
			sendErrorResponse(w, http.StatusNotFound, "Message not found", ctx)
		} else if strings.Contains(err.Error(), "not authorized") {
			sendErrorResponse(w, http.StatusForbidden, "Unauthorized to react to message", ctx)
		} else {
			sendErrorResponse(w, http.StatusInternalServerError, "Failed to create reaction", ctx)
		}
		return
	}

	// 7. Convert to response format
	response := CommentResponse{
		ID:        reaction.ID,
		UserID:    reaction.UserID,
		Username:  reaction.Username,
		Emoticon:  reaction.Emoticon,
		Timestamp: reaction.CreatedAt,
	}

	// 8. Return created/updated reaction as JSON response
	if err := sendJSONResponse(w, http.StatusCreated, response); err != nil {
		ctx.Logger.WithError(err).Error("failed to send reaction response")
	}

	ctx.Logger.Info("Message reaction created successfully", "reactionID", reaction.ID, "messageID", messageID)
}

// uncommentMessage handles removing a reaction/comment from a message
func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Get messageId and commentId from URL path parameters
	messageID := ps.ByName("messageId")
	commentID := ps.ByName("commentId")

	// 2. Validate both IDs format
	if err := validateID(messageID, "messageId"); err != nil {
		ctx.Logger.Error("Invalid message ID", "error", err)
		sendErrorResponse(w, http.StatusBadRequest, err.Error(), ctx)
		return
	}

	if err := validateID(commentID, "commentId"); err != nil {
		ctx.Logger.Error("Invalid comment ID", "error", err)
		sendErrorResponse(w, http.StatusBadRequest, err.Error(), ctx)
		return
	}

	// 3. Get current user from context/token
	userID := ctx.UserID

	// 4. Delete reaction using database operation (it handles ownership validation)
	err := rt.db.DeleteMessageReaction(commentID, userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to delete message reaction")
		if strings.Contains(err.Error(), "not found") {
			sendErrorResponse(w, http.StatusNotFound, "Reaction not found", ctx)
		} else if strings.Contains(err.Error(), "unauthorized") {
			sendErrorResponse(w, http.StatusForbidden, "Unauthorized to delete reaction", ctx)
		} else {
			sendErrorResponse(w, http.StatusInternalServerError, "Failed to delete reaction", ctx)
		}
		return
	}

	// 5. Return 204 No Content response
	w.WriteHeader(http.StatusNoContent)
	ctx.Logger.Info("Message reaction deleted successfully", "commentID", commentID, "messageID", messageID, "userID", userID)
}
