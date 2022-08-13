package gosu

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"

	_ "image/jpeg"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/draw"
)

// Sprites that are independent of key count.
var GeneralSkin *GeneralSkinStruct

type GeneralSkinStruct struct { // Singleton
	// DefaultBackgrounds []Sprite
	DefaultBackground Sprite
	ComboSprites      []Sprite
	ScoreSprites      []Sprite
	JudgmentSprites   []Sprite
	CursorSprites     [2]Sprite // 0: cursor // 1: additive cursor
	// CursorTailSprite   Sprite
	TimingMeterSprite       Sprite
	TimingMeterUnitSprite   Sprite
	TimingMeterAnchorSprite Sprite
}

var (
	dark = color.NRGBA{0, 0, 0, 128}

	white  = color.NRGBA{255, 255, 255, 192}
	red    = color.NRGBA{255, 0, 0, 128}
	purple = color.NRGBA{213, 0, 242, 128}

	gray   = color.NRGBA{109, 120, 134, 255}
	yellow = color.NRGBA{244, 177, 0, 255}
	lime   = color.NRGBA{51, 255, 40, 255}
	sky    = color.NRGBA{85, 251, 255, 255}
	blue   = color.NRGBA{0, 170, 242, 255}
)

// Todo: should each skin has own skin settings?
// Todo: BarLine color settings
type Skin struct {
	*GeneralSkinStruct
	KeyUpSprites   []Sprite
	KeyDownSprites []Sprite
	NoteSprites    []Sprite
	HeadSprites    []Sprite
	TailSprites    []Sprite
	BodySprites    [][]Sprite // Binary-building method

	FieldSprite   Sprite
	HintSprite    Sprite
	BarLineSprite Sprite // Seperator of each bar (aka measure)
}

var SkinMap = make(map[int]Skin)

func LoadSkin() {
	g := &GeneralSkinStruct{
		// DefaultBackgrounds: make([]Sprite, 0, 10),
		ComboSprites:    make([]Sprite, 10),
		ScoreSprites:    make([]Sprite, 10),
		JudgmentSprites: make([]Sprite, 5),
	}
	g.DefaultBackground = Sprite{
		I:      NewImage("skin/default-bg.jpg"),
		Filter: ebiten.FilterLinear,
	}
	g.DefaultBackground.SetWidth(screenSizeX)
	g.DefaultBackground.SetCenterY(screenSizeY / 2)
	for i := 0; i < 10; i++ {
		s := Sprite{
			I:      NewImage(fmt.Sprintf("skin/combo/%d.png", i)),
			Filter: ebiten.FilterLinear,
		}
		s.ApplyScale(ComboScale)
		// ComboSprite's x value is not fixed.
		s.SetCenterY(ComboPosition)
		g.ComboSprites[i] = s
	}
	for i := 0; i < 10; i++ {
		s := Sprite{
			I:      NewImage(fmt.Sprintf("skin/score/%d.png", i)),
			Filter: ebiten.FilterLinear,
		}
		s.ApplyScale(ScoreScale)
		// ScoreSprite's x value is not fixed.
		// ScoreSprite's y value is always 0.
		g.ScoreSprites[i] = s
	}
	for i, name := range []string{"kool", "cool", "good", "bad", "miss"} {
		s := Sprite{
			I:      NewImage(fmt.Sprintf("skin/judgment/%s.png", name)),
			Filter: ebiten.FilterLinear,
		}
		s.ApplyScale(JudgmentScale)
		s.SetCenterX(screenSizeX / 2)
		s.SetCenterY(JudgmentPosition)
		g.JudgmentSprites[i] = s
	}
	for i, name := range []string{"menu-cursor.png", "menu-cursor-additive.png"} {
		s := Sprite{
			I:      NewImage(fmt.Sprintf("skin/cursor/%s", name)),
			Filter: ebiten.FilterLinear,
		}
		s.ApplyScale(CursorScale)
		g.CursorSprites[i] = s
	}
	{ // Timing meter. the height of colored rectangle is 1/4 of meter's.
		meterW := 1 + 2*int(TimingMeterWidth)*int(Miss.Window)
		meterH := int(TimingMeterHeight)
		meter := image.NewRGBA(image.Rect(0, 0, meterW, meterH))
		draw.Draw(meter, meter.Bounds(), &image.Uniform{dark}, image.Point{}, draw.Src)
		y1, y2 := int(float64(meterH)*0.375), int(float64(meterH)*0.625)
		for i, color := range []color.NRGBA{gray, yellow, lime, sky, blue} {
			j := Judgments[4-i]
			w := 1 + 2*int(TimingMeterWidth)*int(j.Window)
			x1 := int(TimingMeterWidth) * (int(Miss.Window - j.Window))
			x2 := x1 + w
			rect := image.Rect(x1, y1, x2, y2)
			draw.Draw(meter, rect, &image.Uniform{color}, image.Point{}, draw.Src)
		}
		// rect := image.Rect(meterW/2-int(TimingMeterWidth/2), 0, meterW/2+int(TimingMeterWidth/2), meterH)
		// draw.Draw(meter, rect, &image.Uniform{color.Black}, image.Point{}, draw.Src) // Todo: need a check draw.Src
		g.TimingMeterSprite = Sprite{
			I: ebiten.NewImageFromImage(meter),
			W: float64(meterW),
			H: float64(meterH),
			Y: screenSizeY - TimingMeterHeight,
		}
		g.TimingMeterSprite.SetCenterX(screenSizeX / 2)

		// Draw middle anchor of timing meter.
		anchor := ebiten.NewImage(int(TimingMeterWidth), int(TimingMeterHeight))
		anchor.Fill(red)
		g.TimingMeterAnchorSprite = Sprite{
			I: anchor,
			W: TimingMeterWidth,
			H: TimingMeterHeight,
			Y: screenSizeY - TimingMeterHeight,
		}
		g.TimingMeterAnchorSprite.SetCenterX(screenSizeX / 2)

		unit := ebiten.NewImage(int(TimingMeterWidth), int(TimingMeterHeight))
		unit.Fill(white)
		g.TimingMeterUnitSprite = Sprite{
			I: unit,
			W: TimingMeterWidth,
			H: TimingMeterHeight,
			// TimingMeterUnitSprite's x value is not fixed.
			Y: screenSizeY - TimingMeterHeight,
		}
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
		noteImages[i] = NewImage(fmt.Sprintf("skin/note/note/%d.png", kind))
		headImages[i] = NewImage(fmt.Sprintf("skin/note/head/%d.png", kind))
		tailImages[i] = NewImage(fmt.Sprintf("skin/note/tail/%d.png", kind))
		// bodyImages[i] = NewImage(fmt.Sprintf("skin/note/body/%d.png", kind))
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
	keyUpImage = NewImage("skin/key/up.png")
	keyDownImage = NewImage("skin/key/down.png")
	hintImage = NewImage("skin/hint.png")

	// Todo: Key count 1~3, KeyCount + scratch
	for keyCount := 4; keyCount <= 10; keyCount++ {
		noteKinds := NoteKindsMap[keyCount]
		noteWidths := NoteWidthsMap[keyCount]
		s := Skin{
			GeneralSkinStruct: GeneralSkin,
			KeyUpSprites:      make([]Sprite, keyCount&ScratchMask),
			KeyDownSprites:    make([]Sprite, keyCount&ScratchMask),
			NoteSprites:       make([]Sprite, keyCount&ScratchMask),
			HeadSprites:       make([]Sprite, keyCount&ScratchMask),
			TailSprites:       make([]Sprite, keyCount&ScratchMask),
			BodySprites:       make([][]Sprite, keyCount&ScratchMask),
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
			s.KeyUpSprites[k] = Sprite{
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
			s.NoteSprites[k] = Sprite{
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
				s.BodySprites[k] = append(s.BodySprites[k], Sprite{
					I: ebiten.NewImageFromImage(dst.SubImage(rect)),
					W: float64(w),
					H: float64(h),
					X: float64(x),
					// BodySprites's y value is not fixed.
				})
			}
			x += int(noteWidths[kind])
		}

		field := ebiten.NewImage(wsum, screenSizeY)
		field.Fill(color.RGBA{0, 0, 0, uint8(255 * FieldDark)})
		s.FieldSprite = Sprite{
			I: field,
			W: float64(wsum),
			H: screenSizeY,
		}
		s.FieldSprite.SetCenterX(screenSizeX / 2)
		// FieldSprite's y value is always 0.

		s.HintSprite = Sprite{
			I: hintImage,
			W: float64(wsum),
			H: HintHeight,
		}
		s.HintSprite.SetCenterX(screenSizeX / 2)
		s.HintSprite.Y = HitPosition - HintHeight

		barLine := ebiten.NewImage(wsum, 1)
		barLine.Fill(color.RGBA{255, 255, 255, 255})
		s.BarLineSprite = Sprite{
			I: barLine,
			W: float64(wsum),
			H: 1,
		}
		s.BarLineSprite.SetCenterX(screenSizeX / 2)
		// BarLineSprite's y value is not fixed.
		SkinMap[keyCount] = s
	}
}

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

// NewImage returns nil when fails to load image from the path.
func NewImage(path string) *ebiten.Image {
	f, err := os.Open(path)
	if err != nil {
		return nil
		// panic(err)
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	if err != nil {
		return nil
		// panic(err)
	}
	return ebiten.NewImageFromImage(i)
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

// May be useful for animation
// {
// 	fs, err := os.ReadDir("skin/bg")
// 	if err != nil {
// 		panic(err)
// 	}
// 	for _, f := range fs {
// 		if f.IsDir() || !strings.HasPrefix(f.Name(), "bg") {
// 			continue
// 		}
// 		sprite := Sprite{
// 			I: NewImage(filepath.Join("skin/bg", f.Name())),
// 		}
// 		sprite.SetFullscreen()
// 		g.DefaultBackgrounds = append(g.DefaultBackgrounds, sprite)
// 	}
// 	r := int(rand.Float64() * float64(len(g.DefaultBackgrounds)))
// 	RandomDefaultBackground = g.DefaultBackgrounds[r]
// }
