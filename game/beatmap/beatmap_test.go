package beatmap

import (
	"fmt"
	"github.com/hndada/gosu/tools"
	"log"
	"testing"
)

func TestParseBeatmap(t *testing.T) {
	b, err := ParseBeatmap("../test_beatmap/test.osu")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(b.HitObjects), len(b.TimingPoints), b.General, b.Metadata)
	// fmt.Printf("%+v, %+v\n", b.General, b.Metadata)
}

func BenchmarkParseBeatmap(b *testing.B) {
	for i:=0; i< b.N; i++ {
		_, err := ParseBeatmap("../test_beatmap/test.osu")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func BenchmarkParseBeatmapParellel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := ParseBeatmap("../test_beatmap/test.osu")
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
	var beatmaps []Beatmap
	fmt.Println(songs)
	for _, song := range songs {
		b, err:= ParseBeatmap(song)
		if err!=nil { log.Fatal(err) }
		beatmaps=append(beatmaps, b)
	}
	for _, b := range beatmaps {
		fmt.Println(b.Metadata)
	}
}


