package main

import (
	"fmt"
	"rbgame/proto/pb"

	"google.golang.org/protobuf/proto"
)

func main() {
	t := &pb.Regist{
		Name:     "xxxxxxxxxxxxxxx",
		Password: "xxxxxxxxxxxxxxxxxxxxxx",
	}
	fmt.Println("pre size", proto.Size(t))
	bs, _ := proto.Marshal(t)
	fmt.Println("now size", len(bs))
}
