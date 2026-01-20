package main

import (
	"SteamInfoScrapper/helpers"
	"SteamInfoScrapper/report"
	"SteamInfoScrapper/steam"
	"SteamInfoScrapper/telegram"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Print("Starting SteamInfoScrapper")

	appConfig, err := LoadConfig()

	if err != nil {
		helpers.ProcessError(err)

		os.Exit(1)
	}

	appConfig.Print()

	steamAdapter := steam.NewSteamApiAdapter(appConfig.Steam.ApiKey, appConfig.Steam.OneTaskTimeoutSec)
	csvReportAdapter := report.NewCSVAdapter()
	telegramAdapter, err := telegram.NewTelegramAdapter(appConfig.Telegram.BotToken)

	if err != nil {
		helpers.ProcessError(err)

		os.Exit(1)
	}

	telegramService := telegram.NewService(telegramAdapter, steamAdapter, csvReportAdapter)

	// Graceful shutdown on SIGINT/SIGTERM.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	go func() {
		<-ctx.Done()
		log.Print("Shutting down...")
		telegramAdapter.Stop()
		steamAdapter.Close()
		os.Exit(0)
	}()

	telegramService.ProcessLoop()

}
