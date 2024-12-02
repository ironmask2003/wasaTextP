package database

// Query used to get all messages from a conversation
var queryGetConversation = `SELECT MessageId, Message, SendTime, SendUserId, ConversationId FROM message WHERE ConversationId = ? AND SendUserId = ?`

func (db *appdbimpl) GetConversationGroup(conv Conversation) ([]Message, error) {
	// Get all members of the group
	members, err := db.GetMembers(conv.GroupId)
	if err != nil {
		return nil, err
	}

	// Get all messages from the conversation with the group
	var messages []Message
	for _, member := range members {
		// Get the conversation from the group and the user id
		convId, err := db.GetConversationsByUserGroup(conv.GroupId, member.UserId)
		if err != nil {
			return nil, err
		}
		rows, err := db.c.Query(queryGetConversation, convId, member.UserId)
		if err != nil {
			return nil, err
		}
		defer func() { err = rows.Close() }()
		for rows.Next() {
			if rows.Err() != nil {
				return nil, err
			}
			var msg Message
			err = rows.Scan(&msg.MessageId, &msg.Text, &msg.SendTime, &msg.SenderUserId)
			if err != nil {
				return nil, err
			}
			messages = append(messages, msg)
		}
	}
	return messages, nil
}

func (db *appdbimpl) GetConversation(userId int, convId int) ([]Message, error) {
	var msg []Message
	// Get conversation by id
	conv, err := db.GetConversationById(convId, userId)
	if err != nil {
		return nil, err
	}

	// Check if the conversation is with a group
	if conv.GroupId != 0 {
		msg, err = db.GetConversationGroup(conv)
		if err != nil {
			return nil, err
		}
	} else {
		// Get all messages from the conversation
		rcvId := conv.SenderUserId
		convId, err := db.GetConversationsBySender(userId, rcvId)
		if err != nil {
			return nil, err
		}
		usr := userId
		cv := conv.ConversationId
		for i := 0; i < 2; i++ {
			rows, err := db.c.Query(queryGetConversation, cv, usr)
			if err != nil {
				return nil, err
			}
			defer func() { err = rows.Close() }()
			for rows.Next() {
				if rows.Err() != nil {
					return nil, err
				}
				var m Message
				err = rows.Scan(&m.MessageId, &m.Text, &m.SendTime, &m.SenderUserId, &m.ConversationId)
				if err != nil {
					return nil, err
				}
				msg = append(msg, m)
			}
			usr = rcvId
			cv = convId
		}
	}

	return msg, nil
}
