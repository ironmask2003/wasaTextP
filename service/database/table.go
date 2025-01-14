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
  Status TEXT,
  ConversationId INTEGER NOT NULL,
  SenderUserId INTEGER NOT NULL,
  PRIMARY KEY(MessageId, ConversationId),
  CONSTRAINT fk_message
    FOREIGN KEY (SenderUserId) REFERENCES user(UserId)
      ON DELETE CASCADE
    FOREIGN KEY (ConversationId) REFERENCES conversation(ConversationId)
      ON DELETE CASCADE
);`

/*
* GROUP TABLE SQL
*	- GroupId: int (PK) Unique ID for each group
*	- GroupName: string not null, Name of the group
 */
var groupTableSQL = `CREATE TABLE IF NOT EXISTS group_t (
	GroupId INTEGER NOT NULL,
	GroupName STRING NOT NULL,
	PRIMARY KEY(GroupId)
);`

var userGroupTableSQL = `CREATE TABLE IF NOT EXISTS user_group (
	GroupId INTEGER NOT NULL,
	UserId INTEGER NOT NULL,
  PRIMARY KEY(GroupId, UserId),
	CONSTRAINT fk_user_group
		FOREIGN KEY (GroupId) REFERENCES group_t(GroupId)
		  ON DELETE CASCADE
		FOREIGN KEY (UserId) REFERENCES user(UserId)
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
	ConversationId INTEGER NOT NULL,
  GroupId INTEGER NOT NULL,
  LastMessageId INTEGER,
  PRIMARY KEY(ConversationId)
  CONSTRAINT fk_conversation
    FOREIGN KEY (LastMessageId, ConversationId) REFERENCES message(MessageId, ConversationId)
      ON DELETE CASCADE
    FOREIGN KEY (GroupId) REFERENCES group_t(GroupId)
      ON DELETE CASCADE
);`

var conversationUsersSQL = `CREATE TABLE IF NOT EXISTS conversation_user (
  ConversationId INTEGER NOT NULL,
  UserId INTEGER NOT NULL,
  PRIMARY KEY(ConversationId, UserId),
  CONSTRAINT fk_conversation_user
    FOREIGN KEY (ConversationId) REFERENCES conversation(ConversationId)
      ON DELETE CASCADE
    FOREIGN KEY (UserId) REFERENCES user(UserId)
      ON DELETE CASCADE
);`

/*
* COMMENT TABLE
* - CommentId: int (PK) Unique ID for each comment
* - Comment: STRING
* - MessageId: int (FK) ID for the message
* - UserId: int (FK) ID for the user
* - ConversationId: int (FK)
* */
var commentTableSQL = `CREATE TABLE IF NOT EXISTS comment (
  CommentId INTEGER NOT NULL,
  CommentUserId INTEGER NOT NULL,
  Comment STRING,
  MessageId INTEGER NOT NULL,
  ConversationId INTEGER NOT NULL,
  PRIMARY KEY(CommentId, MessageId, ConversationId),
  CONSTRAINT fk_comment
    FOREIGN KEY (MessageId, ConversationId) REFERENCES message(MessageId, ConversationId)
      ON DELETE CASCADE
    FOREIGN KEY (CommentUserId) REFERENCES user(UserId)
      ON DELETE CASCADE
);`
