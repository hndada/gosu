package gosu

import (
	"fmt"
	"image"
	"image/color"
	"os"

	_ "image/jpeg"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

type Skin struct {
	DefaultBackground Sprite
	ComboSprites      []Sprite
	ScoreSprites      []Sprite
	JudgmentSprites   []Sprite

	NoteSprites []Sprite
	BodySprites []Sprite
	HeadSprites []Sprite
	TailSprites []Sprite
	FieldSprite Sprite
	HintSprite  Sprite
}

var SkinMap = make(map[int]Skin)

const DefaultBackgroundPath = "skin/bg.jpg"

func LoadSkin() {
	// Following sprites are independent of key count.
	var (
		defaultBackground Sprite
		comboSprites      []Sprite = make([]Sprite, 10)
		scoreSprites      []Sprite = make([]Sprite, 10)
		judgmentSprites   []Sprite = make([]Sprite, 5)
	)
	defaultBackground.I = NewImage(DefaultBackgroundPath)
	defaultBackground.W = screenSizeX
	defaultBackground.H = screenSizeY
	for i := 0; i < 10; i++ {
		s := Sprite{
			I: NewImage(fmt.Sprintf("skin/combo/%d.png", i)),
		}
		s.ApplyScale(ComboScale)
		// ComboSprite's x value is not fixed.
		s.SetCenterXY(0, ComboPosition)
		comboSprites[i] = s
	}
	for i := 0; i < 10; i++ {
		s := Sprite{
			I: NewImage(fmt.Sprintf("skin/score/%d.png", i)),
		}
		s.ApplyScale(ScoreScale)
		// ScoreSprite's x value is not fixed.
		// ScoreSprite's y value is always 0.
		scoreSprites[i] = s
	}
	for i, name := range []string{"kool", "cool", "good", "bad", "miss"} {
		s := Sprite{
			I: NewImage(fmt.Sprintf("skin/judgment/%s.png", name)),
		}
		s.ApplyScale(JudgmentScale)
		s.SetCenterXY(screenSizeX/2, JudgmentPosition)
		judgmentSprites[i] = s
	}

	// Following sprites are dependent of key count.
	// Todo: Key 1 ~ 3, scratch
	for keyCount := 4; keyCount <= 10; keyCount++ {
		s := Skin{
			DefaultBackground: defaultBackground,
			ComboSprites:      comboSprites,
			ScoreSprites:      scoreSprites,
			JudgmentSprites:   judgmentSprites,

			NoteSprites: make([]Sprite, keyCount&ScratchMask),
			BodySprites: make([]Sprite, keyCount&ScratchMask),
			HeadSprites: make([]Sprite, keyCount&ScratchMask),
			TailSprites: make([]Sprite, keyCount&ScratchMask),
		}
		// Todo: 4th note image
		var wsum int
		for k, kind := range NoteKindsMap[keyCount] {
			s.NoteSprites[k] = Sprite{
				I: NewImage("skin/note/" + fmt.Sprintf("n%d.png", []int{1, 2, 3, 3}[kind])),
				W: NoteWidthsMap[keyCount][kind],
				H: NoteHeigth,
			}
			// Each w should be integer, since it is actual sprite's width.
			wsum += int(NoteWidthsMap[keyCount][kind])
		}
		// NoteSprite's x value should be integer as well as w.
		// Todo: Scratch should be excluded to width sum.
		x := (screenSizeX - wsum) / 2
		for k, kind := range NoteKindsMap[keyCount] {
			s.NoteSprites[k].X = float64(x)
			// NoteSprites's y value is not fixed.
			x += int(NoteWidthsMap[keyCount][kind])
		}
		x = (screenSizeX - wsum) / 2
		for k, kind := range NoteKindsMap[keyCount] {
			s.BodySprites[k] = Sprite{
				I: NewImage("skin/note/" + fmt.Sprintf("l%d.png", []int{1, 2, 3, 3}[kind])),
				W: NoteWidthsMap[keyCount][kind],
				H: NoteHeigth, // Fyi, long note body's height doesn't need to be scaled.
			}
			s.BodySprites[k].X = float64(x)
			// BodySprites's y value is not fixed.
			x += int(NoteWidthsMap[keyCount][kind])
		}
		copy(s.HeadSprites, s.NoteSprites)
		copy(s.TailSprites, s.NoteSprites)
		field := ebiten.NewImage(wsum, screenSizeY)
		field.Fill(color.RGBA{0, 0, 0, uint8(255 * FieldDark)})
		s.FieldSprite = Sprite{
			I: field,
			W: float64(wsum),
			H: screenSizeY,
			X: float64(screenSizeX-wsum) / 2,
			Y: 0,
		}
		s.HintSprite = Sprite{
			I: NewImage("skin/play/hint.png"),
			W: float64(wsum),
			H: HintHeight,
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

func NewImage(path string) *ebiten.Image {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	if err != nil {
		panic(err)
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
