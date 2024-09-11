package test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/vmihailenco/msgpack"
)

type msgJson struct {
	Name string `json:"name"`
	Age  int    `json:"Age"`
}

func TestMsgPack(t *testing.T) {
	man := msgJson{"xxx", 20}
	bs, err := msgpack.Marshal(man)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(bs))
	js, ksErr := json.Marshal(man)
	if ksErr != nil {
		panic(ksErr)
	}
	fmt.Println(len(js))
}
