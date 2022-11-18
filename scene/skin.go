package scene

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/mode"
)

type Skin struct {
	Kind mode.SkinKind
	// Number1 [13]draws.Sprite // number and sign(. , %)
	// Number2 [10]draws.Sprite // number only
	Cursor [3]draws.Sprite
}

const (
	CursorBase = iota
	CursorAdditive
	CursorTrail
)

var (
	DefaultSkin = Skin{Kind: mode.SkinKindDefault}
	UserSkin    = Skin{Kind: mode.SkinKindUser}
	// PlaySkin    = Skin{Kind: SkinKindPlay}
)

func (skin *Skin) Load(fsys fs.FS) {
	for i, name := range []string{"base", "additive", "trail"} {
		s := draws.NewSprite(fsys, fmt.Sprintf("cursor/%s.png", name))
		s.ApplyScale(Settings.CursorScale)
		// Cursor should be at CenterMiddle in circle mode (in far future)
		s.Locate(ScreenSizeX/2, ScreenSizeY/2, draws.CenterMiddle)
		skin.Cursor[i] = s
	}
	base := []Skin{Skin{}, DefaultSkin, UserSkin}[skin.Kind]
	skin.fillBlank(base)
}
func (skin *Skin) fillBlank(base Skin) {
	for _, s := range skin.Cursor {
		if !s.IsValid() {
			skin.Cursor = base.Cursor
			break
		}
	}
}
