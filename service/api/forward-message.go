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

	// Get the id of the destination user
	convId, err := strconv.Atoi(ps.ByName("conv"))
	if err != nil {
		BadRequest(w, err, ctx, "Error getting the conversation id")
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
	if check, err := rt.db.CheckUserConv(userId, conv.ConversationId); !check || err != nil {
		BadRequest(w, err, ctx, "The user isn't in the conversation")
		return
	}

	// Query conversation
	var destConv structs.Conversation

	// Get the conversation where forward the message
	if r.URL.Query().Has("dest_user") {
		destUser, err := strconv.Atoi(r.URL.Query().Get("dest_user"))
		if err != nil {
			BadRequest(w, err, ctx, "Can't get the conversation id from the query")
			return
		}
		// Check if the conversation between userId and destUser exist
		if check, err := rt.db.CheckIfExistConv(userId, destUser); !check {
			if err != nil {
				BadRequest(w, err, ctx, "Can't check")
				return
			}
			destConv.GroupId = 0
			// Create the conversation
			destConv, err = rt.db.CreateConversation(destConv)
			if err != nil {
				BadRequest(w, err, ctx, "Error creating the covnersation")
				return
			}

			// Adding the link of the user and the conversation
			if rt.db.AddUserConv(destConv.ConversationId, userId) != nil {
				BadRequest(w, err, ctx, "Error adding in the conversation_user table")
				return
			}

			if rt.db.AddUserConv(destConv.ConversationId, destUser) != nil {
				BadRequest(w, err, ctx, "Error adding the receiver in the conversation_user table")
				return
			}
		} else {
			// Get the conversation if exist
			destConvId, err := rt.db.GetConversation(userId, destUser)
			if err != nil {
				BadRequest(w, err, ctx, "Error getting the conversation id with the receiver")
				return
			}
			destConv, err = rt.db.GetConversationById(destConvId)
			if err != nil {
				BadRequest(w, err, ctx, "Erro getting the conversation with the receiver")
				return
			}
		}
	} else {
		BadRequest(w, err, ctx, "Missing the user id from query")
		return
	}

	// Check if the userId is in the other conversation
	if check, err := rt.db.CheckUserConv(userId, destConv.ConversationId); !check || err != nil {
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
	newMsg.Photo = msg.Photo

	// Create the message in the new conversation
	newMsg, err = rt.db.CreateMessage(newMsg)
	if err != nil {
		BadRequest(w, err, ctx, "Error insert message in the db")
		return
	}

	// Update last message in the conversation
	err = rt.db.UpdateLastMessage(newMsg.MessageId, destConv.ConversationId)
	if err != nil {
		BadRequest(w, err, ctx, "Error updating last message of the conversation")
		return
	}

	// Response
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(newMsg); err != nil {
		InternalServerError(w, err, "Error encoding resposne", ctx)
		return
	}
}
