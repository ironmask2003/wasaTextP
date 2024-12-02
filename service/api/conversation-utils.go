package api

import "wasa.project/service/database"

type Conversation struct {
	ConversationId int `json:"conversationId"`
	UserId         int `json:"userId"`
	SenderUserId   int `json:"senderUserId"`
	GroupId        int `json:"groupId"`
	LastMessageId  int `json:"lastMessageId"`
}

func (c *Conversation) ConvertConversationForDB() database.Conversation {
	return database.Conversation{
		ConversationId: c.ConversationId,
		UserId:         c.UserId,
		SenderUserId:   c.SenderUserId,
		GroupId:        c.GroupId,
		LastMessageId:  c.LastMessageId,
	}
}

func (c *Conversation) ConvertConversationFromDB(convDB database.Conversation) error {
	c.ConversationId = convDB.ConversationId
	c.UserId = convDB.UserId
	c.SenderUserId = convDB.SenderUserId
	c.GroupId = convDB.GroupId
	c.LastMessageId = convDB.LastMessageId
	return nil
}
