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

func (p *Processor) getSatus(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: get server status", err) }()
	//TODO: Сделать нормальный пинг, чтобы отличать выключенный сервер от ошики аутентификации
	if err := p.rcon.Connect(); err != nil {
		return p.tg.SendMessage(chatID, msgServerOffline)
	}
	return p.tg.SendMessage(chatID, msgServerOnline)
}

func (p *Processor) getPlayers(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: get players list", err) }()
	if err := p.rcon.Connect(); err != nil {
		return p.tg.SendMessage(chatID, msgFailedToConn)
	}

	resp, err := p.rcon.Execute("list")

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
	defer func() { err = e.Wrap("can't do command: start tracking players", err) }()
	// TODO: implemet function
	return nil
}

func (p *Processor) stopTracking(chatID int) (err error) {
	defer func() { err = e.Wrap("can't do command: stop tracking players", err) }()
	// TODO: implemet function
	return nil

}
