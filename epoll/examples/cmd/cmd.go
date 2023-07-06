package cmd

const (
	GateBegin    = 1000 // Gate
	ReqGateLogin = 1001 // 登录网关
	ResGateLogin = 1002 // 登录成功
	GateKick     = 1003 // 网关踢人
	GateMulti    = 1004 // 多播
	HeartBeat    = 1005 // 心跳
	Regist       = 1006 // 注册
	ReqGateLeave = 1007 // 离开网关
	ResGateLeave = 1008 // 离开网关
	Offline      = 1009 // 断线
	GateEnd      = 1999 // Gate
)
const (
	HallBegin      = 2000 // Hall
	ReqRoomList    = 2001 // 获取房间列表
	ResRoomList    = 2002 // 获取房间列表
	ReqEnterRoom   = 2003 // 进房请求
	ResEnterRoom   = 2004 // 进房答复
	ReportRoomData = 2005 // 上报桌子信息
	ReqLeaveRoom   = 2006 // 请求离开房间
	ResLeaveRoom   = 2007 // 答复离开房间
	HallEnd        = 2999 // Hall
)
const (
	GameBegin = 3000 // Game
	GameStart = 3001 // 进入游戏
	SyncData  = 3002 // 同步数据
	Tap       = 3003 // 点击游戏
	Round     = 3004 // 回合
	GameOver  = 3005 // 游戏结束
	Reconnect = 3006 // 重连
	CountDown = 3007 // 倒计时
	GameEnd   = 3999 // Game
)
