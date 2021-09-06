package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/hndada/gosu/mania"
)

func main() {
	const dirname = "./bms"
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		s := filepath.Join(dirname, f.Name())
		c, _ := mania.NewChart(s)
		c.CalcDifficulty()
		fmt.Printf("%s %1.2f: %s\n", getManualLV(f.Name()), c.Level/18, f.Name())
	}
}

func getManualLV(name string) string {
	strs := strings.Split(name, " ")
	last := strs[len(strs)-1]
	return last[0 : len(last)-5]
}
