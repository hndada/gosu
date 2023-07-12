package audios

// SoundPlayer is an interface for playing sound.
// It is implemented by Sound and SoundPod.
type SoundPlayer interface {
	Play(vol float64)
}
