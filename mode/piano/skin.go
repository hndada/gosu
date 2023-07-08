package piano

import (
	"fmt"
	"image/color"
	"io/fs"

	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/mode"
)

type Skin struct {
	// Independent of key number
	DefaultBackground draws.Sprite
	Score             [13]draws.Sprite // number + sign(. , %)
	Combo             [10]draws.Sprite // number only
	Judgment          [5]draws.Animation
	Sound             []byte
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
	Skins map[int]*Skin // Key is KeyMode: scratch mode + key count

	defaultBackground draws.Sprite
	score             [13]draws.Sprite // number + sign(. , %)
	Combo             [10]draws.Image  // number only
	Judgment          [5][]draws.Image
	// Bar: generated per skin
	Hint draws.Image
	// Field: generated per skin
	Sound []byte

	Note         [4][4][]draws.Image // Key type, note type
	Key          [2]draws.Image
	KeyLighting  draws.Image
	HitLighting  []draws.Image
	HoldLighting []draws.Image
}

var (
	DefaultSkins = &Skins{Skins: make(map[int]*Skin)}
	UserSkins    = &Skins{Skins: make(map[int]*Skin)}
)

func (skins *Skins) Load(fsys fs.FS) {
	defer skins.fillBlank(DefaultSkins)
	skins.defaultBackground = mode.UserSkin.DefaultBackground
	skins.score = mode.UserSkin.Score
	for i := 0; i < 10; i++ {
		skins.Combo[i] = draws.LoadImage(fsys, fmt.Sprintf("combo/%d.png", i))
	}
	for i, name := range []string{"kool", "cool", "good", "bad", "miss"} {
		skins.Judgment[i] = draws.LoadImages(fsys, fmt.Sprintf("piano/judgment/%s", name))
	}
	skins.Hint = draws.LoadImage(fsys, "piano/stage/hint.png")
	{
		name := "piano/sound.wav"
		s := audios.NewSound(fsys, name)
		if !s.IsValid() {
			s = skins.Sound
		}
		skins.Sound = s
	}

	keyTypes := []KeyType{One, Two, Mid, Tip}
	for k, ktype := range keyTypes {
		for n, ntype := range []string{"normal", "head", "tail", "body"} {
			ktype := []string{"one", "two", "mid", "mid"}[ktype] // Todo: "tip"
			name := fmt.Sprintf("piano/note/%s/%s", ktype, ntype)
			skins.Note[k][n] = draws.LoadImages(fsys, name)
		}
	}
	for i, name := range []string{"up", "down"} {
		skins.Key[i] = draws.LoadImage(fsys, fmt.Sprintf("piano/key/%s.png", name))
	}
	skins.KeyLighting = draws.LoadImage(fsys, "piano/key/lighting.png")
	skins.HitLighting = draws.LoadImages(fsys, "piano/lighting/hit")
	skins.HoldLighting = draws.LoadImages(fsys, "piano/lighting/hold")

	skins.loadSkin(4)
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
		s.MultiplyScale(S.ComboScale)
		s.Locate(S.FieldPosition, S.ComboPosition, draws.CenterMiddle)
		skin.Combo[i] = s
	}
	for i, images := range skins.Judgment {
		a := draws.NewAnimationFromImages(images)
		for frame := range a {
			a[frame].MultiplyScale(S.JudgmentScale)
			a[frame].Locate(S.FieldPosition, S.JudgmentPosition, draws.CenterMiddle)
		}
		skin.Judgment[i] = a
	}
	skin.Sound = skins.Sound

	// Note, Bar, Hint are at bottom of HitPosition.
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
		s.Locate(S.FieldPosition, S.HitPosition, draws.CenterBottom)
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
			skin.Note[k][i] = a
		}
		// Keys are drawn below Hint, which bottom is along with HitPosition.
		for i, image := range skins.Key {
			s := draws.NewSpriteFromSource(image)
			s.SetSize(w, ScreenSizeY-S.HitPosition)
			s.Locate(x, S.HitPosition, draws.CenterTop)
			skin.Key[k][i] = s
		}
		{
			s := draws.NewSpriteFromSource(skins.KeyLighting)
			s.SetScaleToW(w)
			s.Locate(x, S.HitPosition, draws.CenterBottom) // -HintHeight
			skin.KeyLighting[k] = s
		}
		{
			a := draws.NewAnimationFromImages(skins.HitLighting)
			for i := range a {
				a[i].MultiplyScale(S.LightingScale)
				a[i].Locate(x, S.HitPosition, draws.CenterMiddle) // -HintHeight
			}
			skin.HitLighting[k] = a
		}
		{
			a := draws.NewAnimationFromImages(skins.HoldLighting)
			for frame := range a {
				a[frame].MultiplyScale(S.LightingScale)
				a[frame].Locate(x, S.HitPosition-S.HintHeight/2, draws.CenterMiddle)
			}
			skin.HoldLighting[k] = a
		}
		x += w / 2
	}
}
func (skins *Skins) fillBlank(base *Skins) {
	for _, img := range skins.Combo {
		if !img.IsValid() {
			skins.Combo = base.Combo
			break
		}
	}
	for i, imgs := range skins.Judgment {
		for _, img := range imgs {
			if !img.IsValid() {
				skins.Judgment[i] = base.Judgment[i]
				break
			}
		}
	}
	if !skins.Hint.IsValid() {
		skins.Hint = base.Hint
	}
	for k, ktypes := range skins.Note {
		for n := range ktypes {
			switch n {
			case Normal: // Use default's.
				for _, img := range skins.Note[k][n] {
					if !img.IsValid() {
						skins.Note[k][n] = base.Note[k][n]
						break
					}
				}
			case Head: // Use user's.
				skins.Note[k][n] = skins.Note[k][Normal]
			case Tail: // Skip validation for tail.
				continue
			case Body: // Use user's.
				skins.Note[k][n] = skins.Note[k][Normal]
			}
		}
	}
	for _, img := range skins.Key {
		if !img.IsValid() {
			skins.Key = base.Key
			break
		}
	}
	if !skins.KeyLighting.IsValid() {
		skins.KeyLighting = base.KeyLighting
	}
	for _, img := range skins.HitLighting {
		if !img.IsValid() {
			skins.HitLighting = base.HitLighting
			break
		}
	}
	for _, img := range skins.HoldLighting {
		if !img.IsValid() {
			skins.HoldLighting = base.HoldLighting
			break
		}
	}
}
