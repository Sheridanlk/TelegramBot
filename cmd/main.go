package main

import (
	"TelegramBot/internal/clients/rconclient"
	"TelegramBot/internal/clients/tgclient"
	"TelegramBot/internal/config"
	"TelegramBot/internal/consumer/eventconsumer"
	"TelegramBot/internal/events/telegram"
	"log"
	"time"
)

const (
	tgBotHost = "api.telegram.org"
	bachSize  = 100
)

func main() {

	eventsProcessor := telegram.New(
		tgclient.New(tgBotHost, config.MustLoad().Token),
		rconclient.New(config.MustLoad().Address, config.MustLoad().Password, 5*time.Second),
	)

	log.Print("service started")

	consumer := eventconsumer.New(eventsProcessor, eventsProcessor, bachSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}
