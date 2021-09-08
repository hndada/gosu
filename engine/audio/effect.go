package audio

import "github.com/faiface/beep/speaker"

func NewSEPlayer(path string, vol int) func() {
	ap := NewPlayer(path)
	ap.SetVolume(float64(vol) / 100)
	return func() {
		ap.Play()
		ap.Rewind()
	}
}

func NewSEPlayer2(path string) func() {
	ap := NewPlayer2(path)
	return func() {
		speaker.Play(ap)
	}
}
