package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasa.project/service/api/reqcontext"
	"wasa.project/service/api/structs"
)

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the userId from the endpoint
	userId, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Can't take the user id from the endpoint")
		return
	}

	// Check if the user is authorized
	if checkAuth(w, userId, ctx) != nil {
		Forbidden(w, nil, ctx, "The user is not authorized")
		return
	}

	// Get all conversations of the user from the db
	var convs []structs.Conversation
	convs, err = rt.db.GetUserConversations(userId)
	if err != nil {
		BadRequest(w, err, ctx, "Error taking the conversations from the db")
		return
	}

	// Fornire anche informazioni riguardo gruppi o utenti

	// Write the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(convs); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
