package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// Port defines Telegram I/O used by the service layer.
type Port interface {
	ReceiveMessages() <-chan tgbotapi.Update
	SendMessage(chatId int64, text string) error
	SendFile(chatId int64, text string, data []byte) error
	Stop()
}
