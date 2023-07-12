package local

import (
	"fmt"
	"net"
	"rbgame/epoll/examples/cmd"
	"rbgame/epoll/examples/pb"
	"rbgame/epoll/network"
	"time"
)

type IUser interface {
	UserID() uint32
	GameID() uint32
	GateID() uint32
	net.Conn
}

type BaseLocal struct {
	*network.ServerConfig
	m_users   map[uint32]IUser
	m_servers map[network.ServerType][]*network.Conn
	m_route   map[uint16]network.Call
	m_hook    map[uint16]network.Call
	*Timer
}

func NewBase(sconf *network.ServerConfig) *BaseLocal {
	return &BaseLocal{
		ServerConfig: sconf,
		m_users:      make(map[uint32]IUser),
		m_servers:    make(map[network.ServerType][]*network.Conn),
		m_route:      make(map[uint16]network.Call),
		Timer:        NewTimer(1024),
		m_hook:       make(map[uint16]network.Call),
	}
}

func (l *BaseLocal) Init() {
	fmt.Println("base.Init")
	l.AddRoute(cmd.HeartBeat, l.HeartBeat)
	l.AddRoute(cmd.Regist, l.Regist)
	l.AddRoute(cmd.Test, l.TestRequest)
	l.StartTimer(time.Minute, l.TimerHeartBeat, true)
}

func (l *BaseLocal) StartTimer(dur time.Duration, f func(), loop bool) {
	e := &TimerEvent{
		duration:    dur,
		triggerTime: time.Now().Add(dur),
		Loop:        loop,
		event:       f,
	}
	l.Timer.Push(e)
}

func (l *BaseLocal) TimerHeartBeat() {
	for _, t := range l.ServerConfig.Subs {
		if servers, ok := l.m_servers[t]; ok {
			for _, s := range servers {
				msg := network.NewMessage(s.ServerId, t)
				bs, _ := msg.Pack(cmd.HeartBeat, &pb.HeartBeat{
					ServerType: uint32(t),
					ServerId:   s.ServerId,
				})
				// fmt.Printf("send heart beat addr:%s type:%s id:%d\n", s.Addr, s.ServerType, s.ServerId)
				s.Write(bs)
			}
		}
	}
}

func (l *BaseLocal) OnConnect(conn *network.Conn) {
	if sli, ok := l.m_servers[conn.ServerType]; ok {
		// 更新server
		for i := range sli {
			if sli[i].ServerId == conn.ServerId {
				fmt.Printf("AddConn origin:%+v new:%+v\n", sli[i], conn.ServerConfig)
				sli[i] = conn
				return
			}
		}
	}
	// 尾部加入server
	fmt.Printf("AddConn new:%+v\n", conn.ServerConfig)
	l.m_servers[conn.ServerType] = append(l.m_servers[conn.ServerType], conn)

	// 注册server
	msg := network.NewMessage(0, conn.ServerType)
	bs, _ := msg.Pack(cmd.Regist, &pb.Regist{
		ServerId:   l.ServerId,
		ServerType: uint32(l.ServerType),
	})
	conn.Write(bs)
}

func (l *BaseLocal) OnAccept(conn *network.Conn) {

}

func (l *BaseLocal) Close(conn *network.Conn) {
	fmt.Printf("DelConn:%+v\n", conn)
	// 如果是用户连接删除用户数据即可
	if conn.ServerConfig == nil {
		u, ok := conn.Context().(IUser)
		if ok {
			l.DelUser(u.UserID())
		}
		return
	}

	if sli, ok := l.m_servers[conn.ServerType]; ok {
		index := -1
		for i := range sli {
			// 找到对应的元素，移到最后
			if sli[i].ServerId == conn.ServerId {
				fmt.Printf("DelConn index:%d new:%+v\n", i, conn)
				sli[i] = nil
				index = i
			} else if index >= 0 {
				sli[i], sli[index] = sli[index], sli[i]
				index = i
			}
		}
		if index >= 0 {
			// 长度-1
			l.m_servers[conn.ServerType] = sli[:index]
		}
	}
}

func (l *BaseLocal) RangeUser(iter func(u IUser)) {
	for _, u := range l.m_users {
		iter(u)
	}
}

func (l *BaseLocal) Regist(conn *network.Conn, msg *network.Message) error {
	data := new(pb.Regist)
	msg.UnPack(data)
	st := network.ServerType(data.ServerType)
	if sli, ok := l.m_servers[st]; ok {
		for i := range sli {
			if sli[i].ServerConfig != nil && sli[i].ServerId == data.ServerId {
				return fmt.Errorf("re register config: %+v", sli[i])
			}
		}
	}
	conf := &network.ServerConfig{
		Addr:       conn.RemoteAddr().String(),
		ServerType: st,
		ServerId:   data.ServerId,
	}
	conn.ServerConfig = conf
	l.m_servers[st] = append(l.m_servers[st], conn)
	fmt.Printf("regist serverId:%d serverType:%s addr:%s\n", conf.ServerId, conf.ServerType, conf.Addr)
	return nil
}

func (l *BaseLocal) HeartBeat(conn *network.Conn, msg *network.Message) error {
	data := new(pb.HeartBeat)
	msg.UnPack(data)
	// fmt.Println("receive HeartBeat", data.String())
	return nil
}

func (l *BaseLocal) TestRequest(conn *network.Conn, msg *network.Message) error {
	data := new(pb.Test)
	msg.UnPack(data)

	l.AddUser(&UserImplement{
		userId: data.Uid,
		Conn:   conn,
	})

	m := network.NewMessage(data.Uid, network.ST_Client)
	b, _ := m.Pack(cmd.Test, &pb.Test{
		Uid:       data.Uid,
		StartTime: data.StartTime,
	})

	return l.SendToClient(data.Uid, b)
}

func (l *BaseLocal) Tick() {
	l.FrameCheck()
}

func (l *BaseLocal) AddHook(cmd uint16, h network.Call) {
	if _, ok := l.m_hook[cmd]; !ok {
		l.m_hook[cmd] = h
	} else {
		fmt.Printf("err: hook.Add cmd:%d\n", cmd)
	}
}

func (l *BaseLocal) AddRoute(cmd uint16, h network.Call) {
	if _, ok := l.m_route[cmd]; !ok {
		l.m_route[cmd] = h
	} else {
		fmt.Printf("err: handler.Add cmd:%d\n", cmd)
	}
}

func (l *BaseLocal) Route(conn *network.Conn, msg *network.Message) error {

	if msg.UserID() == 0 && msg.Cmd() != cmd.ReqGateLogin && msg.Cmd() != cmd.HeartBeat && msg.Cmd() != cmd.Regist {
		return fmt.Errorf("msg wrong")
	}
	u := l.GetUser(msg.UserID())

	// 优先调用钩子
	if f, ok := l.m_hook[msg.Cmd()]; ok {
		f(conn, msg)
	}

	switch msg.Dest() {
	// 优先调用与本服务
	case l.ServerType:
		if f, ok := l.m_route[msg.Cmd()]; ok {
			return f(conn, msg)
		} else {
			return fmt.Errorf("call: not find cmd %d", msg.Cmd())
		}
	case network.ST_Client:
		return l.SendToClient(msg.UserID(), msg.Bytes())
	case network.ST_Hall:
		return l.SendToHall(msg.UserID(), msg.Bytes())
	case network.ST_Game:
		return l.SendToGame(u.GameID(), msg.Bytes())
	case network.ST_Gate:
		return l.SendToGate(u.GateID(), msg.Bytes())
	}
	return fmt.Errorf("call: not find cmd %d", msg.Cmd())
}

func (l *BaseLocal) SendModUid(uid uint32, buf []byte, t network.ServerType) error {
	if conns, ok := l.m_servers[t]; ok {
		if len(conns) > 0 {
			conn := conns[uid%uint32(len(conns))]
			conn.Write(buf)
			return nil
		} else {
			return fmt.Errorf("server %s size = 0", t)
		}
	} else {
		return fmt.Errorf("error not find server %s", t)
	}
}

func (l *BaseLocal) SendToSid(serverId uint32, buf []byte, t network.ServerType) error {
	if servers, ok := l.m_servers[t]; ok {
		for i := range servers {
			if servers[i].ServerId == serverId {
				servers[i].Write(buf)
				return nil
			}
		}
	}
	return fmt.Errorf("SendToSid: serverType=%s serverId=%d not found", t, serverId)
}

// attention: gateway use this function, other server should be careful
func (l *BaseLocal) SendToClient(uid uint32, buf []byte) error {
	if u, ok := l.m_users[uid]; ok {
		u.Write(buf)
		return nil
	} else {
		return fmt.Errorf("user:%d not find", uid)
	}
}

func (l *BaseLocal) SendToGame(gameId uint32, buf []byte) error {
	if conns, ok := l.m_servers[network.ST_Game]; ok {
		for _, s := range conns {
			if s.ServerId == gameId {
				s.Write(buf)
				return nil
			}
		}
		return fmt.Errorf("game server %d not find", gameId)
	} else {
		return fmt.Errorf("error not find game server")
	}
}

func (l *BaseLocal) SendToGate(gateId uint32, buf []byte) error {
	if conns, ok := l.m_servers[network.ST_Gate]; ok {
		for _, s := range conns {
			if s.ServerId == gateId {
				s.Write(buf)
				return nil
			}
		}
		return fmt.Errorf("game server %d not find", gateId)
	} else {
		return fmt.Errorf("error not find gate server")
	}
}

func (l *BaseLocal) SendToHall(uid uint32, buf []byte) error {
	return l.SendModUid(uid, buf, network.ST_Hall)
}

func (l *BaseLocal) GetUser(uid uint32) IUser {
	if u, ok := l.m_users[uid]; ok {
		return u
	}
	return nil
}

func (l *BaseLocal) AddUser(user IUser) {
	l.m_users[user.UserID()] = user
}

func (l *BaseLocal) DelUser(uid uint32) {
	delete(l.m_users, uid)
}
