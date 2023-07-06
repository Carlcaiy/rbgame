package main

import (
	"flag"
	"fmt"
	"rbgame/epoll/examples/hall/route"
	"rbgame/epoll/network"
	"time"
)

func main() {
	port := 6676
	sid := 1
	flag.IntVar(&port, "p", 6676, "-p 6676")
	flag.IntVar(&sid, "s", 1, "-s 1")
	flag.Parse()

	serverConfig := &network.ServerConfig{
		Addr:       fmt.Sprintf(":%d", port),
		ServerType: network.ST_Hall,
		ServerId:   uint32(sid),
		Subs:       []network.ServerType{network.ST_Game},
	}

	pollConfig := &network.PollConfig{
		HeartBeat: time.Millisecond * 100,
		MaxConn:   1000,
	}

	network.Serve(serverConfig, pollConfig, route.NewLocal(serverConfig))
}
