package network

// #include <pthread.h>

import (
	"errors"
	"fmt"
	"net"
	_ "net/http/pprof"
	"reflect"
	"syscall"
	"time"

	reuseport "github.com/kavu/go_reuseport"
	"golang.org/x/sys/unix"
)

type Handler interface {
	Init()
	Route(conn *Conn, msg *Message) error
	Close(conn *Conn)
	OnConnect(conn *Conn)
	OnAccept(conn *Conn)
	Tick()
}

type PollConfig struct {
	HeartBeat time.Duration // 心跳间隔时间
	MaxConn   int           // 最大连接数
}

type Conn struct {
	*ServerConfig             // 信息
	*net.TCPConn              // 连接
	Fd            int         // 文件描述符
	ActiveTime    time.Time   // 活跃时间
	ctx           interface{} // 该链接附带信息
}

func (c *Conn) Context() interface{} {
	return c.ctx
}

func (c *Conn) SetContext(d interface{}) {
	c.ctx = d
}

type Poll struct {
	epollFd  int
	eventFd  int
	listenFd int
	listener *net.TCPListener
	fdconns  map[int]*Conn
	conn_num int
	ticker   *time.Ticker
	config   *PollConfig
	queue    esqueue
	handle   Handler
}

func NewPoll(conf *PollConfig) *Poll {
	epollFd, err := unix.EpollCreate1(0)
	must(err)
	eventFd, err := unix.Eventfd(0, unix.EFD_CLOEXEC) // unix.EFD_NONBLOCK|unix.EFD_CLOEXEC
	must(err)
	err = unix.EpollCtl(epollFd, unix.EPOLL_CTL_ADD, eventFd, &unix.EpollEvent{Events: unix.EPOLLIN, Fd: int32(eventFd)})
	must(err)

	return &Poll{
		fdconns: make(map[int]*Conn),
		epollFd: epollFd,
		eventFd: eventFd,
		ticker:  time.NewTicker(conf.HeartBeat),
		config:  conf,
	}
}

func (p *Poll) Close() {
	fmt.Println("poll close")
	if p.epollFd > 0 {
		unix.Close(p.epollFd)
	}
	if p.eventFd > 0 {
		unix.Close(p.eventFd)
	}
	if p.listenFd > 0 {
		unix.Close(p.listenFd)
	}
	for _, c := range p.fdconns {
		c.Close()
	}
}

func (p *Poll) LoopRun() {
	wg.Add(1)
	defer func() {
		wg.Done()
	}()

	go func() {
		for {
			<-p.ticker.C
			// fmt.Println("ticker", t)
			p.Trigger(ET_Timer)
		}
	}()

	events := make([]unix.EpollEvent, 64)
	for {
		n, err := unix.EpollWait(p.epollFd, events, 100)
		if err != nil && err != unix.EINTR {
			return
		}

		if err := p.queue.ForEach(func(note interface{}) error {
			switch t := note.(type) {
			case EventType:
				if t == ET_Timer {
					p.handle.Tick()
				} else if t == ET_Close {
					return errors.New("signal close")
				} else if t == ET_Error {
					return errors.New("error")
				}
			case *ServerConfig:
				p.AddConnector(t)
			}
			return nil
		}); err != nil {
			return
		}

		for i := 0; i < n; i++ {
			// fmt.Println("i:", i, "Fd:", events[i].Fd, "Events", events[i].Events, "pad", events[i].Pad, "listenFd", p.epollFd)
			fd := int(events[i].Fd)
			if p.eventFd == fd {
				data := make([]byte, 8)
				unix.Read(fd, data)
			} else if p.listenFd == fd {
				conn, err := p.listener.AcceptTCP()
				if err != nil {
					fmt.Println("AcceptTCP", err)
					continue
				}
				p.Add(conn)
			} else {
				conn := p.fdconns[fd]
				if conn == nil {
					fmt.Println("Get conn", err)
					continue
				}

				msg, err := Parse(conn)
				if err != nil {
					p.Del(fd)
					continue
				}
				fmt.Printf("Route uid:%d cmd:%d dst:%s\n", msg.UserID(), msg.Cmd(), msg.Dest())
				if p.handle != nil {
					if err := p.handle.Route(conn, msg); err != nil {
						fmt.Println(err)
					}
				} else {
					fmt.Println("handle nil, can't deal message")
				}
			}
		}
	}
}

func (p *Poll) AddListener(conf *ServerConfig) {
	if conf.Addr == "" {
		fmt.Println("error addr", conf.Addr)
		return
	}
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: conf.IP(), Port: conf.Port()})
	reuseport.Listen("tcp", conf.Addr)
	must(err)
	p.listenFd = listenFD(listener)
	p.listener = listener
	fmt.Printf("AddListener fd:%d conf:%+v\n", p.listenFd, conf)
	unix.EpollCtl(p.epollFd, syscall.EPOLL_CTL_ADD, p.listenFd, &unix.EpollEvent{Events: unix.EPOLLIN, Fd: int32(p.listenFd)})
}

func (p *Poll) AddConnector(conf *ServerConfig) {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: conf.IP(), Port: conf.Port()})
	if err != nil {
		fmt.Println(err)
		return
	}
	fd := socketFD(conn)
	ptr := &Conn{
		TCPConn:      conn,
		ServerConfig: conf,
		Fd:           fd,
	}
	p.fdconns[fd] = ptr
	p.handle.OnConnect(ptr)

	fmt.Printf("AddConnector fd:%d conf:%+v\n", fd, conf)
	unix.EpollCtl(p.epollFd, syscall.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: unix.EPOLLIN, Fd: int32(fd)})
}

func (p *Poll) Trigger(tri interface{}) {
	p.queue.Add(tri)
	unix.Write(p.eventFd, []byte{1, 1, 1, 1, 1, 1, 1, 1})
}

func (p *Poll) Del(fd int) {
	err := unix.EpollCtl(p.epollFd, syscall.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		fmt.Println("Del", err)
		return
	}
	server := p.fdconns[fd]
	fmt.Printf("Del fd:%d addr:%v conn_num=%d\n", fd, server.RemoteAddr(), p.conn_num)
	p.handle.Close(server)
	delete(p.fdconns, fd)
	p.conn_num--
	server.Close()
}

func (p *Poll) Add(conn *net.TCPConn) {
	if p.conn_num >= p.config.MaxConn {
		fmt.Println("conn num too much.", p.config.MaxConn)
		return
	}
	fd := socketFD(conn)
	err := unix.EpollCtl(p.epollFd, syscall.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: unix.EPOLLIN, Fd: int32(fd)})
	if err != nil {
		fmt.Println("Add", err)
		return
	}
	c := &Conn{
		TCPConn: conn,
		Fd:      fd,
	}
	p.fdconns[fd] = c
	p.conn_num++
	p.handle.OnAccept(c)
	fmt.Printf("Add fd:%d addr:%v conn_num=%d\n", fd, conn.RemoteAddr(), p.conn_num)
}

func socketFD(conn *net.TCPConn) int {
	tcpConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}

func listenFD(conn net.Listener) int {
	fdVal := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
