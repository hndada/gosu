package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/hndada/gosu/game/piano"
)

func main() {
	mpath := `C:\Users\Muang\AppData\Local\osu!\Songs\237212 Paitan - LEMON SUMMER`
	f := os.DirFS(mpath)

	entries, err := fs.ReadDir(f, ".")
	if err != nil {
		panic(err)
	}

	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".osu" {
			fmt.Println("Found:", e.Name())
			c, err := piano.NewChart(f, e.Name(), piano.Mods{})
			if err != nil {
				panic(err)
			}
			lv := c.StandardLevel()
			fmt.Printf("Level: %2.1f\n\n", lv)
		}
	}

}
