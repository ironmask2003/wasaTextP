package api

import (
	"unicode/utf8"
	"wasa.project/service/database"
)

// -- Structs for the Comment -- //
type Comment struct {
	CommentId      int    `json:"commentId"`
	Comment        string `json:"comment"`
	CommentUserId  int    `json:"commentUserId"`
	MessageId      int    `json:"messageId"`
	UserId         int    `json:"userId"`
	ConversationId int    `json:"conversationId"`
}

func (c *Comment) ConvertCommentForDB() database.Comment {
	return database.Comment{
		CommentId:      c.CommentId,
		Comment:        c.Comment,
		CommentUserId:  c.CommentUserId,
		MessageId:      c.MessageId,
		UserId:         c.UserId,
		ConversationId: c.ConversationId,
	}
}

func (c *Comment) ConvertCommentFromDB(dbComment database.Comment) error {
	c.CommentId = dbComment.CommentId
	c.Comment = dbComment.Comment
	c.CommentUserId = dbComment.CommentUserId
	c.MessageId = dbComment.MessageId
	c.UserId = dbComment.UserId
	c.ConversationId = dbComment.ConversationId
	return nil
}

// Control with regexp if the comment is a emoji
func (c *Comment) IsEmoji() bool {
	// Decodifica il primo carattere Unicode nella stringa
	r, _ := utf8.DecodeRuneInString(c.Comment)
	if r == utf8.RuneError {
		return false
	}

	// Verifica se il carattere appartiene a uno dei range degli emoji
	return (r >= 0x1F600 && r <= 0x1F64F) || // Emoticon
		(r >= 0x1F300 && r <= 0x1F5FF) || // Simboli e pittogrammi vari
		(r >= 0x1F680 && r <= 0x1F6FF) || // Trasporti e simboli mappa
		(r >= 0x1F700 && r <= 0x1F77F) || // Simboli alchemici
		(r >= 0x2600 && r <= 0x26FF) || // Simboli vari
		(r >= 0x2700 && r <= 0x27BF) || // Dingbats
		(r >= 0xFE00 && r <= 0xFE0F) || // Variazioni selettori
		(r >= 0x1F900 && r <= 0x1F9FF) || // Simboli supplementari
		(r >= 0x1FA70 && r <= 0x1FAFF) || // Emoji aggiuntivi
		(r >= 0x1F1E6 && r <= 0x1F1FF) // Bandiere (regioni)
}
