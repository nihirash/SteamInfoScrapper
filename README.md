# SteamInfoScrapper

[![Go](https://github.com/nihirash/SteamInfoScrapper/actions/workflows/go.yml/badge.svg)](https://github.com/nihirash/SteamInfoScrapper/actions/workflows/go.yml)

A Telegram bot that collects Steam game information and generates CSV reports. The bot receives Steam game URLs via Telegram messages and returns detailed game information in CSV format.

## Features

- Fetches game information from Steam API
- Web scraping for tags and Steam Deck compatibility status
- Generates CSV reports with game details
- Telegram bot interface for easy interaction
- Graceful shutdown on SIGINT/SIGTERM

## Requirements

- Go 1.25 or later
- Steam API key (optional, currently not required)
- Telegram bot token

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd SteamInfoScrapper
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o steam-info-scrapper
```

## Configuration

1. Copy the example configuration file:
```bash
cp config.yml.example config.yml
```

2. Edit `config.yml` with your settings:
```yaml
steam:
  api_key: "YOUR_STEAM_API_KEY"
  one_task_timeout_sec: 30

telegram:
  bot_token: "YOUR_TELEGRAM_BOT_TOKEN"
```

Alternatively, you can use environment variables:
- `STEAM_API_KEY` - Steam API key
- `TELEGRAM_BOT_TOKEN` - Telegram bot token

The application will first try to load `config.yml` from the current directory, then fall back to `/etc/app_config.yml`.

## Usage

1. Start the bot:
```bash
./steam-info-scrapper
```

2. Send Steam game URLs to your Telegram bot (one per line):
```
https://store.steampowered.com/app/730/CounterStrike_Global_Offensive/
https://store.steampowered.com/app/271590/Grand_Theft_Auto_V/
```

3. The bot will process the URLs and send back a CSV file with game information.

## Project Structure

```
SteamInfoScrapper/
├── main.go              # Application entry point
├── config.go            # Configuration loading
├── config.yml.example   # Example configuration
├── helpers/
│   └── utils.go         # Helper functions
├── steam/
│   ├── adapter.go       # Steam API and web scraping adapter
│   ├── domain.go        # Domain models and utilities
│   └── port.go          # Steam port interface
├── telegram/
│   ├── adapter.go       # Telegram bot adapter
│   ├── domain.go        # Telegram domain models
│   ├── port.go          # Telegram port interface
│   └── service.go       # Telegram service logic
└── report/
    ├── csvadapter.go    # CSV report generation
    └── port.go          # Report port interface
```

## How It Works

1. The bot listens for messages on Telegram
2. When it receives Steam game URLs, it extracts the game IDs
3. For each game, it fetches information from:
   - Steam API (game details, reviews, platforms)
   - Steam store pages (tags, Steam Deck compatibility)
4. All information is compiled into a CSV report
5. The report is sent back to the user as a file

## Shutdown

The application supports graceful shutdown. Send SIGINT (Ctrl+C) or SIGTERM to stop the bot cleanly. All resources will be properly closed.

## License

MIT. See `LICENSE`.
