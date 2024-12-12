package structs

type Conversation struct {
	ConversationId int  `json:"conversationId"`
	IsGroup        bool `json:"isGroup"`
	LastMessageId  int  `json:"lastMessageId"`
}
