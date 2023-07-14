package osr

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestNewFormat(t *testing.T) {
	f, err := os.Open("test.osr")
	if err != nil {
		log.Fatal(err)
	}

	r, err := NewFormat(f)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s's replay. The score is %d\n", r.PlayerName, r.Score)
	var time int64
	for _, rd := range r.ReplayData[len(r.ReplayData)-100:] {
		time += rd.W
		fmt.Printf("%d: %+v\n", time, rd)
	}
}

func TestMD5(t *testing.T) {
	f, err := os.Open("test.osr")
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
