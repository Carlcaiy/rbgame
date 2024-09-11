package test

import (
	"fmt"
	"strconv"
	"testing"
)

const (
	_  = iota             // iota=0
	KB = 1 << (10 * iota) // iota=1
	MB = 1 << (10 * iota) // iota=2
	GB = 1 << (10 * iota) // iota=3
	TB = 1 << (10 * iota) // iota=4
)

func TestIotaB(t *testing.T) {
	t.Logf("11:%d\n", 1<<2)
	t.Logf("KB=%d MB=%d GB=%d TB=%d\n", KB, MB, GB, TB)
}

const (
	FlagNone  = 0
	FlagRead  = 1
	FlagWrite = 2
	FlagExec  = 1 << iota // iota=3
)

func TestFlag(t *testing.T) {
	t.Logf("Read=%d Write=%d Exec=%d\n", FlagRead, FlagWrite, FlagExec)
}

func TestParseF(t *testing.T) {
	str := "1.234234"
	f, err := strconv.ParseFloat(str, 64)
	fmt.Println(f, err)
	p := int64(f * 100)
	fmt.Println(p)
}
