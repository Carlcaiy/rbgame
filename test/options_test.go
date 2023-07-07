package test

import (
	"fmt"
	"testing"
	"time"
)

type options struct {
	network   string
	addr      string
	readTime  time.Duration
	writeTime time.Duration
	userName  string
	password  string
}

type optfunc func(o *options)

func Network(n string) optfunc {
	return func(o *options) {
		o.network = n
	}
}

func Addr(s string) optfunc {
	return func(o *options) {
		o.addr = s
	}
}

func ReadTime(t time.Duration) optfunc {
	return func(o *options) {
		o.readTime = t
	}
}

func WriteTime(t time.Duration) optfunc {
	return func(o *options) {
		o.writeTime = t
	}
}

func UserName(s string) optfunc {
	return func(o *options) {
		o.userName = s
	}
}

func Password(s string) optfunc {
	return func(o *options) {
		o.password = s
	}
}

func StartOption(sli ...optfunc) {
	opt := new(options)
	for _, f := range sli {
		f(opt)
	}
	fmt.Printf("%+v\n", opt)
}

func TestOptions(t *testing.T) {
	StartOption(UserName("caiyunfeng"), Password("123123"), WriteTime(time.Second), ReadTime(time.Second), Network("tcp"), Addr("123.1.1.22:8860"))
}
