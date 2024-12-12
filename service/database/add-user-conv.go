package database

// Query used to add in the conversation_user table all the user in the conversation table
var queryAddUserConv = `INSERT INTO conversation_user (ConversationId, UserId) VALUES (?, ?)`

func (db *appdbimpl) AddUserConv(conversationId int, userId int) error {
	_, err := db.c.Exec(queryAddUserConv, conversationId, userId)
	if err != nil {
		return err
	}
	return nil
}
