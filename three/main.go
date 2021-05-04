package main

import (
	"fmt"
	"goadv/three/srvg"
	"goadv/three/srvg/srv"
	"net/http"
	"time"
)

func main() {
	sg := srvg.New()

	sg.AddServer(server1("server1", ":8081"))
	sg.AddServer(server1("server2", ":8082"))
	sg.AddServer(server1("server3", ":8083"))
	sg.AddServer(server2("server4", ":8084"))

	sg.Run()

	fmt.Println(sg.GetErrors())
	fmt.Println("quit")
}

func server1(name string, addr string) srv.Server {
	serv := server(addr, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(name))
	}))

	return srv.NewHttpServer(name, serv,
		srv.BeforeStart(beforeStart(name)),
		srv.AfterStop(afterStop(name)),
	)
}

func server2(name string, addr string) srv.Server {
	serv := server(addr, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(name))
	}))

	return srv.NewHttpServer(name, serv,
		srv.BeforeStart(beforeStart2(name)),
		srv.AfterStop(afterStop(name)),
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
