package osr

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func TestParse(t *testing.T) {
	dat, err := ioutil.ReadFile("test.osr")
	if err != nil {
		log.Fatal(err)
	}
	r, err := Parse(dat)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s's replay. The score is %d\n", r.PlayerName, r.Score)
	var time int64
	for _, rd := range r.ReplayData {
		time += rd.W
		fmt.Printf("%d: %+v\n", time, rd)
	}
}
