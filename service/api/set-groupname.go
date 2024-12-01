package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasa.project/service/api/reqcontext"
)

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the id of the user who want change the name of the group
	userId, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad Request")
		return
	}

	// Check if the user is authorized
	if userId != ctx.UserId {
		Forbidden(w, nil, ctx, "Forbidden")
		return
	}

	// Get the id of the group
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

	// New name of the group
	var group Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		BadRequest(w, err, ctx, "Bad request")
		return
	}
	// Check if the new username respect the regex
	if !group.IsValid() {
		BadRequest(w, err, ctx, "Invalid Name")
		return
	}

	// Change username
	if err := rt.db.SetGroupName(groupId, group.GroupName); err != nil {
		BadRequest(w, err, ctx, "Name not setted")
		return
	}

	// Set the id of the group from the endpoint
	group.GroupId = groupId

	// Take the group from the db
	err = group.ConvertGroupFromDB(group.ConvertGroupForDB())
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	// Username changed, resposne 200
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(group); err != nil {
		InternalServerError(w, err, ctx)
		return
	}

}
