package scene

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
)

// asset.go: functions that load assets from fs.FS.
// draw.go: functions that draw.
func NewBackgroundDrawer(cfg *Config, asset *Asset, fsys fs.FS, name string) func(draws.Image) {
	s := draws.NewSpriteFromFile(fsys, name)
	if s.IsEmpty() {
		s = asset.DefaultBackgroundSprite
	}
	s.MultiplyScale(cfg.ScreenSize.X / s.W())
	s.Locate(cfg.ScreenSize.X/2, cfg.ScreenSize.Y/2, draws.CenterMiddle)

	return func(dst draws.Image) {
		op := draws.Op{}
		value := cfg.BackgroundBrightness
		op.ColorM.ChangeHSV(0, 1, value)
		s.Draw(dst, op)
	}
}
