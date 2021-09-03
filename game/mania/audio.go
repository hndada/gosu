package mania

import (
	"path/filepath"

	"github.com/hndada/gosu/engine/audio"
)

func SEPlayer(cwd string) func() {
	dir := filepath.Join(cwd, "skin")
	name := "soft-slidertick.wav"
	path := filepath.Join(dir, name)
	ap := audio.NewPlayer(path)
	return func() {
		ap.Play()
		ap.Rewind()
	}
}
