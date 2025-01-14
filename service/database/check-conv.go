package database

// Query
var queryCheckConv = `
		SELECT COUNT(convp.ConversationId)
		FROM conversation_user convp, conversation_user conv, conversation c
		WHERE convp.ConversationId = conv.ConversationId AND convp.UserId = ? AND conv.UserId = ? 
      AND c.ConversationId = convp.ConversationId AND c.GroupId IS NULL
	`
var queryCheckUserConv = `
  SELECT ConversationId
  FROM conversation_user
  WHERE ConversationId = ? AND UserId = ?
`

func (db *appdbimpl) CheckIfExistConv(sender int, receiver int) (bool, error) {
	var count int
	err := db.c.QueryRow(queryCheckConv, sender, receiver).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

func (db *appdbimpl) CheckUserConv(userId int, convId int) (bool, error) {
	var id int
	err := db.c.QueryRow(queryCheckUserConv, convId, userId).Scan(&id)
	if err != nil {
		return false, err
	}
	return id != 0, nil
}
