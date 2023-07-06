package main

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type TroopObj struct {
	Delay  int     `yaml:"delay" json:"delay"`
	Speed  float32 `yaml:"speed" json:"speed"`
	Space  int     `yaml:"apace" json:"apace"`
	Offset int32   `yaml:"offset" json:"offset"`
	RoadId uint32  `yaml:"road" json:"road"`
	FishId int32   `yaml:"fish" json:"fish"`
	Count  int     `yaml:"num" json:"num"`
}

func main() {
	obj := make([]TroopObj, 0)
	str := []byte(`[{offset: 11,road: 608,fish: 11,delay:0,speed:0.105,apace:550,num:10},{offset:0,road:608,fish:11,delay:0,speed:0.105,apace:550,num:10}]`)
	if err := yaml.Unmarshal(str, &obj); err != nil {
		panic(err)
	}
	for i := range obj {
		fmt.Println(obj[i])
	}
}
