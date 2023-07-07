package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Greet struct {
	UnimplementedGreeterServer
}

func (p *Greet) OneToOne(ctx context.Context, req *Request) (*Reply, error) {
	if req.Name == "hello" {
		return &Reply{Message: "hi"}, nil
	} else if req.Name == "good" {
		return &Reply{Message: "bye"}, nil
	} else {
		return nil, errors.New("nothing")
	}
}

func (p *Greet) OneToMul(req *Request, reply Greeter_OneToMulServer) error {
	if req.Name == "yoyo" {
		reply.Send(&Reply{Message: "hey guys"})
		reply.Send(&Reply{Message: "hey guys"})
		reply.Send(&Reply{Message: "hey guys"})
		reply.Send(&Reply{Message: "hey guys"})
		return nil
	}
	return errors.New("nothing")
}

func (p *Greet) MulToOne(req Greeter_MulToOneServer) error {
	for i := 0; i < 4; i++ {
		request, err := req.Recv()
		fmt.Println(request.Name)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
	req.SendAndClose(&Reply{Message: "completed"})
	return nil
}

func (p *Greet) MulToMul(mtm Greeter_MulToMulServer) error {
	for {
		recv, err := mtm.Recv()
		fmt.Println("Recv", recv.Name, err)
		if err != nil {
			break
		}
		err = mtm.Send(&Reply{Message: recv.Name + "1"})
		if err != nil {
			break
		}
	}
	return nil
}

func main() {
	grpcServer := grpc.NewServer()
	RegisterGreeterServer(grpcServer, &Greet{})

	ln, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}
	fmt.Println("listen success", ln.Addr().String())
	grpcServer.Serve(ln)
}
