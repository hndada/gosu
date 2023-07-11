package mode

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/draws"
)

// func NewScoreSprites(fsys fs.FS, screenSize draws.Vector2, scale float64) [13]draws.Sprite {
func NewScoreSprites(fsys fs.FS, cfg Config) [13]draws.Sprite {
	var sprites [13]draws.Sprite
	for i := 0; i < 10; i++ {
		s := draws.LoadSprite(fsys, fmt.Sprintf("score/%d.png", i))
		s.MultiplyScale(cfg.ScoreScale())
		// Score needs to set same base line,
		// since each number might have different height.
		if i == 0 {
			s.Locate(cfg.ScreenSize().X, 0, draws.RightTop)
		} else {
			s.Locate(cfg.ScreenSize().X, sprites[0].H()-s.H(), draws.RightTop)
		}
		sprites[i] = s
	}
	for i, name := range []string{"dot", "comma", "percent"} {
		s := draws.LoadSprite(fsys, fmt.Sprintf("score/%s.png", name))
		s.MultiplyScale(cfg.ScoreScale())
		s.Locate(cfg.ScreenSize().X, 0, draws.RightTop)
		sprites[i+10] = s
	}
	return sprites
}
