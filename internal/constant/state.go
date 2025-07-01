package constant

type ServerState string

const (
	ServerStateConnect  ServerState = "connect"
	ServerStateServing  ServerState = "serving"
	ServerStateStopped  ServerState = "stopped"
	ServerStateNoServer ServerState = "no_server"
)
