package osu

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

// Todo: map my own charts
var tests = []string{"testo.osu", "testt.osu", "testc.osu", "testm.osu"}

func TestParse(t *testing.T) {
	for _, s := range tests {
		dat, err := ioutil.ReadFile(s)
		if err != nil {
			log.Fatal(err)
		}
		o, err := Parse(dat)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%+v\n%+v\n", o.General, o.Metadata)
	}
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dat, err := ioutil.ReadFile("testm.osu")
		if err != nil {
			log.Fatal(err)
		}
		_, err = Parse(dat)
		if err != nil {
			log.Fatal(err)
		}
	}
}
