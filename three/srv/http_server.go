package srv

import (
	"context"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

func NewHttpServer(name string, server *http.Server, opts ...HttpOption) *HttpServer {
	httpServ := &HttpServer{
		name:    name,
		Server:  server,
		timeout: time.Second * 5,
	}

	for _, opt := range opts {
		opt(httpServ)
	}

	return httpServ
}

type HttpServer struct {
	name    string
	Server  *http.Server
	timeout time.Duration

	BeforeStart func(chan<- struct{})
	AfterStop   func(chan<- struct{})
}

func (s *HttpServer) Name() string {
	return s.name
}

func (s *HttpServer) Start() error {
	if s.BeforeStart != nil {
		done := make(chan struct{}, 1)
		go s.BeforeStart(done)

		select {
		case <-done:
		case <-time.After(s.timeout):
			return errors.New("before start timeout")
		}
	}

	if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *HttpServer) Stop() error {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	err = s.Server.Shutdown(ctx)

	if s.AfterStop != nil {
		done := make(chan struct{}, 1)
		go s.AfterStop(done)

		select {
		case <-done:
		case <-time.After(s.timeout):
			return errors.New("after stop timeout")
		}
	}

	return err
}
