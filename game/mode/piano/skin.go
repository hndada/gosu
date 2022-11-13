package piano

import (
	"fmt"
	"image/color"
	"io/fs"
	"math"

	_ "image/jpeg"
	_ "image/png"

	"github.com/hndada/gosu/framework/draws"
	"github.com/hndada/gosu/game"
)

type NoteKind int

const (
	One NoteKind = iota
	Two
	Mid
	Tip = Mid
)

var NoteKindsMap = map[int][]NoteKind{
	0:  {},
	1:  {Mid},
	2:  {One, One},
	3:  {One, Mid, One},
	4:  {One, Two, Two, One},
	5:  {One, Two, Mid, Two, One},
	6:  {One, Two, One, One, Two, One},
	7:  {One, Two, One, Mid, One, Two, One},
	8:  {Tip, One, Two, One, One, Two, One, Tip},
	9:  {Tip, One, Two, One, Mid, One, Two, One, Tip},
	10: {Tip, One, Two, One, Mid, Mid, One, Two, One, Tip},
}

// LeftScratch and RightScratch are bits for indicating scratch mode.
// For example, when key count is 40 = 32 + 8, it is 8-key with left scratch.
const (
	LeftScratch  = 32
	RightScratch = 64
	ScratchMask  = ^(LeftScratch | RightScratch)
)

func init() { // I'm proud of the following code.
	for k := 2; k <= 8; k++ {
		NoteKindsMap[k|LeftScratch] = append([]NoteKind{Tip}, NoteKindsMap[k-1]...)
		NoteKindsMap[k|RightScratch] = append(NoteKindsMap[k-1], Tip)
	}
}

// GeneralSkin is a singleton.
var GeneralSkin struct {
	ComboSprites    [10]draws.Sprite
	JudgmentSprites [5]draws.Animation
}

// Todo: should each skin has own skin settings?
type Skin struct {
	// Sprites which are independent of key count.
	ScoreSprites    [10]draws.Sprite
	SignSprites     [3]draws.Sprite
	ComboSprites    [10]draws.Sprite
	JudgmentSprites [5]draws.Animation
	// Sprites which are dependent of key count.
	KeySprites          [][2]draws.Sprite
	KeyLightingSprites  []draws.Sprite
	HitLightingSprites  []draws.Animation
	HoldLightingSprites []draws.Animation
	NoteSprites         [][4]draws.Animation
	FieldSprite         draws.Sprite
	HintSprite          draws.Sprite
	BarSprite           draws.Sprite
}

var Skins = make(map[int]Skin)

// fsys = "skin"
// Todo: need a test whether fsys is immutable
func LoadSkin(fsys fs.FS) {
	// defer func() { fmt.Printf("%+v\n", Skins[7]) }()
	for i := 0; i < 10; i++ {
		sprite := draws.NewSprite(fsys, fmt.Sprintf("combo/%d.png", i))
		sprite.ApplyScale(ComboScale)
		sprite.Locate(FieldPosition, ComboPosition, draws.CenterMiddle)
		GeneralSkin.ComboSprites[i] = sprite
	}
	for i, name := range []string{"kool", "cool", "good", "bad", "miss"} {
		animation := draws.NewAnimation(fsys, fmt.Sprintf("piano/judgment/%s", name))
		for i := range animation {
			animation[i].ApplyScale(JudgmentScale)
			animation[i].Locate(FieldPosition, JudgmentPosition, draws.CenterMiddle)
		}
		GeneralSkin.JudgmentSprites[i] = animation
	}
	var (
		keyImages          [2]draws.Image
		keyLightingImage   draws.Image
		hitLightingImages  []draws.Image
		holdLightingImages []draws.Image
		noteImages         [4][4][]draws.Image // First 4 is a type, next 4 is a kind.
		hintImage          draws.Image
	)
	for i, name := range []string{"up", "down"} {
		keyImages[i] = draws.LoadImage(fsys, fmt.Sprintf("piano/key/%s.png", name))
	}
	keyLightingImage = draws.LoadImage(fsys, "piano/key/lighting.png")
	hitLightingImages = draws.LoadImages(fsys, "piano/lighting/hit")
	holdLightingImages = draws.LoadImages(fsys, "piano/lighting/hold")
	for i, _type := range []string{"normal", "head", "tail", "body"} {
		for j, kind := range []int{1, 2, 3, 3} { // Todo: 4th note image with custom color settings?
			noteImages[i][j] = draws.LoadImages(fsys, fmt.Sprintf("piano/note/%s/%d", _type, kind))
		}
	}
	hintImage = draws.LoadImage(fsys, "piano/hint.png")

	// Todo: Key count 1, 2, 3 and with scratch
	for keyCount := 4; keyCount <= 10; keyCount++ {
		noteKinds := NoteKindsMap[keyCount]
		noteWidths := NoteWidthsMap[keyCount]
		skin := Skin{
			ScoreSprites:        game.ScoreSprites,
			SignSprites:         game.SignSprites,
			ComboSprites:        GeneralSkin.ComboSprites,
			JudgmentSprites:     GeneralSkin.JudgmentSprites,
			KeySprites:          make([][2]draws.Sprite, len(noteKinds)),
			KeyLightingSprites:  make([]draws.Sprite, len(noteKinds)),
			HitLightingSprites:  make([]draws.Animation, len(noteKinds)),
			HoldLightingSprites: make([]draws.Animation, len(noteKinds)),
			NoteSprites:         make([][4]draws.Animation, len(noteKinds)),
		}
		// Keys are drawn below Hint, which bottom is along with HitPosition.
		// Each w should be integer, since it is a width of independent sprite.
		// Todo: Scratch should be excluded to width sum.
		var wsum float64
		for _, kind := range noteKinds {
			wsum += math.Ceil(noteWidths[kind])
		}
		x := FieldPosition - wsum/2
		for k, kind := range noteKinds {
			w := math.Ceil(noteWidths[kind])
			x += w / 2
			for i, img := range keyImages {
				sprite := draws.NewSpriteFromSource(img)
				sprite.SetSize(w, ScreenSizeY-HitPosition)
				sprite.Locate(x, HitPosition, draws.CenterTop)
				skin.KeySprites[k][i] = sprite
			}
			{
				sprite := draws.NewSpriteFromSource(keyLightingImage)
				sprite.SetScaleToW(w)
				sprite.Locate(x, HitPosition, draws.CenterBottom) // -HintHeight
				skin.KeyLightingSprites[k] = sprite
			}
			{
				animation := draws.NewAnimationFromImages(hitLightingImages)
				for i := range animation {
					animation[i].ApplyScale(LightingScale)
					animation[i].Locate(x, HitPosition, draws.CenterMiddle) // -HintHeight
				}
				skin.HitLightingSprites[k] = animation
			}
			{
				animation := draws.NewAnimationFromImages(holdLightingImages)
				for i := range animation {
					animation[i].ApplyScale(LightingScale)
					animation[i].Locate(x, HitPosition-HintHeight/2, draws.CenterMiddle)
				}
				skin.HoldLightingSprites[k] = animation
			}
			for _type, images := range noteImages {
				animation := draws.NewAnimationFromImages(images[kind])
				for i := range animation {
					animation[i].SetSize(w, NoteHeigth)
					animation[i].Locate(x, HitPosition, draws.CenterBottom)
				}
				skin.NoteSprites[k][_type] = animation
			}
			x += w / 2
		}
		{
			src := draws.NewImage(wsum, ScreenSizeY)
			src.Fill(color.NRGBA{0, 0, 0, uint8(255 * FieldDarkness)})
			sprite := draws.NewSpriteFromSource(src)
			sprite.Locate(FieldPosition, 0, draws.CenterTop)
			skin.FieldSprite = sprite
		}
		{
			sprite := draws.NewSpriteFromSource(hintImage)
			sprite.SetSize(wsum, HintHeight)
			sprite.Locate(FieldPosition, HitPosition-HintHeight, draws.CenterTop)
			skin.HintSprite = sprite
		}
		{
			src := draws.NewImage(wsum, 1)
			src.Fill(color.White)
			sprite := draws.NewSpriteFromSource(src)
			sprite.Locate(FieldPosition, HitPosition, draws.CenterBottom)
			skin.BarSprite = sprite
		}
		Skins[keyCount] = skin
	}
}
