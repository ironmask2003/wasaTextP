package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasa.project/service/api/reqcontext"
)

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get userId from the endpoint
	userId, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad request, can't get the userId")
		return
	}

	// Check if is authorized
	if userId != ctx.UserId {
		Forbidden(w, err, ctx, "Forbidden, you are not authorized to send a message to this conversation")
		return
	}

	// Get the conversationId from the endpoint
	convId, err := strconv.Atoi(ps.ByName("conversation"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad request, can't get the conversationId")
		return
	}

	messageId, err := strconv.Atoi(ps.ByName("message"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad request, can't get the messageId")
		return
	}

	msg, err := rt.db.GetMessage(userId, convId, messageId)
	if err != nil {
		BadRequest(w, err, ctx, "Bad request, error getting the message")
		return
	}

	// Getting the conversationId where forward message
	if !r.URL.Query().Has("dest_conv") {
		BadRequest(w, err, ctx, "Bad request, error missing the conversation id")
		return
	}

	sendConvId, err := strconv.Atoi(r.URL.Query().Get("dest_conv"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad request, error getting the conversation id where forward message")
		return
	}

	newMsg, err := rt.CreateMessageDB(Message{Text: msg.Text, SenderUserId: userId, ConversationId: sendConvId})
	if err != nil {
		BadRequest(w, err, ctx, "Bad request, can't create the message")
		return
	}

	conv, err := rt.db.GetConversationById(sendConvId, userId)
	if err != nil {
		BadRequest(w, err, ctx, "Bad request, conversation not found")
		return
	}

	// Check if the message was send to a group
	if conv.GroupId != 0 {
		convs, err := rt.db.GetConversationsByGroup(conv.GroupId)
		if err != nil {
			BadRequest(w, err, ctx, "Bad Request, failed to get conversations")
			return
		}
		for i := 0; i < len(convs); i++ {
			err = rt.db.UpdateLastMessage(convs[i].ConversationId, convs[i].UserId, newMsg.MessageId)
			if err != nil {
				BadRequest(w, err, ctx, "Bad Request, failed to update last message")
				return
			}
		}
	} else {
		convRcv, err := rt.db.GetConversationsBySender(userId, conv.SenderUserId)
		if err != nil {
			BadRequest(w, err, ctx, "Bad Request, failed to get conversations")
			return
		}
		err = rt.db.UpdateLastMessage(convRcv, conv.SenderUserId, newMsg.MessageId)
		if err != nil {
			BadRequest(w, err, ctx, "Bad Request, failed to update last message 1")
			return
		}
		err = rt.db.UpdateLastMessage(conv.ConversationId, userId, newMsg.MessageId)
		if err != nil {
			BadRequest(w, err, ctx, "Bad Request, failed to update last message")
			return
		}
	}
	w.WriteHeader(http.StatusCreated)

	// Response
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(newMsg); err != nil {
		InternalServerError(w, err, ctx)
		return
	}
}
