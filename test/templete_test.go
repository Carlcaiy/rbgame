package test

import (
	"bytes"
	"fmt"
	"testing"
	"text/template"
)

type Templete struct {
	Name string
	Age  int
	Sex  int
}

func TestTemplete(t *testing.T) {
	temp := template.New("yes")
	temp.Parse("名字是{{.Name}}年龄{{.Age}}性别{{.Sex}}")
	buf := bytes.NewBuffer(make([]byte, 0))
	temp.Execute(buf, &Templete{
		Name: "dada",
		Age:  11,
		Sex:  11,
	})
	fmt.Println(buf.String())
}
