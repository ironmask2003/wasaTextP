package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/julienschmidt/httprouter"
	"wasa.project/service/api/reqcontext"
)

func (rt *_router) searchUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the search query from the request
	query_search := r.URL.Query().Get("username")
	validQuerySearch := regexp.MustCompile(`^[a-z0-9]{1,13}$`)
	if !validQuerySearch.MatchString(query_search) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	dbUsers, err := rt.db.SearchUsers(query_search)
	if err != nil {
		ctx.Logger.Error("Error searching users ", err)
		http.Error(w, "Error searching users", http.StatusInternalServerError)
		return
	}

	users := make([]User, len(dbUsers))
	for i, u := range dbUsers {
		var user User
		err := user.ConvertUserFromDB(u)
		if err != nil {
			ctx.Logger.Error("Error converting users ", err)
			http.Error(w, "Error converting users ", http.StatusInternalServerError)
			return
		}
		users[i] = user
	}

	// Write the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		ctx.Logger.Error("Error encoding users ", err)
		http.Error(w, "Error encoding response ", http.StatusInternalServerError)
		return
	}
}
