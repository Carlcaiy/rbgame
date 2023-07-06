package main

import (
	"fmt"
	"testing"
)

// go类型别名 类型定义
type NewInt = int
type InitAlias int

type new_name_string = string
type new_type_string string

//go test -timeout 30s -run ^Test$ rbgame/practise
func Test(t *testing.T) {
	var nns new_name_string = "nnstr"
	var nnt new_type_string = "ntstr"
	if _, ok := interface{}(nns).(string); ok {
		t.Log("type new_name_string=string 类型相等")
	} else {
		t.Log("type new_name_string!=string 类型不相等")
	}

	if _, ok := interface{}(nnt).(string); ok {
		t.Log("type new_type_string=string 类型相等")
	} else {
		t.Log("type new_type_string!=string 类型不相等")
	}
}

// 类型别名，可以相同类型直接进行计算，自定义类型不同
// 自定义类型可以为类型添加方法

func TestType(t *testing.T) {
	var it NewInt = 100
	var val int = 100
	var nit InitAlias = 100
	if it == val {
		fmt.Println("it == nit")
	}
	if it == int(nit) {
		fmt.Println("需要进行强转")
	}
}
