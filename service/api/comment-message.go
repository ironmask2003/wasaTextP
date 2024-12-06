package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"wasa.project/service/api/reqcontext"
)

// Used for create the conversation in the db
func (rt *_router) CreateCommentDB(c Comment) (Comment, error) {
	// Create the user in the db
	convDB, err := rt.db.CreateComment(c.ConvertCommentForDB())
	if err != nil {
		return c, err
	}

	// Convert the user from the db to the user used in the api
	err = c.ConvertCommentFromDB(convDB)
	if err != nil {
		return c, err
	}

	return c, nil
}

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the id of the user who want comment the message
	userId, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad request")
		return
	}

	// Check if the user is authorized
	if userId != ctx.UserId {
		InternalServerError(w, err, ctx)
		return
	}

	// Get the conversation from the endpoint
	convId, err := strconv.Atoi(ps.ByName("receiver"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad request")
		return
	}

	// Get the message from the endpoint
	msgId, err := strconv.Atoi(ps.ByName("message"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad request")
		return
	}

	// Getting the emoji from the body request
	var c Comment
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		BadRequest(w, err, ctx, "Bad request -> can't take the body request, check the struct")
		return
	}
	if !c.IsEmoji() {
		BadRequest(w, errors.New("The comment is not an emoji"), ctx, "Bad request -> the comment is not an emoji")
		return
	}
	c.CommentUserId = userId
	c.MessageId = msgId
	c.UserId = userId

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
		c.UserId = conv.SenderUserId
	}
	c.ConversationId = convId

	// Se esiste va modificato
	if check, err := rt.db.ExistComment(c.MessageId, c.UserId, c.ConversationId, c.CommentUserId); !check && err != nil {
		// Insert the comment in the db
		c, err = rt.CreateCommentDB(c)
		if err != nil {
			BadRequest(w, err, ctx, "Bad request, can't create the comment")
			return
		}
		w.WriteHeader(http.StatusCreated)
	} else {
		// Get the comment from the db
		cDB, err := rt.db.GetComment(c.CommentUserId, c.MessageId, c.ConversationId, c.UserId)
		if err != nil {
			BadRequest(w, err, ctx, "Bad request, can't get the comment")
			return
		}
		// Set the comment
		err = rt.db.SetComment(cDB.CommentId, c.CommentUserId, c.MessageId, c.ConversationId, c.UserId, c.Comment)
		if err != nil {
			BadRequest(w, err, ctx, "Can't update the comment")
			return
		}
		c.CommentId = cDB.CommentId
	}

	// Response
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(c); err != nil {
		InternalServerError(w, err, ctx)
		return
	}
}
