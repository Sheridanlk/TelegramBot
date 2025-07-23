package rconpoller

import (
	"TelegramBot/internal/clients/rconclient"
	"TelegramBot/internal/mcparse"
	"TelegramBot/internal/probe"
	"TelegramBot/internal/storage"
	"context"
	"log"
	"time"
)

type Poller struct {
	rc       *rconclient.Client
	storage  storage.Storage
	isActive bool

	stop chan struct{}
}

func New(rc *rconclient.Client, storage storage.Storage) *Poller {
	return &Poller{
		rc:       rc,
		storage:  storage,
		isActive: false,
		stop:     make(chan struct{}),
	}
}

func (p *Poller) Start(command string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	online := false
	log.Printf("poller started")
	go func() {
		for {
			if !online {
				select {
				case <-p.stop:
					return
				case <-ticker.C:
					online = probe.IsOnline(p.rc.GetAddres())
				}
				continue
			}

			select {
			case <-p.stop:
				return
			case <-ticker.C:
				resp, err := p.rc.Execute(command)
				if err != nil {
					online = false
				}

				names, isRightCommand := mcparse.ParsePlayersNames(resp)
				if isRightCommand && names != nil {
					for _, name := range names {
						exists, err := p.storage.IsExists(context.Background(), name)
						if err != nil {
							log.Print(err)
							return
						}

						player := &storage.Player{
							Name:      name,
							LastVisit: time.Now().Local(),
						}

						if exists {
							err := p.storage.Update(context.Background(), player)
							if err != nil {
								log.Print(err)
								return
							}
						} else {
							err := p.storage.Save(context.Background(), player)
							if err != nil {
								log.Print(err)
								return
							}
						}
					}
				}
			}
		}
	}()
}

func (p *Poller) Stop() {
	if p.isActive {
		return
	}
	log.Printf("poller is stopped")
	p.isActive = false
	p.stop <- struct{}{}
}

func (p *Poller) IsRun() bool {
	return p.isActive
}
