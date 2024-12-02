package database

// Query used to get all conversation of a user
var queryGetConversations = `SELECT ConversationId FROM conversation WHERE UserId = ?`

func (db *appdbimpl) GetConversations(userId int) ([]Conversation, error) {
	// Get all conversations of a user from the db
	rows, err := db.c.Query(queryGetConversations, userId)
	if err != nil {
		return nil, err
	}
	defer func() { rows.Close() }()

	// All conversation
	var convs []Conversation

	for rows.Next() {
		if rows.Err() != nil {
			return nil, err
		}

		var conv Conversation
		err = rows.Scan(&conv.ConversationId)
		if err != nil {
			return nil, err
		}

		convs = append(convs, conv)
	}
	return convs, nil
}
