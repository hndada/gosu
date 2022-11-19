package mode

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/draws"
)

const (
	// Currently, TPS should be 1000 or greater.
	// TPS supposed to be multiple of 1000, since only one speed value
	// goes passed per Update, while unit of TransPoint's time is 1ms.
	// TPS affects only on Update(), not on Draw().
	// Todo: add lower TPS support
	TPS = 1000

	// ScreenSize is a logical size of in-game screen.
	ScreenSizeX = 1600
	ScreenSizeY = 900
)

type SkinType int

const (
	SkinTypeDefault SkinType = iota
	SkinTypeUser
	SkinTypePlay // refreshes every play
)
const (
	ScoreDot = iota
	ScoreComma
	ScorePercent
)

// In narrow meaning, Skin stands for a set of Sprites.
// In wide meaning, Skin also includes a set of sounds.
// Package defaultskin has a set of sounds.
type Skin struct {
	Type              SkinType
	DefaultBackground draws.Sprite
	Score             [13]draws.Sprite // number + sign(. , %)
	// Combo [10]draws.Sprite // number only
}

var (
	DefaultSkin = Skin{Type: SkinTypeDefault}
	UserSkin    = Skin{Type: SkinTypeUser}
	// PlaySkin    = Skin{Type: SkinTypePlay}
)

// s stands for sprite.
// a stands for animation.
// S stands for UserSettings.
func (skin *Skin) Load(fsys fs.FS) {
	S := UserSettings // abbreviation
	if skin.Type == SkinTypePlay {
		skin.Reset()
	}
	skin.DefaultBackground = NewBackground(fsys, "interface/default-bg.jpg")
	for i := 0; i < 10; i++ {
		s := draws.NewSprite(fsys, fmt.Sprintf("score/%d.png", i))
		s.ApplyScale(S.ScoreScale)
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
		s.ApplyScale(S.ScoreScale)
		s.Locate(ScreenSizeX, 0, draws.RightTop)
		skin.Score[10+i] = s
	}
	base := []Skin{{}, DefaultSkin, UserSkin}[skin.Type]
	skin.fillBlank(base)
}
func (skin *Skin) fillBlank(base Skin) {
	if !skin.DefaultBackground.IsValid() {
		skin.DefaultBackground = base.DefaultBackground
	}
	for _, s := range skin.Score {
		if !s.IsValid() {
			skin.Score = base.Score
			break
		}
	}
}
func (skin *Skin) Reset() {
	kind := skin.Type
	switch kind {
	case SkinTypeUser:
		*skin = DefaultSkin
	case SkinTypePlay:
		*skin = UserSkin
	}
	skin.Type = kind
}

func NewBackground(fsys fs.FS, name string) draws.Sprite {
	s := draws.NewSprite(fsys, name)
	s.ApplyScale(ScreenSizeX / s.W())
	s.Locate(ScreenSizeX/2, ScreenSizeY/2, draws.CenterMiddle)
	return s
}
