package game

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestOneReplay(t *testing.T) {
	r := ParseOsuReplay("test.osr")
	var time int64
	for _, rd := range r.ReplayData {
		time += rd.W
		fmt.Printf("%d: %+v\n", time, rd)
	}
}

func TestAllReplay(t *testing.T) {
	const rDir = "../test/Replays/"
	rs, err := ioutil.ReadDir(rDir)
	if err != nil {
		panic(err)
	}
	for _, rp := range rs {
		func() {
			path := rDir + rp.Name()
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(path)
					fmt.Println(err)
				}
			}()
			_ = ParseOsuReplay(path)
		}()
	}
}
