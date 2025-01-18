package structs

import (
	"unicode/utf8"
)

type Comment struct {
	CommentId      int    `json:"commentId"`
	Comment        string `json:"comment"`
	MessageId      int    `json:"messageId"`
	ConversationId int    `json:"conversationId"`
	CommentUserId  int    `json:"commnetUserId"`
}

type RspComment struct {
	CommentId       int    `json:"commentId"`
	Comment         string `json:"comment"`
	MessageId       int    `json:"messageId"`
	ConversationId  int    `json:"conversationId"`
	CommentUserId   int    `json:"commentUserId"`
	CommentUsername string `json:"commentUsername"`
}

// Function used to check if the comment is valid
func (c Comment) IsValid() bool {
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
