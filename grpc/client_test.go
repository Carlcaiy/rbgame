package main

import (
	"context"
	"fmt"
	"rbgame/grpc/pb"
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
	key := struct{ Name string }{Name: "xxxx"}
	value := struct {
	}{}
	ctx := context.WithValue(context.Background(), key, value)
	reply, err := client().OneToOne(ctx, &pb.Request{Name: "hello"})
	fmt.Println(reply, err)
	reply, err = client().OneToOne(context.Background(), &pb.Request{Name: "good"})
	fmt.Println(reply, err)
	reply, err = client().OneToOne(context.Background(), &pb.Request{Name: ""})
	fmt.Println(reply, err)
}

func TestStreamClient(t *testing.T) {
	stream, err := client().MulToOne(context.Background())
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	for i := 0; i < 5; i++ {
		if err := stream.Send(&pb.Request{Name: "客户端流式RPC"}); err != nil {
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
	steam, err := client().OneToMul(context.Background(), &pb.Request{Name: "yoyo"})
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
	cc, err := client().MulToMul(context.Background())
	if err != nil {
		panic(err)
	}
	go func() {
		for i := 0; i < 10; i++ {
			cc.Send(&pb.Request{Name: fmt.Sprintf("擦儿啊i=%d", i)})
		}
	}()
	for i := 0; i < 10; i++ {
		xx, err := cc.Recv()
		fmt.Println("recv", xx.String(), err)
		if err != nil {
			break
		}
	}
	cc.CloseSend()
}
