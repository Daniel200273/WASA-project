package api

import (
	"net/http"

	"github.com/Daniel200273/WASA-project/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// getMyConversations handles getting all conversations for the current user
func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// TODO: Implementation needed
	// 1. Get current user from context/token
	// 2. Retrieve all conversations for the user from database
	// 3. For each conversation:
	//    - Get last message details
	//    - Calculate unread count
	//    - For direct conversations, get other participant info
	//    - For group conversations, get group details
	// 4. Sort conversations by last message timestamp (most recent first)
	// 5. Format response as JSON with conversations array

	ctx.Logger.Info("getMyConversations endpoint called - TODO: implement")
	http.Error(w, "Not implemented", http.StatusNotImplemented)
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
