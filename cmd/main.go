package main

import (
	"TelegramBot/internal/clients/rconclient"
	"TelegramBot/internal/clients/tgclient"
	"TelegramBot/internal/config"
	"TelegramBot/internal/consumer/eventconsumer"
	"TelegramBot/internal/events/telegram"
	"TelegramBot/internal/storage/sqlite"
	"context"
	"log"
	"os"
	"time"
)

const (
	tgBotHost = "api.telegram.org"
	bachSize  = 100
)

func main() {

	storage, err := sqlite.New(os.Getenv("DATABASE_PATH"))
	if err != nil {
		log.Fatal("can't conneсt storage: ", err)
	}

	if err := storage.Init(context.TODO()); err != nil {
		log.Fatal("can't init storage: ", err)
	}

	eventsProcessor := telegram.New(
		tgclient.New(tgBotHost, config.MustLoad().Token),
		rconclient.New(config.MustLoad().Address, config.MustLoad().Password, 5*time.Second),
		storage,
	)

	log.Print("service started")

	consumer := eventconsumer.New(eventsProcessor, eventsProcessor, bachSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}
