package test

import (
	"fmt"
	"testing"
)

func TestSlice(t *testing.T) {
	slice := make([]byte, 20)
	fmt.Println(len(slice), cap(slice))
	slice = slice[:1]
	fmt.Println(len(slice), cap(slice))
	part1 := slice[:4]
	part1[0] = 1
	part1[1] = 2
	part1[2] = 3
	part1[3] = 4
	fmt.Println(len(part1), cap(part1))
	part2 := slice[4:8]
	part2[0] = 11
	part2[1] = 21
	part2[2] = 31
	part2[3] = 41
	fmt.Println(len(part2), cap(part2))
	fmt.Println(slice)
	slice = slice[:cap(slice)]
	fmt.Println(slice)
	fmt.Println(len(slice), cap(slice))
}
