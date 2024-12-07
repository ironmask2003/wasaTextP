package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasa.project/service/api/reqcontext"
)

func (rt *_router) deleteUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	userID := ctx.UserId

	// Check if the user is authorized
	if userId != userID {
		Forbidden(w, err, ctx, "Forbidden, the user is not authorized")
		return
	}

	// Delete the user from the database
	err = rt.db.DeleteUser(userId)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
