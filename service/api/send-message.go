package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"wasa.project/service/api/reqcontext"
)

// Used for create the message in the db
func (rt *_router) CreateMessageDB(m Message) (Message, error) {
	// Create the message in the db
	msgDB, err := rt.db.CreateMessage(m.ConvertMessageForDB())
	if err != nil {
		return m, err
	}

	// Convert the message from the db to the message used in the api
	err = m.ConvertMessageFromDB(msgDB)
	if err != nil {
		return m, err
	}

	return m, nil
}

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
			err = rt.db.UpdateLastMessage(convs[i].ConversationId, convs[i].UserId, msg.MessageId)
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
		err = rt.db.UpdateLastMessage(convRcv, conv.SenderUserId, msg.MessageId)
		if err != nil {
			BadRequest(w, err, ctx, "Bad Request, failed to update last message 1")
			return
		}
		err = rt.db.UpdateLastMessage(conv.ConversationId, userId, msg.MessageId)
		if err != nil {
			BadRequest(w, err, ctx, "Bad Request, failed to update last message")
			return
		}
	}
	w.WriteHeader(http.StatusCreated)

	// Get the id of the message of the query
	msgId, err := check_query(r.URL.Query())
	if err != nil {
		BadRequest(w, err, ctx, "Bad Request, have the msg response but can't get it")
		return
	}

	if msgId != 0 {
		// Check if the msgId is present in the conversation
		if _, err := rt.db.CheckMessageConv(msgId, conv.ConversationId, userId); err != nil {
			BadRequest(w, err, ctx, "Bad Request, the message isn't in the conversation")
			return
		}
		type Response struct {
			MsgResponse Message `json:"msgResponse"`
			IdResponse  int     `json:"idResponse"`
		}
		w.Header().Set("content-type", "application/json")
		if err := json.NewEncoder(w).Encode(Response{MsgResponse: msg, IdResponse: msgId}); err != nil {
			InternalServerError(w, err, ctx)
			return
		}
		return
	}

	// Response
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(msg); err != nil {
		InternalServerError(w, err, ctx)
		return
	}

}
