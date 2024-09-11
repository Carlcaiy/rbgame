package main

import (
	"log"
	"net"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	tcp_ws()
}
func http_ws() {
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			log.Println(err)
			w.Write([]byte(err.Error()))
			return
		}
		go func() {
			defer conn.Close()

			for {
				msg, op, err := wsutil.ReadClientData(conn)
				if err != nil {
					log.Println(err)
					w.Write([]byte(err.Error()))
					break
				}
				log.Println(string(msg), op)
				err = wsutil.WriteServerMessage(conn, op, msg)
				if err != nil {
					log.Println(err)
					w.Write([]byte(err.Error()))
					break
				}
			}

			log.Println("close")
		}()
	}))
}

func tcp_ws() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println(err)
		return
	}

	u := ws.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		OnHeader: func(key, value []byte) (err error) {
			log.Printf("non-websocket header: %q=%q", key, value)
			return
		},
		Protocol: func(b []byte) bool {
			log.Println(string(b))
			return true
		},
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
			return
		}
		_, err = u.Upgrade(conn)
		if err != nil {
			log.Printf("upgrade error: %s", err)
			return
		}
		log.Println("accept")
		go func() {
			defer conn.Close()

			for {
				msg, err := wsutil.ReadServerText(conn)
				if err != nil {
					log.Println(err)
					break
				}
				log.Println(string(msg))
				err = wsutil.WriteServerText(conn, msg)
				if err != nil {
					log.Println(err)
					break
				}
			}

			log.Println("close")
		}()
	}
}
