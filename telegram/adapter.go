package telegram

import (
	"SteamInfoScrapper/helpers"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Adapter is a thin wrapper around the Telegram Bot API client.
type Adapter struct {
	API *tgbotapi.BotAPI
}

// NewTelegramAdapter creates a Telegram bot client.
func NewTelegramAdapter(token string) (*Adapter, error) {
	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		return nil, err
	}

	return &Adapter{API: bot}, nil
}

// ReceiveMessages returns a channel of updates from Telegram long polling.
func (t *Adapter) ReceiveMessages() <-chan tgbotapi.Update {
	u := tgbotapi.UpdateConfig{
		Offset:  0,
		Limit:   0,
		Timeout: 60,
	}

	updates := t.API.GetUpdatesChan(u)

	return updates
}

// SendMessage sends a text message to a chat.
func (t *Adapter) SendMessage(chatId int64, message string) error {
	msg := tgbotapi.NewMessage(chatId, message)

	_, err := t.API.Send(msg)

	return err
}

// SendFile sends a message and a CSV document to a chat.
func (t *Adapter) SendFile(chatId int64, text string, data []byte) error {
	msg := tgbotapi.NewMessage(chatId, text)
	file := tgbotapi.FileBytes{
		Name:  "report.csv",
		Bytes: data,
	}

	doc := tgbotapi.NewDocument(chatId, file)

	if _, err := t.API.Send(msg); err != nil {
		helpers.ProcessError(err)
	}
	_, err := t.API.Send(doc)

	return err
}

// Stop stops receiving updates. Use it during graceful shutdown.
func (t *Adapter) Stop() {
	t.API.StopReceivingUpdates()
}
