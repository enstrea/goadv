package srv

import "time"

type HttpOption func(serv *HttpServer)

func Timeout(timeout time.Duration) HttpOption {
	return func(serv *HttpServer) {
		serv.timeout = timeout
	}
}

func BeforeStart(beforeStart func(chan<- struct{})) HttpOption {
	return func(serv *HttpServer) {
		serv.BeforeStart = beforeStart
	}
}

func AfterStop(afterStop func(chan<- struct{})) HttpOption {
	return func(serv *HttpServer) {
		serv.AfterStop = afterStop
	}
}
