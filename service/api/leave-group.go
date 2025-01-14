package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasa.project/service/api/reqcontext"
)

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Take the user id of the user who want leave a group
	userId, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad Request")
		return
	}
	// Check if the user is authorized
	if userId != ctx.UserId {
		Forbidden(w, err, ctx, "Forbidden, the user is not authorized")
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
	// Delete the user from the group
	err = rt.db.LeaveGroup(userId, groupId)
	if err != nil {
		InternalServerError(w, err, "Can't delete the user from the user_group table", ctx)
		return
	}

	// Get the id of the group conversation
	convId, err := rt.db.GetConvGroup(groupId)
	if err != nil {
		InternalServerError(w, err, "Error while getting the conversation id", ctx)
		return
	}

	// Delete user from the conversation
	err = rt.db.DeleteUserConv(userId, convId)
	if err != nil {
		InternalServerError(w, err, "Error while deleting the user from the conversation", ctx)
		return
	}

	users, err := rt.db.GetMembers(groupId)
	if err != nil {
		BadRequest(w, err, ctx, "Error while getting the members of the group")
		return
	}

	if len(users) == 0 {
		err = rt.db.DeleteGroup(groupId)
		if err != nil {
			InternalServerError(w, err, "Error deleting the group", ctx)
			return
		}

		// Delete conversation
		err = rt.db.DeleteConv(convId)
		if err != nil {
			InternalServerError(w, err, "Error deleting the conversation", ctx)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
