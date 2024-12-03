package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasa.project/service/api/reqcontext"
)

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get userId from the endpoint
	userId, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad request, can't get the userId")
		return
	}
	// Check if is authorized
	if userId != ctx.UserId {
		Forbidden(w, err, ctx, "Forbidden, you are not authorized to delete a message in this conversation")
		return
	}
	// Get the conversationId from the endpoint
	convId, err := strconv.Atoi(ps.ByName("conversation"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad request, can't get the conversationId")
		return
	}
	messageId, err := strconv.Atoi(ps.ByName("message"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad request, can't get the messageId")
		return
	}

	err = rt.db.DeleteMessage(userId, convId, messageId)
	if err != nil {
		BadRequest(w, err, ctx, "Bad request, error deleting the message")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
