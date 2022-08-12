package gosu

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/draw"
)

// Sprites that are independent of key count.
var GeneralSkin *GeneralSkinStruct

type GeneralSkinStruct struct { // Singleton
	DefaultBackgrounds []Sprite
	ComboSprites       []Sprite
	ScoreSprites       []Sprite
	JudgmentSprites    []Sprite
	CursorSprites      [2]Sprite // 0: cursor // 1: additive cursor
	// CursorTailSprite   Sprite
}

// Todo: should each skin has own skin settings?
// Todo: BarLine color settings
type Skin struct {
	*GeneralSkinStruct
	NoteSprites   []Sprite
	BodySprites   [][]Sprite // Binary-building method
	HeadSprites   []Sprite
	TailSprites   []Sprite
	FieldSprite   Sprite
	HintSprite    Sprite
	BarLineSprite Sprite // Seperator of each bar (aka measure)
}

var SkinMap = make(map[int]Skin)
var RandomDefaultBackground Sprite

func LoadSkin() {
	g := &GeneralSkinStruct{
		DefaultBackgrounds: make([]Sprite, 0, 10),
		ComboSprites:       make([]Sprite, 10),
		ScoreSprites:       make([]Sprite, 10),
		JudgmentSprites:    make([]Sprite, 5),
	}
	{
		fs, err := os.ReadDir("skin/bg")
		if err != nil {
			panic(err)
		}
		for _, f := range fs {
			if f.IsDir() || !strings.HasPrefix(f.Name(), "bg") {
				continue
			}
			sprite := Sprite{
				I: NewImage(filepath.Join("skin/bg", f.Name())),
			}
			sprite.SetFullscreen()
			g.DefaultBackgrounds = append(g.DefaultBackgrounds, sprite)
		}
		r := int(rand.Float64() * float64(len(g.DefaultBackgrounds)))
		RandomDefaultBackground = g.DefaultBackgrounds[r]
	}
	for i := 0; i < 10; i++ {
		s := Sprite{
			I: NewImage(fmt.Sprintf("skin/combo/%d.png", i)),
		}
		s.ApplyScale(ComboScale)
		// ComboSprite's x value is not fixed.
		s.SetCenterXY(0, ComboPosition)
		g.ComboSprites[i] = s
	}
	for i := 0; i < 10; i++ {
		s := Sprite{
			I: NewImage(fmt.Sprintf("skin/score/%d.png", i)),
		}
		s.ApplyScale(ScoreScale)
		// ScoreSprite's x value is not fixed.
		// ScoreSprite's y value is always 0.
		g.ScoreSprites[i] = s
	}
	for i, name := range []string{"kool", "cool", "good", "bad", "miss"} {
		s := Sprite{
			I: NewImage(fmt.Sprintf("skin/judgment/%s.png", name)),
		}
		s.ApplyScale(JudgmentScale)
		s.SetCenterXY(screenSizeX/2, JudgmentPosition)
		g.JudgmentSprites[i] = s
	}
	for i, name := range []string{"menu-cursor.png", "menu-cursor-additive.png"} {
		s := Sprite{
			I: NewImage(fmt.Sprintf("skin/cursor/%s", name)),
		}
		s.ApplyScale(CursorScale)
		g.CursorSprites[i] = s
	}
	GeneralSkin = g

	// Following sprites are dependent of key count.
	// Todo: Key 1 ~ 3, scratch
	// Todo: 4th note image
	for keyCount := 4; keyCount <= 10; keyCount++ {
		noteKinds := NoteKindsMap[keyCount]
		noteWidths := NoteWidthsMap[keyCount]
		s := Skin{
			GeneralSkinStruct: GeneralSkin,
			NoteSprites:       make([]Sprite, keyCount&ScratchMask),
			// BodySprites:       make([]Sprite, keyCount&ScratchMask),
			BodySprites: make([][]Sprite, keyCount&ScratchMask),
			HeadSprites: make([]Sprite, keyCount&ScratchMask),
			TailSprites: make([]Sprite, keyCount&ScratchMask),
		}
		var wsum int
		for k, kind := range noteKinds {
			s.NoteSprites[k] = Sprite{
				I: NewImage("skin/note/" + fmt.Sprintf("n%d.png", []int{1, 2, 3, 3}[kind])),
				W: noteWidths[kind],
				H: NoteHeigth,
			}
			// Each w should be integer, since it is actual sprite's width.
			wsum += int(noteWidths[kind])
		}
		// NoteSprite's x value should be integer as well as w.
		// Todo: Scratch should be excluded to width sum.
		x := screenSizeX/2 - wsum/2
		for k, kind := range noteKinds {
			s.NoteSprites[k].X = float64(x)
			// NoteSprites's y value is not fixed.
			x += int(noteWidths[kind])
		}
		// Draw max length of long note body sprite in advance.
		// Todo: change l%d.png to b%d.png (with animation)
		x = screenSizeX/2 - wsum/2
		for k, kind := range noteKinds {
			f, err := os.Open("skin/note/" + fmt.Sprintf("l%d.png", []int{1, 2, 3, 3}[kind]))
			if err != nil {
				panic(err)
			}
			defer f.Close()
			src, _, err := image.Decode(f)
			if err != nil {
				panic(err)
			}
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
		copy(s.HeadSprites, s.NoteSprites)
		copy(s.TailSprites, s.NoteSprites)
		field := ebiten.NewImage(wsum, screenSizeY)
		field.Fill(color.RGBA{0, 0, 0, uint8(255 * FieldDark)})
		s.FieldSprite = Sprite{
			I: field,
			W: float64(wsum),
			H: screenSizeY,
			X: float64(screenSizeX)/2 - float64(wsum)/2,
			Y: 0,
		}
		s.HintSprite = Sprite{
			I: NewImage("skin/play/hint.png"),
			W: float64(wsum),
			H: HintHeight,
		}
		barLine := ebiten.NewImage(wsum, 1)
		barLine.Fill(color.RGBA{255, 255, 255, 255})
		s.BarLineSprite = Sprite{
			I: barLine,
			W: float64(wsum),
			H: 1,
			X: float64(screenSizeX)/2 - float64(wsum)/2,
			Y: 0,
		}
		s.HintSprite.SetCenterXY(screenSizeX/2, HintPosition)
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
