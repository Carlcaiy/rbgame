package main

import (
	"context"
	"log"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	conn, _, _, err := ws.Dial(context.Background(), "ws://localhost:8080")
	if err != nil {
		log.Println(err)
		return
	}
	for {
		if err := wsutil.WriteServerText(conn, []byte("PING")); err != nil {
			log.Println(err)
			break
		}
		time.Sleep(time.Second)
		bs, err := wsutil.ReadServerText(conn)
		log.Println(string(bs), err)
	}
}
