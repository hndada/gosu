package piano

import (
	"fmt"
	"image/color"
	"math"

	_ "image/jpeg"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu"
	"github.com/hndada/gosu/draws"
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
	KeySprites  [][2]draws.Sprite
	NoteSprites [][4]draws.Animation
	FieldSprite draws.Sprite
	HintSprite  draws.Sprite
	BarSprite   draws.Sprite
}

var Skins = make(map[int]Skin)

func LoadSkin() {
	for i := 0; i < 10; i++ {
		s := draws.NewSprite(fmt.Sprintf("skin/combo/%d.png", i))
		s.ApplyScale(ComboScale)
		s.SetPoint(FieldPosition, ComboPosition, draws.CenterMiddle)
		GeneralSkin.ComboSprites[i] = s
	}
	for i, name := range []string{"kool", "cool", "good", "bad", "miss"} {
		a := draws.NewAnimation(fmt.Sprintf("skin/piano/judgment/%s", name))
		for i := range a {
			a[i].ApplyScale(JudgmentScale)
			a[i].SetPoint(FieldPosition, JudgmentPosition, draws.CenterMiddle)
		}
		GeneralSkin.JudgmentSprites[i] = a
	}
	var (
		keyImages  [2]*ebiten.Image
		noteImages [4][4][]*ebiten.Image // First 4 is a type, next 4 is a kind.
		hintImage  *ebiten.Image
	)
	for i, name := range []string{"up", "down"} {
		keyImages[i] = draws.NewImage(fmt.Sprintf("skin/piano/key/%s.png", name))
	}
	hintImage = draws.NewImage("skin/piano/hint.png")
	for i, _type := range []string{"normal", "head", "tail", "body"} {
		for j, kind := range []int{1, 2, 3, 3} { // Todo: 4th note image with custom color settings?
			noteImages[i][j] = draws.NewImages(fmt.Sprintf("skin/piano/note/%s/%d", _type, kind))
		}
	}
	// Todo: Key count 1, 2, 3 and with scratch
	for keyCount := 4; keyCount <= 10; keyCount++ {
		noteKinds := NoteKindsMap[keyCount]
		noteWidths := NoteWidthsMap[keyCount]
		skin := Skin{
			ScoreSprites:    gosu.ScoreSprites,
			SignSprites:     gosu.SignSprites,
			ComboSprites:    GeneralSkin.ComboSprites,
			JudgmentSprites: GeneralSkin.JudgmentSprites,
			KeySprites:      make([][2]draws.Sprite, len(noteKinds)),
			NoteSprites:     make([][4]draws.Animation, len(noteKinds)),
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
			for i, img := range keyImages {
				s := draws.NewSpriteFromImage(img)
				s.SetSize(w, screenSizeY-HitPosition)
				s.SetPoint(x, HitPosition, draws.LeftTop)
				skin.KeySprites[k][i] = s
			}
			for _type, images := range noteImages {
				animation := draws.NewAnimationFromImages(images[kind])
				for i := range animation {
					animation[i].SetSize(w, NoteHeigth)
					animation[i].SetPoint(x, HitPosition, draws.LeftBottom)
				}
				skin.NoteSprites[k][_type] = animation
			}
			x += w
		}
		{
			src := ebiten.NewImage(int(wsum), screenSizeY)
			src.Fill(color.NRGBA{0, 0, 0, uint8(255 * FieldDarkness)})
			s := draws.NewSpriteFromImage(src)
			s.SetPoint(FieldPosition, 0, draws.CenterTop)
			skin.FieldSprite = s
		}
		{
			s := draws.NewSpriteFromImage(hintImage)
			s.SetSize(wsum, HintHeight)
			s.SetPoint(FieldPosition, HitPosition-HintHeight, draws.CenterTop)
			skin.HintSprite = s
		}
		{
			src := ebiten.NewImage(int(wsum), 1)
			src.Fill(color.NRGBA{255, 255, 255, 255}) // White
			s := draws.NewSpriteFromImage(src)
			s.SetPoint(FieldPosition, HitPosition, draws.CenterBottom)
			skin.BarSprite = s
		}
		Skins[keyCount] = skin
	}
}
