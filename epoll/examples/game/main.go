package main

import (
	"rbgame/epoll/examples/game/route"
	"rbgame/epoll/network"
	"time"
)

func main() {
	serverConfig := &network.ServerConfig{
		Addr:       ":6686",
		ServerType: network.ST_Game,
		ServerId:   1,
		Subs:       []network.ServerType{network.ST_Gate},
	}
	pollConfig := &network.PollConfig{
		HeartBeat: time.Millisecond * 100,
		MaxConn:   1000,
	}
	network.Serve(serverConfig, pollConfig, route.NewLocal(serverConfig))
}
