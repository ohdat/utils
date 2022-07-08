package utils

import (
	"fmt"
	"testing"
)

func TestStructAssign(t *testing.T) {
	//fmt.Println("price", price)
	type S1 struct {
		SSS string
		DDD string
	}

	type S2 struct {
		SSS int
		DDD string
	}

	var s1 = &S1{"123", "123"}
	var s2 = new(S2)

	StructAssign(s2, s1)
	fmt.Println("s2:", s2)
}

func TestStrToArrInt(t *testing.T) {
	arr, err := Str2ArrInt("1,2,3,4,5")
	fmt.Println("arr:", arr, "err:", err)
}
func TestDuplicationArrInt(t *testing.T) {
	arr := []int{1, 5, 2, 5, 2, 2, 52, 3, 4, 5}
	arr = DuplicationArrInt(arr)
	fmt.Println("arr:", arr)
}
func TestContainsInt(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	fmt.Println("ContainsInt:", ContainsInt(arr, 3))
}
