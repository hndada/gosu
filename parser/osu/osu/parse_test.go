package osu

import (
	"fmt"
	"log"
	"testing"
)

func TestParse(t *testing.T) {
	for _, s := range []string{"testo.osu", "testt.osu", "testc.osu", "testm.osu"} {
		o, err := Parse(s)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%+v", o)
	}
}
