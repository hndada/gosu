package mania

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestFinger(t *testing.T) {
	for keys := 0; keys <= 10; keys++ {
		s := make([]string, 0, keys)
		for k := 0; k < keys; k++ {
			s = append(s, strconv.Itoa(finger(keys, k)))
		}
		fmt.Printf("fingers[%d] = []int{%s}\n", keys, strings.Join(s, ", "))
	}
}

func TestScratchFinger(t *testing.T) {
	fs := make(map[int][]int)
	for keys := 1; keys <= 8; keys++ {
		fs[keys] = make([]int, keys)
		for k := 0; k < keys; k++ {
			fs[keys][k] = finger(keys, k)
		}
		fmt.Printf("%+v\n", fs[keys])
	}
	for keys := 2; keys <= 8; keys++ {
		fs[keys|ScratchLeft] = append([]int{fingers[keys-1][0] + 1}, fingers[keys-1]...)
		fmt.Printf("%+v\n", fs[keys|ScratchLeft])
		fs[keys|ScratchRight] = append(fingers[keys-1], fingers[keys-1][keys-2]+1)
		fmt.Printf("%+v\n", fs[keys|ScratchRight])
	}
	fmt.Println("fingers:")
	for keys := 2; keys <= 8; keys++ {
		fmt.Printf("%+v\n", fingers[keys|ScratchLeft])
		fmt.Printf("%+v\n", fingers[keys|ScratchRight])
	}
}
