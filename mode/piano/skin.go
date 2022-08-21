package piano

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"

	_ "image/jpeg"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
	"golang.org/x/image/draw"
)

// var colors = []color.NRGBA{gosu.Gray, gosu.Yellow, gosu.Lime, gosu.Sky, gosu.Blue}

type NoteKind int

const (
	One NoteKind = iota
	Two
	Mid
	Tip
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

// Sprites that are independent of key count.
var GeneralSkin *GeneralSkinStruct

type GeneralSkinStruct struct { // Singleton
	ComboSprites    []draws.Sprite
	ScoreSprites    []draws.Sprite
	JudgmentSprites []draws.Sprite
}

// Todo: should each skin has own skin settings?
type Skin struct {
	*GeneralSkinStruct
	KeyUpSprites   []draws.Sprite
	KeyDownSprites []draws.Sprite
	NoteSprites    []draws.Sprite
	HeadSprites    []draws.Sprite
	TailSprites    []draws.Sprite
	BodySprites    [][]draws.Sprite // Binary-building method

	FieldSprite   draws.Sprite
	HintSprite    draws.Sprite
	BarLineSprite draws.Sprite // Seperator of each bar (aka measure)

	// BodySpritesTest []draws.Sprite
}

var SkinMap = make(map[int]Skin)

func LoadSkin() {
	g := &GeneralSkinStruct{
		// DefaultBackgrounds: make([]draws.Sprite, 0, 10),
		ComboSprites:    make([]draws.Sprite, 10),
		ScoreSprites:    make([]draws.Sprite, 10),
		JudgmentSprites: make([]draws.Sprite, 5),
	}

	for i := 0; i < 10; i++ {
		s := draws.Sprite{
			I:      draws.NewImage(fmt.Sprintf("skin/combo/%d.png", i)),
			Filter: ebiten.FilterLinear,
		}
		s.ApplyScale(ComboScale)
		// ComboSprite's x value is not fixed.
		s.SetCenterY(ComboPosition)
		g.ComboSprites[i] = s
	}
	for i := 0; i < 10; i++ {
		s := draws.Sprite{
			I:      draws.NewImage(fmt.Sprintf("skin/score/%d.png", i)),
			Filter: ebiten.FilterLinear,
		}
		s.ApplyScale(ScoreScale)
		// ScoreSprite's x value is not fixed.
		// ScoreSprite's y value is always 0.
		g.ScoreSprites[i] = s
	}
	for i, name := range []string{"kool", "cool", "good", "bad", "miss"} {
		s := draws.Sprite{
			I:      draws.NewImage(fmt.Sprintf("skin/judgment/%s.png", name)),
			Filter: ebiten.FilterLinear,
		}
		s.ApplyScale(JudgmentScale)
		s.SetCenterX(screenSizeX / 2)
		s.SetCenterY(JudgmentPosition)
		g.JudgmentSprites[i] = s
	}

	GeneralSkin = g

	// Following sprites are dependent of key count.
	// Todo: animated sprite support. Starting with [4][]*ebiten.Image will help.
	var (
		keyUpImage   *ebiten.Image
		keyDownImage *ebiten.Image
		noteImages   [4]*ebiten.Image
		headImages   [4]*ebiten.Image
		tailImages   [4]*ebiten.Image
		bodyImages   [4]image.Image // binary-building method
		// bodyImages   [4]*ebiten.Image
		hintImage *ebiten.Image
	)
	// Todo: 4th note image
	// Currently head and tail use note's image.
	for i, kind := range []int{1, 2, 3, 3} {
		noteImages[i] = draws.NewImage(fmt.Sprintf("skin/note/note/%d.png", kind))
		headImages[i] = draws.NewImage(fmt.Sprintf("skin/note/head/%d.png", kind))
		tailImages[i] = draws.NewImage(fmt.Sprintf("skin/note/tail/%d.png", kind))
		// bodyImages[i] = draws.NewImage(fmt.Sprintf("skin/note/body/%d.png", kind))
		{
			f, err := os.Open(fmt.Sprintf("skin/note/body/%d.png", kind))
			if err != nil {
				panic(err)
			}
			defer f.Close()
			src, _, err := image.Decode(f)
			if err != nil {
				panic(err)
			}
			bodyImages[i] = src
		}
	}
	keyUpImage = draws.NewImage("skin/key/up.png")
	keyDownImage = draws.NewImage("skin/key/down.png")
	hintImage = draws.NewImage("skin/hint.png")

	// Todo: Key count 1~3, KeyCount + scratch
	for keyCount := 4; keyCount <= 10; keyCount++ {
		noteKinds := NoteKindsMap[keyCount]
		noteWidths := NoteWidthsMap[keyCount]
		s := Skin{
			GeneralSkinStruct: GeneralSkin,
			KeyUpSprites:      make([]draws.Sprite, keyCount&ScratchMask),
			KeyDownSprites:    make([]draws.Sprite, keyCount&ScratchMask),
			NoteSprites:       make([]draws.Sprite, keyCount&ScratchMask),
			HeadSprites:       make([]draws.Sprite, keyCount&ScratchMask),
			TailSprites:       make([]draws.Sprite, keyCount&ScratchMask),
			BodySprites:       make([][]draws.Sprite, keyCount&ScratchMask),
			// BodySpritesTest:   make([]draws.Sprite, keyCount&ScratchMask),
		}

		var wsum int
		for _, kind := range noteKinds {
			// Each w should be integer, since it is actual sprite's width.
			wsum += int(noteWidths[kind])
		}

		// Todo: Scratch should be excluded to width sum.
		// KeyUp and KeyDown are drawn below Hint, which bottom is along with HitPosition.
		x := screenSizeX/2 - wsum/2
		for k, kind := range noteKinds {
			s.KeyUpSprites[k] = draws.Sprite{
				I: keyUpImage,
				X: float64(x),
				Y: HitPosition,
			}
			s.KeyUpSprites[k].SetWidth(noteWidths[kind])
			s.KeyDownSprites[k] = s.KeyUpSprites[k]
			s.KeyDownSprites[k].I = keyDownImage
			s.KeyUpSprites[k].SetWidth(noteWidths[kind])
			x += int(noteWidths[kind])
		}

		x = screenSizeX/2 - wsum/2 // x should be integer like w as well.
		for k, kind := range noteKinds {
			s.NoteSprites[k] = draws.Sprite{
				I: noteImages[kind],
				W: noteWidths[kind],
				H: NoteHeigth,
				X: float64(x),
				// NoteSprites's y value is not fixed.
			}
			s.HeadSprites[k] = s.NoteSprites[k]
			s.HeadSprites[k].I = headImages[kind]
			s.TailSprites[k] = s.NoteSprites[k]
			s.TailSprites[k].I = tailImages[kind]
			x += int(noteWidths[kind])
		}

		// Draw max length of long note body sprite in advance.
		x = screenSizeX/2 - wsum/2
		for k, kind := range noteKinds {
			src := bodyImages[kind]

			w := int(noteWidths[kind])
			scale := float64(w) / float64(src.Bounds().Dx())
			h := int(scale * float64(src.Bounds().Dy()))
			dst := image.NewRGBA(image.Rect(0, 0, w, screenSizeY))
			switch BodySpriteStyle {
			case BodySpriteStyleStretch:
				draw.BiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
			case BodySpriteStyleAttach:
				for rect := image.Rect(0, 0, w, h); rect.Min.Y < dst.Bounds().Dy(); {
					draw.BiLinear.Scale(dst, rect, src, src.Bounds(), draw.Over, nil)
					rect.Min.Y += h
					rect.Max.Y += h
				}
			}
			for pow := 0; pow < int(math.Log2(screenSizeY))+1; pow++ {
				h := 1 << pow
				rect := image.Rect(0, 0, w, h)
				s.BodySprites[k] = append(s.BodySprites[k], draws.Sprite{
					I: ebiten.NewImageFromImage(dst.SubImage(rect)),
					W: float64(w),
					H: float64(h),
					X: float64(x),
					// BodySprites's y value is not fixed.
				})
			}

			// This is for test.
			// {
			// 	w := int(noteWidths[kind])
			// 	scale := float64(w) / float64(src.Bounds().Dx())
			// 	h := int(scale * float64(src.Bounds().Dy()))
			// 	dst := image.NewRGBA(image.Rect(0, 0, int(w), screenSizeY))
			// 	for rect := image.Rect(0, 0, w, h); rect.Min.Y < dst.Bounds().Dy(); {
			// 		draw.BiLinear.Scale(dst, rect, src, src.Bounds(), draw.Over, nil)
			// 		rect.Min.Y += h
			// 		rect.Max.Y += h
			// 	}
			// 	s.BodySpritesTest[k] = draws.Sprite{
			// 		I: ebiten.NewImageFromImage(dst),
			// 		W: float64(dst.Bounds().Dx()), // noteWidths[kind]
			// 		H: float64(dst.Bounds().Dy()), // screenSizeY
			// 	}
			// 	s.BodySpritesTest[k].X = float64(x)
			// 	// BodySprites's y value is not fixed.
			// }

			x += int(noteWidths[kind])
		}

		field := ebiten.NewImage(wsum, screenSizeY)
		field.Fill(color.RGBA{0, 0, 0, uint8(255 * FieldDark)})
		s.FieldSprite = draws.Sprite{
			I: field,
			W: float64(wsum),
			H: screenSizeY,
		}
		s.FieldSprite.SetCenterX(screenSizeX / 2)
		// FieldSprite's y value is always 0.

		s.HintSprite = draws.Sprite{
			I: hintImage,
			W: float64(wsum),
			H: HintHeight,
		}
		s.HintSprite.SetCenterX(screenSizeX / 2)
		s.HintSprite.Y = HitPosition - HintHeight

		barLine := ebiten.NewImage(wsum, 1)
		barLine.Fill(color.RGBA{255, 255, 255, 255})
		s.BarLineSprite = draws.Sprite{
			I: barLine,
			W: float64(wsum),
			H: 1,
		}
		s.BarLineSprite.SetCenterX(screenSizeX / 2)
		// BarLineSprite's y value is not fixed.
		SkinMap[keyCount] = s
	}
}
