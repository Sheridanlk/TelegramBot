package rconpoller

import (
	"TelegramBot/internal/clients/rconclient"
	"time"
)

type Poller struct {
	rc       *rconclient.Client
	isActive bool
}

func New(rc *rconclient.Client) *Poller {
	return &Poller{
		rc:       rc,
		isActive: false,
	}
}

func (p *Poller) Start(command string, interval time.Duration) error {

	return nil
}

func (p *Poller) Stop() error {
	return nil
}

func (p *Poller) IsRun() (bool, error) {
	return true, nil
}
