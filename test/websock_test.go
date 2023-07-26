package test

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func TestWSServer(t *testing.T) {
	ln, _ := net.Listen("tcp", ":8080")

	for {
		conn, _ := ln.Accept()
		ws.Upgrade(conn)
		go func() {
			for {
				text, err := wsutil.ReadClientBinary(conn)
				if err != nil {
					fmt.Println("wsutil", err)
					break
				}
				fmt.Printf("receive:%+v\n", string(text))
			}
		}()
	}
}

func TestWSClient(t *testing.T) {
	conn, _, _, err := ws.Dial(context.Background(), ":8080")
	if err != nil {
		t.Log(err)
		return
	}
	for {
		err = wsutil.WriteServerBinary(conn, []byte("good"))
		if err != nil {
			t.Log(err)
			break
		}
		bs, _ := wsutil.ReadServerBinary(conn)
		fmt.Println(string(bs))
	}
}
