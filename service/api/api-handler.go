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
	rt.router.PUT("/session/:user/username", rt.wrap(rt.setMyUserName, true))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
