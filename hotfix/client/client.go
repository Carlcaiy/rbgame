package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

var cmd []byte = make([]byte, 2)
var ring []byte = make([]byte, 1024*64)

func main() {
	conn, err := net.Dial("tcp", ":9999")
	if err != nil {
		panic(err)
	}

	go func() {
		for {

		}
	}()

	c := 0
	buf := make([]byte, 1000)
	for {
		time.Sleep(time.Millisecond * 1000)
		msg := fmt.Sprintf("messsage%d", c)
		binary.BigEndian.PutUint16(buf[:2], uint16(len(msg)))
		copy(buf[2:], []byte(msg))
		_, err := conn.Write(buf[:2+len(msg)])
		if err != nil {
			break
		}
		println(c)
		c++

		n, err := io.ReadFull(conn, cmd)
		if err != nil {
			fmt.Println("ReadFull", err)
			break
		}
		if n != 2 {
			fmt.Println("n != 2")
			break
		}
		length := binary.BigEndian.Uint16(cmd)
		n, err = io.ReadFull(conn, ring[:length])
		if err != nil {
			fmt.Println(err)
			break
		}
		if n != int(length) {
			fmt.Println("n != int(len)")
			break
		}
		fmt.Println(string(ring[:length]))
	}
}
