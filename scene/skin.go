package scene

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/draws"
)

// Skin is a set of Sprites and sounds.
type Skin struct {
	DefaultBackground draws.Sprite
	Cursor            [3]draws.Sprite
	Score             [10]draws.Sprite
	Sign              [3]draws.Sprite
	// ChartItemBoxSprite draws.Sprite
	// ChartLevelBoxSprite draws.Sprite
	Sounds
}

var (
	defaultSkin *Skin
	skin        *Skin
)

func CurrentSkin() *Skin { return skin }
func (s *Skin) Load(fsys fs.FS, base *Skin) {
	skin.DefaultBackground = NewBackground(fsys, "default-bg.jpg")
	names := []string{"menu-cursor", "menu-cursor-additive", "cursortrail"}
	for i, name := range names {
		s := draws.NewSprite(fsys, fmt.Sprintf("cursor/%s.png", name))
		s.ApplyScale(settings.CursorScale)
		s.Locate(ScreenSizeX/2, ScreenSizeY/2, draws.CenterMiddle)
		skin.Cursor[i] = s
	}
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
	s.Sounds.Load(fsys, &base.Sounds)
}

// {
// 	s := draws.NewSprite(fsys, "box-mask.png")
// 	s.SetSize(ChartInfoBoxWidth, ChartInfoBoxHeight)
// 	s.Locate(ScreenSizeX+chartInfoBoxshrink, ScreenSizeY/2, draws.RightMiddle)
// 	ChartItemBoxSprite = s
// }
