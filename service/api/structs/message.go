package structs

import (
	"time"
)

type Message struct {
	MessageId      int       `json:"messageId"`
	Text           string    `json:"text"`
	SendTime       time.Time `json:"sendTime"`
	Status         string    `json:"status"`
	SenderUserId   int       `json:"senderUserId"`
	ConversationId int       `json:"conversationId"`
	Photo          string    `json:"photo"`
}
