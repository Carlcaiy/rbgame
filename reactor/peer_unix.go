//go:build !windows
// +build !windows

package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

type Poller struct {
	epfd  int
	wfd   int
	ln    int
	queue []func()
}

func NewPoller() *Poller {
	epfd, _ := unix.EpollCreate1(unix.EPOLL_CLOEXEC)
	wfd, _ := unix.Eventfd(0, unix.EFD_NONBLOCK|unix.EFD_CLOEXEC)
	unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, wfd, &unix.EpollEvent{Fd: int32(wfd), Events: unix.EPOLLPRI | unix.EPOLLIN})
	return &Poller{
		epfd: epfd,
		wfd:  wfd,
	}
}

func (p *Poller) Add(fd int) {
	err := unix.EpollCtl(p.epfd, unix.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Fd: int32(fd), Events: unix.EPOLLPRI | unix.EPOLLIN})
	fmt.Println("Add", p.epfd, fd, err)
}

func (p *Poller) Del(fd int) {
	err := unix.EpollCtl(p.epfd, unix.EPOLL_CTL_DEL, fd, &unix.EpollEvent{Fd: int32(fd), Events: unix.EPOLLPRI | unix.EPOLLIN})
	fmt.Println("Del", p.epfd, fd, err)
}

var a int32 = 0
var b = []byte{1, 1, 1, 1, 1, 1, 1, 1}

func (p *Poller) Loop() {
	go func() {
		events := make([]unix.EpollEvent, 100)
		timeout := -1
		for {
			wake := false
			n, err := unix.EpollWait(p.epfd, events, timeout)
			fmt.Println("loop", p.epfd, n, err)
			if n == 0 || (n < 1 && err == unix.EINTR) {
				timeout = -1
				runtime.Gosched()
				continue
			} else if err != nil {
				fmt.Println(err)
				return
			}
			timeout = 0
			for i := 0; i < n; i++ {
				fmt.Println("wait fd", events[i].Fd)
				if events[i].Fd == int32(p.wfd) {
					bs := make([]byte, 8)
					unix.Read(p.wfd, bs)
					fmt.Println("wake subreactor", bs, p.wfd)
					wake = true
				} else {
					if events[i].Events&(unix.EPOLLIN|unix.EPOLLPRI) > 0 {
						fd := int(events[i].Fd)
						fmt.Println("epoll in", fd)
						p.queue = append(p.queue, func() {
							bs := make([]byte, 128)
							n, _ := unix.Read(fd, bs)
							if n > 0 {
								fmt.Println(string(bs[:n]))
							} else {
								p.Del(fd)
							}
							unix.Write(fd, bs[:n])
						})

						if atomic.CompareAndSwapInt32(&a, 0, 1) {
							fmt.Println("wake msg")
							unix.Write(p.wfd, b)
						}
					}
					if events[i].Events&unix.EPOLLOUT > 0 {
						fmt.Println("can out")
					}
				}
			}

			atomic.StoreInt32(&a, 0)
			if wake {
				fmt.Println("wake")
				for i := range p.queue {
					p.queue[i]()
				}
				p.queue = p.queue[:0]
			}
		}
	}()
}

func server() {
	cur := 0
	subs := make([]*Poller, 3)
	for i := 0; i < 3; i++ {
		subs[i] = NewPoller()
		subs[i].Loop()
	}

	main := NewPoller()
	fd, _ := unix.Socket(unix.AF_INET, unix.SOCK_STREAM|unix.SOCK_NONBLOCK|unix.SOCK_CLOEXEC, unix.IPPROTO_TCP)
	unix.Bind(fd, &unix.SockaddrInet4{Port: 8080, Addr: [4]byte{127, 0, 0, 1}})
	unix.Listen(fd, 10)
	main.ln = fd
	main.Add(fd)
	events := make([]unix.EpollEvent, 100)

	timeout := -1
	for {
		n, err := unix.EpollWait(main.epfd, events, timeout)
		fmt.Println("accept", n, err)
		if n == 0 || (n < 1 && err == unix.EINTR) {
			timeout = -1
			runtime.Gosched()
			continue
		} else if err != nil {
			fmt.Println(err)
			return
		}
		timeout = 0
		for i := 0; i < n; i++ {
			if events[i].Fd == int32(main.wfd) {
				fmt.Println("wake main reactor")
			} else {
				nfd, sa, err := unix.Accept(int(events[i].Fd))
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("new connnect socket", sa)
				subs[cur].Add(nfd)
				cur++
			}
		}
	}
}

func client() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal("client", err)
	}

	fmt.Println("connect success")
	for {
		time.Sleep(1 * time.Second) // wait for close on the server side
		if _, err := conn.Write([]byte("hello")); err != nil {
			log.Printf("client: %v", err)
		}
		fmt.Println("write success")
		data := make([]byte, 10)
		n, err := conn.Read(data)
		if err != nil {
			log.Printf("client: %v", err)
			if errors.Is(err, syscall.ECONNRESET) {
				log.Print("This is connection reset by peer error")
			}
			continue
		}
		fmt.Println("recv", string(data[:n]))
	}
}

func main() {
	server()
	// time.Sleep(time.Second * 3)
	// client()
}
