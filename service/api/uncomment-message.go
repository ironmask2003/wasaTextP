package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasa.project/service/api/reqcontext"
	"wasa.project/service/api/structs"
)

func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the user id from the endpoint
	userId, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Can't get the user id from the endpoint")
		return
	}

	// Check if the user is authorized
	if checkAuth(w, userId, ctx) != nil {
		return
	}

	// Get the conversation id
	convId, err := strconv.Atoi(ps.ByName("conv"))
	if err != nil {
		BadRequest(w, err, ctx, "Can't get the conversation id from the endpoint")
		return
	}

	// Get the conversation
	var conv structs.Conversation
	conv, err = rt.db.GetConversationById(convId)
	if err != nil {
		BadRequest(w, err, ctx, "Can't get the conversation from the db")
		return
	}

	// Check if the user is in the conversation
	if check, err := rt.db.CheckUserConv(userId, conv.ConversationId); check || err != nil {
		BadRequest(w, err, ctx, "The user isn't in the conversation")
		return
	}

	// Get the message id from the endpoint
	msgId, err := strconv.Atoi(ps.ByName("message"))
	if err != nil {
		BadRequest(w, err, ctx, "Can't get the message id from the endpoint")
		return
	}

	// Get the message
	var msg structs.Message
	msg, err = rt.db.GetMessageById(msgId, conv.ConversationId)
	if err != nil {
		BadRequest(w, err, ctx, "Can't get the message from the db")
		return
	}

	// Get the comment id
	comId, err := strconv.Atoi(ps.ByName("comment"))
	if err != nil {
		BadRequest(w, err, ctx, "Can't get the commen id from the endpoint")
		return
	}

	// Get the comment
	var comment structs.Comment
	comment, err = rt.db.GetCommentByUser(userId, msg.MessageId, conv.ConversationId)
	if err != nil || comment.CommentId != comId {
		BadRequest(w, err, ctx, "Bad request")
		return
	}

	// Delete the comment
	err = rt.db.DeleteComment(comment.CommentId, msg.MessageId, conv.ConversationId)
	if err != nil {
		BadRequest(w, err, ctx, "Can't delete the commet")
		return
	}

	// Response
	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("content-type", "plain/text")
	if err := json.NewEncoder(w).Encode("User deleted"); err != nil {
		InternalServerError(w, err, "Error encoding the response", ctx)
		return
	}
}
