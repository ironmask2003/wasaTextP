package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// User routes
	//
	// -- Do Login -- //
	rt.router.POST("/session", rt.wrap(rt.doLogin, false))
	//
	// -- Search User -- //
	rt.router.GET("/profiles", rt.wrap(rt.searchUsers, true))
	//
	// -- Set my UserName -- //
	rt.router.PUT("/profiles/:user/username", rt.wrap(rt.setMyUserName, true))
	//
	// -- Set my Photo -- //
	rt.router.PUT("/profiles/:user/photo", rt.wrap(rt.setMyPhoto, true))
	//
	// -- Delete User -- //
	rt.router.DELETE("/profiles/:user", rt.wrap(rt.deleteUser, true))

	// Group routes
	//
	// -- Create Group -- //
	rt.router.POST("/profiles/:user/groups", rt.wrap(rt.createGroup, true))
	//
	// -- Leave Group -- //
	rt.router.DELETE("/profiles/:user/groups/:group", rt.wrap(rt.leaveGroup, true))
	//
	// -- Add to Group -- //
	rt.router.PUT("/profiles/:user/groups/:group", rt.wrap(rt.addToGroup, true))
	//
	// -- Set new name to the group -- //
	rt.router.PUT("/profiles/:user/groups/:group/groupname", rt.wrap(rt.setGroupName, true))
	//
	// -- Set new photo group -- //
	rt.router.PUT("/profiles/:user/groups/:group/grouphoto", rt.wrap(rt.setGroupPhoto, true))

	// -- Conversation Routes -- //
	//
	rt.router.PUT("/profiles/:user/conversations/:dest", rt.wrap(rt.createConversation, true))
	//
	// -- Get User Conversations -- //
	rt.router.GET("/profiles/:user/conversations", rt.wrap(rt.getMyConversations, true))
	//
	// -- Get conversation -- //
	rt.router.GET("/profiles/:user/conversations/:conv", rt.wrap(rt.getConversation, true))

	// Message Routes
	//
	// -- Send Message -- //
	rt.router.POST("/profiles/:user/conversations/:conv/messages", rt.wrap(rt.sendMessage, true))
	//
	// -- Forward Message -- //
	rt.router.POST("/profiles/:user/conversations/:conv/messages/:message", rt.wrap(rt.forwardMessage, true))
	//
	// -- Delete the message -- //
	rt.router.DELETE("/profiles/:user/conversations/:conv/messages/:message", rt.wrap(rt.deleteMessage, true))

	// Comment Routes
	//
	// -- Comment Message -- //
	rt.router.PUT("/profiles/:user/conversations/:dest/messages/:message/reactions", rt.wrap(rt.commentMessage, true))
	//
	// -- Uncomment Message -- //
	rt.router.DELETE("/profiles/:user/conversations/:conv/messages/:message/reactions/:comment", rt.wrap(rt.uncommentMessage, true))

	return rt.router
}
