package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// Message is an internal representation of a Telegram message used by the service.
type Message struct {
	ChatId     int64
	SenderName string
	Text       string
}

// BuildFromTgMessage converts a Telegram message into internal Message.
func BuildFromTgMessage(message *tgbotapi.Message) Message {
	return Message{
		ChatId:     message.Chat.ID,
		SenderName: message.Chat.UserName,
		Text:       message.Text,
	}
}
