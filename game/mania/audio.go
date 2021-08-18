package mania

import "github.com/hndada/gosu/game"

func SEPlayer() func() {
	var path string
	skinPath := `E:\gosu\skin\`
	fname := "soft-slidertick.wav"
	path = skinPath + fname
	ap := game.NewAudioPlayer(path)
	return func() {
		ap.Play()
		ap.Rewind()
	}
}
