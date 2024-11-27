package api

import (
	"time"
	"wasa.project/service/database"
)

type Message struct {
	MessageId    int       `json:"messageId"`
	Text         string    `json:"text"`
	SendTime     time.Time `json:"sendTime"`
	Photo        string    `json:"photo"`
	Status       string    `json:"status"`
	SenderUserId int       `json:"senderUserId"`
}

func (m *Message) ConvertMessageForDB() database.Message {
	return database.Message{
		MessageId:    m.MessageId,
		Text:         m.Text,
		SendTime:     m.SendTime,
		Status:       m.Status,
		SenderUserId: m.SenderUserId,
	}
}

func (m *Message) ConvertMessageFromDB(message database.Message) {
	m.MessageId = message.MessageId
	m.Text = message.Text
	m.SendTime = message.SendTime
	m.Status = message.Status
	m.SenderUserId = message.SenderUserId
}
