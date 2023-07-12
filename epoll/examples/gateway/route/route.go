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
}

func NewLocal(st *network.ServerConfig) *Local {
	return &Local{
		BaseLocal: local.NewBase(st),
	}
}

func (l *Local) Init() {
	l.BaseLocal.Init()
	l.AddRoute(cmd.ReqGateLogin, l.reqGateLogin)
	l.AddRoute(cmd.GateMulti, l.multi)
	l.AddRoute(cmd.ReqGateLeave, l.reqLeaveGate)
	l.AddHook(cmd.SyncData, l.sync)
	l.AddHook(cmd.ResGateLogin, l.resGateLogin)
}

func (l *Local) reqGateLogin(c *network.Conn, msg *network.Message) error {
	data := new(pb.ReqGateLogin)
	msg.UnPack(data)
	user, ok := l.GetUser(msg.UserID()).(*User)
	if ok {
		if c != user.Conn {
			buf, _ := network.Pack(msg.UserID(), network.ST_Client, cmd.GateKick, &pb.GateKick{
				Type: pb.KickType_Squeeze,
			})
			user.Write(buf)
			user.Conn = c
		} else {
			buf, _ := network.Pack(msg.UserID(), network.ST_Client, cmd.ResGateLogin, &pb.ResGateLogin{
				Ret: 1,
			})
			c.Write(buf)
		}
	} else {
		u := &User{
			uid:  msg.UserID(),
			Conn: c,
		}
		c.SetContext(u)
		l.AddUser(u)

		buf, _ := network.Pack(msg.UserID(), network.ST_Hall, cmd.ReqGateLogin, &pb.ReqGateLogin{
			GateId: l.ServerId,
		})
		l.SendModUid(u.uid, buf, network.ST_Hall)
	}
	return nil
}

func (l *Local) resGateLogin(c *network.Conn, msg *network.Message) error {
	pack := new(pb.ResGateLogin)
	msg.UnPack(pack)

	user, ok := l.GetUser(msg.UserID()).(*User)
	if !ok {
		return fmt.Errorf("resGateLogin not find user:%d", msg.UserID())
	}
	user.hallId = pack.HallId
	// user.gameId = pack.GameID
	return nil
}

func (l *Local) multi(c *network.Conn, msg *network.Message) error {
	pack := new(pb.MultiMsg)
	msg.UnPack(pack)
	for _, uid := range pack.Uids {
		if user := l.GetUser(uid); user != nil {
			user.Write(pack.Data)
		}
	}
	return nil
}

// 记录游戏服务ID
func (l *Local) sync(c *network.Conn, msg *network.Message) error {
	pack := new(pb.SyncData)
	msg.UnPack(pack)
	fmt.Printf("sync userId:%d gameId:%d roomId:%d\n", msg.UserID(), pack.GameId, pack.RoomId)
	if user, ok := l.GetUser(msg.UserID()).(*User); ok && user != nil {
		user.gameId = pack.GameId
	}
	return nil
}

// 离开网关
func (l *Local) reqLeaveGate(c *network.Conn, msg *network.Message) error {
	u, ok := c.Context().(*User)
	if ok {
		b, _ := network.Pack(u.UserID(), network.ST_Client, cmd.ResGateLeave, &pb.Empty{})
		l.SendToClient(u.UserID(), b)
		l.DelUser(u.UserID())
		return nil
	}
	return fmt.Errorf("reqLeaveGate not find user: %d", c.Fd)
}

func (l *Local) Close(conn *network.Conn) {
	l.BaseLocal.Close(conn)
	// 1、游戏服，踢出游戏内所有玩家
	if conn.ServerConfig == nil {
		u, ok := conn.Context().(*User)
		if ok {
			if u.gameId > 0 {
				b, _ := network.Pack(u.UserID(), network.ST_Game, cmd.Offline, &pb.Offline{})
				l.SendToSid(u.gameId, b, network.ST_Game)
			} else if u.hallId > 0 {
				b, _ := network.Pack(u.UserID(), network.ST_Hall, cmd.Offline, &pb.Offline{})
				l.SendToSid(u.hallId, b, network.ST_Hall)
			}
		}
	} else if conn.ServerType == network.ST_Game {
		l.RangeUser(func(u local.IUser) {
			if u.GameID() == conn.ServerId {
				b, _ := network.Pack(u.UserID(), network.ST_Client, cmd.GateKick, &pb.GateKick{
					Type: pb.KickType_GameNotFound,
				})
				l.SendToClient(u.UserID(), b)
			}
		})
	} else if conn.ServerType == network.ST_Client {
		u, ok := conn.Context().(*User)
		if ok {
			if u.gameId > 0 {
				b, _ := network.Pack(u.UserID(), network.ST_Game, cmd.Offline, &pb.Offline{})
				l.SendToSid(u.gameId, b, network.ST_Game)
			} else if u.hallId > 0 {
				b, _ := network.Pack(u.UserID(), network.ST_Hall, cmd.Offline, &pb.Offline{})
				l.SendToSid(u.hallId, b, network.ST_Hall)
			}
		}
	}
}
