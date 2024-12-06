package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasa.project/service/api/reqcontext"
)

func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Getting the userId from the endpoint
	userId, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad request")
		return
	}

	// Check if the user is authorized
	if userId != ctx.UserId {
		Forbidden(w, err, ctx, "Forbidden -> the user is not authorized")
		return
	}

	// Getting the conversationId from the endpoint
	convId, err := strconv.Atoi(ps.ByName("conversation"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad request")
		return
	}

	// Get the messageId
	msgId, err := strconv.Atoi(ps.ByName("message"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad request")
		return
	}

	// Get the comment id
	commentId, err := strconv.Atoi(ps.ByName("comment"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad request")
		return
	}

	// Get the conversation
	// Get the message from the db
	_, err = rt.db.GetMessage(userId, convId, msgId)
	if err != nil {
		// get the conversation from the db
		conv, err := rt.db.GetConversationById(convId, userId)
		if err != nil {
			BadRequest(w, err, ctx, "Bad request, can't get the conversation")
			return
		}
		convId, err = rt.db.GetConversationsBySender(userId, conv.SenderUserId)
		if err != nil {
			BadRequest(w, err, ctx, "Bad request, can't get the conversation")
			return
		}
		_, err = rt.db.GetMessage(conv.SenderUserId, convId, msgId)
		if err != nil {
			BadRequest(w, err, ctx, "Bad request, can't get the message")
			return
		}
		// Check if the comment exist
		exist, err := rt.db.ExistCommentWithId(commentId, msgId, conv.SenderUserId, convId, userId)
		if !exist && err != nil {
			BadRequest(w, err, ctx, "Bad request, can't check if the comment exist")
			return
		}
		err = rt.db.DeleteComment(commentId, userId, msgId, convId, conv.SenderUserId)
		if err != nil {
			BadRequest(w, err, ctx, "Bad request, can't delete the comment")
			return
		}
	} else {
		// Check if the comment exist
		exist, err := rt.db.ExistCommentWithId(commentId, msgId, userId, convId, userId)
		if !exist && err != nil {
			BadRequest(w, err, ctx, "Bad request, can't check if the comment exist")
			return
		}

		err = rt.db.DeleteComment(commentId, userId, msgId, convId, userId)
		if err != nil {
			BadRequest(w, err, ctx, "Bad request, can't delete the comment")
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
