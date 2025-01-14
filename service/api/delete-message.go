package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"wasa.project/service/api/reqcontext"
	"wasa.project/service/api/structs"
)

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the id of the user who want delete message
	userId, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Can't get the user id from the endpoint check it")
		return
	}

	// Check if the user is authorized
	if checkAuth(w, userId, ctx) != nil {
		return
	}

	// Get the id of the conversation
	convId, err := strconv.Atoi(ps.ByName("conv"))
	if err != nil {
		BadRequest(w, err, ctx, "Can't get the conversation id from the endpoint check URL")
		return
	}

	// Check if the conversation exist taking it from the db
	var conv structs.Conversation
	conv, err = rt.db.GetConversationById(convId)
	if err != nil {
		BadRequest(w, err, ctx, "Can't get the conversation from the db")
		return
	}

	// Check if the user is in the conversation
	if check, err := rt.db.CheckUserConv(userId, conv.ConversationId); !check || err != nil {
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

	// Get the max id of the message table
	maxId, err := rt.db.GetMaxMessageId(msg.ConversationId)
	if err != nil {
		BadRequest(w, err, ctx, "Can't get the max id of the message table")
		return
	}

	if maxId == msg.MessageId {
		// Check if the maxId is the first message
		if maxId == 1 {
			err = rt.db.UpdateLastMessage(0, msg.ConversationId)
			if err != nil {
				BadRequest(w, err, ctx, "Error updating last message with NULL")
				return
			}
		} else {
			// Update
			maxId -= 1
			err = rt.db.UpdateLastMessage(maxId, msg.ConversationId)
			for err != nil {
				maxId -= 1
				err = rt.db.UpdateLastMessage(maxId, msg.ConversationId)
			}
		}
	}

	// Delete the message from the db
	err = rt.db.DeleteMessage(msg.MessageId, msg.ConversationId)
	if err != nil {
		BadRequest(w, err, ctx, "Can't delete the message "+strconv.Itoa(msg.MessageId)+" "+strconv.Itoa(msg.ConversationId))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "plain/text")
	if err := json.NewEncoder(w).Encode("Message deleted"); err != nil {
		InternalServerError(w, err, "Error encoding the response", ctx)
		return
	}

	// Delete the message from the db
	err = rt.db.DeleteMessage(msg.MessageId, conv.ConversationId)
	if err != nil {
		BadRequest(w, err, ctx, "Can't delete the message")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Header().Set("content-type", "plain/text")
	if err := json.NewEncoder(w).Encode("Message deleted"); err != nil {
		InternalServerError(w, err, "Error encoding the response", ctx)
		return
	}
}
