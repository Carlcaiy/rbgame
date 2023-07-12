package local

import (
	"fmt"
	"rbgame/epoll/examples/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPC struct {
	pb.RPCClient
	client *grpc.ClientConn
	addr   string
}

func NewGRPC(addr string) *GRPC {
	return &GRPC{
		addr: addr,
	}
}

func (g *GRPC) Init() {
	client, err := grpc.Dial("addr", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		return
	}
	g.client = client
	g.RPCClient = pb.NewRPCClient(client)
}

func (g *GRPC) Close() {
	g.client.Close()
}
