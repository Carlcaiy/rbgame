package local

import "rbgame/epoll/network"

type UserImplement struct {
	userId uint32
	gameId uint32
	gateId uint32
	*network.Conn
}

func (u *UserImplement) UserID() uint32 {
	return u.userId
}

func (u *UserImplement) GameID() uint32 {
	return u.gameId
}

func (u *UserImplement) GateID() uint32 {
	return u.gateId
}
