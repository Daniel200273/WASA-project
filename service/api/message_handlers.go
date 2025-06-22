package api

import (
	"net/http"

	"github.com/Daniel200273/WASA-project/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// sendMessage handles sending a new message to a conversation
func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// TODO: Implementation needed
	// 1. Get conversationId from URL path parameters
	// 2. Validate conversationId format
	// 3. Get current user from context/token
	// 4. Check if user is participant in the conversation
	// 5. Parse request body/form data:
	//    - For text messages: get content and optional replyTo
	//    - For photo messages: get photo file and optional replyTo
	// 6. Validate message content (text length, photo size/format)
	// 7. If replyTo is provided, validate it exists in the conversation
	// 8. Save photo file if it's a photo message
	// 9. Create message in database
	// 10. Return created message as JSON response

	ctx.Logger.Info("sendMessage endpoint called - TODO: implement")
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// forwardMessage handles forwarding an existing message to another conversation
func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// TODO: Implementation needed
	// 1. Get messageId from URL path parameters
	// 2. Validate messageId format
	// 3. Parse request body to get target conversationId
	// 4. Validate conversationId format
	// 5. Get current user from context/token
	// 6. Check if user has access to the original message
	// 7. Check if user is participant in target conversation
	// 8. Get original message details
	// 9. Create forwarded message in target conversation
	// 10. Return created message as JSON response

	ctx.Logger.Info("forwardMessage endpoint called - TODO: implement")
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// deleteMessage handles deleting a message sent by the current user
func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// TODO: Implementation needed
	// 1. Get messageId from URL path parameters
	// 2. Validate messageId format
	// 3. Get current user from context/token
	// 4. Check if message exists and was sent by current user
	// 5. Delete message from database (consider cascade effects on reactions)
	// 6. Return 204 No Content response

	ctx.Logger.Info("deleteMessage endpoint called - TODO: implement")
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// commentMessage handles adding a reaction/comment to a message
func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// TODO: Implementation needed
	// 1. Get messageId from URL path parameters
	// 2. Validate messageId format
	// 3. Parse request body to get emoticon
	// 4. Validate emoticon format and length
	// 5. Get current user from context/token
	// 6. Check if user has access to the message
	// 7. Check if user already reacted to this message (update vs create)
	// 8. Create or update reaction in database
	// 9. Return created/updated reaction as JSON response

	ctx.Logger.Info("commentMessage endpoint called - TODO: implement")
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// uncommentMessage handles removing a reaction/comment from a message
func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// TODO: Implementation needed
	// 1. Get messageId and commentId from URL path parameters
	// 2. Validate both IDs format
	// 3. Get current user from context/token
	// 4. Check if comment exists and was created by current user
	// 5. Delete reaction from database
	// 6. Return 204 No Content response

	ctx.Logger.Info("uncommentMessage endpoint called - TODO: implement")
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
