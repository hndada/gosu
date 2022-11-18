package new

import (
	"io/fs"

	"github.com/hndada/gosu/defaultskin"
	"github.com/hndada/gosu/draws"
)

type SkinKind int

const (
	SkinKindDefault SkinKind = iota
	SkinKindUser
	SkinKindPlay
)

type Skin struct {
	Kind  SkinKind
	Score [13]draws.Sprite
	Combo [10]draws.Sprite
}

var (
	DefaultSkin = Skin{Kind: SkinKindDefault}
	UserSkin    = Skin{Kind: SkinKindUser}
	PlaySkin    = Skin{Kind: SkinKindPlay}
)

func (skin *Skin) Load(fsys fs.FS) {
	if skin.Kind == SkinKindPlay {
		skin.Reset()
	}
	// loads skin data
	skin.fillBlank(fsys)
}
func (skin *Skin) fillBlank(fsys fs.FS) {
	base := []Skin{Skin{}, DefaultSkin, UserSkin}[skin.Kind]
	// fills blank data with base
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
func init() {
	DefaultSkin.Load(defaultskin.FS)
}

// game.go
func load(fsys fs.FS) {
	mode.UserSkin.Load(fsys)
}
