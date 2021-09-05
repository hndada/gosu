package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hndada/gosu/common"
	"github.com/hndada/gosu/mania"
)

func main() {
	const instability = 15

	path := `E:\gosu\Music\932851 Kurokotei - wtf\Kurokotei - wtf (FAMoss) [easy].osu`
	switch strings.ToLower(filepath.Ext(path)) {
	case ".osu":
		switch common.OsuMode(path) {
		case common.ModeMania:
			c, err := mania.NewChart(path)
			if err != nil {
				panic(err)
			}
			autoGen := c.GenAutoKeyEvents(instability)
			for t := 0; t < int(c.EndTime()); t += 60 {
				fmt.Println(autoGen(int64(t)))
			}
		default:
			panic("not reach")
		}
	}
}
