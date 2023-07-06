package test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/nats-io/nats.go"
)

func ts(s string) []byte {
	return []byte(s)
}

func TestServer(t *testing.T) {
	wg := sync.WaitGroup{}
	server, err := nats.Connect("nats://for:bar@localhost:4222")
	if err != nil {
		panic(err)
	}
	_, err = server.Subscribe("foo", func(msg *nats.Msg) {
		t.Logf("%+v data=%s\n", msg.Subject, string(msg.Data))
		wg.Done()
	})
	wg.Add(10)
	for i := 1; i < 11; i++ {
		server.Publish("foo", ts(fmt.Sprintf("dasdasdasdad%d", i)))
		t.Log("publish foo")
	}
	if err != nil {
		panic(err)
	}
	wg.Wait()
	server.Close()
}
