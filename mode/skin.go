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

type SkinKind int

const (
	SkinKindDefault SkinKind = iota
	SkinKindUser
	SkinKindPlay // refreshes every play
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
	Kind              SkinKind
	DefaultBackground draws.Sprite
	Score             [13]draws.Sprite // number + sign(. , %)
	// Combo [10]draws.Sprite // number only
}

var (
	DefaultSkin = Skin{Kind: SkinKindDefault}
	UserSkin    = Skin{Kind: SkinKindUser}
	// PlaySkin    = Skin{Kind: SkinKindPlay}
)

func (skin *Skin) Load(fsys fs.FS) {
	if skin.Kind == SkinKindPlay {
		skin.Reset()
	}
	skin.DefaultBackground = NewBackground(fsys, "interface/default-bg.jpg")
	for i := 0; i < 10; i++ {
		s := draws.NewSprite(fsys, fmt.Sprintf("score/%d.png", i))
		s.ApplyScale(UserSettings.ScoreScale)
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
		s.ApplyScale(UserSettings.ScoreScale)
		s.Locate(ScreenSizeX, 0, draws.RightTop)
		skin.Score[10+i] = s
	}
	base := []Skin{{}, DefaultSkin, UserSkin}[skin.Kind]
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
	switch skin.Kind {
	case SkinKindUser:
		*skin = DefaultSkin
		skin.Kind = SkinKindUser
	case SkinKindPlay:
		*skin = UserSkin
		skin.Kind = SkinKindPlay
	}
}

func NewBackground(fsys fs.FS, name string) draws.Sprite {
	s := draws.NewSprite(fsys, name)
	s.ApplyScale(ScreenSizeX / s.W())
	s.Locate(ScreenSizeX/2, ScreenSizeY/2, draws.CenterMiddle)
	return s
}
