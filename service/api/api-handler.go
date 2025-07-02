package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// Register all API endpoints here
	// Liveness endpoint for health checks
	rt.router.GET("/liveness", rt.wrap(rt.liveness, false))

	// Authentication endpoints
	rt.router.POST("/session", rt.wrap(rt.doLogin, false)) // ❌ NO auth (è il login!)

	// User Management endpoints - consistent pattern with userId
	rt.router.GET("/users", rt.wrap(rt.searchUsers, true))
	rt.router.GET("/users/:userId", rt.wrap(rt.getUserProfile, true))
	rt.router.PUT("/users/:userId/username", rt.wrap(rt.setMyUserName, true))
	rt.router.PUT("/users/:userId/photo", rt.wrap(rt.setMyPhoto, true))

	// Conversations endpoints - consistent with user-centric pattern
	rt.router.POST("/users/:userId/conversations", rt.wrap(rt.startConversation, true))
	rt.router.GET("/users/:userId/conversations", rt.wrap(rt.getMyConversations, true))
	rt.router.GET("/users/:userId/conversations/:conversationId", rt.wrap(rt.getConversation, true))

	// Messages endpoints - nested under conversations
	rt.router.POST("/users/:userId/conversations/:conversationId/messages", rt.wrap(rt.sendMessage, true))
	rt.router.POST("/users/:userId/messages/:messageId/forward", rt.wrap(rt.forwardMessage, true))
	rt.router.DELETE("/users/:userId/messages/:messageId", rt.wrap(rt.deleteMessage, true))
	rt.router.POST("/users/:userId/messages/:messageId/comments", rt.wrap(rt.commentMessage, true))
	rt.router.DELETE("/users/:userId/messages/:messageId/comments/:commentId", rt.wrap(rt.uncommentMessage, true))

	// Groups endpoints - consistent with user-centric pattern
	rt.router.POST("/users/:userId/groups", rt.wrap(rt.createGroup, true))
	rt.router.POST("/users/:userId/groups/:groupId/members", rt.wrap(rt.addToGroup, true))
	rt.router.DELETE("/users/:userId/groups/:groupId/members", rt.wrap(rt.leaveGroup, true))
	rt.router.DELETE("/users/:userId/groups/:groupId/members/:memberId", rt.wrap(rt.removeMemberFromGroup, true))
	rt.router.PUT("/users/:userId/groups/:groupId/name", rt.wrap(rt.setGroupName, true))
	rt.router.PUT("/users/:userId/groups/:groupId/photo", rt.wrap(rt.setGroupPhoto, true))

	// Static file serving for uploaded images (temporary storage)
	rt.router.ServeFiles("/uploads/*filepath", http.Dir("tmp/uploads"))

	return rt.router
}
