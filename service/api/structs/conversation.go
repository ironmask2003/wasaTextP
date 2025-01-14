package structs

type Conversation struct {
	ConversationId int `json:"conversationId"`
	GroupId        int `json:"GroupId"`
	LastMessageId  int `json:"lastMessageId"`
}
