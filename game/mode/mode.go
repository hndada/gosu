package mode

import (
	"github.com/hndada/gosu/framework/input"
)

const (
	ModeNone = iota - 1
	ModePiano
	ModeDrum
	ModeSing
)

// Prop stands for Mode properties.
type Prop struct {
	LoadSkin   func()
	SpeedScale *float64
	Settings   map[string]*float64
	// NewScenePlay func(cpath string, rf *osr.Format) (scene.Scene, error)
	ExposureTime func(float64) float64
	KeySettings  map[int][]input.Key
}
