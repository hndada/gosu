package drum

import (
	"fmt"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu"
	draws "github.com/hndada/gosu/draws2"
)

var (
	ColorRed    = color.NRGBA{235, 69, 44, 255}
	ColorBlue   = color.NRGBA{68, 141, 171, 255}
	ColorYellow = color.NRGBA{230, 170, 0, 255} // 252, 83, 6
	ColorPurple = color.NRGBA{150, 100, 200, 255}
)

const (
	LeftBlue = iota
	LeftRed
	RightRed
	RightBlue
)

const (
	DancerIdle = iota
	DancerYes
	DancerNo
	DancerHigh
)

var DefaultSkin Skin

// Order of fields of Skin is roughly consistent with drawing order.
// https://osu.ppy.sh/wiki/en/Skinning/osu%21taiko
type Skin struct {
	FieldSprite     draws.Sprite
	HintSprites     [2]draws.Sprite
	BarSprite       draws.Sprite
	JudgmentSprites [2][3]draws.Sprite

	NoteSprites       [2][4]draws.Sprite
	HeadSprites       [2]draws.Sprite
	TailSprites       [2]draws.Sprite
	OverlaySprites    [2][]draws.Sprite
	BodySprites       [2]draws.Sprite
	DotSprite         draws.Sprite
	ShakeBorderSprite draws.Sprite
	ShakeSprite       draws.Sprite

	KeySprites     [4]draws.Sprite
	KeyFieldSprite draws.Sprite
	DancerSprites  [4][]draws.Sprite
	ScoreSprites   [10]draws.Sprite
	ComboSprites   [10]draws.Sprite
}

// Todo: embed default skins to code for preventing panic when files are missing
func LoadSkin() {
	var skin Skin
	defer func() { DefaultSkin = skin }()
	var noteImage = draws.NewImage("skin/drum/note/note.png")
	{
		s := draws.NewSprite("skin/drum/field.png")
		s.SetSize(screenSizeX, FieldHeight)
		s.SetPoint(0, FieldPosition, draws.LeftMiddle)
		skin.FieldSprite = s
	}
	for i := range skin.HintSprites {
		const (
			padScale   = 1.1
			outerScale = 1.2
		)
		// sw, sh := noteImage.Size()
		srcSize := draws.IntPt(noteImage.Size())
		outer := draws.NewScaledImage(noteImage, outerScale)
		pad := draws.NewScaledImage(noteImage, padScale)
		inner := noteImage
		a := uint8(255 * FieldDarkness)
		img := ebiten.NewImage(outer.Size())
		{
			op := &ebiten.DrawImageOptions{}
			op.ColorM.ScaleWithColor(color.NRGBA{128, 128, 128, a})
			op.GeoM.Translate(0, 0)
			img.DrawImage(outer, op)
		}
		{
			op := &ebiten.DrawImageOptions{}
			op.ColorM.ScaleWithColor(color.NRGBA{255, 255, 0, a})
			if i == 0 { // Blank for idle, Yellow for highlight.
				op.CompositeMode = ebiten.CompositeModeDestinationOut
			}
			op.GeoM.Translate(srcSize.Mul(draws.Scalar((outerScale - padScale) / 2)).XY())
			// op.GeoM.Translate(0.05*float64(sw), 0.05*float64(sh))
			img.DrawImage(pad, op)
		}
		{
			op := &ebiten.DrawImageOptions{}
			op.ColorM.ScaleWithColor(color.NRGBA{60, 60, 60, a})
			op.GeoM.Translate(srcSize.Mul(draws.Scalar((outerScale - 1) / 2)).XY())
			// op.GeoM.Translate(0.1*float64(sw), 0.1*float64(sh))
			img.DrawImage(inner, op)
		}
		s := draws.NewSpriteFromImage(img)
		s.SetScaleToH(1.2 * regularNoteHeight)
		s.SetPoint(HitPosition, FieldPosition, draws.CenterMiddle)
		skin.HintSprites[i] = s
	}
	{
		src := ebiten.NewImage(1, int(FieldInnerHeight))
		src.Fill(color.NRGBA{255, 255, 255, 255})
		s := draws.NewSpriteFromImage(src)
		s.SetPoint(HitPosition, FieldPosition, draws.CenterMiddle)
		skin.BarSprite = s
	}
	var (
		rollEndImage = draws.NewImage("skin/drum/note/roll/end.png")
		rollMidImage = draws.NewImage("skin/drum/note/roll/mid.png")
	)
	for i, sname := range []string{"regular", "big"} {
		noteHeight := regularNoteHeight
		if i == Big {
			noteHeight = bigNoteHeight
		}
		for j, jname := range []string{"cool", "good", "miss"} {
			var path string
			if jname == "miss" {
				path = "skin/drum/judgment/miss.png"
			} else {
				path = fmt.Sprintf("skin/drum/judgment/%s/%s.png", sname, jname)
			}
			s := draws.NewSprite(path)
			s.SetScale(draws.Scalar(JudgmentScale))
			s.SetPoint(HitPosition, FieldPosition, draws.CenterMiddle)
			skin.JudgmentSprites[i][j] = s
		}
		for j, clr := range []color.NRGBA{ColorRed, ColorBlue, ColorYellow, ColorPurple} {
			img := ebiten.NewImage(noteImage.Size())
			op := &ebiten.DrawImageOptions{}
			op.ColorM.ScaleWithColor(clr)
			img.DrawImage(noteImage, op)

			s := draws.NewSpriteFromImage(img)
			s.SetScaleToH(noteHeight)
			s.SetPoint(HitPosition, FieldPosition, draws.CenterMiddle)
			skin.NoteSprites[i][j] = s
		}
		{
			s := draws.NewSpriteFromImage(rollEndImage)
			s.SetScaleToH(noteHeight)
			s.SetPoint(HitPosition, FieldPosition, draws.LeftMiddle)
			skin.TailSprites[i] = s
		}
		{
			s := draws.NewSpriteFromImage(draws.NewXFlippedImage(rollEndImage))
			s.SetScaleToH(noteHeight)
			s.SetPoint(HitPosition, FieldPosition, draws.RightMiddle)
			skin.HeadSprites[i] = s
		}
		{
			paths := gosu.Paths(fmt.Sprintf("skin/drum/note/overlay/%s", sname))
			skin.OverlaySprites[i] = make([]draws.Sprite, len(paths))
			for j, path := range paths {
				s := draws.NewSprite(path)
				s.SetScaleToH(noteHeight)
				s.SetPoint(HitPosition, FieldPosition, draws.CenterMiddle)
				skin.OverlaySprites[i][j] = s
			}
		}
		{
			s := draws.NewSpriteFromImage(rollMidImage)
			s.SetScaleToH(noteHeight)
			s.SetPoint(HitPosition, FieldPosition, draws.LeftMiddle)
			s.Filter = ebiten.FilterNearest
			skin.BodySprites[i] = s
		}
	}
	{
		s := draws.NewSprite("skin/drum/note/roll/dot.png")
		s.SetScale(draws.Scalar(DotScale))
		s.SetPoint(HitPosition, FieldPosition, draws.CenterMiddle)
		skin.DotSprite = s
	}
	{
		sw, sh := noteImage.Size()
		inner := draws.NewScaledImage(noteImage, 4)
		shake := ebiten.NewImage(inner.Size())
		{
			op := &ebiten.DrawImageOptions{}
			color := ColorPurple
			color.A = 128
			op.ColorM.ScaleWithColor(color)
			shake.DrawImage(inner, op)
		}
		{
			s := draws.NewSpriteFromImage(shake)
			s.SetScaleToH(4 * regularNoteHeight)
			s.SetPoint(HitPosition, FieldPosition, draws.CenterMiddle)
			skin.ShakeSprite = s
		}

		outer := draws.NewScaledImage(noteImage, 4.1)
		border := ebiten.NewImage(outer.Size())
		{
			op := &ebiten.DrawImageOptions{}
			op.ColorM.ScaleWithColor(color.NRGBA{255, 255, 255, 255})
			op.GeoM.Translate(0, 0)
			border.DrawImage(outer, op)
		}
		{
			op := &ebiten.DrawImageOptions{}
			op.ColorM.ScaleWithColor(color.NRGBA{255, 255, 255, 255})
			op.CompositeMode = ebiten.CompositeModeDestinationOut
			op.GeoM.Translate(0.05*float64(sw), 0.05*float64(sh))
			border.DrawImage(inner, op)
		}
		{
			s := draws.NewSpriteFromImage(border)
			s.SetScaleToH(4.1 * regularNoteHeight)
			s.SetPoint(HitPosition, FieldPosition, draws.CenterMiddle)
			skin.ShakeBorderSprite = s
		}
	}
	// Position of combo is dependent on widths of key sprite.
	// Key sprites are overlapped at each side.
	{
		s := draws.NewSprite("skin/drum/key/in.png")
		s.SetScaleToH(FieldInnerHeight)
		s.SetPoint(0, FieldPosition, draws.LeftMiddle)
		keyCenter = s.W()
		skin.KeySprites[LeftRed] = s
	}
	{
		s := draws.NewSprite("skin/drum/key/out.png")
		s.SetScaleToH(FieldInnerHeight)
		s.SetPoint(keyCenter, FieldPosition, draws.LeftMiddle)
		skin.KeySprites[RightBlue] = s
	}
	{
		src := draws.NewImage("skin/drum/key/out.png")
		s := draws.NewSpriteFromImage(draws.NewXFlippedImage(src))
		s.SetScaleToH(FieldInnerHeight)
		s.SetPoint(0, FieldPosition, draws.LeftMiddle)
		skin.KeySprites[LeftBlue] = s
	}
	{
		src := draws.NewImage("skin/drum/key/in.png")
		s := draws.NewSpriteFromImage(draws.NewXFlippedImage(src))
		s.SetScaleToH(FieldInnerHeight)
		s.SetPoint(keyCenter, FieldPosition, draws.LeftMiddle)
		skin.KeySprites[RightRed] = s
	}
	{
		w := keyCenter + skin.KeySprites[RightBlue].W()
		h := skin.KeySprites[RightBlue].H()
		src := ebiten.NewImage(int(w), int(h))
		src.Fill(color.NRGBA{0, 0, 0, uint8(255 * FieldDarkness)})
		s := draws.NewSpriteFromImage(src)
		s.SetPoint(0, FieldPosition, draws.LeftMiddle)
		skin.KeyFieldSprite = s
	}
	for i, name := range []string{"idle", "yes", "no", "high"} {
		fs, err := os.ReadDir(fmt.Sprintf("skin/drum/dancer/%s", name))
		if err != nil {
			continue
		}
		skin.DancerSprites[i] = make([]draws.Sprite, len(fs))
		for j := range fs {
			path := fmt.Sprintf("skin/drum/dancer/%s/%d.png", name, j)
			s := draws.NewSprite(path)
			s.SetScale(draws.Scalar(DancerScale))
			s.SetPoint(DancerPositionX, DancerPositionY, draws.CenterMiddle)
			skin.DancerSprites[i][j] = s
		}
	}
	skin.ScoreSprites = gosu.ScoreSprites
	var comboImages [10]*ebiten.Image
	for i := 0; i < 10; i++ {
		comboImages[i] = draws.NewImage(fmt.Sprintf("skin/combo/%d.png", i))
	}
	for i := 0; i < 10; i++ {
		s := draws.NewSpriteFromImage(comboImages[i])
		s.SetScale(draws.Scalar(ComboScale))
		s.SetPoint(keyCenter, FieldPosition, draws.CenterMiddle)
		skin.ComboSprites[i] = s
	}
}
