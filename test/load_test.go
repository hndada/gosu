package test

import (
	"fmt"
	"github.com/hndada/gosu/game/beatmap"
	"github.com/hndada/gosu/game/tools"
	"log"
	"testing"
)

func TestLoadSongList(t *testing.T) {
	songs, err := tools.LoadSongList("../test_beatmap/", ".osu")
	if err != nil {
		log.Fatal(err)
	}
	var beatmaps []beatmap.Beatmap
	fmt.Println(songs)
	for _, song := range songs {
		b, err:=beatmap.ParseBeatmap(song)
		if err!=nil { log.Fatal(err) }
		beatmaps=append(beatmaps, b)
	}
	for _, b := range beatmaps {
		fmt.Println(b.Metadata)
	}
}
