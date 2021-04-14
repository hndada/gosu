package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/hndada/gosu/game/mania"
	"github.com/hndada/rg-parser/osugame/osu"
)

func main() {
	const dirname = "./charts"
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		s := filepath.Join(dirname, f.Name())
		o, err := osu.Parse(s)
		if err != nil {
			log.Fatal(err)
		}
		c, _ := mania.NewChartFromOsu(o, s)
		c.CalcDifficulty()
		fmt.Printf("%1.2f: %s\n", c.Level/50, f.Name())
	}
}
