package route

import (
	"fmt"
	"rbgame/epoll/examples/cmd"
	"rbgame/epoll/examples/pb"
	"rbgame/epoll/local"
	"rbgame/epoll/network"
)

type Local struct {
	*local.BaseLocal
	rooms map[uint32]*Room // 房间配置
	users map[uint32]*Room
}

func NewLocal(st *network.ServerConfig) *Local {
	return &Local{
		BaseLocal: local.NewBase(st),
		rooms:     make(map[uint32]*Room),
		users:     make(map[uint32]*Room),
	}
}

func (l *Local) Init() {
	l.BaseLocal.Init()
	l.AddRoute(cmd.GameStart, l.startGame)
	l.AddRoute(cmd.Tap, l.tapGame)
	l.AddRoute(cmd.Reconnect, l.reconnect)
	l.AddRoute(cmd.Offline, l.offline)
}

func (l *Local) offline(c *network.Conn, msg *network.Message) error {
	if room, ok := l.users[msg.UserID()]; ok {
		room.Offline(msg.UserID())
	}
	return nil
}

func (l *Local) reconnect(c *network.Conn, msg *network.Message) error {
	pack := new(pb.Reconnect)
	msg.UnPack(pack)
	fmt.Println("reconnect", pack.String())
	if pack.RoomId > 0 {
		room, ok := l.rooms[pack.RoomId]
		if ok {
			room.Reconnect(msg.UserID(), pack.GateId)
		}
	}

	return nil
}

func (l *Local) startGame(c *network.Conn, msg *network.Message) error {
	data := new(pb.StartGame)
	msg.UnPack(data)
	fmt.Println("startGame", data.String())
	room, ok := l.rooms[data.RoomId]
	if !ok {
		room = &Room{
			hall:   c,
			hallId: data.HallId,
			roomId: data.RoomId,
			tempId: data.TempId,
			l:      l,
		}
		l.rooms[data.RoomId] = room
	}
	room.Users = make([]*User, len(data.Uids))
	for i := range room.Users {
		room.Users[i] = &User{
			uid:    data.Uids[i],
			gateId: data.Gates[i],
			tap:    0,
		}
		l.users[data.Uids[i]] = room
	}
	room.Start()
	return nil
}

func (l *Local) tapGame(c *network.Conn, msg *network.Message) error {
	data := new(pb.Tap)
	msg.UnPack(data)
	fmt.Println("tap game")
	if room, ok := l.rooms[data.RoomId]; ok {
		room.Tap(msg.UserID(), data.Tap)
	}

	return nil
}
