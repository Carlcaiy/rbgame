package route

import (
	"fmt"
	"rbgame/epoll/examples/cmd"
	"rbgame/epoll/examples/pb"
	"rbgame/epoll/local"
	"rbgame/epoll/network"
	"time"
)

type Local struct {
	templetes map[uint32]*RoomTemplete
	rooms     map[uint32]*RoomInstance
	*local.BaseLocal
}

func NewLocal(st *network.ServerConfig) *Local {
	return &Local{
		BaseLocal: local.NewBase(st),
		templetes: make(map[uint32]*RoomTemplete),
		rooms:     make(map[uint32]*RoomInstance),
	}
}

func (l *Local) Init() {
	l.load_room_templete()

	l.BaseLocal.Init()
	l.AddRoute(cmd.ReqRoomList, l.getRoomList)
	l.AddRoute(cmd.ReqEnterRoom, l.reqEnterRoom)
	l.AddRoute(cmd.ReqLeaveRoom, l.reqLeaveRoom)
	l.AddRoute(cmd.GameOver, l.gameOver)
	l.AddRoute(cmd.ReqGateLogin, l.login)
	l.AddRoute(cmd.Offline, l.offline)
}

func (l *Local) login(c *network.Conn, msg *network.Message) error {
	data := new(pb.ReqGateLogin)
	msg.UnPack(data)
	u, ok := l.GetUser(msg.UserID()).(*User)
	if !ok {
		u = &User{
			userID: msg.UserID(),
			gateID: data.GateId,
			hallID: l.ServerId,
		}
		l.AddUser(u)
		fmt.Printf("login add uid:%d game:%d room:%d gate:%d\n", msg.UserID(), u.gameID, u.roomID, data.GateId)
	} else {
		fmt.Printf("login contain uid:%d game:%d room:%d gate:%d\n", msg.UserID(), u.gameID, u.roomID, data.GateId)
	}

	if u.gameID > 0 && u.roomID > 0 {
		buf, _ := network.Pack(msg.UserID(), network.ST_Game, cmd.Reconnect, &pb.Reconnect{
			GateId: u.gateID,
			RoomId: u.roomID,
		})
		l.SendToGame(u.gameID, buf)
	} else {
		buf, _ := network.Pack(msg.UserID(), network.ST_Client, cmd.ResGateLogin, &pb.ResGateLogin{
			GameId: u.gameID,
			RoomId: u.roomID,
			HallId: l.ServerId,
		})
		c.Write(buf)
	}
	return nil
}

func (l *Local) offline(c *network.Conn, msg *network.Message) error {
	l.DelUser(msg.UserID())
	return nil
}

func (l *Local) load_room_templete() {
	for i := uint32(0); i < 5; i++ {
		l.templetes[i] = &RoomTemplete{
			TempId:    i,
			UserCount: 2,
			GameId:    1,
		}
	}
}

func (l *Local) gameOver(c *network.Conn, msg *network.Message) error {
	data := new(pb.GameOver)
	msg.UnPack(data)

	var room *RoomInstance
	if r, ok := l.rooms[data.RoomId]; ok {
		room = r
	}

	// 先发消息后处理数据
	if room != nil {
		if room.sitCount == room.UserCount {
			for _, user := range room.users {
				bs, _ := network.Pack(user.UserId(), network.ST_Client, cmd.CountDown, &pb.Empty{})
				l.SendToGate(user.gateID, bs)
			}
			l.Start(room.tevent)
		}

		for _, u := range room.users {
			bs, _ := network.Pack(u.userID, network.ST_Client, cmd.GameOver, data)
			l.SendToGate(u.gateID, bs)
		}
	}

	return nil
}

func (l *Local) getRoomList(c *network.Conn, msg *network.Message) error {
	fmt.Println("getRoomList")
	data := new(pb.ReqRoomList)
	msg.UnPack(data)
	res := new(pb.ResRoomList)
	res.Rooms = make([]*pb.RoomInfo, len(l.templetes))
	for i, temp := range l.templetes {
		res.Rooms[i] = &pb.RoomInfo{
			ServerId: 0,
			RoomId:   temp.TempId,
			Tag:      1,
		}
	}
	if buf, err := network.Pack(msg.UserID(), network.ST_Client, cmd.ResRoomList, res); err == nil {
		c.Write(buf)
	}
	return nil
}

func (l *Local) reqEnterRoom(c *network.Conn, msg *network.Message) error {
	data := new(pb.ReqEnterRoom)
	msg.UnPack(data)
	fmt.Printf("reqEnterRoom uid:%d tempId:%d\n", msg.UserID(), data.TempleteId)

	user, ok := l.GetUser(msg.UserID()).(*User)
	if !ok {
		return fmt.Errorf("not found user:%d", user.userID)
	}

	var room *RoomInstance

	// 寻找一个空的房间
	for _, r := range l.rooms {
		if r.TempId == data.TempleteId && r.status == 0 && r.sitCount < r.UserCount {
			room = r
		}
	}

	// 没找到房间，新建一个房间
	if room == nil {
		temp := l.templetes[data.TempleteId]
		room = &RoomInstance{
			RoomTemplete: temp,
			status:       0,
			users:        make([]*User, temp.UserCount),
			conn:         c,
			roomID:       uint32(len(l.rooms)) + 1000,
			tevent: local.NewDelayEvent(time.Second*10, func() {
				if room.sitCount < room.UserCount {
					return
				}
				room.status = 1
				greq := &pb.StartGame{
					TempId: room.TempId,
					RoomId: room.roomID,
					HallId: l.ServerId,
					Uids:   make([]uint32, room.sitCount),
					Gates:  make([]uint32, room.sitCount),
				}
				for i := range greq.Uids {
					greq.Uids[i] = room.users[i].userID
					greq.Gates[i] = room.users[i].gateID
					fmt.Printf("i:%d uid:%d gateid:%d\n", i, room.users[i].userID, room.users[i].gateID)
				}
				bs, _ := network.Pack(msg.UserID(), network.ST_Game, cmd.GameStart, greq)
				l.SendToGame(room.GameId, bs)
			}),
		}
		l.rooms[room.roomID] = room
	}

	if room != nil {
		user.roomID = room.roomID
		user.gameID = room.GameId

		for i := range room.users {
			if room.users[i] == nil {
				room.users[i] = user
				room.sitCount++
				break
			}
		}

		res := &pb.ResEnterRoom{}
		res.Uids = make([]uint32, 0, room.sitCount)
		for i := range room.users {
			if room.users[i] != nil {
				res.Uids = append(res.Uids, room.users[i].userID)
			}
		}

		bs, _ := network.Pack(msg.UserID(), network.ST_Client, cmd.ResEnterRoom, res)
		c.Write(bs)

		if room.sitCount == room.UserCount {
			for _, user := range room.users {
				bs, _ := network.Pack(user.UserId(), network.ST_Client, cmd.CountDown, &pb.Empty{})
				l.SendToGate(user.gateID, bs)
			}

			l.Start(room.tevent)
		}
	}

	return nil
}

func (l *Local) reqLeaveRoom(c *network.Conn, msg *network.Message) error {
	data := new(pb.ReqLeaveRoom)
	msg.UnPack(data)

	if room, ok := l.rooms[data.RoomId]; ok {
		if room.status == 0 {
			for i, u := range room.users {
				if u.UserId() == msg.UserID() {
					room.users[i] = nil
					room.sitCount -= 1
					l.Stop(room.tevent)
					return nil
				}
			}
		}
	}

	return fmt.Errorf("reqLeaveRoom error %s", data)
}

func (l *Local) Close(conn *network.Conn) {
	l.BaseLocal.Close(conn)
	// 大厅服，清理所有相关桌子
	if conn.ServerType == network.ST_Game {
		for _, room := range l.rooms {
			if room.GameId == conn.ServerId {
				room.status = 0
				room.sitCount = 0
				for i := range room.users {
					room.users[i] = nil
				}
			}
		}
	}
}
