package network

import (
	"bytes"
	"fmt"
	"rbgame/epoll/examples/pb"
	"testing"
)

func TestParser(t *testing.T) {
	msg1 := NewMessage(101, ST_Hall)
	bs1, _ := msg1.Pack(1000, &pb.HeartBeat{
		ServerType: uint32(ST_Client),
		ServerId:   100,
	})
	r := bytes.NewReader(bs1)
	msg2, err := Parse(r)

	fmt.Println(msg1)
	fmt.Println(msg2, err)

	bs2, _ := Pack(101, ST_Hall, 1000, &pb.HeartBeat{
		ServerType: uint32(ST_Client),
		ServerId:   100,
	})
	r1 := bytes.NewReader(bs2)
	msg3, err := Parse(r1)

	fmt.Println(msg3, err)
}
