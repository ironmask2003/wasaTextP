package database

var queryFindConversationById = `SELECT ConversationId, GroupId FROM conversation WHERE ConversationId = ?`
var queryFindConversationUserById = `SELECT ConversationId, SenderUserId FROM conversation WHERE ConversationId = ?`

func (db *appdbimpl) GetConversationById(convId int) (Conversation, error) {
	var conv Conversation
	err := db.c.QueryRow(queryFindConversationById, convId).Scan(&conv.ConversationId, &conv.GroupId)
	if err != nil {
		err = db.c.QueryRow(queryFindConversationUserById, convId).Scan(&conv.ConversationId, &conv.SenderUserId)
		return conv, err
	}
	return conv, nil
}
