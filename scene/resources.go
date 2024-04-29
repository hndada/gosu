package scene

import "github.com/hndada/gosu/game/piano"

type Resources struct {
	Piano *piano.Resources
}

// Todo: deal with two kinds of values: file path and directory path.
const (
	SoundToggleOff      = "interface/sound/toggle/off.wav"
	SoundToggleOn       = "interface/sound/toggle/on.wav"
	SoundTransitionDown = "interface/sound/transition/down.wav"
	SoundTransitionUp   = "interface/sound/transition/up.wav"
	SoundTaps           = "interface/sound/tap/"
	SoundSwipes         = "interface/sound/swipe/"
)
