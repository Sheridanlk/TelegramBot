package telegram

const msgHelp = `I can tell you about the list of players on the server, as well as the server status.
	list of commands:
	/list
	/status
	/start_track
	/stop_track`

const msgHello = "Hello! \n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command"
	msgServerOffline  = "Server is offline now"
)
