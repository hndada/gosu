package mania

import "github.com/hndada/gosu/game"

func SEPlayer() func() {
	var path string
	skinPath := `E:\gosu\Skin\`
	fname := "soft-slidertick.wav"
	path = skinPath + fname
	ap := game.NewAudioPlayer(path, 16) // 매번 새로 load?
	return func() {
		ap.Play()
		ap.Rewind()
	}
}
