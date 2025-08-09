package eventconsumer

import (
	"TelegramBot/internal/events"
	"log"
	"sync"
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

func (c *Consumer) handleEvents(evenList []events.Event) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(evenList))

	for _, event := range evenList {
		log.Printf("got new event: %s", event.Text)
		wg.Add(1)

		go func(e events.Event) {
			defer wg.Done()
			if err := c.processor.Process(e); err != nil {
				errCh <- err
			}
		}(event)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			log.Printf("can't hundle event: %s", err.Error())
			continue
		}
	}
	return nil
}
