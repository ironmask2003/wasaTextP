package database

// Query used ty insert a new conversation in the conversation table
var queryAddConversation = `INSERT INTO conversation (ConversationId, GroupId, UserId, LastMessageId) VALUES (?, ?, ?, ?);`

// Query for take the max id in the conversation
var queryMaxConversationID = `SELECT MAX(ConversationId) FROM conversation`

// Funzione nel caso in cui si voglia creare una conversazione tra utenti
func CreateForUsers() {

}

// Funzione nel caso in cui si voglia creare una conversazione tra gruppi

func (db *appdbimpl) CreateConversation(c Conversation, m Message) (Conversation, error) {
	var conv Conversation
	if c.GroupId != 0 {
		conv.GroupId = c.GroupId
	} else {
		conv.UserId = c.UserId
	}

	// Getting the max id in the conversation table
	maxID, err := db.GetLastElem(queryMaxConversationID)
	if err != nil {
		return conv, err
	}
	// Setting the id of the new conversation
	conv.ConversationId += maxID + 1

	return conv, nil
}
