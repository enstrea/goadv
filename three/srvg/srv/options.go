package srv

type Option func(serv *Server)

func BeforeStart(beforeStart func(chan<- struct{})) Option {
	return func(serv *Server) {
		serv.BeforeStart = beforeStart
	}
}

func AfterStop(afterStop func(chan<- struct{})) Option {
	return func(serv *Server) {
		serv.AfterStop = afterStop
	}
}
