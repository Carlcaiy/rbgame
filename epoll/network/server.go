package network

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var wg sync.WaitGroup

// 控制所有使用select阻塞功能的组件结束阻塞
var closech = make(chan struct{})

func Serve(sconf *ServerConfig, pconf *PollConfig, handle Handler) {

	handle.Init()

	poll := NewPoll(pconf)
	poll.AddListener(sconf)

	poll.handle = handle
	go poll.LoopRun()

	etcd := NewEtcd(sconf, poll.Trigger)
	etcd.Init()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	sig := <-ch
	if sig == syscall.SIGQUIT || sig == syscall.SIGTERM || sig == syscall.SIGINT {
		poll.Trigger(ET_Close)
		close(closech)
	}

	wg.Wait()
	poll.Close()
	etcd.Close()
}

func Client(sconf *ServerConfig, pconf *PollConfig) {

	poll := NewPoll(pconf)
	poll.AddConnector(sconf)
	go poll.LoopRun()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	sig := <-ch
	if sig == syscall.SIGQUIT || sig == syscall.SIGTERM || sig == syscall.SIGINT {
		poll.Trigger(ET_Close)
	}

	poll.Close()
}
