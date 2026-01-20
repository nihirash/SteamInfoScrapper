package telegram

import (
	"SteamInfoScrapper/helpers"
	"SteamInfoScrapper/report"
	"SteamInfoScrapper/steam"
	"fmt"
	"log"
	"strings"
)

// Service implements the Telegram bot workflow.
type Service struct {
	TgPort     Port
	SteamPort  steam.Port
	ReportPort report.Port
}

// NewService wires Telegram, Steam and Report ports into a bot service.
func NewService(telegramPort Port, steamPort steam.Port, reportPort report.Port) *Service {
	return &Service{telegramPort, steamPort, reportPort}
}

// ProcessMessage parses Steam URLs from a message, builds a CSV report, and replies back.
func (svc *Service) ProcessMessage(message Message) error {
	log.Printf("Building report for %s", message.SenderName)

	_ = svc.TgPort.SendMessage(message.ChatId, fmt.Sprintf("%s, взяли работу!", message.SenderName))

	var ids []uint

	for line := range strings.Lines(message.Text) {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		id, err := steam.GameUrl(line).ExtractID()
		if err != nil {
			_ = svc.TgPort.SendMessage(message.ChatId, fmt.Sprintf("Битый адрес: %s ", line))
			helpers.ProcessError(err)

			continue
		}

		ids = append(ids, id)
	}

	pages, err := svc.SteamPort.GetGamePages(ids)
	if err != nil {
		_ = svc.TgPort.SendMessage(message.ChatId, fmt.Sprintf("Случилась ошибка при сборе отчета: %s ", err))

		return err
	}

	data, err := svc.ReportPort.Store(pages)

	if err != nil {
		_ = svc.TgPort.SendMessage(message.ChatId, fmt.Sprintf("Ошибка при создании файла отчета: %s", err))

		return err
	}

	err = svc.TgPort.SendFile(message.ChatId, "Отчет готов!", data.Bytes())

	if err != nil {
		return err
	}

	log.Printf("Report for %s built!", message.SenderName)

	return nil
}

// StartCommand sends a short usage hint.
func (svc *Service) StartCommand(chatId int64) {
	_ = svc.TgPort.SendMessage(chatId, "Привет!\n"+
		"Я ОЧЕНЬ ПРОСТОЙ телеграм бот для сбора информации по стим играм.\n"+
		"Присылай мне только ссылки на игры одним сообщением и получишь отчет в виде csv файла!")
}

// ProcessLoop consumes Telegram updates and dispatches work.
func (svc *Service) ProcessLoop() {
	channel := svc.TgPort.ReceiveMessages()

	log.Print("Processing telegram messages!")

	for update := range channel {
		if update.Message == nil {
			continue
		}

		if update.Message.Text == "/start" {
			svc.StartCommand(update.Message.Chat.ID)

			continue
		}

		msg := BuildFromTgMessage(update.Message)

		go func(msg Message) {
			err := svc.ProcessMessage(msg)
			if err != nil {
				helpers.ProcessError(err)
			}
		}(msg)

	}
}
