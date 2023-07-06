//+build !debug

package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type mstring = string
type tstring string

func test() {
	fmt.Println("REALSE MODE")

	strs := strings.Split("100000-200000", "-")
	for _, str := range strs {
		strKm, err := strconv.Atoi(str)
		if err != nil {
			panic(err)
		}
		fmt.Println(strKm)
	}
}
func conf() error {
	// 使用匿名结构体，防止污染项目内结构体
	mc := struct {
		MakeCards map[int][]int `yaml:"make_cards"`
	}{
		MakeCards: make(map[int][]int),
	}

	{
		file, err := ioutil.ReadFile("make_cards.yaml")
		if err != nil {
			fmt.Println("ReadFile", err)
		}
		err = yaml.Unmarshal(file, mc)
		if err != nil {
			fmt.Println("Unmarshal", err)
		}
		fmt.Println("ucards", mc)
	}
	return nil
}
