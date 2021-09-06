package audio

func SEPlayer(path string) func() {
	ap := NewPlayer(path)
	// ap.SetVolume(common.Settings.MasterVolume * common.Settings.SFXVolume)
	ap.SetVolume(0.25) // temp
	return func() {
		ap.Play()
		ap.Rewind()
	}
}
