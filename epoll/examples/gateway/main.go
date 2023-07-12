package main

import (
	"flag"
	"fmt"
	"rbgame/epoll/examples/gateway/route"
	"rbgame/epoll/network"
	"time"
)

func main() {

	port := 6666
	sid := 1
	flag.IntVar(&port, "p", 6666, "-p 6666")
	flag.IntVar(&sid, "s", 1, "-s 1")
	flag.Parse()

	serverConfig := &network.ServerConfig{
		Addr:       fmt.Sprintf(":%d", port),
		ServerType: network.ST_Gate,
		ServerId:   uint32(sid),
		Subs:       []network.ServerType{network.ST_Hall},
	}
	pollConfig := &network.PollConfig{
		HeartBeat: time.Millisecond * 100,
		MaxConn:   50000,
	}
	network.Serve(serverConfig, pollConfig, route.NewLocal(serverConfig))
}
