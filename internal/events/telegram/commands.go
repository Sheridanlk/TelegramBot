package telegram

import (
	"TelegramBot/internal/probe"
	"TelegramBot/lib/e"
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	StartCmd          = "/start"
	HelpCmd           = "/help"
	StartTrackingCmd  = "/start_track"
	StopTrackingCmd   = "/stop_track"
	GetPlayersCmd     = "/list"
	GetSastusCmd      = "/status"
	GetPalyersStatCmd = "/stat"
)

const (
	GetPlayerListRcon = "/list"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	switch text {
	case StartCmd:
		return p.sendHello(chatID)
	case HelpCmd:
		return p.sendHelp(chatID)
	case GetSastusCmd:
		return p.getSatus(chatID)
	case GetPlayersCmd:
		return p.getPlayers(chatID)
	case StartTrackingCmd:
		return p.startTracking(chatID)
	case StopTrackingCmd:
		return p.stopTracking(chatID)
	case GetPalyersStatCmd:
		return p.getPlayersStat(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)

	}
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func (p *Processor) getSatus(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: get server status", err) }()
	if probe.IsOnline(p.rcon.GetAddres()) {
		return p.tg.SendMessage(chatID, msgServerOnline)
	}
	return p.tg.SendMessage(chatID, msgServerOffline)
}

func (p *Processor) getPlayers(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: get players list", err) }()
	if !probe.IsOnline(p.rcon.GetAddres()) {
		p.rcon.Close()
		return p.tg.SendMessage(chatID, msgFailedToConn)
	}

	resp, err := p.rcon.Execute(GetPlayerListRcon)

	if err != nil {
		return err
	}

	err = p.tg.SendMessage(chatID, resp)

	if err != nil {
		return err
	}

	return nil
}

func (p *Processor) startTracking(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: start tracking players", err) }()
	if p.rconpoller.IsRun() {
		return p.tg.SendMessage(chatID, msgPollerIsRun)
	}

	p.rconpoller.Start(GetPlayerListRcon, 30*time.Second)
	return p.tg.SendMessage(chatID, msgPollerStart)
}

func (p *Processor) stopTracking(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: stop tracking players", err) }()
	if p.rconpoller.IsRun() {
		return p.tg.SendMessage(chatID, msgPollerIsStopped)
	}
	p.rconpoller.Stop()
	return p.tg.SendMessage(chatID, msgPollerStop)
}

func (p *Processor) getPlayersStat(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: get players status", err) }()
	players, err := p.storage.GetPlayersLastLogin(context.Background())

	if err != nil {
		return err
	}

	var output strings.Builder
	for _, player := range players {
		fmt.Fprintf(&output, "%s - %s\n", player.Name, player.LastVisit.Format("02.01.2006 15:04"))
	}
	return p.tg.SendMessage(chatID, "Players:\n"+output.String())
}
