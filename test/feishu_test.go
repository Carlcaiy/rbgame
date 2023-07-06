package test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func sendMsg(apiUrl, msg string) {
	// json
	contentType := "application/json"
	// data
	sendData := `{
		"msg_type": "text",
		"content": {"text": "` + "消息通知:" + msg + `"}
	}`
	// request
	result, err := http.Post(apiUrl, contentType, strings.NewReader(sendData))
	if err != nil {
		fmt.Printf("post failed, err:%v\n", err)
		return
	}
	defer result.Body.Close()

}
func TestFeishu(t *testing.T) {
	var v = decimal.NewFromInt(1000)
	v.Div(decimal.NewFromInt(33))
	fmt.Println(v.Div(decimal.NewFromInt(33)))

	stamp := time.Now().Unix()
	fmt.Println(stamp, time.Now().Day(), time.Unix(stamp, 1231231).Unix())

	var num int64 = 1000000
	sum := decimal.NewFromInt(num).DivRound(decimal.NewFromFloat(100), 2)
	fmt.Println(sum)

	sendMsg("https://open.feishu.cn/open-apis/bot/v2/hook/ae447214-c3c6-4187-ac3f-1868c0349599", "xxxxxxxxx")

	m := make(map[int]int)
	m[1] = 100
	m[2] = 200
	m[3] = 300
	for i, c := range "加大大书法水电费沙发上服阿斯顿发sdf" {
		fmt.Printf("%d,%U \n", i, c)
	}
	fmt.Println(Unhex("123GoodStudy"))

	src := []int{1, 2, 3, 4, 5, 6, 7, 78, 8}
	size := 1
Loop:
	for n := 0; n < len(src); n += size {
		switch {
		case src[n] < 5:
			break Loop
			// if src[n] < 5 {
			// 	break
			// }
			// fmt.Println("1", src[n])
		case src[n] < 8:
			fmt.Println("2", src[n])
			break Loop
		}
	}
	fmt.Println("dadasd")

	var _type interface{} = 10
	switch t := _type.(type) {
	default:
		fmt.Printf("unexpected type %T\n", t) // %T prints whatever type t has
	case bool:
		fmt.Printf("boolean %t\n", t) // t has type bool
	case int:
		fmt.Printf("integer %d\n", t) // t has type int
	case *bool:
		fmt.Printf("pointer to boolean %t\n", *t) // t has type *bool
	case *int:
		fmt.Printf("pointer to integer %d\n", *t) // t has type *int
	}
}

// if b= [0,9] get [0,9]
// if b= [a,z] get [a,z]
// if b= [A,Z] get [a,z]
func unhex(b byte) byte {
	if b >= '0' && b <= '9' {
		return b
	}
	if b >= 'a' && b <= 'z' {
		return b
	}
	if b >= 'A' && b <= 'Z' {
		return ('a' - 'A') + b
	}
	return 0
}

func Unhex(s string) string {
	sli := []byte(s)
	for i, c := range sli {
		sli[i] = unhex(c)
	}
	return string(sli)
}
