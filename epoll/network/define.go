package network

type ServerType int32

const (
	ST_InValid   ServerType = iota // 无效的服务器
	ST_Client                      // 客户端
	ST_Gate                        // 网关服
	ST_Login                       // 登录服
	ST_Hall                        // 大厅服
	ST_Broadcast                   // 广播服
	ST_Money                       // 金币服
	ST_Game                        // 游戏服
)

func (s ServerType) String() string {
	switch s {
	case ST_Client:
		return "客户端"
	case ST_Gate:
		return "网关服"
	case ST_Login:
		return "登录服"
	case ST_Hall:
		return "大厅服"
	case ST_Broadcast:
		return "广播服"
	case ST_Money:
		return "金币服"
	case ST_Game:
		return "游戏服"
	default:
		return "未知的服务"
	}
}

type EventType int

const (
	ET_Close EventType = iota
	ET_Error
	ET_Timer
)

type Compoment interface {
	Init()
	Run()
	Close()
}

type Call func(*Conn, *Message) error
