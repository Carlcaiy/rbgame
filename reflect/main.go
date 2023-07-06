package main

import (
	"fmt"
	"rbgame/reflect/tests"
	"reflect"
)


func main() {
	iface := &tests.Xxxx{
		Num: 100,
		Val: 100,
	}
	value := reflect.ValueOf(iface)
	ttype := reflect.TypeOf(iface)
	v := reflect.Indirect(value)
	fmt.Println(value.CanAddr(), v.CanAddr())
	fmt.Println(v.NumField()) 
	for i := 0; i < 2; i++ {
		fmt.Println(v.Field(i).Int())
	}

	fmt.Println(value.MethodByName("GetNum").Call(nil)[0].Int()) 
	for i := 0; i < value.NumMethod(); i++ {
		fmt.Println(value.Method(i), ttype.Method(i).Name, ttype.Method(i).Type.PkgPath()+"--") 
	}
}

func Add[V int](a V, b V) V {
	return a + b
}

func AA[V int](b V) V {
	return b
}
