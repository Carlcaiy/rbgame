package main

import "fmt"

// go build -buildmode=plugin -o myplugin-v1.19-3.so plugin.go

func init() {
	fmt.Println("version 1.2")
}

func MyPrin(a, b int) int {
	fmt.Println("MyPrin+")
	return a + b
}
