package osr

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

// Observation:
// 1. There might be some more ways to check whether the replay is auto or not,
// but they varies among modes.
// 2. There is no difference in replays whether a player skipped the intro or not.
// 3. Auto always hit precisely, without 1ms error on every notes.
func TestNewFormat(t *testing.T) {
	for _, name := range []string{"4k_auto.osr"} {
		path := filepath.Join("testdata", name)
		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}

		r, err := NewFormat(f)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s's replay. The score is %d\n", r.PlayerName, r.Score)
		var time int64
		for _, rd := range r.ReplayData[len(r.ReplayData)-20:] {
			time += rd.W
			// fmt.Printf("%d: %+v\n", time, rd)
			fmt.Printf("%+v\n", rd)
		}
	}
}

func TestMD5(t *testing.T) {
	for _, name := range []string{"4k.osr", "7k.osr", "4k_nc.osr"} {
		f, err := os.Open(name)
		if err != nil {
			log.Fatal(err)
		}

		r, err := NewFormat(f)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(r.BeatmapMD5)
		fmt.Println(r.MD5())
	}
}
