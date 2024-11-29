package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"wasa.project/service/api/reqcontext"
	"wasa.project/service/database"
)

// Used for create the conversation in the db if his dosen't exist
func (rt *_router) CreateConversationDB(c Conversation, m Message) (Conversation, error) {
	if m.MessageId == 0 {
		// Create the user in the db
		dbConversation, err := rt.db.CreateConversation(c.ConvertConversationForDB(), database.Message{})
		if err != nil {
			return c, err
		}
		// Convert the user from the db to the user used in the api
		c.ConvertConversationFromDB(dbConversation)
	} else {
		dbConversation, err := rt.db.CreateConversation(c.ConvertConversationForDB(), m.ConvertMessageForDB())
		if err != nil {
			return c, err
		}
		// Convert the user from the db to the user used in the api
		c.ConvertConversationFromDB(dbConversation)
	}
	return c, nil
}

// Function used to check if the rcv is a user
func (rt *_router) CheckIfRcvUser(rcv int) (User, error) {
	var user User

	userDB, err := rt.db.GetUserById(rcv)
	if err != nil {
		return user, err
	}

	user.ConvertUserFromDB(userDB)

	return user, nil
}

// Function used to check if the rcv is a group
func (rt *_router) CheckIfRcvGroup(rcv int) (Group, error) {
	var group Group

	groupDB, err := rt.db.GetGroupById(rcv)
	if err != nil {
		return group, err
	}

	group.ConvertGroupFromDB(groupDB)

	return group, nil
}

func (rt *_router) CreateConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Se nella request body non c'e un messaggio allora l'id del receiver sarà di un gruppo
	var m Message
	err_m := json.NewDecoder(r.Body).Decode(&m)

	// Get the id of the user who want to create a conversation
	user_id, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad request")
		return
	}

	// Get the id of the receiver
	receiver, err := strconv.Atoi(ps.ByName("dest_user_id"))
	if err != nil {
		BadRequest(w, err, ctx, "Bad request")
		return
	}

	// Check if the user_id is the same of the user who called the action
	if user_id != ctx.UserId {
		BadRequest(w, nil, ctx, "Bad request")
		return
	}

	// Check if the receiver is different from the user
	if user_id == receiver {
		BadRequest(w, nil, ctx, "Bad request")
		return
	}

	// Conversation to return
	var conv Conversation

	// Check if the receiver is a user or a group
	if user, err := rt.CheckIfRcvUser(receiver); err == nil {
		// Check if the conversation exist
		exist, err := rt.db.CheckIfExistConversationWithUser(user_id, user.UserId)
		if err != nil {
			InternalServerError(w, err, ctx)
			return
		}

		if !exist {
			if err_m != nil {
				BadRequest(w, err_m, ctx, "Bad request")
			}

			// Create the conversation
			conv, err = rt.CreateConversationDB(Conversation{UserId: user_id, SenderUserId: receiver}, m)

			if err != nil {
				InternalServerError(w, err, ctx)
				return
			}
		}
	} else if group, err := rt.CheckIfRcvGroup(receiver); err == nil {
		// Check if the conversation exist
		exist, err := rt.db.CheckIfExistConversationWithGroup(user_id, group.GroupId)
		if err != nil {
			InternalServerError(w, err, ctx)
			return
		}

		if !exist {
			// Creare la conversazione senza messaggio
			conv, err = rt.CreateConversationDB(Conversation{UserId: user_id, GroupId: group.GroupId}, Message{})

			if err != nil {
				InternalServerError(w, err, ctx)
				return
			}
		}
	} else {
		BadRequest(w, err, ctx, "Bad Request")
		return
	}

	// Response
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(conv); err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	return
}
