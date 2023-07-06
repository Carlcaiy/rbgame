package test

import (
	"fmt"
	"reflect"
	"testing"
)

type School struct {
	Math   int
	Chines int
}

type Info struct {
	Name  string
	Age   int
	Is    bool
	Func  func(a int, b int) int
	Grade *School
}

func TestDeepcopy(t *testing.T) {
	cc := &Info{
		Name: "caiyunfeng",
		Age:  33,
		Is:   false,
		Grade: &School{
			Math:   100,
			Chines: 200,
		},
		Func: func(a int, b int) int {
			return a + b
		},
	}

	tc := reflect.TypeOf(cc)
	if tc.Kind() == reflect.Pointer {
		fmt.Println("pointer")
		ele := tc.Elem()
		fmt.Println("num field", ele.NumField())
		for i := 0; i < ele.NumField(); i++ {
			field := ele.Field(i)
			fmt.Println(field.Name, field.Type, field.Offset, field.Index, field.Anonymous)
		}
	} else {
		fmt.Println("value")
		fmt.Println("num field", tc.NumField())
		for i := 0; i < tc.NumField(); i++ {
			fmt.Println(tc.Field(i))
		}
	}

	vc := reflect.ValueOf(cc)
	fmt.Println("Can Addr", vc.CanAddr(), "Can Set", vc.CanSet())
	if vc.CanSet() {
		for i := 0; i < vc.NumField(); i++ {
			fmt.Println(vc.Field(i))
		}
	} else {
		vvc := vc.Elem()
		for i := 0; i < vvc.NumField(); i++ {
			fmt.Println(vvc.Field(i), vvc.Field(i).Type(), vvc.Field(i).CanAddr(), vvc.Field(i).CanSet())
		}
	}

	fmt.Println("--------------reflect of func--------------")
	rf := reflect.ValueOf(cc.Func)
	if rf.Kind() == reflect.Func {
		fmt.Println("is func")
		ret := rf.Call([]reflect.Value{reflect.ValueOf(2), reflect.ValueOf(10)})
		for i := range ret {
			if ret[i].CanInt() {
				fmt.Println(ret[i].Int())
			} else if ret[i].CanFloat() {
				fmt.Println(ret[i].Float())
			} else if ret[i].CanUint() {
				fmt.Println(ret[i].Uint())
			} else {
				fmt.Println(ret[i])
			}
		}
	} else {
		fmt.Println("not func")
	}

	bb := new(Info)
	copy(bb, cc)
	fmt.Println("bb", bb, bb.Grade)
	cc.Grade.Math = 1011
	fmt.Println("bb", bb, bb.Grade)
	fmt.Println("cc", cc, cc.Grade)
}
func copy(dst interface{}, src interface{}) {
	deepcopy(reflect.ValueOf(dst), reflect.ValueOf(src))
}

func deepcopy(v1 reflect.Value, v2 reflect.Value) {
	if v1.Kind() != v2.Kind() {
		fmt.Println("kind not same")
		return
	}

	if v1.Kind() == reflect.Pointer {
		v1 = v1.Elem()
		v2 = v2.Elem()
	}

	if v1.NumField() != v2.NumField() {
		fmt.Println("field num not same")
		return
	}
	for i := 0; i < v1.NumField(); i++ {
		if v1.Field(i).Kind() == v2.Field(i).Kind() {
			if v1.Field(i).Kind() == reflect.Pointer && v1.Field(i).IsNil() && !v2.Field(i).IsNil() {
				v1.Field(i).Set(reflect.New(v2.Field(i).Elem().Type()))
				deepcopy(v1.Field(i).Elem(), v2.Field(i).Elem())
			} else if v1.Field(i).CanSet() && v1.Field(i).CanAddr() {
				v1.Field(i).Set(v2.Field(i))
			} else {
				fmt.Println(i, "can set not", "can addr not")
			}
		}
	}
}
