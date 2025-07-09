package telegram

import (
	"TelegramBot/lib/e"
	"log"
	"strings"
)

const (
	// start: /start + start tracking:
	StartCmd = "/start"
	// help:
	HelpCmd = "/help"
	// start tracking Players:
	StartTrackingCmd = "/start_track"
	// stop tracking Players:
	StopTrackingCmd = "/stop_track"
	// get players list:
	GetPlayersCmd = "/list"
	// get server status:
	GetSastusCmd = "/status"
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

func (p *Processor) getPlayers(chatID int) (err error) {
	defer func() { err = e.Wrap("can't do comand: get players list", err) }()
	// TODO: implemet function
	return nil
}

func (p *Processor) getSatus(chatID int) (err error) {
	defer func() { err = e.Wrap("can't do comand: get server status", err) }()
	// TODO: implemet function
	return nil
}

func (p *Processor) startTracking(chatID int) (err error) {
	defer func() { err = e.Wrap("can't do comand: start tracking players", err) }()
	// TODO: implemet function
	return nil
}

func (p *Processor) stopTracking(chatID int) (err error) {
	defer func() { err = e.Wrap("can't do comand: stop tracking players", err) }()
	// TODO: implemet function
	return nil
}
