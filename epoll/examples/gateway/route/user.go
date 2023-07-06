package route

import "rbgame/epoll/network"

type User struct {
	uid           uint32 // 玩家uid
	gameId        uint32 // GameUid
	hallId        uint32 // HallId
	*network.Conn        // 连接
}

func (u *User) UserId() uint32 {
	return u.uid
}

func (u *User) GameId() uint32 {
	return u.gameId
}

func (u *User) GateId() uint32 {
	return 0
}
