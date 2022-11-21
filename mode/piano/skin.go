package piano

import (
	"fmt"
	"image/color"
	"io/fs"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/mode"
)

// KeyMode int // scratch mode + key count
type Skin struct {
	// Independent of key number
	DefaultBackground draws.Sprite
	Score             [13]draws.Sprite // number + sign(. , %)
	Combo             [10]draws.Sprite // number only
	Judgment          [5]draws.Animation
	// Dependent of key number
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
// Fields starting with lowercase are derived from mode package.
type Skins struct {
	Type  int
	Skins map[int]*Skin

	defaultBackground draws.Sprite
	score             [13]draws.Sprite // number + sign(. , %)
	Combo             [10]draws.Image  // number only
	Judgment          [5][]draws.Image
	// Bar: generated per skin
	Hint draws.Image
	// Field: generated per skin
	Note         [][4][]draws.Image
	Key          [][2]draws.Image
	KeyLighting  []draws.Image
	HitLighting  [][]draws.Image
	HoldLighting [][]draws.Image
}

var (
	DefaultSkins = Skins{Type: mode.Default, Skins: make(map[int]*Skin)}
	UserSkins    = Skins{Type: mode.User, Skins: make(map[int]*Skin)}
)

func (skins *Skins) Load(fsys fs.FS) {
	skins.defaultBackground = mode.UserSkin.DefaultBackground
	skins.score = mode.UserSkin.Score
	for i := 0; i < 10; i++ {
		skins.Combo[i] = draws.LoadImage(fsys, fmt.Sprintf("combo/%d.png", i))
	}
	for i, name := range []string{"kool", "cool", "good", "bad", "miss"} {
		skins.Judgment[i] = draws.LoadImages(fsys, fmt.Sprintf("piano/judgment/%s", name))
	}
	skins.Hint = draws.LoadImage(fsys, "piano/stage/hint.png")

	keyTypes := []KeyType{One, Two, Mid, Tip}
	keyCount := len(keyTypes)
	skins.Note = make([][4][]draws.Image, keyCount)
	skins.Key = make([][2]draws.Image, keyCount)
	skins.KeyLighting = make([]draws.Image, keyCount)
	skins.HitLighting = make([][]draws.Image, keyCount)
	skins.HoldLighting = make([][]draws.Image, keyCount)
	for k, ktype := range keyTypes {
		for i, ntype := range []string{"normal", "head", "tail", "body"} {
			ktype := []string{"one", "two", "mid", "mid"}[ktype] // Todo: "tip"
			name := fmt.Sprintf("piano/note/%s/%s", ntype, ktype)
			skins.Note[k][i] = draws.LoadImages(fsys, name)
		}
		for i, name := range []string{"up", "down"} {
			skins.Key[k][i] = draws.LoadImage(fsys, fmt.Sprintf("piano/key/%s.png", name))
		}
		skins.KeyLighting[k] = draws.LoadImage(fsys, "piano/key/lighting.png")
		skins.HitLighting[k] = draws.LoadImages(fsys, "piano/lighting/hit")
		skins.HoldLighting[k] = draws.LoadImages(fsys, "piano/lighting/hold")
	}
	// for keyMode := 4; keyMode <= 9; keyMode++ {
	// 	skins.loadSkin(keyMode)
	// }
	skins.loadSkin(7)
}

// Todo: Should each w be math.Ceil()?
// Todo: should Scratch be excluded from fieldWidth?
func (skins *Skins) loadSkin(keyMode int) {
	keyCount := len(KeyTypes[keyMode])
	skin := &Skin{
		Note:         make([][4]draws.Animation, keyCount),
		Key:          make([][2]draws.Sprite, keyCount),
		KeyLighting:  make([]draws.Sprite, keyCount),
		HitLighting:  make([]draws.Animation, keyCount),
		HoldLighting: make([]draws.Animation, keyCount),
	}
	defer func() { skins.Skins[keyMode] = skin }()
	skin.DefaultBackground = skins.defaultBackground
	skin.Score = skins.score
	for i := 0; i < 10; i++ {
		s := draws.NewSpriteFromSource(skins.Combo[i])
		s.ApplyScale(S.ComboScale)
		s.Locate(S.FieldPosition, S.ComboPosition, draws.CenterMiddle)
		skin.Combo[i] = s
	}
	for i, images := range skins.Judgment {
		a := draws.NewAnimationFromImages(images)
		for frame := range a {
			a[frame].ApplyScale(S.JudgmentScale)
			a[frame].Locate(S.FieldPosition, S.JudgmentPosition, draws.CenterMiddle)
		}
		skin.Judgment[i] = a
	}

	var fieldWidth float64
	widths := S.NoteWidths[keyMode]
	for _, ktype := range KeyTypes[keyMode] {
		fieldWidth += widths[ktype] // Todo: math.Ceil()?
	}
	{
		src := draws.NewImage(fieldWidth, 1)
		src.Fill(color.White)
		s := draws.NewSpriteFromSource(src)
		s.Locate(S.FieldPosition, S.HitPosition, draws.CenterBottom)
		skin.Bar = s
	}
	{
		s := draws.NewSpriteFromSource(skins.Hint)
		s.SetSize(fieldWidth, S.HintHeight)
		s.Locate(S.FieldPosition, S.HitPosition-S.HintHeight, draws.CenterTop)
		skin.Hint = s
	}
	{
		src := draws.NewImage(fieldWidth, ScreenSizeY)
		src.Fill(color.NRGBA{0, 0, 0, uint8(255 * S.FieldOpaque)})
		s := draws.NewSpriteFromSource(src)
		s.Locate(S.FieldPosition, 0, draws.CenterTop)
		skin.Field = s
	}
	x := S.FieldPosition - fieldWidth/2
	for k, ktype := range KeyTypes[keyMode] {
		w := S.NoteWidths[keyMode][ktype]
		x += w / 2
		for i, images := range skins.Note[ktype] {
			a := draws.NewAnimationFromImages(images)
			for frame := range a {
				a[frame].SetSize(w, S.NoteHeigth)
				a[frame].Locate(x, S.HitPosition, draws.CenterBottom)
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
		// Keys are drawn below Hint, which bottom is along with HitPosition.
		for i, image := range skins.Key[ktype] {
			s := draws.NewSpriteFromSource(image)
			s.SetSize(w, ScreenSizeY-S.HitPosition)
			s.Locate(x, S.HitPosition, draws.CenterTop)
			skin.Key[k][i] = s
		}
		{
			s := draws.NewSpriteFromSource(skins.KeyLighting[ktype])
			s.SetScaleToW(w)
			s.Locate(x, S.HitPosition, draws.CenterBottom) // -HintHeight
			skin.KeyLighting[k] = s
		}
		{
			a := draws.NewAnimationFromImages(skins.HitLighting[ktype])
			for i := range a {
				a[i].ApplyScale(S.LightingScale)
				a[i].Locate(x, S.HitPosition, draws.CenterMiddle) // -HintHeight
			}
			skin.HitLighting[k] = a
		}
		{
			a := draws.NewAnimationFromImages(skins.HoldLighting[ktype])
			for frame := range a {
				a[frame].ApplyScale(S.LightingScale)
				a[frame].Locate(x, S.HitPosition-S.HintHeight/2, draws.CenterMiddle)
			}
			skin.HoldLighting[k] = a
		}
		x += w / 2
	}
}
