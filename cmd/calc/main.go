package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/hndada/gosu/game/mania"
	"github.com/hndada/rg-parser/osugame/osu"
)

func main() {
	const dirname = "./bms"
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
		fmt.Printf("%s %1.2f: %s\n", getManualLV(f.Name()), c.Level/18, f.Name())
	}
}

func getManualLV(name string) string {
	strs:= strings.Split(name, " ")
	last:=strs[len(strs)-1]
	return last[0:len(last)-5]
}