package database

import "time"

// -- Structs used in the database for the user containning the userId and the username -- //
type User struct {
	UserId   int    `json:"userId"`
	Username string `json:"username"`
}

// -- Structs for the message -- //
type Message struct {
	MessageId      int       `json:"messageId"`
	Text           string    `json:"text"`
	SendTime       time.Time `json:"sendTime"`
	Status         string    `json:"status"`
	SenderUserId   int       `json:"senderUserId"`
	ConversationId int       `json:"conversationId"`
}

// -- Structs for the group -- //
type Group struct {
	GroupId   int    `json:"groupId"`
	GroupName string `json:"groupName"`
}

// -- Structs for the Conversation -- //
type Conversation struct {
	ConversationId int `json:"conversationId"`
	SenderUserId   int `json:"senderUserId"`
	GroupId        int `json:"groupId"`
	UserId         int `json:"userId"`
	LastMessageId  int `json:"lastMessageId"`
}

// -- Structs for the Comment -- //
type Comment struct {
	CommentId      int    `json:"commentId"`
	Comment        string `json:"comment"`
	CommentUserId  int    `json:"commentUserId"`
	MessageId      int    `json:"messageId"`
	UserId         int    `json:"userId"`
	ConversationId int    `json:"conversationId"`
}
