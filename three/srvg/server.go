package srvg

import "net/http"

type Server struct {
	Name   string
	Server *http.Server

	BeforeStart func(chan<- struct{})
	AfterStop   func(chan<- struct{})
}

func NewServer(name string, server *http.Server, opts ...SrvOption) *Server {
	serv := &Server{
		Name:   name,
		Server: server,
	}

	for _, opt := range opts {
		opt(serv)
	}

	return serv
}
