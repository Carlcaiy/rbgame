package composite

type BaseI interface {
	GetName() string
}

type Base struct {
}

func (p *Base) GetName() string {
	return "base"
}

type Leaf struct {
	Base
}

var sli []BaseI

func xxx() {
	sli = append(sli, &Leaf{})
}
