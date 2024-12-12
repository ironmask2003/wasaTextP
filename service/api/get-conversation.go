package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"wasa.project/service/api/reqcontext"
	"wasa.project/service/api/structs"
)

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the userId from the endpoint
	userId, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Can't take the user id from the endpoint")
		return
	}
	// Check if the user is authorized
	if checkAuth(w, userId, ctx) != nil {
		return
	}

	// Get the id of the conversation
	convId, err := strconv.Atoi(ps.ByName("conv"))
	if err != nil {
		BadRequest(w, err, ctx, "Can't take the conversation id from the endpoint")
		return
	}

	// Get the conversation by id
	var conv structs.Conversation
	conv, err = rt.db.GetConversationById(convId)
	if err != nil {
		BadRequest(w, err, ctx, "Can't take the conversation from the db")
		return
	}

	messages, err := rt.db.GetMessages(conv.ConversationId)
	if err != nil {
		BadRequest(w, err, ctx, "Can't get messages of the conversation")
		return
	}

	for i := 0; i < len(messages); i++ {
		err = rt.db.UpdateStatusMessage(messages[i].MessageId, conv.ConversationId)
		if err != nil {
			BadRequest(w, err, ctx, "Error updating status o message")
			return
		}
	}

	// Getting all message after the update of the status
	messages, err = rt.db.GetMessages(conv.ConversationId)
	if err != nil {
		BadRequest(w, err, ctx, "Can't get messages of the conversation")
		return
	}

	// Write the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(messages); err != nil {
		InternalServerError(w, err, "Error encoding response", ctx)
		return
	}
}
