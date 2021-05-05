package main

import (
	"context"
	"fmt"
	srv2 "goadv/three/srv"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	servers := map[string]srv2.Server{
		"server1": server1("server1", ":8081"),
		"server2": server1("server2", ":8082"),
		"server3": server2("server3", ":8083"),
		"server4": server1("server4", ":8084"),
	}

	group := new(errgroup.Group)
	group, ctx := errgroup.WithContext(context.Background())

	stop := make(chan struct{})
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

	// 监听退出
	go func() {
		select {
		case <-quit:
			close(stop)
		}
	}()

	for _, serv := range servers {
		serv := serv
		group.Go(func() error {
			return serv.Start()
		})
		group.Go(func() error {
			select {
			case <-ctx.Done():
			case <-stop:
			}
			return serv.Stop()
		})
	}

	if err := group.Wait(); err != nil {
		fmt.Println("server err: ", err)
	}
}

func server1(name string, addr string) srv2.Server {
	serv := server(addr, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(name))
	}))

	return srv2.NewHttpServer(name, serv,
		srv2.BeforeStart(beforeStart(name)),
		srv2.AfterStop(afterStop(name)),
	)
}

func server2(name string, addr string) srv2.Server {
	serv := server(addr, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(name))
	}))

	return srv2.NewHttpServer(name, serv,
		srv2.BeforeStart(beforeStart2(name)),
		srv2.AfterStop(afterStop(name)),
	)
}

func server(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadTimeout:       time.Second,
		ReadHeaderTimeout: time.Second,
		WriteTimeout:      2 * time.Second,
		IdleTimeout:       5 * time.Second,
		MaxHeaderBytes:    1024,
	}
}

func beforeStart(name string) func(chan<- struct{}) {
	return func(c chan<- struct{}) {
		fmt.Println(fmt.Sprintf("%s start", name))
		time.Sleep(time.Second)
		c <- struct{}{}
	}
}

func beforeStart2(name string) func(chan<- struct{}) {
	return func(c chan<- struct{}) {
		fmt.Println(fmt.Sprintf("%s start", name))
		time.Sleep(time.Second * 10)
		c <- struct{}{}
	}
}

func afterStop(name string) func(chan<- struct{}) {
	return func(c chan<- struct{}) {
		fmt.Println(fmt.Sprintf("%s close", name))
		time.Sleep(time.Second)
		c <- struct{}{}
	}
}
