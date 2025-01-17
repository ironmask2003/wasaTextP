package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"wasa.project/service/api/reqcontext"
)

// Used for create the user in the db if his dosen't exist
func (rt *_router) CreateUser(u User) (User, error) {
	// Create the user in the db
	dbUser, err := rt.db.CreateUser(u.ConvertUserForDB())
	if err != nil {
		return u, err
	}

	// Convert the user from the db to the user used in the api
	err = u.ConvertUserFromDB(dbUser)
	if err != nil {
		return u, err
	}

	return u, nil
}

// Main function used for the login/registrasion
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User

	// Read the request body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		BadRequest(w, err, ctx, "Bad request")
		return
	}

	// Check if the patter of the username respect the regex
	if !user.IsValid() {
		BadRequest(w, err, ctx, "Invalid username")
		return
	}

	// Check if the username is already used
	exist, err := rt.db.CheckIfExist(user.Username)
	if err != nil {
		InternalServerError(w, err, "Can't check if the user exist", ctx)
		return
	}

	// If doesn't exist create the user
	if !exist {
		user, err = rt.CreateUser(user)
		if err != nil {
			InternalServerError(w, err, "Can't add user in the user table", ctx)
			return
		}
		w.WriteHeader(http.StatusCreated)
	} else {
		// If exist find the user in the db with username and take the info
		// Find the user in the db
		dbUser, err := rt.db.GetUserByName(user.Username)
		if err != nil {
			InternalServerError(w, err, "Error getting user in the request body from user table", ctx)
			return
		}
		err = user.ConvertUserFromDB(dbUser)
		if err != nil {
			InternalServerError(w, err, "Error converting the user from the database struct", ctx)
			return
		}
		w.WriteHeader(http.StatusOK)
	}

	// Respose
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		InternalServerError(w, err, "Error encoding the response", ctx)
		return
	}
}
