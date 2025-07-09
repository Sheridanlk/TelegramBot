package main

import (
	"TelegramBot/internal/clients/tgclient"
	"TelegramBot/internal/config"
	"TelegramBot/internal/consumer/eventconsumer"
	"TelegramBot/internal/events/telegram"
	"log"
)

const (
	tgBotHost = "api.telegram.org"
	bachSize  = 100
)

func main() {

	eventsProcessor := telegram.New(
		tgclient.New(tgBotHost, config.MustLoad().Token),
	)

	log.Print("service started")

	consumer := eventconsumer.New(eventsProcessor, eventsProcessor, bachSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}
