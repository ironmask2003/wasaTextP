package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasa.project/service/api/reqcontext"
)

// Used for create the message in the db
func (rt *_router) CreateMessageDB(m Message) (Message, error) {
	// Create the user in the db
	msgDB, err := rt.db.CreateMessage(m.ConvertMessageForDB())
	if err != nil {
		return m, err
	}

	// Convert the user from the db to the user used in the api
	err = m.ConvertMessageFromDB(msgDB)
	if err != nil {
		return m, err
	}

	return m, nil
}

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the user who want send the message
	userId, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad Request")
		return
	}

	// Check if the user is authorized
	if userId != ctx.UserId {
		Forbidden(w, err, ctx, "Forbidden")
		return
	}

	// Take the conversation id
	convId, err := strconv.Atoi(ps.ByName("conversation"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad Request")
		return
	}

	// Get the conversation
	conv, err := rt.db.GetConversationById(convId, userId)
	if err != nil {
		BadRequest(w, err, ctx, "Bad request, can't get the conversation")
		return
	}

	var group bool

	if conv.GroupId != 0 {
		// Check if the user is member of the group
		isMember, err := rt.db.CheckMember(userId, conv.GroupId)
		if !isMember || err != nil {
			BadRequest(w, err, ctx, "Bad request, the user is not member of the group")
			return
		}
		group = true
	}

	// Take the message
	type TextBody struct {
		Message string `json:"message"`
	}
	var txt TextBody

	if err := json.NewDecoder(r.Body).Decode(&txt); err != nil {
		BadRequest(w, err, ctx, "Bad request -> can't take the body request, check the struct")
		return
	}

	msg, err := rt.CreateMessageDB(Message{Text: txt.Message, SenderUserId: userId, ConversationId: convId})
	if err != nil {
		BadRequest(w, err, ctx, "Bad request, can't create the message")
		return
	}

	// Check if the message was send to a group
	if group {
		convs, err := rt.db.GetConversationsByGroup(conv.GroupId)
		if err != nil {
			BadRequest(w, err, ctx, "Bad Request, failed to get conversations")
			return
		}
		for i := 0; i < len(convs); i++ {
			err = rt.db.UpdateLastMessage(convs[i].ConversationId, msg.MessageId)
			if err != nil {
				BadRequest(w, err, ctx, "Bad Request, failed to update last message")
				return
			}
		}
	}
	w.WriteHeader(http.StatusCreated)

	// Response
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(msg); err != nil {
		InternalServerError(w, err, ctx)
		return
	}

}
