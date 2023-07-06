package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"rbgame/epoll/examples/cmd"
	"rbgame/epoll/examples/pb"
	"rbgame/epoll/network"
	"syscall"
)

var uid int = 123
var roomId uint32 = 0
var gateid int = 6666

var req = make(map[int]func())
var des = make(map[int]string)

func add(cmd int, str string, f func()) {
	req[cmd] = f
	des[cmd] = str
}

func main() {
	flag.IntVar(&uid, "u", 123, "-u 123")
	flag.IntVar(&gateid, "p", 6666, "-p 6666")
	flag.Parse()

	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", gateid))
	if err != nil {
		panic(err)
	}
	add(cmd.ReqGateLogin, "请求登录", func() {
		msg := network.NewMessage(uint32(uid), network.ST_Gate)
		bs, _ := msg.Pack(cmd.ReqGateLogin, &pb.ReqGateLogin{})
		conn.Write(bs)
	})
	add(cmd.ReqRoomList, "请求房间列表", func() {
		msg := network.NewMessage(uint32(uid), network.ST_Hall)
		bs, _ := msg.Pack(cmd.ReqRoomList, &pb.ReqRoomList{})
		conn.Write(bs)
	})
	add(cmd.ReqEnterRoom, "请求进房间", func() {
		fmt.Println("请输入进入的房间")
		roomId := uint32(0)
		fmt.Scanln(&roomId)
		msg := network.NewMessage(uint32(uid), network.ST_Hall)
		bs, _ := msg.Pack(cmd.ReqEnterRoom, &pb.ReqEnterRoom{
			TempleteId: roomId,
		})
		conn.Write(bs)
	})
	add(cmd.Tap, "猜测一个数值", func() {
		fmt.Println("请输入键入的数值：")
		num := int32(0)
		fmt.Scanln(&num)
		msg := network.NewMessage(uint32(uid), network.ST_Game)
		bs, _ := msg.Pack(cmd.Tap, &pb.Tap{
			RoomId: roomId,
			Tap:    num,
		})
		conn.Write(bs)
	})
	add(cmd.ReqLeaveRoom, "请求离开房间", func() {
		msg := network.NewMessage(uint32(uid), network.ST_Hall)
		bs, _ := msg.Pack(cmd.ReqLeaveRoom, &pb.ReqLeaveRoom{
			RoomId: roomId,
		})
		conn.Write(bs)
	})

	go func() {
		show_op()
		for {
			msg, err := network.Parse(conn)
			if err != nil {
				break
			}
			fmt.Println("receive msg:", msg.Cmd())
			switch msg.Cmd() {
			case cmd.ResGateLogin:
				p := new(pb.ResGateLogin)
				err := msg.UnPack(p)
				if err != nil {
					fmt.Println(err)
					continue
				}
				if p.Ret == 0 {
					fmt.Println("login success")
					show_op()
				} else if p.Ret == 1 {
					fmt.Println("error: relogin")
				}
			case cmd.GateKick:
				p := new(pb.GateKick)
				msg.UnPack(p)
				if p.Type == pb.KickType_Unknow {
					fmt.Println("游戏服务关闭")
				} else if p.Type == pb.KickType_Squeeze {
					fmt.Println("挤号")
				} else if p.Type == pb.KickType_GameNotFound {
					fmt.Println("服务未发现")
				} else {
					fmt.Println("踢出：位置错误")
				}
			case cmd.ResRoomList:
				p := new(pb.ResRoomList)
				msg.UnPack(p)
				fmt.Println(p.String())
				show_op()
			case cmd.ResEnterRoom:
				p := new(pb.ResEnterRoom)
				msg.UnPack(p)
				for _, uids := range p.Uids {
					fmt.Println("房间玩家", uids)
				}
			case cmd.GameStart:
				p := new(pb.StartGame)
				msg.UnPack(p)
				fmt.Println(p.String())
			case cmd.SyncData:
				p := new(pb.SyncData)
				msg.UnPack(p)
				fmt.Println(p.String())
				roomId = p.RoomId
			case cmd.GameOver:
				p := new(pb.GameOver)
				msg.UnPack(p)
				fmt.Println(p.String())
			case cmd.Tap:
				p := new(pb.Tap)
				msg.UnPack(p)
				fmt.Println(p.String())
			case cmd.Round:
				show_op()
			case cmd.CountDown:
				fmt.Println("剩余时间10")
				// show_op()
			}
		}
	}()

	sig := make(chan os.Signal, 10)
	signal.Notify(sig, syscall.SIGTERM)
	<-sig

	show_op()
}

func show_op() {
	for {
		for cmd, des := range des {
			fmt.Printf("%d: %s\n", cmd, des)
		}
		fmt.Println("选择想执行的操作")
		cmd := 0
		fmt.Scanln(&cmd)
		if f, ok := req[cmd]; ok {
			f()
			break
		}
	}
}
