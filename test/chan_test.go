package test

import (
	"fmt"
	"testing"
)

func TestChan(t *testing.T) {
	ch := make(chan int)
	l1 := make(chan struct{}, 1)
	l2 := make(chan struct{}, 1)

	go func() {
		for {
			<-l2
			fmt.Println("ch0", <-ch)
			l1 <- struct{}{}
		}
	}()

	go func() {
		for {
			<-l1
			fmt.Println("ch1", <-ch)
			l2 <- struct{}{}
		}
	}()

	l1 <- struct{}{}
	for i := 0; ; i++ {
		ch <- i
	}

}

/*
这个问题的关键在于理解每个小时的价格是如何变化的。 第一个小时是10块，第二个小时是3块，第三个小时是4块，
以此类推，每增加一个小时，价格就增加1块，直到一天的最高价格达到35元。 我们可以用一个循环来计算一天的总价格。
首先，我们初始化总价格为0，然后对于每一天的每一个小时，我们根据上述规则计算价格，并将价格加到总价格上。
注意，因为一天有24小时，所以我们循环24次来计算每一小时的价格。 所以，一天的总价格是756元。
*/
func TestStop(t *testing.T) {
	price := func(min int) int {
		if min <= 15 {
			return 0
		}
		if min <= 60 {
			return 10
		}

		min -= 60
		ext := min % 60
		res := min / 60
		if ext > 0 {
			res += 1
		}

		return 10 + 2 + res
	}

	total := 0
	pre := 0
	for i := 0; i <= 24*60; i++ {
		consume := price(i)
		if pre != consume {
			total += consume
			if i < 60 {
				fmt.Printf("停车时长%dm 消费:%d\n", i, consume)
			} else {
				fmt.Printf("停车时长%dh 消费:%d\n", i/60, consume)
			}
			pre = consume
		}
	}
	fmt.Println(total)
}
