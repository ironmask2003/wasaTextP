package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"wasa.project/service/api/reqcontext"
	"wasa.project/service/api/structs"
)

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the id of the user who want Forwar the message
	userId, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Can't take the user id from the endpoint")
		return
	}

	// Check if the user is authorized
	if checkAuth(w, userId, ctx) != nil {
		return
	}

	// Get the id of the conversation
	convId, err := strconv.Atoi(ps.ByName("conv"))
	if err != nil {
		BadRequest(w, err, ctx, "Can't take the conversation id from the endpoint")
		return
	}

	// Get the conversation from the db
	var conv structs.Conversation
	conv, err = rt.db.GetConversationById(convId)
	if err != nil {
		BadRequest(w, err, ctx, "Can't take the conversation from the db")
		return
	}

	// Check if the user is in the Conversation
<<<<<<< HEAD
	if check, err := rt.db.CheckUserConv(userId, conv.ConversationId); !check && err != nil {
		BadRequest(w, err, ctx, "The user is not in the conversation")
=======
	if check, err := rt.db.CheckUserConv(userId, conv.ConversationId); check || err != nil {
		BadRequest(w, err, ctx, "The conversation already exist")
>>>>>>> 779b51a (Modified table and function)
		return
	}

	// Query conversation
	var destConv structs.Conversation

	// Get the conversation where forward the message
	if r.URL.Query().Has("dest_conv") {
		destConvId, err := strconv.Atoi(r.URL.Query().Get("dest_conv"))
		if err != nil {
			BadRequest(w, err, ctx, "Can't get the conversation id from the query")
			return
		}
		// Get the conversation
		destConv, err = rt.db.GetConversationById(destConvId)
		if err != nil {
			BadRequest(w, err, ctx, "Error getting the conversation from the db")
		}
	} else {
		BadRequest(w, err, ctx, "Missing the conversation id from query")
		return
	}

	// Check if the userId is in the other conversation
<<<<<<< HEAD
	if check, err := rt.db.CheckUserConv(userId, destConv.ConversationId); !check && err != nil {
=======
	if check, err := rt.db.CheckUserConv(userId, destConv.ConversationId); check || err != nil {
>>>>>>> 779b51a (Modified table and function)
		BadRequest(w, err, ctx, "The user isn't in the conversation")
		return
	}

	// Getting the message id from the endpoint
	msgId, err := strconv.Atoi(ps.ByName("message"))
	if err != nil {
		BadRequest(w, err, ctx, "Can't take the message id, check the endpoint")
		return
	}

	// Getting the message from the id
	msg, err := rt.db.GetMessageById(msgId, conv.ConversationId)
	if err != nil {
		BadRequest(w, err, ctx, "Can't get the message from the db")
		return
	}

	// New message
	var newMsg structs.Message

	// Setting the value of the new message
	newMsg.ConversationId = destConv.ConversationId
	newMsg.SenderUserId = userId
	newMsg.Text = msg.Text
	newMsg.Status = "Sended"

	// Create the message in the new conversation
	newMsg, err = rt.db.CreateMessage(newMsg)
	if err != nil {
		BadRequest(w, err, ctx, "Error insert message in the db")
		return
<<<<<<< HEAD
	}

	// Update last message in the conversation
	err = rt.db.UpdateLastMessage(newMsg.MessageId, destConv.ConversationId)
	if err != nil {
		BadRequest(w, err, ctx, "Error updating last message of the conversation")
		return
=======
>>>>>>> 779b51a (Modified table and function)
	}

	// Response
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(newMsg); err != nil {
		InternalServerError(w, err, "Error encoding resposne", ctx)
		return
	}
}
