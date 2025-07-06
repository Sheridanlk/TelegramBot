package main

import (
	"TelegramBot/internal/clients/telegram"
	"TelegramBot/internal/config"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	//load config

	tgClietn := telegram.New(tgBotHost, config.MustLoad().Token)

}
