package api

import (
	"wasa.project/service/database"
)

type Conversation struct {
	ConversationId int `json:"conversationId"`
	SenderUserId   int `json:"senderUserId"`
	GroupId        int `json:"groupId"`
	LastMessageId  int `json:"lasrMessageId"`
}

// Function used to convert the Conversation struct used in the api package in the Conversation struct used in the database package
func (c *Conversation) ConvertConversationForDB() database.Conversation {
	return database.Conversation{
		ConversationId: c.ConversationId,
		SenderUserId:   c.SenderUserId,
		GroupId:        c.GroupId,
		LastMessageId:  c.LastMessageId,
	}
}

// Function used to convert the Conversation struct used in the database package in the Conversation struct used in the api package
func (c *Conversation) ConvertConversationFromDB(conversation database.Conversation) {
	c.ConversationId = conversation.ConversationId
	c.SenderUserId = conversation.SenderUserId
	c.GroupId = conversation.GroupId
	c.LastMessageId = conversation.LastMessageId
}
