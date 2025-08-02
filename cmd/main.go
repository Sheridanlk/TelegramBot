package main

import (
	"TelegramBot/internal/clients/rconclient"
	"TelegramBot/internal/clients/tgclient"
	"TelegramBot/internal/config"
	"TelegramBot/internal/consumer/eventconsumer"
	"TelegramBot/internal/events/telegram"
	rconpoller "TelegramBot/internal/poller/rconPoller"
	"TelegramBot/internal/storage/sqlite"
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

const (
	tgBotHost = "api.telegram.org"
	bachSize  = 100
)

func main() {
	_ = godotenv.Load(filepath.Join("..", ".env"))

	storage, err := sqlite.New(os.Getenv("DATABASE_PATH"))
	if err != nil {
		log.Fatal("can't conneсt storage: ", err)
	}

	if err := storage.Init(context.TODO()); err != nil {
		log.Fatal("can't init storage: ", err)
	}

	telegramClient := tgclient.New(tgBotHost, config.MustLoad().Token)
	rconClietn := rconclient.New(config.MustLoad().Address, config.MustLoad().Password, 5*time.Second)
	rconPoller := rconpoller.New(rconClietn, storage)

	eventsProcessor := telegram.New(
		telegramClient,
		rconClietn,
		rconPoller,
		storage,
	)

	log.Print("service started")

	consumer := eventconsumer.New(eventsProcessor, eventsProcessor, bachSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}
