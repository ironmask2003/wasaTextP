package database

// In this file we define the SQL for the tables we want to create in the database.

/*
* USER TABEL SQL
*	- UserId: int (PK) Unique ID for each user
*	- Username: string (Unique) Username of the user (used for login)
 */
var userTableSQL = `CREATE TABLE IF NOT EXISTS user (
	UserId INTEGER NOT NULL UNIQUE,
	Username STRING NOT NULL UNIQUE,
	PRIMARY KEY(UserId)
);`

/*
* MESSAGE TABLE SQL
*	- MessageId: int not null, (PK) ID for each message
*	- Message: TEXT, text of the message
*	- SendTime: DateTime current time
*	- Status: string
*	- SenderUserID: int (PK) (FK) User that sent the message
 */
var messageTableSQL = `CREATE TABLE IF NOT EXISTS message (
	MessageId INTEGER NOT NULL,
	Message TEXT,
	SendTime DATETIME DEFAULT CURRENT_TIMESTAMP,
	Status STRING,
	SendUserID INTEGER NOT NULL,
  ConversationId INTEGER NOT NULL
	PRIMARY KEY(MessageId, SendUserID, ConversationId),
	CONSTAINT fk_message
		FOREIGN KEY (SendUserID) REFERENCES user(UserId)
		ON DELTE CASCADE
);`

/*
* GROUP TABLE SQL
*	- GroupId: int (PK) Unique ID for each group
*	- GroupName: string not null, Name of the group
 */
var groupTableSQL = `CREATE TABLE IF NOT EXISTS group (
	GroupId INTEGER NOT NULL UNIQUE,
	GroupName STRING NOT NULL,
	PRIMARY KEY(GroupId)
);`

var userGroupTableSQL = `CREATE TABLE IF NOT EXISTS user_group (
	GroupId INTEGER NOT NULL,
	UserID INTEGER NOT NULL,
	PRIMARY KEY(GroupId, UserID),
	CONSTRAINT fk_user_group
		FORIEGN KEY (GroupId) REFERENCES group(GroupId)
		ON DELETE CASCADE
		FOREIGN KEY (UserID) REFERENCES user(UserId)
		ON DELETE CASCADE
)`

/*
* CONVERSATION TABLE SQL
*	- ConversationId: int (PK) Unique ID for each conversation (chat)
*	- GroupId: int (FK)
*	- UserId: int (FK)
*	- LastMessageId: int (FK)
 */
var conversationTableSQL = `CREATE TABLE IF NOT EXISTS conversation (
	ConversationId INTEGER NOT NULL UNIQUE,
	UserId INTEGER UNIQUE NOT NULL,
  GroupId INTEGER UNIQUE,
	SenderUserId INTEGER UNIQUE,
	LastMessageId INTEGER,
	PRIMARY KEY(ConversationId, UserId),
	CONSTRAINT fk_conversation
		FOREIGN KEY (GroupId) REFERENCES group(GroupId)
		ON DELETE CASCADE
		FOREIGN KEY (UserID) REFERENCES user(UserId)
		ON DELETE CASCADE
		FOREIGN KEY (LastMessageId) REFERENCES message(MessageId)
		ON DELETE CASCADE
);`
