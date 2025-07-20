package storage

import (
	"context"
	"time"
)

type Storage interface {
	Save(ctx context.Context, p *Player) error
	Update(ctx context.Context, p *Player) error
	Remove(ctx context.Context, p *Player) error
	IsExists(ctx context.Context, p *Player) (bool, error)
	GetPlayersLastLogin(ctx context.Context) ([]Player, error)
}

type Player struct {
	Name      string
	LastVisit time.Time
}
