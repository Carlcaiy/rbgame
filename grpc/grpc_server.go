package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"rbgame/grpc/pb"

	"google.golang.org/grpc"
)

type Greet struct {
	pb.UnimplementedGreeterServer
}

func (p *Greet) OneToOne(ctx context.Context, req *pb.Request) (*pb.Reply, error) {
	fmt.Println(ctx.Value("name"))
	if req.Name == "hello" {
		return &pb.Reply{Message: "hi"}, nil
	} else if req.Name == "good" {
		return &pb.Reply{Message: "bye"}, nil
	} else {
		return nil, errors.New("nothing")
	}
}

func (p *Greet) OneToMul(req *pb.Request, reply pb.Greeter_OneToMulServer) error {
	if req.Name == "yoyo" {
		reply.Send(&pb.Reply{Message: "hey guys"})
		reply.Send(&pb.Reply{Message: "hey guys"})
		reply.Send(&pb.Reply{Message: "hey guys"})
		reply.Send(&pb.Reply{Message: "hey guys"})
		return nil
	}
	return errors.New("nothing")
}

func (p *Greet) MulToOne(req pb.Greeter_MulToOneServer) error {
	for i := 0; i < 4; i++ {
		request, err := req.Recv()
		fmt.Println(request.Name)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
	req.SendAndClose(&pb.Reply{Message: "completed"})
	return nil
}

func (p *Greet) MulToMul(mtm pb.Greeter_MulToMulServer) error {
	for {
		recv, err := mtm.Recv()
		fmt.Println("Recv", recv.Name, err)
		if err != nil {
			break
		}
		err = mtm.Send(&pb.Reply{Message: recv.Name + "1"})
		if err != nil {
			break
		}
	}
	return nil
}

func XXX(srv interface{}, stream grpc.ServerStream) error {

	return nil
}

func main() {
	grpcServer := grpc.NewServer(
		grpc.NumStreamWorkers(1),
		grpc.UnknownServiceHandler(XXX),
	)
	pb.RegisterGreeterServer(grpcServer, &Greet{})

	ln, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}
	fmt.Println("listen success", ln.Addr().String())
	grpcServer.Serve(ln)
}
