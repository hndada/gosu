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
	for _, name := range []string{"circles(7k).osr"} {
		path := filepath.Join("testdata", name)
		dat, err := os.ReadFile(path)
		if err != nil {
			log.Fatal(err)
		}

		r, err := NewFormat(dat)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s's replay. The score is %d\n", r.PlayerName, r.Score)
		var time int64
		for _, rd := range r.ReplayData[:500] {
			time += rd.W
			// fmt.Printf("%d: %+v\n", time, rd)
			fmt.Printf("%+v\n", rd)
		}
	}
}

func TestMD5(t *testing.T) {
	for _, name := range []string{"4k.osr", "7k.osr", "4k_nc.osr"} {
		dat, err := os.ReadFile(name)
		if err != nil {
			log.Fatal(err)
		}

		r, err := NewFormat(dat)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(r.BeatmapMD5)
		// fmt.Println(r.MD5())
	}
}
