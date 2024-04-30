//go:build !js || !wasm

package gosu

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"

	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/scene"
	"github.com/hndada/gosu/util"
)

func (Game) createOptionsFile(fname string) {
	options := scene.NewOptions()
	data, err := json.Marshal(options)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(fname, data, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Generated %s with default values.\n", fname)
}

func (g *Game) loadOptions() {
	const fname = "options.json"

	// Try to create the file if it doesn't exist.
	_, err := os.Stat(fname)
	if os.IsNotExist(err) {
		g.createOptionsFile(fname)
	}

	// Read the file.
	data, err := os.ReadFile(fname)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, g.options)
	if err != nil {
		panic(err)
	}

	// g.options.Normalize()
}

// load functions should not load the entire file system into memory.
// Instead, they should load the path to the file and read the file when needed.

// 1. Parse the whole file system.
// 2. Save header data to map with file path as a key.
func loadReplays(fsys fs.FS) map[string]game.Replay {
	replays := make(map[string]game.Replay)

	fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		r, err := game.NewReplay(fsys, path, 4)
		replays[util.MD5(dat)] = r
		return nil
	})
	return replays
}
