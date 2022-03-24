package utils

import (
	"fmt"
	"testing"
)

func TestStructAssign(t *testing.T) {
	//fmt.Println("price", price)
	type S1 struct {
		SSS string
	}

	type S2 struct {
		SSS int
	}

	var s1 = &S1{"123"}
	var s2 = new(S2)

	StructAssign(s2, s1)
	fmt.Println("s2:", s2)
}
