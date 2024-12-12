package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"wasa.project/service/api/reqcontext"
	"wasa.project/service/api/structs"
)

func (rt *_router) createConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the userId from the endpoint
	userId, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Can't take the user id from the endpoint")
		return
	}

	// Check if the user is authorized
	if checkAuth(w, userId, ctx) != nil {
		return
	}

	// Take the destination user
	destId, err := strconv.Atoi(ps.ByName("dest"))
	if err != nil {
		BadRequest(w, err, ctx, "Can't take te id of the user or the group where send the message")
		return
	}

	// Check if the conversation exist
	if check, err := rt.db.CheckIfExistConv(userId, destId); check || err != nil {
		BadRequest(w, err, ctx, "The conversation already exist")
		return
	}

	// Take the message to sent from the Request
	var msgRequest structs.Message
	err = json.NewDecoder(r.Body).Decode(&msgRequest)
	if err != nil {
		BadRequest(w, err, ctx, "Error reading the request body")
		return
	}
	msgRequest.SenderUserId = userId
	msgRequest.Status = "Sended"

	// get user by id
	destUserDB, err := rt.db.GetUserById(destId)
	if err != nil {
		BadRequest(w, err, ctx, "Error taking the user from the user table")
		return
	}

	// Convert the user from the db
	var destUser User
	err = destUser.ConvertUserFromDB(destUserDB)
	if err != nil {
		BadRequest(w, err, ctx, "Bad request")
		return
	}

	// New conversation
	var c structs.Conversation
	c.IsGroup = false

	// Create the conversation in the db
	c, err = rt.db.CreateConversation(c)
	if err != nil {
		BadRequest(w, err, ctx, "Error adding the conversation in the db")
		return
	}

	// Adding the link of the user and the conversation
	if rt.db.AddUserConv(c.ConversationId, userId) != nil {
		BadRequest(w, err, ctx, "Error adding in the conversation_user table")
		return
	}

	if rt.db.AddUserConv(c.ConversationId, destUser.UserId) != nil {
		BadRequest(w, err, ctx, "Error adding the receiver in the conversation_user table")
		return
	}

	// Setting the conversation id of the message and adding in the db
	msgRequest.ConversationId = c.ConversationId
	msgRequest, err = rt.db.CreateMessage(msgRequest)
	if err != nil {
		BadRequest(w, err, ctx, "Error adding message in the db")
		return
	}

	// Update the last message id in the conversation
	err = rt.db.UpdateLastMessage(msgRequest.MessageId, c.ConversationId)
	if err != nil {
		BadRequest(w, err, ctx, "Error updating the last message id of the conversation")
		return
	}
	c.LastMessageId = msgRequest.MessageId

	w.WriteHeader(http.StatusCreated)

	// Response
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(c); err != nil {
		InternalServerError(w, err, "Error encoding resposne", ctx)
		return
	}
}
