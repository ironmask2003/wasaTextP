package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasa.project/service/api/reqcontext"
)

// Used for create the group in the db
func (rt *_router) CreateGroupDB(g Group, userId int) (Group, error) {
	// Create the user in the db
	groupDB, err := rt.db.CreateGroup(g.ConvertGroupForDB(), userId)
	if err != nil {
		return g, err
	}

	// Convert the user from the db to the user used in the api
	err = g.ConvertGroupFromDB(groupDB)
	if err != nil {
		return g, err
	}

	return g, nil
}

func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the user who want create the group
	userId, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad Request")
		return
	}

	// Check if the user is authorized
	if userId != ctx.UserId {
		InternalServerError(w, err, ctx)
		return
	}

	// Group to return on the Respose
	var g Group

	// Struct used for the body request
	type RequestBodyCG struct {
		GroupName string `json:"groupName"`
		Users     []User `json:"users"`
	}
	var body RequestBodyCG

	// Take the gorup name
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		BadRequest(w, err, ctx, "Bad request -> can't take the body request, check the struct")
		return
	}

	g.GroupName = body.GroupName

	// Creation of the group in the db
	g, err = rt.CreateGroupDB(g, userId)
	if err != nil {
		BadRequest(w, err, ctx, "Bad request, can't create the group")
		return
	}
	w.WriteHeader(http.StatusCreated)

	// List of users to add in the group
	user := body.Users

	for i := 0; i < len(user); i++ {
		userDB, err := rt.db.GetUserByName(user[i].Username)
		if err != nil {
			InternalServerError(w, err, ctx)
			return
		}

		err = user[i].ConvertUserFromDB(userDB)
		if err != nil {
			InternalServerError(w, err, ctx)
			return
		}

		// Controlla se l'utente è gia stato aggiunto al gruppo
		if check, err := rt.db.CheckMember(user[i].UserId, g.GroupId); !check && err != nil {
			// Check if the username is already used
			exist, err := rt.db.CheckIfExist(user[i].Username)
			if err != nil {
				InternalServerError(w, err, ctx)
				return
			}

			if !exist {
				BadRequest(w, err, ctx, "The user dosn't exist")
				return
			}

			// Add
			err = rt.db.AddUserGroup(user[i].UserId, g.GroupId)
			if err != nil {
				InternalServerError(w, err, ctx)
				return
			}
		}
	}

	// Respose
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(g); err != nil {
		InternalServerError(w, err, ctx)
		return
	}

}
