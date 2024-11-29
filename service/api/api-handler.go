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
	// -- Set my UserName -- //
	rt.router.PUT("/profiles/:user/username", rt.wrap(rt.setMyUserName, true))
	//
	// -- Set my Photo -- //
	rt.router.PUT("/profiles/:user/photo", rt.wrap(rt.setMyPhoto, true))

	// Conversation routes
	//
	// -- Create Conversation -- //
	rt.router.PUT("/profiles/:user/conversations/:dest_user_id", rt.wrap(rt.CreateConversation, true))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
