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
	// ap.SetVolume(game.Settings.MasterVolume * game.Settings.SFXVolume)
	ap.SetVolume(0.25) // temp
	return func() {
		ap.Play()
		ap.Rewind()
	}
}
