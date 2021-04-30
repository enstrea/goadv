package srv

import "net/http"

type Server struct {
	Name   string
	Server *http.Server

	BeforeStart func(chan<- struct{})
	AfterStop   func(chan<- struct{})
}

func New(name string, server *http.Server, opts ...Option) *Server {
	serv := &Server{
		Name:   name,
		Server: server,
	}

	for _, opt := range opts {
		opt(serv)
	}

	return serv
}
