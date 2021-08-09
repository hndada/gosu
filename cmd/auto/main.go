package main

import (
	"fmt"

	"github.com/hndada/gosu/game/mania"
	"github.com/hndada/rg-parser/osugame/osu"
)

func main() {
	// fpath := `/home/hndadada/gosu/cmd/gosu/Music/665131 volthi - beyond the clouds/volthi - beyond the clouds (Hydria) [beginner].osu`
	// fpath := `E:\gosu\Music\665131 volthi - beyond the clouds\volthi - beyond the clouds (Hydria) [beginner].osu`
	fpath := `E:\gosu\Music\932851 Kurokotei - wtf\Kurokotei - wtf (FAMoss) [easy].osu`
	o, err := osu.Parse(fpath)
	if err != nil {
		panic(err) // todo: log and continue
	}
	c, err := mania.NewChartFromOsu(o, fpath)
	if err != nil {
		panic(err) // todo: log and continue
	}
	autoGen := c.GenAutoKeyEvents()
	for t := 0; t < int(c.EndTime()); t += 60 {
		fmt.Println(autoGen(int64(t)))
	}
}
