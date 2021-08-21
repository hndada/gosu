package mania

import (
	"path/filepath"

	"github.com/hndada/gosu/game"
)

func SEPlayer(cwd string) func() {
	dir := filepath.Join(cwd, "skin")
	name := "soft-slidertick.wav"
	path := filepath.Join(dir, name)
	ap := game.NewAudioPlayer(path)
	return func() {
		ap.Play()
		ap.Rewind()
	}
}
