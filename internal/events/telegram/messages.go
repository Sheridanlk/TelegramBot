package telegram

const msgHelp = `I can tell you about the list of players on the server, as well as the server status.
	list of commands:
	/list - outputs online players list
	/status - outputs server status
	/start_track - collects stasistic about players activity
	/stop_track
	/stat - outputs players statisitc`

const msgHello = "Hello! \n" + msgHelp
const (
	msgUnknownCommand     = "Unknown command"
	msgFailedToConn       = "Failed to connect server"
	msgServerOnline       = "Server online"
	msgServerOffline      = "Server offline"
	msgAuthenticationFail = "Authorization faild Check the address and password and try again."
	msgExecuteFail        = "Unable to execute command"
	msgPollerIsRun        = "The poller is already running"
	msgPollerIsStopped    = "The poller is already stopped"
	msgPollerStart        = "Poller started"
	msgPollerStop         = "Poller is stopped"
)
