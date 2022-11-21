package piano

import (
	"fmt"
	"image/color"
	"io/fs"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/mode"
)

type Skin struct {
	Type    int
	KeyMode int // scratch mode + key count

	// independent of key number
	DefaultBackground draws.Sprite
	Score             [13]draws.Sprite // number + sign(. , %)
	Combo             [10]draws.Sprite // number only
	Judgment          [5]draws.Animation
	// dependent of key number
	Bar          draws.Sprite
	Hint         draws.Sprite
	Field        draws.Sprite
	Note         [][4]draws.Animation
	Key          [][2]draws.Sprite
	KeyLighting  []draws.Sprite
	HitLighting  []draws.Animation
	HoldLighting []draws.Animation
}

// Each piano's sub mode has different skin.
type Skins map[int]*Skin

const general = 0

var (
	DefaultSkins = Skins{general: NewSkin(mode.Default, general)}
	UserSkins    = Skins{general: NewSkin(mode.User, general)}
	// PlaySkin doesn't have to be slice, since it is one-time struct.
	// PlaySkin *Skin
)

func NewSkin(_type, k int) *Skin {
	keyCount := len(KeyTypes[k])
	return &Skin{
		Type:         _type,
		KeyMode:      k,
		Note:         make([][4]draws.Animation, keyCount),
		Key:          make([][2]draws.Sprite, keyCount),
		KeyLighting:  make([]draws.Sprite, keyCount),
		HitLighting:  make([]draws.Animation, keyCount),
		HoldLighting: make([]draws.Animation, keyCount),
	}
}
func (s Skin) isGeneral() bool { return s.KeyMode == general }
func (skins Skins) Load(fsys fs.FS) {
	skins[general].Load(fsys)
	for k := 4; k <= 9; k++ {
		skins.loadKeyMode(fsys, k)
	}
}

// LoadKeyMode loads particular key mode. Useful for lazy load.
func (skins Skins) loadKeyMode(fsys fs.FS, k int) {
	skins[k] = NewSkin(skins[general].Type, k)
	skins[k].Load(fsys)
	// fmt.Printf("%d: %+v\n", k, skins[k].KeyLighting)
}

func (skin *Skin) Load(fsys fs.FS) {
	var generalSkin *Skin
	switch skin.Type {
	case mode.Default:
		generalSkin = DefaultSkins[general]
		defer func() { DefaultSkins[skin.KeyMode] = skin }()
	case mode.User:
		generalSkin = UserSkins[general]
		defer func() { UserSkins[skin.KeyMode] = skin }()
	case mode.Play:
		skin.Reset()
	}
	skin.DefaultBackground = mode.UserSkin.DefaultBackground
	skin.Score = mode.UserSkin.Score
	for i := 0; i < 10; i++ {
		var s draws.Sprite
		if skin.isGeneral() {
			s = draws.NewSprite(fsys, fmt.Sprintf("combo/%d.png", i))
		} else {
			s = generalSkin.Combo[i]
		}
		s.ApplyScale(S.ComboScale)
		s.Locate(S.FieldPosition, S.ComboPosition, draws.CenterMiddle)
		skin.Combo[i] = s
	}
	for i, name := range []string{"kool", "cool", "good", "bad", "miss"} {
		var a draws.Animation
		if skin.isGeneral() {
			a = draws.NewAnimation(fsys, fmt.Sprintf("piano/judgment/%s", name))
		} else {
			a = make(draws.Animation, len(generalSkin.Judgment[i]))
			copy(a, generalSkin.Judgment[i])
			for i := range a {
				a[i].ApplyScale(S.JudgmentScale)
				a[i].Locate(S.FieldPosition, S.JudgmentPosition, draws.CenterMiddle)
				fmt.Println(a[i].Scale)
			}
		}
		skin.Judgment[i] = a
		if skin.KeyMode == 7 {
			fmt.Printf("%+v\n", skin.Judgment)
		}
	}
	// Keys are drawn below Hint, which bottom is along with HitPosition.
	// Each w should be integer, since it is a width of independent sprite.
	// Todo: should Scratch be excluded from fw?
	fw := skin.fieldWidth()
	{
		src := draws.NewImage(fw, 1)
		src.Fill(color.White)
		s := draws.NewSpriteFromSource(src)
		s.Locate(S.FieldPosition, S.HitPosition, draws.CenterBottom)
		skin.Bar = s
	}
	{
		var s draws.Sprite
		if skin.isGeneral() {
			s = draws.NewSprite(fsys, "piano/stage/hint.png")
		} else {
			s = generalSkin.Hint
		}
		s.SetSize(fw, S.HintHeight)
		s.Locate(S.FieldPosition, S.HitPosition-S.HintHeight, draws.CenterTop)
		skin.Hint = s
	}
	{
		src := draws.NewImage(fw, ScreenSizeY)
		src.Fill(color.NRGBA{0, 0, 0, uint8(255 * S.FieldOpaque)})
		s := draws.NewSpriteFromSource(src)
		s.Locate(S.FieldPosition, 0, draws.CenterTop)
		skin.Field = s
	}
	x := S.FieldPosition - fw/2
	// keyCount := len(KeyTypes[skin.KeyMode])
	for k, ktype := range KeyTypes[skin.KeyMode] {
		w := S.NoteWidths[skin.KeyMode][ktype] // Todo: math.Ceil()?
		x += w / 2
		// skin.Note = make([][4]draws.Animation, keyCount)
		for i, ntype := range []string{"normal", "head", "tail", "body"} {
			var a draws.Animation
			if skin.isGeneral() {
				ktype := []string{"one", "two", "mid", "mid"}[ktype] // Todo: "tip"
				name := fmt.Sprintf("piano/note/%s/%s", ntype, ktype)
				a = draws.NewAnimation(fsys, name)
			} else {
				a = generalSkin.Note[ktype][i]
			}
			// if !skin.isGeneral() && ktype == Tip && !a.IsValid() {
			// 	for frame := range a {
			// 		a[frame] = skin.Note[k][0][frame]
			// 		op := draws.Op{}
			// 		op.ColorM.ScaleWithColor(S.scratchColor)
			// 		i := a[frame].Source.(draws.Image) // Todo: looks weird usage to me
			// 		skin.Note[k][0][frame].Draw(i, op)
			// 		a[frame].SetSize(w, S.NoteHeigth)
			// 		a[frame].Locate(x, S.HitPosition, draws.CenterBottom)
			// 	}
			// }
			skin.Note[k][i] = a
		}
		// skin.Key = make([][2]draws.Sprite, keyCount)
		for i, name := range []string{"up", "down"} {
			var s draws.Sprite
			if skin.isGeneral() {
				s = draws.NewSprite(fsys, fmt.Sprintf("piano/key/%s.png", name))
			} else {
				s = generalSkin.Key[0][i]
			}
			s.SetSize(w, ScreenSizeY-S.HitPosition)
			s.Locate(x, S.HitPosition, draws.CenterTop)
			skin.Key[k][i] = s
		}
		{
			// skin.KeyLighting = make([]draws.Sprite, keyCount)
			var s draws.Sprite
			if skin.isGeneral() {
				s = draws.NewSprite(fsys, "piano/key/lighting.png")
			} else {
				s = generalSkin.KeyLighting[ktype]
			}
			s.SetScaleToW(w)
			s.Locate(x, S.HitPosition, draws.CenterBottom) // -HintHeight
			skin.KeyLighting[k] = s
		}
		{
			// skin.HitLighting = make([]draws.Animation, keyCount)
			var a draws.Animation
			if skin.isGeneral() {
				a = draws.NewAnimation(fsys, "piano/lighting/hit")
			} else {
				a = generalSkin.HitLighting[ktype]
			}
			for i := range a {
				a[i].ApplyScale(S.LightingScale)
				a[i].Locate(x, S.HitPosition, draws.CenterMiddle) // -HintHeight
			}
			skin.HitLighting[k] = a
		}
		{
			// skin.HoldLighting = make([]draws.Animation, keyCount)
			var a draws.Animation
			if skin.isGeneral() {
				a = draws.NewAnimation(fsys, "piano/lighting/hold")
			} else {
				a = generalSkin.HoldLighting[ktype]
			}
			for i := range a {
				a[i].ApplyScale(S.LightingScale)
				a[i].Locate(x, S.HitPosition-S.HintHeight/2, draws.CenterMiddle)
			}
			skin.HoldLighting[k] = a
		}
		x += w / 2
	}
}
func (skin Skin) fieldWidth() float64 {
	var fw float64
	widths := S.NoteWidths[skin.KeyMode]
	for _, ktype := range KeyTypes[skin.KeyMode] {
		fw += widths[ktype] // Todo: math.Ceil()?
	}
	return fw
}

func (skin *Skin) Reset() {
	kind := skin.Type
	switch kind {
	case mode.User:
		*skin = *DefaultSkins[skin.KeyMode]
	case mode.Play:
		*skin = *UserSkins[skin.KeyMode]
	}
	skin.Type = kind
}
