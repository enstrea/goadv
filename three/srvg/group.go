package srvg

import (
	"goadv/three/srvg/srv"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func New(opts ...Option) *serverGroup {
	servGroup := &serverGroup{
		wg:      new(sync.WaitGroup),
		servers: make(map[string]srv.Server),
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

	servers map[string]srv.Server
	errors  sync.Map
	run     bool
	stopped bool
	stop    chan struct{}
}

func (sg *serverGroup) AddServer(server srv.Server) {
	sg.mu.Lock()
	defer sg.mu.Unlock()

	if sg.run {
		return
	}

	if _, ok := sg.servers[server.Name()]; !ok {
		sg.servers[server.Name()] = server
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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)

	// 监听退出
	go func() {
		for {
			select {
			case <-sg.stop:
				sg.stopAll()
				return
			case <-quit:
				sg.Shutdown()
			}
		}
	}()

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
		serv := serv
		sg.wg.Add(1)

		go func() {
			if err := serv.Start(); err != nil {
				sg.errors.Store(serv.Name(), err)
				sg.Shutdown()
			}
		}()
	}
}

// 关闭所有服务
func (sg *serverGroup) stopAll() {
	for _, serv := range sg.servers {
		serv := serv

		go func() {
			defer sg.wg.Done()

			if err := serv.Stop(); err != nil {
				sg.errors.Store(serv.Name(), err)
			}
		}()
	}

	sg.wg.Wait()
}

func (sg *serverGroup) GetErrors() map[string]error {
	errMap := make(map[string]error)
	sg.errors.Range(func(key, value interface{}) bool {
		errMap[key.(string)] = value.(error)
		return true
	})
	return errMap
}
