package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"wasa.project/service/api/reqcontext"
)

func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Take the user id of the user who want add member to a group
	userId, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad Request")
		return
	}

	// Check if the user is authorized
	if checkAuth(w, userId, ctx) != nil {
		return
	}

	// Group ID
	groupId, err := strconv.Atoi(ps.ByName("group"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad Request")
		return
	}

	// Check if the group exist
	_, err = rt.db.GetGroupById(groupId)
	if err != nil {
		BadRequest(w, err, ctx, "The group doesn't exist")
		return
	}

	// Check if the user is a member of the group
	isMember, err := rt.db.CheckMember(userId, groupId)
	if !isMember || err != nil {
		BadRequest(w, err, ctx, "The user is not a member of the group")
		return
	}

	// Struct to take the body of the request
	type UserToAdd struct {
		Users []User `json:"users"`
	}
	var contentRequest UserToAdd

	// Take the content of the request
	if err := json.NewDecoder(r.Body).Decode(&contentRequest); err != nil {
		BadRequest(w, err, ctx, "Bad request, can't take the body of the request, check the struct or the http reequst")
		return
	}

	// List of users
	user := contentRequest.Users

	for i := 0; i < len(user); i++ {
		userDB, err := rt.db.GetUserByName(user[i].Username)
		if err != nil {
			InternalServerError(w, err, "Error getting user in the request body from user table", ctx)
			return
		}

		err = user[i].ConvertUserFromDB(userDB)
		if err != nil {
			InternalServerError(w, err, "Error converting the user from the database struct", ctx)
			return
		}

		// Check if the user is a member of the group
		isMember, err := rt.db.CheckMember(user[i].UserId, groupId)
		if !isMember || err != nil {
			// Check if the username is already used
			exist, err := rt.db.CheckIfExist(user[i].Username)
			if err != nil {
				InternalServerError(w, err, "Can't check if the user exist", ctx)
				return
			}

			if !exist {
				BadRequest(w, err, ctx, "The user doesn't exist")
				return
			}

			// Add
			err = rt.db.AddUserGroup(user[i].UserId, groupId)
			if err != nil {
				InternalServerError(w, err, "Error adding the user in the group", ctx)
				return
			}
		}
	}

	// Response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "plain/text")
	if err := json.NewEncoder(w).Encode("Member added"); err != nil {
		InternalServerError(w, err, "Error econding the response", ctx)
		return
	}

}
