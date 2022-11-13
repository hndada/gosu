package game

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/framework/audios"
	"github.com/hndada/gosu/framework/draws"
)

// Skin is a set of Sprites and sounds.
var (
	DefaultBackground  draws.Sprite
	CursorSprites      [3]draws.Sprite
	ChartItemBoxSprite draws.Sprite
	// ChartLevelBoxSprite draws.Sprite

	ScoreSprites [10]draws.Sprite
	SignSprites  [3]draws.Sprite
)

const (
	CursorSpriteBase = iota
	CursorSpriteAdditive
	CursorSpriteTrail
)

var (
	SelectSound      []byte
	SwipeSound       []byte
	TapSound         []byte
	ToggleSounds     [2][]byte
	TransitionSounds [2][]byte
)

// fsys = "skin"
func LoadGeneralSkin(fsys fs.FS) {
	DefaultBackground = NewBackground(fsys, "default-bg.jpg")
	// Todo: cursor may have outer circle
	names := []string{"menu-cursor", "menu-cursor-additive", "cursortrail"}
	for i, name := range names {
		s := draws.NewSprite(fsys, fmt.Sprintf("cursor/%s.png", name))
		s.ApplyScale(CursorScale)
		s.Locate(ScreenSizeX/2, ScreenSizeY/2, draws.CenterMiddle)
		CursorSprites[i] = s
	}
	{
		s := draws.NewSprite(fsys, "box-mask.png")
		s.SetSize(ChartInfoBoxWidth, ChartInfoBoxHeight)
		s.Locate(ScreenSizeX+chartInfoBoxshrink, ScreenSizeY/2, draws.RightMiddle)
		ChartItemBoxSprite = s
	}
	// Todo: ChartLevelBoxSprite
	for i := 0; i < 10; i++ {
		s := draws.NewSprite(fsys, fmt.Sprintf("score/%d.png", i))
		s.ApplyScale(ScoreScale)
		if i == 0 {
			s.Locate(ScreenSizeX, 0, draws.RightTop)
		} else { // Need to set same base line, since each number has different height.
			s.Locate(ScreenSizeX, ScoreSprites[0].H()-s.H(), draws.RightTop)
		}
		ScoreSprites[i] = s
	}
	for i, name := range []string{"dot", "comma", "percent"} {
		s := draws.NewSprite(fsys, fmt.Sprintf("score/%s.png", name))
		s.ApplyScale(ScoreScale)
		s.Locate(ScreenSizeX, 0, draws.RightTop)
		SignSprites[i] = s
	}
	LoadSound(fsys)
}

// Todo: need a test whether fsys is immutable
func LoadSound(fsys fs.FS) {
	fsys, err := fs.Sub(fsys, "sound")
	if err != nil {
		return
	}
	TapSound, _ = audios.NewBytes(fsys, "tap/0.wav")
	SelectSound, _ = audios.NewBytes(fsys, "old/restart.wav")
	SwipeSound, _ = audios.NewBytes(fsys, "swipe.wav")
	for i, name := range []string{"off", "on"} {
		name := fmt.Sprintf("toggle/%s.wav", name)
		ToggleSounds[i], _ = audios.NewBytes(fsys, name)
	}
	for i, name := range []string{"down", "up"} {
		name := fmt.Sprintf("transition/%s.wav", name)
		TransitionSounds[i], _ = audios.NewBytes(fsys, name)
	}
}
func NewBackground(fsys fs.FS, name string) draws.Sprite {
	s := draws.NewSprite(fsys, name)
	s.ApplyScale(ScreenSizeX / s.W())
	s.Locate(ScreenSizeX/2, ScreenSizeY/2, draws.CenterMiddle)
	return s
}

// func Paths(base string) (paths []string) {
// 	if _, err := os.Stat(base); os.IsNotExist(err) { // Only one path.
// 		return []string{fmt.Sprintf("%s.png", base)}
// 	}
// 	fs, err := os.ReadDir(base)
// 	if err != nil {
// 		return
// 	}
// 	paths = make([]string, len(fs))
// 	for i := range paths {
// 		paths[i] = fmt.Sprintf("%s/%d.png", base, i)
// 	}
// 	return
// }
