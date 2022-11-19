package piano

import (
	"fmt"
	"image/color"
	"io/fs"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/mode"
)

type Skin struct {
	Kind    mode.SkinKind
	KeyMode int // scratch mode + key count

	// independent of key number
	Score    [13]draws.Sprite
	Combo    [10]draws.Sprite
	Judgment [5]draws.Animation
	// dependent of key number
	Key          [2][]draws.Sprite
	KeyLighting  []draws.Sprite
	HitLighting  []draws.Animation
	HoldLighting []draws.Animation
	Note         [4][]draws.Animation
	Field        draws.Sprite
	Hint         draws.Sprite
	Bar          draws.Sprite
}
type Skins map[int]Skin

const base = 0

func (s Skin) isBase() bool { return s.KeyMode == base }

// Each piano's sub mode has different skin.
// PlaySkin doesn't have to be slice, since it is one-time struct.
var (
	DefaultSkins = Skins{base: {Kind: mode.SkinKindDefault}}
	UserSkins    = Skins{base: {Kind: mode.SkinKindUser}}
	PlaySkin     = Skin{Kind: mode.SkinKindPlay}
)

func (skins Skins) Load(fsys fs.FS) {
	baseSkin := skins[base]
	baseSkin.Load(fsys)
	skins[base] = baseSkin
	for k := 4; k <= 9; k++ {
		skins.load(fsys, k)
	}
}

// load lazy loads less popular key mode.
func (skins Skins) load(fsys fs.FS, k int) {
	skin := Skin{Kind: skins[base].Kind, KeyMode: k}
	skin.Load(fsys)
	skins[k] = skin
}

func (skin *Skin) Load(fsys fs.FS) {
	var baseSkin Skin
	switch skin.Kind {
	case mode.SkinKindDefault:
		baseSkin = DefaultSkins[base]
	case mode.SkinKindUser:
		baseSkin = UserSkins[base]
	case mode.SkinKindPlay:
		skin.Reset()
	}
	S := UserSettings
	for i := 0; i < 10; i++ {
		s := baseSkin.Combo[i]
		if skin.isBase() {
			s = draws.NewSprite(fsys, fmt.Sprintf("combo/%d.png", i))
		}
		s.ApplyScale(S.ComboScale)
		s.Locate(S.FieldPosition, S.ComboPosition, draws.CenterMiddle)
		skin.Combo[i] = s
	}
	for i, name := range []string{"kool", "cool", "good", "bad", "miss"} {
		a := baseSkin.Judgment[i]
		if skin.isBase() {
			a = draws.NewAnimation(fsys, fmt.Sprintf("piano/judgment/%s", name))
		}
		for i := range a {
			a[i].ApplyScale(S.JudgmentScale)
			a[i].Locate(S.FieldPosition, S.JudgmentPosition, draws.CenterMiddle)
		}
		skin.Judgment[i] = a
	}
	// Keys are drawn below Hint, which bottom is along with HitPosition.
	// Each w should be integer, since it is a width of independent sprite.
	// Todo: should Scratch be excluded from fw?
	fw := skin.fieldWidth()
	x := S.FieldPosition - fw/2
	keyCount := len(KeyTypes[skin.KeyMode])
	for k, ktype := range KeyTypes[skin.KeyMode] {
		w := S.NoteWidths[skin.KeyMode][ktype] // Todo: math.Ceil()?
		if skin.isBase() {
			w = 0.06 * ScreenSizeX
		}
		x += w / 2
		for i, name := range []string{"up", "down"} {
			skin.Key[i] = make([]draws.Sprite, keyCount)
			s := baseSkin.Key[i][0]
			if skin.isBase() {
				s = draws.NewSprite(fsys, fmt.Sprintf("piano/key/%s.png", name))
			}
			s.SetSize(w, ScreenSizeY-S.HitPosition)
			s.Locate(x, S.HitPosition, draws.CenterTop)
			skin.Key[i][k] = s
		}
		{
			skin.KeyLighting = make([]draws.Sprite, keyCount)
			s := baseSkin.KeyLighting[ktype]
			if skin.isBase() {
				s = draws.NewSprite(fsys, "piano/key/lighting.png")
			}
			s.SetScaleToW(w)
			s.Locate(x, S.HitPosition, draws.CenterBottom) // -HintHeight
			skin.KeyLighting[k] = s
		}
		{
			skin.HitLighting = make([]draws.Animation, keyCount)
			a := baseSkin.HitLighting[ktype]
			if skin.isBase() {
				a = draws.NewAnimation(fsys, "piano/lighting/hit")
			}
			for i := range a {
				a[i].ApplyScale(S.LightingScale)
				a[i].Locate(x, S.HitPosition, draws.CenterMiddle) // -HintHeight
			}
			skin.HitLighting[k] = a
		}
		{
			skin.HoldLighting = make([]draws.Animation, keyCount)
			a := baseSkin.HoldLighting[ktype]
			if skin.isBase() {
				a = draws.NewAnimation(fsys, "piano/lighting/hold")
			}
			for i := range a {
				a[i].ApplyScale(S.LightingScale)
				a[i].Locate(x, S.HitPosition-S.HintHeight/2, draws.CenterMiddle)
			}
			skin.HoldLighting[k] = a
		}
		for i, nTypeName := range []string{"normal", "head", "tail", "body"} {
			skin.Note[i] = make([]draws.Animation, keyCount)
			a := baseSkin.Note[i][ktype]
			if skin.isBase() {
				kTypeName := []string{"one", "two", "mid", "tip"}[ktype]
				name := fmt.Sprintf("piano/note/%s/%s", nTypeName, kTypeName)
				a = draws.NewAnimation(fsys, name)
			}
			for frame := range a {
				if ktype == Tip && !a.IsValid() {
					a[frame] = skin.Note[0][k][frame]
					op := draws.Op{}
					op.ColorM.ScaleWithColor(S.scratchColor)
					i := a[frame].Source.(draws.Image) // Todo: looks weird usage to me
					skin.Note[0][k][frame].Draw(i, op)
				}
				a[frame].SetSize(w, S.NoteHeigth)
				a[frame].Locate(x, S.HitPosition, draws.CenterBottom)
			}
			skin.Note[i][k] = a
		}
		x += w / 2
	}
	{
		s := baseSkin.Field
		if skin.isBase() {
			src := draws.NewImage(fw, ScreenSizeY)
			src.Fill(color.NRGBA{0, 0, 0, uint8(255 * S.FieldOpaque)})
			s = draws.NewSpriteFromSource(src)
		}
		s.Locate(S.FieldPosition, 0, draws.CenterTop)
		skin.Field = s
	}
	{
		s := baseSkin.Hint
		if skin.isBase() {
			s = draws.NewSprite(fsys, "piano/stage/hint.png")
		}
		s.SetSize(fw, S.HintHeight)
		s.Locate(S.FieldPosition, S.HitPosition-S.HintHeight, draws.CenterTop)
		skin.Hint = s
	}
	{
		s := baseSkin.Bar
		if skin.isBase() {
			src := draws.NewImage(fw, 1)
			src.Fill(color.White)
			s = draws.NewSpriteFromSource(src)
		}
		s.Locate(S.FieldPosition, S.HitPosition, draws.CenterBottom)
		skin.Bar = s
	}
}
func (skin Skin) fieldWidth() float64 {
	var fw float64
	widths := UserSettings.NoteWidths[skin.KeyMode]
	for _, ktype := range KeyTypes[skin.KeyMode] {
		fw += widths[ktype] // Todo: math.Ceil()?
	}
	return fw
}

func (skin *Skin) Reset() {
	kind := skin.Kind
	switch kind {
	case mode.SkinKindUser:
		*skin = DefaultSkins[skin.KeyMode]
	case mode.SkinKindPlay:
		*skin = UserSkins[skin.KeyMode]
	}
	skin.Kind = kind
}
