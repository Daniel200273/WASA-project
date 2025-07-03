package api

import (
	"net/http"

	"github.com/Daniel200273/WASA-project/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// Constants for conversation types
const (
	ConversationTypeDefault = "Conversation"
)

// startConversation creates or gets a direct conversation with another user
func (rt *_router) startConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Get user ID from URL parameter and validate authorization
	userID := ps.ByName("userId")
	if userID == "" {
		sendErrorResponse(w, http.StatusBadRequest, "User ID is required", ctx)
		return
	}

	// Authorization check - user can only create conversations for themselves
	if userID != ctx.UserID {
		sendErrorResponse(w, http.StatusForbidden, "You can only create conversations for yourself", ctx)
		return
	}

	// 2. Parse request body to get target user ID
	var req StartConversationRequest
	if err := parseJSONRequest(r, &req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request body", ctx)
		return
	}

	// 3. Validate target user ID format
	if err := validateID(req.UserID, "userId"); err != nil {
		ctx.Logger.Error("Invalid user ID", "error", err)
		sendErrorResponse(w, http.StatusBadRequest, err.Error(), ctx)
		return
	}

	// 4. Validate that user is not trying to start conversation with themselves
	if ctx.UserID == req.UserID {
		sendErrorResponse(w, http.StatusBadRequest, "Cannot start conversation with yourself", ctx)
		return
	}

	// 5. Check if target user exists
	_, err := rt.db.GetUserByID(req.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Target user not found")
		sendErrorResponse(w, http.StatusNotFound, "User not found", ctx)
		return
	}

	// 6. Get or create direct conversation using existing database method
	conversation, err := rt.db.GetOrCreateDirectConversation(ctx.UserID, req.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to create conversation")
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to create conversation", ctx)
		return
	}

	// 7. Return conversation details
	response := ConversationDetailResponse{
		ID:       conversation.ID,
		Type:     conversation.Type,
		Members:  make([]UserResponse, len(conversation.Participants)),
		Messages: []MessageResponse{}, // Empty for new conversations
	}

	// Set conversation name to other participant's name
	if conversation.OtherParticipant != nil {
		response.Name = conversation.OtherParticipant.Username
	} else {
		response.Name = ConversationTypeDefault
	}

	// Convert participants to response format
	for i, participant := range conversation.Participants {
		response.Members[i] = UserResponse{
			ID:       participant.ID,
			Username: participant.Username,
			PhotoURL: participant.PhotoURL,
		}
	}

	if err := sendJSONResponse(w, http.StatusCreated, response); err != nil {
		ctx.Logger.WithError(err).Error("failed to send conversation response")
	}

	ctx.Logger.Info("Conversation created/retrieved successfully", "conversationID", conversation.ID)
}

// getMyConversations handles getting all conversations for the current user
func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Get user ID from URL parameter and validate authorization
	userID := ps.ByName("userId")
	if userID == "" {
		sendErrorResponse(w, http.StatusBadRequest, "User ID is required", ctx)
		return
	}

	// Authorization check - user can only get their own conversations
	if userID != ctx.UserID {
		sendErrorResponse(w, http.StatusForbidden, "You can only access your own conversations", ctx)
		return
	}

	// 2. Retrieve all conversations for the user from database
	dbConversations, err := rt.db.GetUserConversations(ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("failed to retrieve user conversations")
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve conversations", ctx)
		return
	}

	// 3. Map database models to API response format
	response := ConversationsResponse{
		Conversations: make([]ConversationResponse, len(dbConversations)),
	}

	for i, dbConv := range dbConversations {
		// Convert each ConversationPreview to ConversationResponse
		convResp := ConversationResponse{
			ID:          dbConv.ID,
			Type:        dbConv.Type,
			UnreadCount: dbConv.UnreadCount,
		}

		// Handle name - check if direct conversation and use other participant's name if available
		switch {
		case dbConv.Type == "direct" && dbConv.OtherParticipant != nil:
			convResp.Name = dbConv.OtherParticipant.Username
		case dbConv.Name != nil:
			convResp.Name = *dbConv.Name
		default:
			convResp.Name = ConversationTypeDefault // Fallback name
		}

		// Handle photo URL
		convResp.PhotoURL = dbConv.PhotoURL

		// Convert last message if present
		if dbConv.LastMessage != nil {
			convResp.LastMessage = &MessagePreview{
				ID:             dbConv.LastMessage.ID,
				Content:        dbConv.LastMessage.Content,
				Timestamp:      dbConv.LastMessage.Timestamp,
				SenderUsername: dbConv.LastMessage.SenderUsername,
				HasPhoto:       dbConv.LastMessage.HasPhoto,
			}
		}

		response.Conversations[i] = convResp
	}

	// 4. Return the response as JSON
	if err := sendJSONResponse(w, http.StatusOK, response); err != nil {
		ctx.Logger.WithError(err).Error("failed to send conversations response")
	}

	ctx.Logger.Info("Retrieved conversations successfully", "count", len(dbConversations))
}

// getConversation handles getting messages in a specific conversation
func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Get user ID from URL parameter and validate authorization
	userID := ps.ByName("userId")
	if userID == "" {
		sendErrorResponse(w, http.StatusBadRequest, "User ID is required", ctx)
		return
	}

	// Authorization check - user can only access their own conversations
	if userID != ctx.UserID {
		sendErrorResponse(w, http.StatusForbidden, "You can only access your own conversations", ctx)
		return
	}

	// 2. Get conversationId from URL path parameters
	conversationID := ps.ByName("conversationId")

	// 3. Validate conversationId format
	if err := validateID(conversationID, "conversationId"); err != nil {
		ctx.Logger.Error("Invalid conversation ID", "error", err)
		sendErrorResponse(w, http.StatusBadRequest, err.Error(), ctx)
		return
	}

	// 4. Check if user is participant in the conversation
	if val, err := rt.db.IsUserInConversation(conversationID, ctx.UserID); !val && err == nil {
		ctx.Logger.Error("User not authorized to access conversation", "userID", ctx.UserID, "conversationID", conversationID)
		sendErrorResponse(w, http.StatusForbidden, "Unauthorized access to conversation", ctx)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("Failed to check user participation in conversation")
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to check conversation participation", ctx)
		return
	}

	// 5. Get conversation details (type, name, photo, members)
	conversationDetails, err := rt.db.GetConversation(conversationID, ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to retrieve conversation details")
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve conversation details", ctx)
		return
	}

	// 6. Mark conversation as read when user opens it
	if err := rt.db.MarkConversationAsRead(conversationID, ctx.UserID); err != nil {
		ctx.Logger.WithError(err).Warn("Failed to mark conversation as read") // Don't fail the request for this
	}

	// 7. Get all messages in conversation with sender info
	messages, err := rt.db.GetConversationMessages(conversationID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to retrieve conversation messages")
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve messages", ctx)
		return
	}

	// 8. Format response as JSON with conversation details and messages
	response := ConversationDetailResponse{
		ID:   conversationDetails.ID,
		Type: conversationDetails.Type,
	}

	// Handle conversation name - for direct conversations, use other participant's name
	switch {
	case conversationDetails.Type == "direct" && conversationDetails.OtherParticipant != nil:
		response.Name = conversationDetails.OtherParticipant.Username
	case conversationDetails.Name != nil:
		response.Name = *conversationDetails.Name
	default:
		response.Name = ConversationTypeDefault // Fallback name
	}

	// Handle photo URL
	response.PhotoURL = conversationDetails.PhotoURL

	// Handle timestamps
	response.CreatedAt = &conversationDetails.CreatedAt
	response.LastMessageAt = &conversationDetails.LastMessageAt

	// Convert participants to members format
	response.Members = make([]UserResponse, len(conversationDetails.Participants))
	for i, participant := range conversationDetails.Participants {
		response.Members[i] = UserResponse{
			ID:       participant.ID,
			Username: participant.Username,
			PhotoURL: participant.PhotoURL,
		}
	}

	// Convert messages to response format
	response.Messages = make([]MessageResponse, len(messages))
	for i, msg := range messages {
		// Convert comments/reactions
		comments := make([]CommentResponse, len(msg.Comments))
		for j, comment := range msg.Comments {
			comments[j] = CommentResponse{
				ID:        comment.ID,
				UserID:    comment.UserID,
				Username:  comment.Username,
				Emoticon:  comment.Emoticon,
				Timestamp: comment.CreatedAt,
			}
		}

		response.Messages[i] = MessageResponse{
			ID:             msg.ID,
			SenderID:       msg.SenderID,
			SenderUsername: msg.SenderUsername,
			Content:        msg.Content,
			PhotoURL:       msg.PhotoURL,
			ReplyToID:      msg.ReplyToID,
			Forwarded:      msg.Forwarded,
			Timestamp:      msg.CreatedAt,
			Status:         msg.Status,
			Comments:       comments,
		}
	}

	// 8. Return the response as JSON
	if err := sendJSONResponse(w, http.StatusOK, response); err != nil {
		ctx.Logger.WithError(err).Error("failed to send conversation response")
	}

	ctx.Logger.Info("Retrieved conversation successfully", "conversationID", conversationID, "messageCount", len(messages))
}
