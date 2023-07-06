package route

import (
	"fmt"
	"math/rand"
	"rbgame/epoll/examples/cmd"
	"rbgame/epoll/examples/pb"
	"rbgame/epoll/network"
)

type User struct {
	uid     uint32
	gateId  uint32 // 网关ID
	tap     int32
	offline bool
}

type RoomTemplete struct {
	Id        int32
	UserCount int32
}

type Room struct {
	hall      *network.Conn
	hallId    uint32
	roomId    uint32
	tempId    uint32
	l         *Local
	Users     []*User
	turn      int
	guess_num int32
}

func (r *Room) Reset() {
	for _, u := range r.Users {
		u.tap = 0
	}
	r.guess_num = rand.Int31n(100) + 1
}

func (r *Room) GetUser(uid uint32) *User {
	for _, u := range r.Users {
		if u.uid == uid {
			return u
		}
	}
	return nil
}

func (r *Room) Offline(uid uint32) {
	if u := r.GetUser(uid); u != nil {
		fmt.Printf("user :%d offline\n", uid)
		u.offline = true
	}
}

func (r *Room) Reconnect(uid uint32, gateId uint32) {
	for i, u := range r.Users {
		if u.uid == uid {
			u.gateId = gateId
			u.offline = false

			fmt.Println(u.uid, u.gateId)
			bs, _ := network.Pack(u.uid, network.ST_Client, cmd.SyncData, &pb.SyncData{
				Data:   "Reconnect game success",
				RoomId: r.roomId,
			})
			r.l.SendToGate(u.gateId, bs)

			fmt.Println("Reconnect", "uid:", uid, "sit", i, "turn", r.turn)
			if i == r.turn {
				bs, _ := network.Pack(uid, network.ST_Client, cmd.Round, &pb.Empty{})
				r.l.SendToGate(u.gateId, bs)
			}
			return
		}
	}
	fmt.Printf("Reconnect error: not find uid:%d\n", uid)
}

func (r *Room) Start() {
	r.Reset()
	fmt.Println("Start", "turn:", r.turn)
	for i, u := range r.Users {
		fmt.Printf("uid:%d seat:%d gateId:%d\n", u.uid, i, u.gateId)
		bs, _ := network.Pack(u.uid, network.ST_Client, cmd.SyncData, &pb.SyncData{
			Data:   "game start",
			RoomId: r.roomId,
			GameId: r.l.ServerId,
		})
		r.l.SendToGate(u.gateId, bs)
	}
	u := r.Users[r.turn]
	bs, _ := network.Pack(u.uid, network.ST_Client, cmd.Round, &pb.Empty{})
	r.l.SendToGate(u.gateId, bs)
}

func (r *Room) Tap(uid uint32, tap int32) {
	u := r.Users[r.turn]
	if uid != u.uid {
		fmt.Printf("tap err, uid:%d should uid:%d\n", uid, u.uid)
		return
	}

	tips := ""
	if tap < r.guess_num {
		tips = "小了"
	} else if tap > r.guess_num {
		tips = "大了"
	} else {
		tips = fmt.Sprintf("答对了%d", r.guess_num)
	}
	for _, u := range r.Users {
		bs, _ := network.Pack(u.uid, network.ST_Client, cmd.Tap, &pb.Tap{
			Uid:    uid,
			RoomId: r.roomId,
			Tap:    tap,
			Tips:   tips,
		})
		r.l.SendToGate(u.gateId, bs)
	}

	if tap == r.guess_num {
		bs, _ := network.Pack(uid, network.ST_Hall, cmd.GameOver, &pb.GameOver{
			TempId: r.tempId,
			RoomId: r.roomId,
			Data:   "game over",
		})
		r.l.SendToSid(r.hallId, bs, network.ST_Hall)
		r.gameOver()
		return
	}

	r.turn = (r.turn + 1) % len(r.Users)
	u = r.Users[r.turn]
	bs, _ := network.Pack(u.uid, network.ST_Client, cmd.Round, &pb.Empty{})
	r.l.SendToGate(u.gateId, bs)
}

func (r *Room) gameOver() {
	fmt.Println("game over")
}

func (r *Room) SendOne(bs []byte) {
	r.hall.Write(bs)
}

func (r *Room) SendOther(uid uint32, bs []byte) {
	multi := &pb.MultiMsg{
		Data: bs,
	}
	for _, u := range r.Users {
		if u.uid != uid {
			multi.Uids = append(multi.Uids, u.uid)
		}
	}
	buf, _ := network.Pack(0, network.ST_Gate, cmd.GateMulti, multi)
	r.hall.Write(buf)
}

func (r *Room) SendAll(bs []byte) {
	multi := &pb.MultiMsg{
		Data: bs,
	}
	for _, u := range r.Users {
		multi.Uids = append(multi.Uids, u.uid)
	}
	buf, _ := network.Pack(0, network.ST_Gate, cmd.GateMulti, multi)
	r.hall.Write(buf)
}
