package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/url"
	"strconv"
	"wasa.project/service/api/reqcontext"
	"wasa.project/service/api/structs"
)

// Function used to check if the query have the message to response
// Check if the query has the msg id
func check_query(query url.Values) (int, error) {
	id_msg_response := 0
	var err error

	if query.Has("msg") {
		id_msg_response, err = strconv.Atoi(query.Get("msg"))
		if err != nil {
			return 0, err
		}
	}

	return id_msg_response, nil
}

// Function used to send a message to a conversation
func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the id of the user who want send the message
	userId, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Can't take the user id from the endpoint")
		return
	}

	// Check if the user is authorized
	if checkAuth(w, userId, ctx) != nil {
		Forbidden(w, err, ctx, "The user is not authorized")
		return
	}

	// Get the conversation id where send the message
	convId, err := strconv.Atoi(ps.ByName("conv"))
	if err != nil {
		BadRequest(w, err, ctx, "Can't take the conversation id from the endpoint")
		return
	}

	// Message
	var msg structs.Message

	// Get the conversation by the id taked from the endpoint
	conv, err := rt.db.GetConversationById(convId)
	if err != nil {
		BadRequest(w, err, ctx, "Conversatino not found")
		return
	}

	// Check if the user is part of the conversation
	if _, err := rt.db.CheckUserConv(userId, convId); err != nil {
		BadRequest(w, err, ctx, "The user is not in the conversation")
		return
	}

	// Take the message from the request body
	err = json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		BadRequest(w, err, ctx, "Error decoding the request body")
		return
	}

	// Set the id of the conversation
	msg.ConversationId = conv.ConversationId
	msg.SenderUserId = userId
	msg.Status = "Sended"

	// query message
	type Response struct {
		MsgQuery  structs.Message `json:"msgQuery"`
		MsgSended structs.Message `json:"msgSended"`
	}
	var response Response

	// Check if the query has the msg id
	id, err := check_query(r.URL.Query())
	if err != nil {
		BadRequest(w, err, ctx, "Error taking the query")
		return
	}
	if id != 0 {
		// Get the message by the id
		msgQuery, err := rt.db.GetMessageById(id, conv.ConversationId)
		if err != nil {
			BadRequest(w, err, ctx, "Error taking the message by the id")
			return
		}
		// Set the message query in the response
		response.MsgQuery = msgQuery
	}

	// Insert the message in the db
	msg, err = rt.db.CreateMessage(msg)
	if err != nil {
		BadRequest(w, err, ctx, "Error insert the message in the db")
		return
	}

	// Update the last message
	err = rt.db.UpdateLastMessage(msg.MessageId, conv.ConversationId)
	if err != nil {
		BadRequest(w, err, ctx, "Error updating last message id")
		return
	}

	// Set the message sended in the response
	response.MsgSended = msg

	// Resposne
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		InternalServerError(w, err, "Error encoding response", ctx)
		return
	}
}
