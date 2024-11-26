package api

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"wasa.project/service/api/reqcontext"
)

func (rt *_router) CreateConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the id of the user who want to create a conversation
	user_id, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad request")
		return
	}

	// Get the id of the receiver
	receiver, err := strconv.Atoi(ps.ByName("dest_user_id"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad request")
		return
	}

	// Check if the user_id is the same of the user who called the action
	if user_id != ctx.UserId {
		BadRequest(w, nil, ctx, "Bad request")
		return
	}

	// Check if the receiver is different from the user
	if user_id == receiver {
		BadRequest(w, nil, ctx, "Bad request")
		return
	}

	return
}
