package rconpoller

import (
	"TelegramBot/internal/clients/rconclient"
	"TelegramBot/internal/mcparse"
	"TelegramBot/internal/probe"
	"TelegramBot/internal/storage"
	"context"
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"
)

type Poller struct {
	rc      *rconclient.Client
	storage storage.Storage

	mu   sync.Mutex
	stop chan struct{}
	done chan struct{}
	run  bool
}

func New(rc *rconclient.Client, storage storage.Storage) *Poller {
	return &Poller{
		rc:      rc,
		storage: storage,
		run:     false,
	}
}

func (p *Poller) Start(command string, interval time.Duration) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.run {
		return fmt.Errorf("the poller is already running")
	}

	p.stop = make(chan struct{})
	p.done = make(chan struct{})

	go p.pollLoop(p.stop, p.done, command, interval)
	p.run = true
	log.Printf("poller started")
	return nil
}

func (p *Poller) Stop() error {
	p.mu.Lock()
	if !p.run {
		p.mu.Unlock()
		return fmt.Errorf("the poller is already stopped")
	}
	close(p.stop)
	done := p.done
	p.run = false
	p.mu.Unlock()

	<-done
	log.Printf("poller is stopped")
	return nil
}

func (p *Poller) pollLoop(stop <-chan struct{}, done chan<- struct{}, command string, interval time.Duration) {
	const probeinterval = 60 * time.Second
	reqticker := time.NewTicker(interval)
	online := false

	defer func() {
		reqticker.Stop()
		close(done)
	}()

	for {
		if !online {
			select {
			case <-stop:
				return
			default:
				online = probe.IsOnline(p.rc.GetAddres())
				if !online {
					time.Sleep(probeinterval)
					continue
				}
			}
		}

		select {
		case <-stop:
			return
		case <-reqticker.C:
			resp, err := p.rc.Execute(command)
			if err != nil {
				online = false
				continue
			}

			prevnames := []string{}
			names, isRightCommand := mcparse.ParsePlayersNames(resp)
			rewrite := false

			if isRightCommand && names != nil {
				rewrite = reflect.DeepEqual(prevnames, names)
				prevnames = names
				for _, name := range names {
					player := &storage.Player{
						Name:      name,
						LastVisit: time.Now().Local(),
					}
					if rewrite {
						err := p.storage.Update(context.Background(), player)
						if err != nil {
							log.Print(err)
							return
						}
						continue
					}
					exists, err := p.storage.IsExists(context.Background(), name)
					if err != nil {
						log.Print(err)
						return
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
}
