//go:build !js || !wasm

package gosu

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hndada/gosu/scene"
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
