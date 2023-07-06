package main

import (
	"fmt"
	"rbgame/gof/pubsub"
	"time"
)

func main() {
	pub := pubsub.NewPublisher(time.Second, 10)
	sub1 := pub.Subscribe()
	sub2 := pub.SubscribeTopic(func(parm interface{}) bool {
		return parm.(string) == "topic"
	})
	go func() {
		// for {
		select {
		case a := <-sub1:
			fmt.Println("ch1", a)
		case b := <-sub2:
			fmt.Println("ch2", b)
		}
		// }
	}()
	pub.Publish("topic")
	time.Sleep(time.Second * 2)
}
