package drum

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
	"golang.org/x/image/draw"
)

// [2] that most sprites have.
const (
	NormalNote = iota
	BigNote
)

// Significant keyword goes ahead, just as number is written: Left
const (
	KeyLeftKat = iota
	KeyLeftDon
	KeyRightDon
	KeyRightKat
)
const (
	Ground = iota
	Overlay1
	Overlay2
)
const (
	ShakeNote = iota
	ShakeBottom
	ShakeTop
)
const (
	DancerIdle = iota
	DancerYes
	DancerNo
	DancerHigh
)

var DefaultSkin Skin

// https://osu.ppy.sh/wiki/en/Skinning/osu%21taiko
type Skin struct {
	FieldSprite draws.Sprite
	HintSprite  draws.Sprite

	// First [2] are for big notes.
	JudgmentSprites    [2][3]draws.Sprite
	RedSprites         [2]draws.Sprite
	BlueSprites        [2]draws.Sprite
	NoteOverlaySprites [2][2]draws.Sprite
	HeadSprites        [2]draws.Sprite
	TailSprites        [2]draws.Sprite
	BodySprites        [2]draws.Sprite

	RollTickSprites draws.Sprite
	ShakeSprites    [3]draws.Sprite
	RollDotSprite   draws.Sprite
	BarSprite       draws.Sprite

	KeyFieldSprite draws.Sprite
	KeySprites     [4]draws.Sprite
	DancerSprites  [4][]draws.Sprite

	ScoreSprites          [10]draws.Sprite
	ComboSprites          [10]draws.Sprite
	RollTickComboSprites  [10]draws.Sprite
	ShakeCountdownSprites [10]draws.Sprite
}

var (
	ColorDon  = color.NRGBA{235, 69, 44, 255}
	ColorKat  = color.NRGBA{68, 141, 171, 255}
	ColorRoll = color.NRGBA{252, 83, 6, 255}
)

// func IsKeyImageFlipped(keyType int) bool {
// 	return keyType == KeyLeftKat || keyType == KeyRightDon
// }

// Todo: embed default skins to code for preventing panic when files are missing
func LoadSkin() {
	var skin Skin
	skin.ScoreSprites = make([]draws.Sprite, 10)
	skin.ComboSprites = make([]draws.Sprite, 10)
	skin.TickComboSprites = make([]draws.Sprite, 10)
	for i := 0; i < 10; i++ {
		s := draws.Sprite{
			I:      draws.NewImage(fmt.Sprintf("skin/score/%d.png", i)),
			Filter: ebiten.FilterLinear,
		}
		s.ApplyScale(ScoreScale)
		// x is not fixed.
		// y is always 0.
		skin.ScoreSprites[i] = s
	}
	for i := 0; i < 10; i++ {
		s := draws.Sprite{
			I:      draws.NewImage(fmt.Sprintf("skin/combo/%d.png", i)),
			Filter: ebiten.FilterLinear,
		}
		s.ApplyScale(ComboScale)
		// x is not fixed.
		s.SetCenterY(FieldPosition)
		skin.ComboSprites[i] = s

		s2 := draws.Sprite{
			I:      s.I,
			Filter: ebiten.FilterLinear,
		}
		s2.ApplyScale(RollTickComboScale)
		// x is not fixed.
		s2.SetCenterY(FieldPosition)
		skin.TickComboSprites[i] = s2
	}

	// Key sprites should be calculated first.
	// Sum of widths of Key sprites is essential to most of other sprites' X position.
	// Or, maybe not.
	leftDon := draws.NewImage("skin/drum/key/in.png")
	rightDon := draws.FlipX(leftDon)
	rightKat := draws.NewImage("skin/drum/key/out.png")
	leftKat := draws.FlipX(rightKat)
	x := 0.0
	for i, img := range []*ebiten.Image{leftDon, rightKat, leftKat, rightDon} {
		s := draws.Sprite{
			I: img,
			X: x,
		}
		s.SetHeight(FieldHeight)
		s.SetCenterY(FieldPosition)
		skin.KeySprites[[]int{KeyLeftDon, KeyRightKat,
			KeyLeftKat, KeyRightDon}[i]] = s
		x += s.W
		if i == 0 {
			comboPosition = s.W
		}
		if i == 1 { // Each side's Don and Kat are overlapped.
			x = 0
		}
	}

	skin.FieldSprite = draws.Sprite{
		I: draws.NewImage("skin/drum/field.png"),
		W: screenSizeX - x,
		H: FieldHeight,
		X: x,
		// Y is centered.
	}
	skin.FieldSprite.SetCenterY(FieldPosition)

	skin.HintSprite = draws.Sprite{
		I: draws.NewImage("skin/drum/hint.png"),
		H: NormalNoteHeight,
	}
	// skin.HintSprite.SetHeight(NormalNoteHeight)
	skin.HintSprite.ApplyScale(skin.HintSprite.ScaleH())
	skin.HintSprite.SetCenterX(HitPosition)
	skin.HintSprite.SetCenterY(FieldPosition)

	barLine := ebiten.NewImage(1, int(FieldInnerHeight))
	barLine.Fill(color.RGBA{255, 255, 255, 255})
	skin.BarSprite = draws.Sprite{
		I: barLine,
		W: 1,
		H: FieldInnerHeight,
	}
	// X is not fixed.
	skin.BarSprite.SetCenterY(FieldPosition)

	for size, sizeWord := range []string{"normal", "big"} {
		height := []float64{NormalNoteHeight, BigNoteHeight}[size]
		for judge, judgeWord := range []string{"cool", "good"} {
			path := fmt.Sprintf("skin/drum/judgment/%s/%s.png", sizeWord, judgeWord)
			s := draws.Sprite{
				I: draws.NewImage(path),
			}
			s.SetHeight(FieldHeight)
			s.ApplyScale(JudgmentScale)
			s.SetCenterX(HitPosition)
			s.SetCenterY(FieldPosition)
			skin.JudgmentSprites[size][judge] = s
		}
		{
			const miss = 2
			s := draws.Sprite{
				I: draws.NewImage("skin/drum/judgment/miss.png"),
			}
			// s.SetHeight(FieldHeight)
			s.ApplyScale(JudgmentScale)
			s.SetCenterX(HitPosition)
			s.SetCenterY(FieldPosition)
			skin.JudgmentSprites[size][miss] = s
		}

		for noteType, clr := range []color.NRGBA{ColorDon, ColorKat} {
			img := draws.NewImage("skin/drum/note/normal/note.png")
			img = draws.ApplyColor(img, clr)
			s := draws.Sprite{I: img, H: height}
			// s.SetHeight(height)
			s.ApplyScale(s.ScaleH())
			s.SetCenterY(FieldPosition)
			if noteType == 0 {
				skin.DonSprites[size][0] = s
			} else {
				skin.KatSprites[size][0] = s
			}
		}
		for noteType, clr := range []color.NRGBA{ColorDon, ColorKat} {
			img := draws.NewImage(fmt.Sprintf("skin/drum/note/%s/note.png", sizeWord))
			img = draws.ApplyColor(img, clr)
			s := draws.Sprite{I: img, H: height}
			// s.SetHeight(height)
			s.ApplyScale(s.ScaleH())
			s.SetCenterY(FieldPosition)
			if noteType == 0 {
				skin.DonSprites[size][0] = s
			} else {
				skin.KatSprites[size][0] = s
			}
		}

		overlayPath := fmt.Sprintf("skin/drum/note/%s/overlay", sizeWord)
		if ok, err := IsDir(overlayPath); err != nil {
			// fmt.Printf("loading %s's overlay occurrs an err: %s\n", sizeWord, err)
			fmt.Printf("%s's overlay has one frame.\n", sizeWord)
			continue
		} else if ok { // 2 frames.
			for i := 0; i < 2; i++ {
				overlayPath += fmt.Sprintf("/%d.png", i)
				s := draws.Sprite{I: draws.NewImage(overlayPath), H: height}
				// s.SetHeight(height)
				s.ApplyScale(s.ScaleH())
				s.SetCenterY(FieldPosition)
				skin.DonSprites[size][i+1] = s
			}
		} else { // 1 frame. Copy 1st frame to 2nd frame.
			overlayPath += ".png"
			s := draws.Sprite{I: draws.NewImage(overlayPath), H: height}
			// s.SetHeight(height)
			s.ApplyScale(s.ScaleH())
			s.SetCenterY(FieldPosition)
			skin.DonSprites[size][1] = s
			skin.DonSprites[size][2] = s // Copy the same one.
		}

		end := draws.NewImage("skin/drum/note/roll/end.png")
		tailSrc := draws.ApplyColor(end, ColorRoll)
		headSrc := draws.FlipX(tailSrc)

		head := draws.Sprite{I: headSrc}
		head.SetHeight(height)
		head.SetCenterY(FieldPosition)
		skin.HeadSprites[size] = head

		tail := draws.Sprite{I: tailSrc}
		tail.SetHeight(height)
		tail.SetCenterY(FieldPosition)
		skin.TailSprites[size] = tail
		{
			src := draws.NewImageSrc("skin/drum/note/roll/mid.png")
			dst := image.NewRGBA(image.Rect(0, 0, screenSizeX, int(height)))
			draw.BiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
			for pow := 0; pow < int(math.Log2(screenSizeX))+1; pow++ {
				w := 1 << pow
				rect := image.Rect(0, 0, w, int(height))
				s := draws.Sprite{
					I: ebiten.NewImageFromImage(dst.SubImage(rect)),
					W: float64(w),
					H: float64(height),
					// X is not fixed.
				}
				s.SetCenterY(FieldPosition)
				skin.BodySprites[size] = append(skin.BodySprites[size], s)
			}
		}
	}
	dot := draws.Sprite{
		I: draws.NewImage("skin/drum/note/roll/dot.png"),
	}
	dot.ApplyScale(RollDotScale)
	dot.SetCenterY(FieldPosition)
	skin.RollDotSprite = dot

	for i, name := range []string{"note", "bottom", "top"} {
		s := draws.Sprite{
			I: draws.NewImage(fmt.Sprintf("skin/drum/note/shake/%s.png", name)),
		}
		if i == ShakeNote {
			s.SetHeight(BigNoteHeight)
		} else {
			s.ApplyScale(ShakeScale)
		}
		s.SetCenterY(FieldPosition)
		skin.ShakeSprites[i] = s
	}

	DefaultSkin = skin
}
func IsDir(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}
