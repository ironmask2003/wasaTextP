package database

var queryGetConversations = `SELECT ConversationId, UserId, GroupId, SendUserId, LastMessageId FROM Conversation WHERE UserId = ?`

func (db *appdbimpl) GetConversations(userId int) ([]Conversation, error) {
	// Get the posts from the database
	rows, err := db.c.Query(queryGetConversations, userId)
	if err != nil {
		return nil, err
	}
	defer func() { err = rows.Close() }()

	// List of conversation
	var convs []Conversation

	for rows.Next() {
		if rows.Err() != nil {
			return nil, err
		}

		var conv Conversation

		// Get conversation data
		err = rows.Scan(&conv.ConversationId, &conv.UserId, &conv.GroupId, &conv.SenderUserId, &conv.LastMessageId)
		if err != nil {
			return nil, err
		}

		// Append sigle conversation in the conversation array
		convs = append(convs, conv)
	}
	return convs, nil
}
