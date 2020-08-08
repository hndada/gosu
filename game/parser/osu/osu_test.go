package osu

import (
	"fmt"
	"github.com/hndada/gosu/tools"
	"log"
	"testing"
)

func TestNewOSU(t *testing.T) {
	b, err := NewOSU("../test_beatmap/test.osu")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(b.HitObjects), len(b.TimingPoints), b.General, b.Metadata)
	// fmt.Printf("%+v, %+v\n", b.General, b.Metadata)
}

func BenchmarkNewOSU(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := NewOSU("../test_beatmap/test.osu")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func BenchmarkNewOSUParellel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := NewOSU("../test_beatmap/test.osu")
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}

func TestLoadSongList(t *testing.T) {
	songs, err := tools.LoadSongList("../test_beatmap/", ".osu")
	if err != nil {
		log.Fatal(err)
	}
	var beatmaps []OSU
	fmt.Println(songs)
	for _, song := range songs {
		b, err := NewOSU(song)
		if err != nil {
			log.Fatal(err)
		}
		beatmaps = append(beatmaps, b)
	}
	for _, b := range beatmaps {
		fmt.Println(b.Metadata)
	}
}
