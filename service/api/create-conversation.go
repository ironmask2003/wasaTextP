package api

import (
	"bytes"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"strconv"
	"wasa.project/service/api/reqcontext"
)

// Used for create the conversation in the db
func (rt *_router) CreateConversationDB(c Conversation) (Conversation, error) {
	// Create the user in the db
	convDB, err := rt.db.CreateConversation(c.ConvertConversationForDB())
	if err != nil {
		return c, err
	}

	// Convert the user from the db to the user used in the api
	err = c.ConvertConversationFromDB(convDB)
	if err != nil {
		return c, err
	}

	return c, nil
}

func (rt *_router) createConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the user who want create the conversation
	userId, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad Request")
		return
	}

	// Check if the user is authorized
	if userId != ctx.UserId {
		InternalServerError(w, err, ctx)
		return
	}

	// Get the receiver
	rcvId, err := strconv.Atoi(ps.ByName("receiver"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad Request")
		return
	}

	// Check if the bosy i empty
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		BadRequest(w, err, ctx, "Bad Request -> error reading the body")
		return
	}

	// Ripristina il body in modo che possa essere letto successivamente
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Creazione della conversazione
	var c Conversation
	c.UserId = userId

	var user User

	// Controlla se il body è vuoto
	userDB, err := rt.db.GetUserById(rcvId)
	if err != nil {
		BadRequest(w, err, ctx, "Bad Request -> error taking the user")
		return
	}
	err = user.ConvertUserFromDB(userDB)
	if err != nil {
		BadRequest(w, err, ctx, "Bad Request -> error converting the user from the db")
	}
	c.SenderUserId = user.UserId

	// Creation of the group in the db
	c, err = rt.CreateConversationDB(c)
	if err != nil {
		BadRequest(w, err, ctx, "Bad request, can't create the conversation with user")
		return
	}

	// Take the message
	type TextBody struct {
		Message string `json:"message"`
	}
	// Prendere il messaggio
	var txt TextBody

	if err := json.NewDecoder(r.Body).Decode(&txt); err != nil {
		BadRequest(w, err, ctx, "Bad request -> can't take the body request, check the struct")
		return
	}

	msg, err := rt.CreateMessageDB(Message{Text: txt.Message, SenderUserId: userId, ConversationId: c.ConversationId})
	if err != nil {
		BadRequest(w, err, ctx, "Bad request, can't create the message")
		return
	}
	err = rt.db.UpdateLastMessage(msg.ConversationId, msg.MessageId)
	if err != nil {
		BadRequest(w, err, ctx, "Bad Request, failed to update last message")
	}
	w.WriteHeader(http.StatusCreated)

	// Response
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(c); err != nil {
		InternalServerError(w, err, ctx)
		return
	}

}
