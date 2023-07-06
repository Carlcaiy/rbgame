package code

import (
	"bufio"
	"io"
	"net"
	"rbgame/websocket/gopool"
	"sync"
	"time"

	"github.com/gobwas/ws"
)

var pool *gopool.Pool
var poller *sync.Pool

func init() {
	pool = gopool.New(128)
}

func wsserve() {
	ln, _ := net.Listen("tcp", ":8080")

	for {
		// Try to accept incoming connection inside free pool worker.
		// If there no free workers for 1ms, do not accept anything and try later.
		// This will help us to prevent many self-ddos or out of resource limit cases.
		err := pool.ScheduleTimeout(time.Millisecond, func() {
			conn, _ := ln.Accept()
			_, _ = ws.Upgrade(conn)

			// Wrap WebSocket connection with our Channel struct.
			// This will help us to handle/send our app's packets.
			ch := NewChannel(conn)

			// Wait for incoming bytes from connection.
			poller.Start(conn, netpoll.EventRead, func() {
				// Do not cross the resource limits.
				pool.Schedule(func() {
					// Read and handle incoming packet(s).
					ch.Recevie()
				})
			})
		})
		if err != nil {
			time.Sleep(time.Millisecond)
		}
	}
}

// getReadBuf, putReadBuf are intended to
// reuse *bufio.Reader (with sync.Pool for example).
func getReadBuf(io.Reader) *bufio.Reader
func putReadBuf(*bufio.Reader)

// readPacket must be called when data could be read from conn.
func readPacket(conn io.Reader) (*Packet, error) {
	buf := getReadBuf(conn)
	defer putReadBuf(buf)

	buf.Reset(conn)
	frame, _ := ReadFrame(buf)
	parsePacket(frame.Payload)
	return nil, nil
}

func writePacket(conn io.Writer, pkt *Packet) error {
	conn.Write(pkt.Bytes())
	return nil
}

func (ch *Channel) Send(p Packet) {
	if ch.noWriterYet() {
		pool.Schedule(ch.writer)
	}
	ch.send <- p
}

func (ch *Channel) noWriterYet() bool {
	return false
}

// Packet represents application level data.
type Packet struct {
	bs []byte
}

func (p *Packet) Bytes() []byte {
	return p.bs
}

// Channel wraps user connection.
type Channel struct {
	conn net.Conn    // WebSocket connection.
	send chan Packet // Outgoing packets queue.
}

func NewChannel(conn net.Conn) *Channel {
	c := &Channel{
		conn: conn,
		send: make(chan Packet, 128),
	}

	go c.reader()
	go c.writer()

	return c
}

func (c *Channel) reader() {
	// We make a buffered read to reduce read syscalls.
	buf := bufio.NewReader(c.conn)

	for {
		pkt, _ := readPacket(buf)
		c.handle(pkt)
	}
}

func (c *Channel) writer() {
	// We make buffered write to reduce write syscalls.
	buf := bufio.NewWriter(c.conn)

	for pkt := range c.send {
		_ := writePacket(buf, pkt)
		buf.Flush()
	}
}

func (c *Channel) handle(pkt *Packet) {

}
