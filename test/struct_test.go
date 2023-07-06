package test

import (
	"fmt"
	"testing"
)

type nils struct {
	Name string
}

func (n *nils) GetMame() string {
	if n == nil {
		return "dada"
	}
	return n.Name
}

func TestStructFunc(t *testing.T) {
	var x *nils
	fmt.Println(x.GetMame())
}
