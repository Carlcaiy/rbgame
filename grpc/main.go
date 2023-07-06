package main

import (
	"context"
	"fmt"
	"log"
	"net"
	pb "rbgame/proto/grpc"
	"time"

	"google.golang.org/grpc"
)

type Greet struct {
	pb.UnimplementedGreeterServer
}

func (p *Greet) SayHello(context.Context, *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "HeyBro"}, nil
}

func (p *Greet) SayHelloST(svr pb.Greeter_SayHelloSTServer) error {
	for {
		request, err := svr.Recv()
		fmt.Println(request, err)
		if err != nil {
			break
		}
	}
	return svr.SendAndClose(&pb.HelloReply{Message: "眼神中飘逸"})
}

func (p *Greet) SayHelloAgain(xx *pb.HelloRequest, bb pb.Greeter_SayHelloAgainServer) error {
	fmt.Println(xx.Name)
	for i := 0; i < 5; i++ {
		bb.Send(&pb.HelloReply{Message: fmt.Sprintf("response %d", i)})
	}
	return nil
}

func (p *Greet) SayHelloSSTT(s pb.Greeter_SayHelloSSTTServer) error {
	go func() {
		for {
			time.Sleep(time.Millisecond * 200)
			msg, err := s.Recv()
			fmt.Println("recv", msg, err)
			if err != nil {
				break
			}
		}
	}()
	for i := 0; i < 15; i++ {
		time.Sleep(time.Millisecond * 300)
		err := s.Send(&pb.HelloReply{Message: fmt.Sprintf("SSTT%d", i)})
		fmt.Println("send", err)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	grpcServer := grpc.NewServer()
	pb.RegisterGreeterServer(grpcServer, &Greet{})

	ln, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}
	fmt.Println("listen success", ln.Addr().String())
	grpcServer.Serve(ln)
}
