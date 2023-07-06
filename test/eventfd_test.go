package test

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

var eventfd_num int64 = 0

func tobytes(val int) []byte {
	return (*(*[]byte)(unsafe.Pointer(&val)))
}

func TestEventfd(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			epfd, err := unix.EpollCreate1(unix.EPOLL_CLOEXEC)
			if err != nil {
				log.Fatalln(err)
			}

			efd, err := unix.Eventfd(0, unix.EFD_NONBLOCK|unix.EFD_CLOEXEC)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(efd)

			unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, efd, &unix.EpollEvent{Events: unix.POLLIN | unix.POLLHUP, Fd: int32(efd)})

			events := make([]unix.EpollEvent, 10)
			buf := []byte{1, 1, 1, 1, 1, 1, 1, 1}

			if id == 4 {
				unix.Write(efd, buf)
			}
			for {
				n, err := unix.EpollWait(epfd, events, 0)
				if err != nil {
					break
				}
				for i := 0; i < n; i++ {
					if events[i].Fd == int32(efd) {
						n, err := unix.Read(efd, buf)
						if n == 0 || err != nil {
							if err == unix.EAGAIN {
								continue
							}
						}

						for i := range buf {
							if buf[i] < 9 {
								buf[i]++
								break
							}
						}
						fmt.Println(id, n, buf, err)
						unix.Write(efd, buf)
						time.Sleep(time.Second)
					}
				}
			}
		}(i)
	}
	wg.Wait()
}
