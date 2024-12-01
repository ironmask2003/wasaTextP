package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasa.project/service/api/reqcontext"
)

func (rt *_router) deleteGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the id of the user who want delete the group
	userId, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad Request")
		return
	}

	// Get the id of the group
	groupId, err := strconv.Atoi(ps.ByName("group"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad Request")
	}

	// Check if the user is authorized
	if userId != ctx.UserId {
		Forbidden(w, nil, ctx, "Forbidden")
		return
	}

	// Check if the group exist
	_, err = rt.db.GetGroupById(groupId)
	if err != nil {
		BadRequest(w, err, ctx, "The group doesn't exist")
		return
	}

	// Check if the user is mebmber of the group
	isMember, err := rt.db.CheckMember(userId, groupId)
	if !isMember || err != nil {
		BadRequest(w, err, ctx, "The user is not a member of the group")
		return
	}

	// Delete the group from the db
	err = rt.db.DeleteGroup(groupId)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error while deleting the user")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
