package main

import (
	"context"
	"fmt"
	"rbgame/pb"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func client() pb.GreeterClient {
	client, err := grpc.Dial("localhost:8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return pb.NewGreeterClient(client)
}

func TestRPC(t *testing.T) {
	reply, err := client().SayHello(context.Background(), &pb.HelloRequest{Name: "普通RPC"})
	fmt.Println(reply, err)
}

func TestStreamClient(t *testing.T) {
	stream, err := client().SayHelloST(context.Background())
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	for i := 0; i < 5; i++ {
		if err := stream.Send(&pb.HelloRequest{Name: "客户端流式RPC"}); err != nil {
			fmt.Println("Send err", err)
			break
		}
		fmt.Println("Send Message Successful", i)
		time.Sleep(time.Millisecond)
	}

	recv, err := stream.CloseAndRecv()
	fmt.Println("RecvMsg", recv, err)
}

func TestStreamServer(t *testing.T) {
	steam, err := client().SayHelloAgain(context.Background(), &pb.HelloRequest{Name: "服务端流式RPC"})
	if err != nil {
		panic(err)
	}
	for {
		data, err := steam.Recv()
		if err != nil {
			fmt.Println("receive", err)
			break
		}
		fmt.Println("receave", data)
	}
	fmt.Println("end", steam.CloseSend())
}

func TestDoubleStream(t *testing.T) {
	cc, err := client().SayHelloSSTT(context.Background())
	if err != nil {
		panic(err)
	}
	go func() {
		for i := 0; i < 10; i++ {
			cc.Send(&pb.HelloRequest{Name: fmt.Sprintf("擦儿啊i=%d", i)})
		}
	}()
	for {
		xx, err := cc.Recv()
		fmt.Println("recv", xx.String(), err)
		if err != nil {
			break
		}
	}
	cc.CloseSend()
}
