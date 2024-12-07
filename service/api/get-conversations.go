package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasa.project/service/api/reqcontext"
)

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get user id of the user who want to get his conversations
	userId, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad Request")
		return
	}

	// Check if the user is authorized
	if userId != ctx.UserId {
		Forbidden(w, err, ctx, "Forbidden, the user is not authorized")
		return
	}

	convsDB, err := rt.db.GetConversations(userId)
	if err != nil {
		BadRequest(w, err, ctx, "Bad Request -> can't take the conversations")
		return
	}

	convs := make([]Conversation, len(convsDB))

	for i, dbConv := range convsDB {
		var conv Conversation
		prova, err := rt.db.GetConversationById(dbConv.ConversationId, userId)
		if err != nil {
			BadRequest(w, err, ctx, "Bad Request -> can't take the conversations")
			return
		}
		err = conv.ConvertConversationFromDB(prova)
		if err != nil {
			InternalServerError(w, err, ctx)
			return
		}
		convs[i] = conv
	}

	// Write the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(convs); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
