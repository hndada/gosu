package tools

import (
	"log"
)

func CheckArgsLen(args []string, l int) {
	if len(args) < l {
		log.Fatalf("not enough args; needs %d", l)
	} else if len(args) > l {
		log.Fatal("too much args; needs %d", l)
	}
}

func CheckValidMode(mode int) {
	keys := make([]int, 0, 4)
	for k := range ModePrefix {
		if mode == k {
			return
		}
		keys = append(keys, k)
	}
	log.Fatalf("mode should be one %v", keys)
}
