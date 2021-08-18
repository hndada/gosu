package mania

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestFinger(t *testing.T) {
	for keyCount := 0; keyCount <= 10; keyCount++ {
		s := make([]string, 0, keyCount)
		for k := 0; k < keyCount; k++ {
			s = append(s, strconv.Itoa(finger(keyCount, k)))
		}
		fmt.Printf("fingers[%d] = []int{%s}\n", keyCount, strings.Join(s, ", "))
	}
}

func TestScratchFinger(t *testing.T) {
	fs := make(map[int][]int)
	for keyCount := 1; keyCount <= 8; keyCount++ {
		fs[keyCount] = make([]int, keyCount)
		for k := 0; k < keyCount; k++ {
			fs[keyCount][k] = finger(keyCount, k)
		}
		fmt.Printf("%+v\n", fs[keyCount])
	}
	for keyCount := 2; keyCount <= 8; keyCount++ {
		fs[keyCount|leftScratch] = append([]int{fingers[keyCount-1][0] + 1}, fingers[keyCount-1]...)
		fmt.Printf("%+v\n", fs[keyCount|leftScratch])
		fs[keyCount|rightScratch] = append(fingers[keyCount-1], fingers[keyCount-1][keyCount-2]+1)
		fmt.Printf("%+v\n", fs[keyCount|rightScratch])
	}
	fmt.Println("fingers:")
	for keyCount := 2; keyCount <= 8; keyCount++ {
		fmt.Printf("%+v\n", fingers[keyCount|leftScratch])
		fmt.Printf("%+v\n", fingers[keyCount|rightScratch])
	}
}
