package tests

type Xxxx struct {
	Num int
	Val int
}

func (x *Xxxx) GetNum() int {
	return x.Num
}

func (x *Xxxx) GetVal() int {
	return x.Val
}

type Ifxx interface {
	GetNum() int
	GetVal() int
}
