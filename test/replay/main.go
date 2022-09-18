package main

import (
	"fmt"
	"os"

	"github.com/hndada/gosu/format/osr"
)

// - Soleily - Renatus [don DON] (2022-09-16) Taiko.osr
// Idle: {W:13 X:320 Y:9999 Z:0}
// Left don: {W:16 X:0 Y:9999 Z:1}
// Right don: {W:15 X:640 Y:9999 Z:20}
// Left kat: {W:12 X:0 Y:9999 Z:2}
// Right kat: {W:3 X:640 Y:9999 Z:8}

// Z value for [K, D, D, K]: [2, 1, 4+16, 8]
// X = 320 when at idle. X = 640 when only right hand is hitting.
// X = 0 when left hand or both hands are hitting.
func main() {
	rd, err := os.ReadFile("test.osr")
	if err != nil {
		panic(err)
	}
	f, err := osr.Parse(rd)
	if err != nil {
		panic(err)
	}
	for i, data := range f.ReplayData[:15] {
		fmt.Printf("%d: %+v\n", i, data)
	}
}
