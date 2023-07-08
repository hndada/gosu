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

// Skin has not only a set of Sprites, but also a set of sounds.
type SkinType struct {
	DefaultBackground draws.Sprite
	Score             [13]draws.Sprite
}

var Skin SkinType

func (skin *SkinType) Load(fsys fs.FS) {
	skin.DefaultBackground = NewBackground(fsys, "interface/default-bg.jpg")
	for i := 0; i < 10; i++ {
		s := draws.LoadSprite(fsys, fmt.Sprintf("score/%d.png", i))
		s.MultiplyScale(Settings.ScoreScale)

		// Need to set same base line, since each number has different height.
		if i == 0 {
			s.Locate(ScreenSizeX, 0, draws.RightTop)
		} else {
			s.Locate(ScreenSizeX, skin.Score[0].H()-s.H(), draws.RightTop)
		}
		skin.Score[i] = s
	}
	for i, name := range []string{"dot", "comma", "percent"} {
		s := draws.LoadSprite(fsys, fmt.Sprintf("score/%s.png", name))
		s.MultiplyScale(Settings.ScoreScale)
		s.Locate(ScreenSizeX, 0, draws.RightTop)
		skin.Score[i+10] = s
	}
}
