package play

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/draws"
)

func LoadSkin(fsys fs.FS, base Skin) {
	for i := 0; i < 10; i++ {
		s := draws.NewSprite(fsys, fmt.Sprintf("score/%d.png", i))
		s.ApplyScale(settings.ScoreScale)
		// Need to set same base line, since each number has different height.
		if i == 0 {
			s.Locate(ScreenSizeX, 0, draws.RightTop)
		} else {
			s.Locate(ScreenSizeX, skin.Score[0].H()-s.H(), draws.RightTop)
		}
		skin.Score[i] = s
	}
	for i, name := range []string{"dot", "comma", "percent"} {
		s := draws.NewSprite(fsys, fmt.Sprintf("score/%s.png", name))
		s.ApplyScale(settings.ScoreScale)
		s.Locate(ScreenSizeX, 0, draws.RightTop)
		skin.Sign[i] = s
	}
}
