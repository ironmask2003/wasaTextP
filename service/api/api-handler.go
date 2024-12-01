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
	//
	// -- Delete User -- //
	rt.router.DELETE("/profiles/:user", rt.wrap(rt.deleteUser, true))

	// Group routes
	//
	// -- Create Group -- //
	rt.router.POST("/profiles/:user/groups", rt.wrap(rt.createGroup, true))
	//
	// -- Delete Group -- //
	rt.router.DELETE("/profiles/:user/groups/:group", rt.wrap(rt.deleteGroup, true))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
