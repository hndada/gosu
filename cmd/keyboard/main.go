package main

import (
	"fmt"
	"time"

	"github.com/eiannone/keyboard"
)

func main() {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Println("Press ESC to quit")
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		t:=time.Now()
		fmt.Printf("You pressed: rune %q, key %X at time: %+v\r\n", char, key, t)
		if key == keyboard.KeyEsc {
			break
		}
	}
}
