package test

import "testing"

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
