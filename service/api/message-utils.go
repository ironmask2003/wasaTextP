package api

import (
	"time"
	"wasa.project/service/database"
)

type Message struct {
	MessageId      int       `json:"messageId"`
	Text           string    `json:"text"`
	SendTime       time.Time `json:"sendTime"`
	Status         string    `json:"status"`
	SenderUserId   int       `json:"senderUserId"`
	ConversationId int       `json:"conversationId"`
}

func (m *Message) ConvertMessageForDB() database.Message {
	return database.Message{
		MessageId:      m.MessageId,
		Text:           m.Text,
		SendTime:       m.SendTime,
		Status:         m.Status,
		SenderUserId:   m.SenderUserId,
		ConversationId: m.ConversationId,
	}
}

func (m *Message) ConvertMessageFromDB(msgDB database.Message) error {
	m.MessageId = msgDB.MessageId
	m.Text = msgDB.Text
	m.SendTime = msgDB.SendTime
	m.Status = msgDB.Status
	m.SenderUserId = msgDB.SenderUserId
	m.ConversationId = msgDB.ConversationId
	return nil
}
