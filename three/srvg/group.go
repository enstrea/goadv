package srvg

import (
	"context"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func New(opts ...Option) *serverGroup {
	servGroup := &serverGroup{
		wg:      new(sync.WaitGroup),
		servers: make(map[string]*Server),
		wait:    time.Second * 5,
		stop:    make(chan struct{}),
	}

	for _, opt := range opts {
		opt(servGroup)
	}

	return servGroup
}

type serverGroup struct {
	mu sync.Mutex
	wg *sync.WaitGroup

	servers map[string]*Server
	errors  sync.Map
	wait    time.Duration
	run     bool
	stopped bool
	stop    chan struct{}
}

func (sg *serverGroup) AddServer(server *Server) {
	sg.mu.Lock()
	defer sg.mu.Unlock()

	if sg.run {
		return
	}

	if _, ok := sg.servers[server.Name]; !ok {
		sg.servers[server.Name] = server
	}
}

func (sg *serverGroup) Run() {
	sg.mu.Lock()
	if sg.run {
		sg.mu.Unlock()
		return
	}
	sg.run = true
	sg.mu.Unlock()

	sg.startAll()

	// 监听退出
	go func() {
		<-sg.stop
		sg.stopAll()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)

	select {
	case <-quit:
		sg.Shutdown()
	case <-sg.stop:
	}

	sg.wg.Wait()
}

func (sg *serverGroup) Shutdown() {
	sg.mu.Lock()
	defer sg.mu.Unlock()

	if sg.stopped {
		return
	}
	sg.stopped = true
	close(sg.stop)
}

// 启动所有服务
func (sg *serverGroup) startAll() {
	for _, serv := range sg.servers {
		sg.wg.Add(1)
		go sg.startServer(serv)
	}
}

// 关闭所有服务
func (sg *serverGroup) stopAll() {
	ctx, cancel := context.WithTimeout(context.Background(), sg.wait)
	defer cancel()

	for _, serv := range sg.servers {
		go sg.stopServer(ctx, serv)
	}

	sg.wg.Wait()
}

func (sg *serverGroup) startServer(serv *Server) {
	if serv.BeforeStart != nil {
		done := make(chan struct{}, 1)
		go serv.BeforeStart(done)

		select {
		case <-done:
		case <-time.After(sg.wait):
			sg.errors.Store(serv.Name, errors.New("before start timeout"))
			sg.Shutdown()
			return
		}
	}

	if err := serv.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		sg.errors.Store(serv.Name, err)
		sg.Shutdown()
	}
}

func (sg *serverGroup) stopServer(ctx context.Context, serv *Server) {
	defer sg.wg.Done()

	if err := serv.Server.Shutdown(ctx); err != nil {
		sg.errors.Store(serv.Name, err)
	}

	if serv.AfterStop != nil {
		done := make(chan struct{}, 1)
		go serv.AfterStop(done)

		select {
		case <-done:
		case <-time.After(sg.wait):
			sg.errors.Store(serv.Name, errors.New("after stop timeout"))
		}
	}
}

func (sg *serverGroup) GetErrors() map[string]error {
	errMap := make(map[string]error)
	sg.errors.Range(func(key, value interface{}) bool {
		errMap[key.(string)] = value.(error)
		return true
	})
	return errMap
}
