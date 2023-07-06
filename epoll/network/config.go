package network

import (
	"strconv"
	"strings"
)

type ServerConfig struct {
	Addr       string       // 服务地址
	ServerType ServerType   // 服务类型
	ServerId   uint32       // 服务ID
	Subs       []ServerType // 订阅的服务类型
}

func (s ServerConfig) IP() []byte {
	strs := strings.Split(s.Addr, ":")
	return []byte(strs[0])
}

func (s ServerConfig) Port() int {
	strs := strings.Split(s.Addr, ":")
	port, _ := strconv.Atoi(strs[1])
	return port
}

// 配置相同
func (s ServerConfig) Equal(a ServerConfig) bool {
	if s.Addr != a.Addr || s.ServerType != a.ServerType || s.ServerId != a.ServerId {
		return false
	}
	if len(s.Subs) != len(a.Subs) {
		return false
	}
	for i := 0; i < len(a.Subs); i++ {
		if s.Subs[i] != a.Subs[i] {
			return false
		}
	}
	return true
}

// 是否被订阅
func (p *ServerConfig) isSub(st ServerType) bool {
	for i := range p.Subs {
		if st == p.Subs[i] {
			return true
		}
	}
	return false
}
