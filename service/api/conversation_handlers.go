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
	// TODO: Implementation needed
	// 1. Get conversationId from URL path parameters
	// 2. Validate conversationId format
	// 3. Get current user from context/token
	// 4. Check if user is participant in the conversation
	// 5. Get conversation details (type, name, photo, members)
	// 6. Get all messages in conversation with sender info
	// 7. Get reactions/comments for each message
	// 8. Mark messages as read for current user
	// 9. Format response as JSON with conversation details and messages

	ctx.Logger.Info("getConversation endpoint called - TODO: implement")
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
