package ui

import (
	"github.com/hndada/gosu/input"
)

type ControlType int

const (
	Decrease ControlType = iota
	Increase
)

const (
	Toggle ControlType = iota
)

type Control struct {
	Key           input.Key
	Type          ControlType
	SoundFilename string
}

// const decreaseSoundFilename = "interface/sound/toggle/off.wav"
// const increaseSoundFilename = "interface/sound/toggle/on.wav"

// var DecreaseControl = Control{
// 	Key:           input.KeyArrowLeft,
// 	Type:          Decrease,
// 	SoundFilename: decreaseSoundFilename,
// }
// var IncreaseControl = Control{
// 	Key:           input.KeyArrowRight,
// 	Type:          Increase,
// 	SoundFilename: increaseSoundFilename,
// }

// var DefaultControls = []Control{DecreaseControl, IncreaseControl}
