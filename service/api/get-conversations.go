package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasa.project/service/api/reqcontext"
	"wasa.project/service/api/structs"
)

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the userId from the endpoint
	userId, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		BadRequest(w, err, ctx, "Can't take the user id from the endpoint")
		return
	}

	// Check if the user is authorized
	if checkAuth(w, userId, ctx) != nil {
		Forbidden(w, nil, ctx, "The user is not authorized")
		return
	}

	// Get all conversations of the user from the db
	var convs []structs.Conversation
	convs, err = rt.db.GetUserConversations(userId)
	if err != nil {
		BadRequest(w, err, ctx, "Error taking the conversations from the db")
		return
	}

	// Struct used for the response
	type Response struct {
		Conversation structs.Conversation `json:"conversation"`
		User         User                 `json:"user"`
		Group        Group                `json:"group"`
		GroupUsers   []User               `json:"groupUsers"`
		Message      structs.Message      `json:"message"`
		SenderUser   User                 `json:"senderUser"`
	}

	// Response
	response := make([]Response, len(convs))

	// Fornire anche informazioni riguardo gruppi o utenti
	for idx, conv := range convs {
		if conv.GroupId == 0 {
			// Get the user from the conversation
			userID_DB, err := rt.db.GetUsersConv(conv.ConversationId, userId)
			if err != nil {
				BadRequest(w, err, ctx, "Error taking user of the conversation")
				return
			}
			userDB, err := rt.db.GetUserById(userID_DB.UserId)
			if err != nil {
				BadRequest(w, err, ctx, "Error taking user of the conversation")
				return
			}
			var user User
			err = user.ConvertUserFromDB(userDB)
			if err != nil {
				BadRequest(w, err, ctx, "Error converting user")
				return
			}

			// Get last message
			message, err := rt.db.GetMessageById(conv.LastMessageId, conv.ConversationId)
			if err != nil {
				BadRequest(w, err, ctx, "Error taking the message")
				return
			}

			SenderUserDB, err := rt.db.GetUserById(message.SenderUserId)
			if err != nil {
				BadRequest(w, err, ctx, "Error taking the sender user from the message")
				return
			}

			var senderUser User
			err = senderUser.ConvertUserFromDB(SenderUserDB)
			if err != nil {
				BadRequest(w, err, ctx, "Error converting senderUser")
				return
			}

			response[idx] = Response{
				Conversation: conv,
				User:         user,
				Message:      message,
				SenderUser:   senderUser,
			}
		} else {
			// Get the group from the conversation
			groupDB, err := rt.db.GetGroupById(conv.GroupId)
			if err != nil {
				BadRequest(w, err, ctx, "Error taking the group from the conversation")
				return
			}
			var group Group
			err = group.ConvertGroupFromDB(groupDB)
			if err != nil {
				BadRequest(w, err, ctx, "Error converting the group")
				return
			}

			// Get last message
			message, err := rt.db.GetMessageById(conv.LastMessageId, conv.ConversationId)
			if err != nil {
				BadRequest(w, err, ctx, "Error taking the message dio")
				return
			}

			SenderUserDB, err := rt.db.GetUserById(message.SenderUserId)
			if err != nil {
				BadRequest(w, err, ctx, "Error taking the sender user from the message")
				return
			}

			var senderUser User
			err = senderUser.ConvertUserFromDB(SenderUserDB)
			if err != nil {
				BadRequest(w, err, ctx, "Error converting senderUser")
				return
			}

			users, err := rt.db.GetMembers(conv.GroupId)
			if err != nil {
				BadRequest(w, err, ctx, "Error taking the members of the group")
				return
			}

			var groupUsers []User
			for i := 0; i < len(users); i++ {
				userDB, err := rt.db.GetUserById(users[i].UserId)
				if err != nil {
					BadRequest(w, err, ctx, "Error taking the user from the user table")
					return
				}
				var user User
				err = user.ConvertUserFromDB(userDB)
				if err != nil {
					BadRequest(w, err, ctx, "Error converting the user from the database struct")
					return
				}
				groupUsers = append(groupUsers, user)
			}

			response[idx] = Response{
				Conversation: conv,
				Group:        group,
				GroupUsers:   groupUsers,
				Message:      message,
				SenderUser:   senderUser,
			}
		}
	}

	// Write the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
