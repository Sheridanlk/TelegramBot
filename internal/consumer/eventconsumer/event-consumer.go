package eventconsumer

import (
	"TelegramBot/internal/events"
	"log"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, bachSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: bachSize,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err)

			continue
		}
	}
}

/*
Исправить потрерю событий: ретраи, возврщение в хранилище, фоллбэк, подтверждение
Обработка всей пачки: останавливаться после первой ошибки, счётчик ошибок
паралельная обработка событий(прям надо) поможет WaitGroup
*/

func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("can't hundle event: %s", err.Error())

			continue
		}
	}
	return nil
}
