package scene

import (
	"github.com/hndada/gosu/draws"
)

// asset.go: functions that load assets from fs.FS.
// draw.go: functions that draw.
func NewDrawBackgroundFunc(s draws.Sprite,
	screenSize draws.Vector2, bgBrightness *float64) func(draws.Image) {

	return func(dst draws.Image) {
		op := draws.Op{}
		value := *bgBrightness
		op.ColorM.ChangeHSV(0, 1, value)
		s.Draw(dst, op)
	}
}
