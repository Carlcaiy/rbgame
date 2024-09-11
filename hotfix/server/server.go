package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

var cmd []byte = make([]byte, 2)
var ring []byte = make([]byte, 1024*64)
var msgch = make(chan string)

func main() {
	re := flag.Bool("re", false, "-re")
	flag.Parse()
	var ln net.Listener
	conns := make([]net.Conn, 0)
	if *re {
		f := os.NewFile(3, "")
		ln, _ = net.FileListener(f)
		cf, _ := net.FileConn(os.NewFile(4, ""))
		conns = append(conns, cf)
		syscall.Kill(os.Getppid(), syscall.SIGTERM)
		time.Sleep(time.Second)
		loop(cf)
	} else {
		ln, _ = net.Listen("tcp", "127.0.0.1:9999")
	}
	fmt.Println("re:", *re)

	go func() {
		for {
			acc, err := ln.Accept()
			if err != nil {
				fmt.Println(err)
				break
			}
			conns = append(conns, acc)
			loop(acc)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTTIN, syscall.SIGTERM)
	for {
		select {
		case s := <-ch:
			if s == syscall.SIGTTIN {
				file1, err := ln.(*net.TCPListener).File()
				fmt.Println(err)
				cmd := exec.Command(os.Args[0], "-re")
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.ExtraFiles = append(cmd.ExtraFiles, file1)
				for _, conn := range conns {
					file2, _ := conn.(*net.TCPConn).File()
					cmd.ExtraFiles = append(cmd.ExtraFiles, file2)
				}
				cmd.Start()
			} else if s == syscall.SIGTERM {
				goto end
			}
		}
	}
end:
	fmt.Println("jies")
	for _, conn := range conns {
		conn.Close()
	}
	ln.Close()
}

func loop(conn net.Conn) {
	for {
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

		var buf = make([]byte, 1000)
		msg := "xx" + string(ring[:length])
		binary.BigEndian.PutUint16(buf, uint16(len(msg)))
		copy(buf[2:], []byte(msg))
		conn.Write(buf[:len(msg)+2])
	}
	fmt.Println("loop end")
}
