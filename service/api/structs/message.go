package structs

import (
	"time"
)

/*
* Struct used to rappsent a message
*
* MessageId is the id of the message
* Text is the body of the message
* SendTime is the time when the message was sent
* Status is the status of the message
* SenderUserId is the id of the user who sent the message
* ConversationId is the id of the conversation where the message was sent
* Photo is the photo of the user who sent the message
 */

type Message struct {
	MessageId      int       `json:"messageId"`
	Text           string    `json:"text"`
	SendTime       time.Time `json:"sendTime"`
	Status         string    `json:"status"`
	SenderUserId   int       `json:"senderUserId"`
	ConversationId int       `json:"conversationId"`
	Photo          string    `json:"photo"`
}
