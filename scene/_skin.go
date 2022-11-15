package scene

import (
	"embed"
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/draws"
)

// const (
// 	CursorBase = iota
// 	CursorAdditive
// 	CursorTrail
// )

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

var skin *Skin
var defaultSkin *Skin

func (Skin) init() {
	//go:embed skin/*
	var fs embed.FS
	data, err := Skin{}.Load(fs)
	if err != nil {
		panic(err)
	}
	defaultSkin.Set(data)
}

func (Skin) Default() Setter { return *defaultSkin }
func (Skin) Current() Setter { return *skin }
func (Skin) Set(s Setter)    { skin = s.(*Skin) }
func (Skin) Load(fsys any) (Setter, error) {
	return Skin{}.load(fsys.(fs.FS))
}
func (Skin) load(fsys fs.FS) (Setter, error) {
	skin := Skin{}.Default().(Skin)
	skin.DefaultBackground = NewBackground(fsys, "default-bg.jpg")
	// Todo: cursor may have outer circle
	names := []string{"menu-cursor", "menu-cursor-additive", "cursortrail"}
	for i, name := range names {
		s := draws.NewSprite(fsys, fmt.Sprintf("cursor/%s.png", name))
		s.ApplyScale(settings.CursorScale)
		s.Locate(ScreenSizeX/2, ScreenSizeY/2, draws.CenterMiddle)
		skin.Cursor[i] = s
	}
	// {
	// 	s := draws.NewSprite(fsys, "box-mask.png")
	// 	s.SetSize(ChartInfoBoxWidth, ChartInfoBoxHeight)
	// 	s.Locate(ScreenSizeX+chartInfoBoxshrink, ScreenSizeY/2, draws.RightMiddle)
	// 	ChartItemBoxSprite = s
	// }
	for i := 0; i < 10; i++ {
		s := draws.NewSprite(fsys, fmt.Sprintf("score/%d.png", i))
		s.ApplyScale(settings.ScoreScale)
		if i == 0 {
			s.Locate(ScreenSizeX, 0, draws.RightTop)
		} else { // Need to set same base line, since each number has different height.
			s.Locate(ScreenSizeX, ScoreSprites[0].H()-s.H(), draws.RightTop)
		}
		ScoreSprites[i] = s
	}
	for i, name := range []string{"dot", "comma", "percent"} {
		s := draws.NewSprite(fsys, fmt.Sprintf("score/%s.png", name))
		s.ApplyScale(settings.ScoreScale)
		s.Locate(ScreenSizeX, 0, draws.RightTop)
		SignSprites[i] = s
	}
	LoadSound(fsys)
}
