package game

import (
	"fmt"

	"github.com/hndada/gosu/framework/audios"
	"github.com/hndada/gosu/framework/draws"
)

// ScreenSize is a logical size of in-game screen.
const (
	screenSizeX = 1600
	screenSizeY = 900

	ScreenSizeX = screenSizeX
	ScreenSizeY = screenSizeY
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

func LoadGeneralSkin() {
	DefaultBackground = NewBackground("skin/default-bg.jpg")
	// Todo: cursor may have outer circle
	names := []string{"menu-cursor", "menu-cursor-additive", "cursortrail"}
	for i, name := range names {
		s := draws.NewSprite(fmt.Sprintf("skin/cursor/%s.png", name))
		s.ApplyScale(CursorScale)
		s.Locate(screenSizeX/2, screenSizeY/2, draws.CenterMiddle)
		CursorSprites[i] = s
	}
	{
		s := draws.NewSprite("skin/box-mask.png")
		s.SetSize(ChartInfoBoxWidth, ChartInfoBoxHeight)
		s.Locate(screenSizeX+chartInfoBoxshrink, screenSizeY/2, draws.RightMiddle)
		ChartItemBoxSprite = s
	}
	// Todo: ChartLevelBoxSprite
	for i := 0; i < 10; i++ {
		s := draws.NewSprite(fmt.Sprintf("skin/score/%d.png", i))
		s.ApplyScale(ScoreScale)
		if i == 0 {
			s.Locate(screenSizeX, 0, draws.RightTop)
		} else { // Need to set same base line, since each number has different height.
			s.Locate(screenSizeX, ScoreSprites[0].H()-s.H(), draws.RightTop)
		}
		ScoreSprites[i] = s
	}
	for i, name := range []string{"dot", "comma", "percent"} {
		s := draws.NewSprite(fmt.Sprintf("skin/score/%s.png", name))
		s.ApplyScale(ScoreScale)
		s.Locate(screenSizeX, 0, draws.RightTop)
		SignSprites[i] = s
	}
	TapSound, _ = audios.NewBytes("skin/sound/tap/0.wav")
	SelectSound, _ = audios.NewBytes("skin/sound/old/restart.wav")
	SwipeSound, _ = audios.NewBytes("skin/sound/swipe.wav")
	for i, name := range []string{"off", "on"} {
		path := fmt.Sprintf("skin/sound/toggle/%s.wav", name)
		ToggleSounds[i], _ = audios.NewBytes(path)
	}
	for i, name := range []string{"down", "up"} {
		path := fmt.Sprintf("skin/sound/transition/%s.wav", name)
		TransitionSounds[i], _ = audios.NewBytes(path)
	}
}
func NewBackground(path string) draws.Sprite {
	s := draws.NewSprite(path)
	s.ApplyScale(screenSizeX / s.W())
	s.Locate(screenSizeX/2, screenSizeY/2, draws.CenterMiddle)
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
