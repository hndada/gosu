package mode

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/draws"
)

const (
	ScoreDot = iota + 10
	ScoreComma
	ScorePercent
)

func NewScoreNumbers(fsys fs.FS, ScreenSize draws.Vector2, scale float64) [13]draws.Sprite {
	var sprites [13]draws.Sprite
	for i := 0; i < 10; i++ {
		s := draws.NewSpriteFromFile(fsys, fmt.Sprintf("score/%d.png", i))
		s.MultiplyScale(scale)
		// Score needs to set same base line,
		// since each number might have different height.
		if i == 0 {
			s.Locate(ScreenSize.X, 0, draws.RightTop)
		} else {
			s.Locate(ScreenSize.X, sprites[0].H()-s.H(), draws.RightTop)
		}
		sprites[i] = s
	}
	for i, name := range []string{"dot", "comma", "percent"} {
		s := draws.NewSpriteFromFile(fsys, fmt.Sprintf("score/%s.png", name))
		s.MultiplyScale(scale)
		s.Locate(ScreenSize.X, 0, draws.RightTop)
		sprites[i+10] = s
	}
	return sprites
}

// type Config interface {
// 	ScreenSize() draws.Vector2
// 	ScoreScale() float64
// }
// func NewScoreNumbers(fsys fs.FS, cfg Config) [13]draws.Sprite {
