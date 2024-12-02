package database

// Query used to get all conversation with a group
var queryGetConversationsByGroup = `SELECT ConversationId FROM conversation WHERE GroupId = ?`

func (db *appdbimpl) GetConversationsByGroup(groupId int) ([]Conversation, error) {
	// Get all conversations of a group from the db
	rows, err := db.c.Query(queryGetConversationsByGroup, groupId)
	if err != nil {
		return nil, err
	}
	defer func() { rows.Close() }()

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
