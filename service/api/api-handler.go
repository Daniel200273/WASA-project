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

	// User Management endpoints
	rt.router.GET("/users/me", rt.wrap(rt.getMyProfile, true))
	rt.router.GET("/users/:userId", rt.wrap(rt.getUserProfile, true))
	rt.router.PUT("/users/me/username", rt.wrap(rt.setMyUserName, true))
	rt.router.PUT("/users/me/photo", rt.wrap(rt.setMyPhoto, true))
	rt.router.GET("/users", rt.wrap(rt.searchUsers, true))
	rt.router.POST("/users/:userId/conversations", rt.wrap(rt.startConversation, true))

	// Conversations endpoints
	rt.router.GET("/conversations", rt.wrap(rt.getMyConversations, true))
	rt.router.GET("/conversations/:conversationId", rt.wrap(rt.getConversation, true))

	// Messages endpoints
	rt.router.POST("/conversations/:conversationId/messages", rt.wrap(rt.sendMessage, true))
	rt.router.POST("/messages/:messageId/forward", rt.wrap(rt.forwardMessage, true))
	rt.router.DELETE("/messages/:messageId", rt.wrap(rt.deleteMessage, true))
	rt.router.POST("/messages/:messageId/comments", rt.wrap(rt.commentMessage, true))
	rt.router.DELETE("/messages/:messageId/comments/:commentId", rt.wrap(rt.uncommentMessage, true))

	// Groups endpoints
	rt.router.POST("/groups", rt.wrap(rt.createGroup, true))
	rt.router.POST("/groups/:groupId/members", rt.wrap(rt.addToGroup, true))
	rt.router.DELETE("/groups/:groupId/members/me", rt.wrap(rt.leaveGroup, true))
	rt.router.PUT("/groups/:groupId/name", rt.wrap(rt.setGroupName, true))
	rt.router.PUT("/groups/:groupId/photo", rt.wrap(rt.setGroupPhoto, true))

	// Static file serving for uploaded images (temporary storage)
	rt.router.ServeFiles("/uploads/*filepath", http.Dir("tmp/uploads"))

	return rt.router
}
