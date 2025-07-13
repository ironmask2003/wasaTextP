package api

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"wasa.project/service/api/reqcontext"
)

// handleWebSocket gestisce le connessioni WebSocket
func (rt *_router) handleWebSocket(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	username := r.URL.Query().Get("username")

	conn, err := rt.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	rt.wsConnMutex.Lock()
	rt.wsConnMap[username] = conn
	rt.wsConnMutex.Unlock()

	defer func() {
		rt.wsConnMutex.Lock()
		delete(rt.wsConnMap, username)
		rt.wsConnMutex.Unlock()
		conn.Close()
	}()

	log.Println("WebSocket connected:", username)

	// Opzionale: leggi messaggi dal client
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket disconnected:", username)
			break
		}
	}
}
