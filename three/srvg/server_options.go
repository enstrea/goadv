package srvg

type SrvOption func(serv *Server)

func BeforeStart(beforeStart func(chan<- struct{})) SrvOption {
	return func(serv *Server) {
		serv.BeforeStart = beforeStart
	}
}

func AfterStop(afterStop func(chan<- struct{})) SrvOption {
	return func(serv *Server) {
		serv.AfterStop = afterStop
	}
}
