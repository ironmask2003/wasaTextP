package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"regexp"
	"wasa.project/service/api/reqcontext"
)

/*
searchUsers is the handler for the GET /users/search endpoint
It returns the users that match the given search query
*/
func (rt *_router) searchUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the search query from the request
	query_search := r.URL.Query().Get("search")
	validQuerySearch := regexp.MustCompile(`^[a-z][a-z0-9]{2,13}$`)
	if !validQuerySearch.MatchString(query_search) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	userID := ctx.UserId

	dbUsers, err := rt.db.SearchUsers(userID, query_search)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	// Write the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dbUsers); err != nil {
		ctx.Logger.Error("Error encoding users ", err)
		http.Error(w, "Error encoding response ", http.StatusInternalServerError)
		return
	}

}
