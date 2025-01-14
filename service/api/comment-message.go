package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"wasa.project/service/api/reqcontext"
	"wasa.project/service/api/structs"
)

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the user id from the endpoint
	userId, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Can't get the user id from the endpoint check URL")
		return
	}

	// Check if the user is authorized
	if checkAuth(w, userId, ctx) != nil {
		return
	}

	// Get the conversation id
	convId, err := strconv.Atoi(ps.ByName("dest"))
	if err != nil {
		BadRequest(w, err, ctx, "Can't get the conversation id from the endpoint check URL")
		return
	}

	// Get the conversation by the id taked from the endpoint
	conv, err := rt.db.GetConversationById(convId)
	if err != nil {
		BadRequest(w, err, ctx, "Conversatino not found")
		return
	}

	// Check if the user is part of the conversation
	if _, err := rt.db.CheckUserConv(userId, conv.ConversationId); err != nil {
		BadRequest(w, err, ctx, "The user is not in the conversation")
		return
	}

	// Get the id of the message
	msgId, err := strconv.Atoi(ps.ByName("message"))
	if err != nil {
		BadRequest(w, err, ctx, "Can't get the id of the message check the endpoint")
		return
	}

	// Take the message from the db
	var msg structs.Message
	msg, err = rt.db.GetMessageById(msgId, conv.ConversationId)
	if err != nil {
		BadRequest(w, err, ctx, "Can't get the message from the db")
		return
	}

	// Comment
	var comment structs.Comment

	// Check if the comment exist
	if comment, err = rt.db.GetCommentByUser(userId, msg.MessageId, conv.ConversationId); err == nil {

		// Take the comment from the request body
		err = json.NewDecoder(r.Body).Decode(&comment)
		if err != nil {
			BadRequest(w, err, ctx, "Error decode request body")
			return
		}

		// Check if the comment is valid
		if !comment.IsValid() {
			BadRequest(w, err, ctx, "The comment isn't a emoji")
			return
		}

		// Update comment
		err = rt.db.UpdateComment(comment.Comment, comment.CommentId, msg.MessageId, conv.ConversationId)
		if err != nil {
			BadRequest(w, err, ctx, "Error updating the comment")
			return
		}

		// Response
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("content-type", "application/json")
		if err := json.NewEncoder(w).Encode(comment); err != nil {
			InternalServerError(w, err, "Errro encode response", ctx)
			return
		}

		// Stop function
		return
	}

	// Take the comment from the request body
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		BadRequest(w, err, ctx, "Error decode request body")
		return
	}

	// Check if the comment is valid
	if !comment.IsValid() {
		BadRequest(w, err, ctx, "The comment isn't a emoji")
		return
	}

	// Set values of comment
	comment.MessageId = msg.MessageId
	comment.ConversationId = msg.ConversationId
	comment.CommentUserId = userId

	// Create the comment in the db
	comment, err = rt.db.CreateComment(comment)
	if err != nil {
		BadRequest(w, err, ctx, "Error add comment in the db")
		return
	}

	// Response
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(comment); err != nil {
		InternalServerError(w, err, "Errro encode response", ctx)
		return
	}
}
