package database

var queryFindConversationById = `SELECT ConversationId, GroupId FROM conversation WHERE ConversationId = ? AND UserId = ?`
var queryFindConversationUserById = `SELECT ConversationId, SenderUserId FROM conversation WHERE ConversationId = ? AND UserId = ?`

func (db *appdbimpl) GetConversationById(convId int, userId int) (Conversation, error) {
	var conv Conversation
	err := db.c.QueryRow(queryFindConversationById, convId, userId).Scan(&conv.ConversationId, &conv.GroupId)
	if err != nil {
		err = db.c.QueryRow(queryFindConversationUserById, convId, userId).Scan(&conv.ConversationId, &conv.SenderUserId)
		return conv, err
	}
	return conv, nil
}
