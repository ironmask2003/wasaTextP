package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"wasa.project/service/api/reqcontext"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the userID from the URL
	profileUserID, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad request")
		return
	}

	userId := ctx.UserId
	// Check if the user is authorized
	if profileUserID != userId {
		Forbidden(w, err, ctx, "Forbidden")
		return
	}

	// Take the user object from the bosy of the request
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		BadRequest(w, err, ctx, "Bad request")
		return
	}
	// Check if the new username respect the regex
	if !user.IsValid() {
		BadRequest(w, err, ctx, "Invalid username")
		return
	}

	// Change username
	if err := rt.db.SetMyUsername(userId, user.Username); err != nil {
		BadRequest(w, err, ctx, "Username already used")
		return
	}

	// Set the id of the user from the endpoint
	user.UserId = userId

	// Take the user from the db
	err = user.ConvertUserFromDB(user.ConvertUserForDB())
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	// Username changed, resposne 200
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		InternalServerError(w, err, ctx)
		return
	}
}
