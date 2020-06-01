package beatmap

import (
	"fmt"
	"log"
	"testing"
)

func TestParseBeatmap(t *testing.T) {
	beatmap, err := ParseBeatmap("../../test/test.osu")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(beatmap.HitObjects), len(beatmap.TimingPoints), beatmap.General, beatmap.Metadata)
	// fmt.Printf("%+v, %+v\n", beatmap.General, beatmap.Metadata)
}

func BenchmarkParseBeatmap(b *testing.B) {
	for i:=0; i< b.N; i++ {
		_, err := ParseBeatmap("../../test/test.osu")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func BenchmarkParseBeatmapParellel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := ParseBeatmap("../../test/test.osu")
			if err != nil {
				log.Fatal(err)
			}
		}
	})
}
