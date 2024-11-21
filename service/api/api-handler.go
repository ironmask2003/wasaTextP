package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// Settare tutti i path con le funzioni che gestiscono le richieste api (es. rt.router.GET("/users/:userID/conversations", rt.wrp (rt.getConversations, true)))

	return rt.router
}
