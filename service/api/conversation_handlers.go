package api

import (
	"net/http"

	"github.com/Daniel200273/WASA-project/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// getMyConversations handles getting all conversations for the current user
func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Get current user from context/token (already authenticated by wrapper)
	userID := ctx.UserID

	// 2. Retrieve all conversations for the user from database
	dbConversations, err := rt.db.GetUserConversations(userID)
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
		if dbConv.Type == "direct" && dbConv.OtherParticipant != nil {
			convResp.Name = dbConv.OtherParticipant.Username
		} else if dbConv.Name != nil {
			convResp.Name = *dbConv.Name
		} else {
			convResp.Name = "Conversation" // Fallback name
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

	// 4. Check if user is participant in the conversation
	if val, err := rt.db.IsUserInConversation(conversationID, userID); !val && err == nil {
		ctx.Logger.Error("User not authorized to access conversation", "userID", userID, "conversationID", conversationID)
		sendErrorResponse(w, http.StatusForbidden, "Unauthorized access to conversation", ctx)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("Failed to check user participation in conversation")
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to check conversation participation", ctx)
		return
	}

	// 5. Get conversation details (type, name, photo, members)
	conversationDetails, err := rt.db.GetConversation(conversationID, userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to retrieve conversation details")
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve conversation details", ctx)
		return
	}

	// 6. Get all messages in conversation with sender info
	messages, err := rt.db.GetConversationMessages(conversationID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to retrieve conversation messages")
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve messages", ctx)
		return
	}

	// 7. Format response as JSON with conversation details and messages
	response := ConversationDetailResponse{
		ID:   conversationDetails.ID,
		Type: conversationDetails.Type,
	}

	// Handle conversation name - for direct conversations, use other participant's name
	if conversationDetails.Type == "direct" && conversationDetails.OtherParticipant != nil {
		response.Name = conversationDetails.OtherParticipant.Username
	} else if conversationDetails.Name != nil {
		response.Name = *conversationDetails.Name
	} else {
		response.Name = "Conversation" // Fallback name
	}

	// Handle photo URL
	response.PhotoURL = conversationDetails.PhotoURL

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
