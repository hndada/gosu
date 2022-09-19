package drum

import (
	"fmt"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu"
	"github.com/hndada/gosu/draws"
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
	OverlaySprites    [2][2]draws.Sprite
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
		ratioW := screenSizeX / s.W()
		ratioH := FieldHeight / s.H()
		s.SetScaleXY(ratioW, ratioH, ebiten.FilterLinear)
		s.SetPosition(0, FieldPosition, draws.OriginLeftMiddle)
		skin.FieldSprite = s
	}
	for i := range skin.HintSprites {
		sw, sh := noteImage.Size()
		outer := draws.NewScaledImage(noteImage, 1.2)
		pad := draws.NewScaledImage(noteImage, 1.1)
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
			op.GeoM.Translate(0.05*float64(sw), 0.05*float64(sh))
			img.DrawImage(pad, op)
		}
		{
			op := &ebiten.DrawImageOptions{}
			op.ColorM.ScaleWithColor(color.NRGBA{60, 60, 60, a})
			op.GeoM.Translate(0.1*float64(sw), 0.1*float64(sh))
			img.DrawImage(inner, op)
		}
		s := draws.NewSpriteFromImage(img)
		s.SetScale(1.2 * regularNoteHeight / s.H())
		s.SetPosition(HitPosition, FieldPosition, draws.OriginCenterMiddle)
		skin.HintSprites[i] = s
	}
	{
		src := ebiten.NewImage(1, int(FieldInnerHeight))
		src.Fill(color.NRGBA{255, 255, 255, 255})
		s := draws.NewSpriteFromImage(src)
		s.SetPosition(HitPosition, FieldPosition, draws.OriginCenterMiddle)
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
			s.SetScale(JudgmentScale)
			s.SetPosition(HitPosition, FieldPosition, draws.OriginCenterMiddle)
			skin.JudgmentSprites[i][j] = s
		}
		for j, clr := range []color.NRGBA{ColorRed, ColorBlue, ColorYellow, ColorPurple} {
			img := ebiten.NewImage(noteImage.Size())
			op := &ebiten.DrawImageOptions{}
			op.ColorM.ScaleWithColor(clr)
			img.DrawImage(noteImage, op)

			s := draws.NewSpriteFromImage(img)
			s.SetScale(noteHeight / s.H())
			s.SetPosition(HitPosition, FieldPosition, draws.OriginCenterMiddle)
			skin.NoteSprites[i][j] = s
		}
		{
			s := draws.NewSpriteFromImage(rollEndImage)
			s.SetScale(noteHeight / s.H())
			s.SetPosition(HitPosition, FieldPosition, draws.OriginLeftMiddle)
			skin.TailSprites[i] = s
		}
		{
			s := draws.NewSpriteFromImage(draws.XFlippedImage(rollEndImage))
			s.SetScale(noteHeight / s.H())
			s.SetPosition(HitPosition, FieldPosition, draws.OriginRightMiddle)
			skin.HeadSprites[i] = s
		}
		for j := 0; j < 2; j++ {
			path := fmt.Sprintf("skin/drum/note/overlay/%s/%d.png", sname, j)
			if _, err := os.Stat(path); os.IsNotExist(err) { // One overlay.
				path = fmt.Sprintf("skin/drum/note/overlay/%s.png", sname)
			}
			s := draws.NewSprite(path)
			s.SetScale(noteHeight / s.W())
			s.SetPosition(HitPosition, FieldPosition, draws.OriginCenterMiddle)
			skin.OverlaySprites[i][j] = s
		}
		{
			s := draws.NewSpriteFromImage(rollMidImage)
			s.SetScale(noteHeight / s.H())
			s.SetPosition(HitPosition, FieldPosition, draws.OriginLeftMiddle)
			skin.BodySprites[i] = s
		}
	}
	{
		s := draws.NewSprite("skin/drum/note/roll/dot.png")
		s.SetScale(DotScale)
		s.SetPosition(HitPosition, FieldPosition, draws.OriginCenterMiddle)
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
			s.SetScale(4 * regularNoteHeight / s.H())
			s.SetPosition(HitPosition, FieldPosition, draws.OriginCenterMiddle)
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
			s.SetScale(4.1 * regularNoteHeight / s.H())
			s.SetPosition(HitPosition, FieldPosition, draws.OriginCenterMiddle)
			skin.ShakeBorderSprite = s
		}
	}
	// Position of combo is dependent on widths of key sprite.
	// Key sprites are overlapped at each side.
	{
		s := draws.NewSprite("skin/drum/key/in.png")
		s.SetScale(FieldInnerHeight / s.H())
		s.SetPosition(0, FieldPosition, draws.OriginLeftMiddle)
		keyCenter = s.W()
		skin.KeySprites[LeftRed] = s
	}
	{
		s := draws.NewSprite("skin/drum/key/out.png")
		s.SetScale(FieldInnerHeight / s.H())
		s.SetPosition(keyCenter, FieldPosition, draws.OriginLeftMiddle)
		skin.KeySprites[RightBlue] = s
	}
	{
		src := draws.NewImage("skin/drum/key/out.png")
		s := draws.NewSpriteFromImage(draws.XFlippedImage(src))
		s.SetScale(FieldInnerHeight / s.H())
		s.SetPosition(0, FieldPosition, draws.OriginLeftMiddle)
		skin.KeySprites[LeftBlue] = s
	}
	{
		src := draws.NewImage("skin/drum/key/in.png")
		s := draws.NewSpriteFromImage(draws.XFlippedImage(src))
		s.SetScale(FieldInnerHeight / s.H())
		s.SetPosition(keyCenter, FieldPosition, draws.OriginLeftMiddle)
		skin.KeySprites[RightRed] = s
	}
	{
		w := keyCenter + skin.KeySprites[RightBlue].W()
		h := skin.KeySprites[RightBlue].H()
		src := ebiten.NewImage(int(w), int(h))
		src.Fill(color.NRGBA{0, 0, 0, uint8(255 * FieldDarkness)})
		s := draws.NewSpriteFromImage(src)
		s.SetPosition(0, FieldPosition, draws.OriginLeftMiddle)
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
			s.SetScale(DancerScale)
			s.SetPosition(DancerPositionX, DancerPositionY, draws.OriginCenterMiddle)
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
		s.SetScale(ComboScale)
		s.SetPosition(keyCenter, FieldPosition, draws.OriginCenterMiddle)
		skin.ComboSprites[i] = s
	}
}
