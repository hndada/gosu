package audio

func NewSEPlayer(path string, vol int) func() {
	ap := NewPlayer(path)
	ap.SetVolume(float64(vol) / 100)
	return func() {
		ap.Play()
		ap.Rewind()
	}
}
