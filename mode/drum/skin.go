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
	FieldSprites    [2]draws.Sprite
	HintSprites     [2]draws.Sprite
	BarSprite       draws.Sprite
	JudgmentSprites [2][3]draws.Sprite

	NoteSprites       [2][4]draws.Sprite
	HeadSprites       [2]draws.Sprite
	TailSprites       [2]draws.Sprite
	OverlaySprites    [2]draws.Animation
	BodySprites       [2]draws.Sprite
	DotSprite         draws.Sprite
	ShakeBorderSprite draws.Sprite
	ShakeSprite       draws.Sprite

	KeySprites     [4]draws.Sprite
	KeyFieldSprite draws.Sprite
	DancerSprites  [4]draws.Animation
	ScoreSprites   [10]draws.Sprite
	ComboSprites   [10]draws.Sprite
}

// Todo: embed default skins to code for preventing panic when files are missing
func LoadSkin() {
	var skin Skin
	defer func() { DefaultSkin = skin }()
	var noteImage = draws.NewImage("skin/drum/note/note.png")
	for i, name := range []string{"idle", "high"} {
		sprite := draws.NewSprite(fmt.Sprintf("skin/drum/field/%s.png", name))
		sprite.SetSize(screenSizeX, FieldHeight)
		sprite.SetPoint(0, FieldPosition, draws.LeftMiddle)
		skin.FieldSprites[i] = sprite
	}
	for i := range skin.HintSprites {
		const (
			padScale   = 1.1
			outerScale = 1.2
		)
		sw, sh := noteImage.Size()
		// srcSize := draws.IntPt(noteImage.Size())
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
			sd := outerScale - padScale // Size difference.
			// op.GeoM.Translate(srcSize.Mul(draws.Scalar((outerScale - padScale) / 2)).XY())
			op.GeoM.Translate(sd/2*float64(sw), sd/2*float64(sh))
			img.DrawImage(pad, op)
		}
		{
			op := &ebiten.DrawImageOptions{}
			op.ColorM.ScaleWithColor(color.NRGBA{60, 60, 60, a})
			sd := outerScale - 1 // Size difference.
			// op.GeoM.Translate(srcSize.Mul(draws.Scalar((outerScale - 1) / 2)).XY())
			op.GeoM.Translate(sd/2*float64(sw), sd/2*float64(sh))
			img.DrawImage(inner, op)
		}
		sprite := draws.NewSpriteFromImage(img)
		sprite.SetScaleToH(1.2 * regularNoteHeight)
		sprite.SetPoint(HitPosition, FieldPosition, draws.CenterMiddle)
		skin.HintSprites[i] = sprite
	}
	{
		src := ebiten.NewImage(1, int(FieldInnerHeight))
		src.Fill(color.NRGBA{255, 255, 255, 255})
		sprite := draws.NewSpriteFromImage(src)
		sprite.SetPoint(HitPosition, FieldPosition, draws.CenterMiddle)
		skin.BarSprite = sprite
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
			sprite := draws.NewSprite(path)
			sprite.ApplyScale(JudgmentScale)
			sprite.SetPoint(HitPosition, FieldPosition, draws.CenterMiddle)
			skin.JudgmentSprites[i][j] = sprite
		}
		for j, clr := range []color.NRGBA{ColorRed, ColorBlue, ColorYellow, ColorPurple} {
			img := ebiten.NewImage(noteImage.Size())
			op := &ebiten.DrawImageOptions{}
			op.ColorM.ScaleWithColor(clr)
			img.DrawImage(noteImage, op)

			sprite := draws.NewSpriteFromImage(img)
			sprite.SetScaleToH(noteHeight)
			sprite.SetPoint(HitPosition, FieldPosition, draws.CenterMiddle)
			skin.NoteSprites[i][j] = sprite
		}
		{
			sprite := draws.NewSpriteFromImage(rollEndImage)
			sprite.SetScaleToH(noteHeight)
			sprite.SetPoint(HitPosition, FieldPosition, draws.LeftMiddle)
			skin.TailSprites[i] = sprite
		}
		{
			sprite := draws.NewSpriteFromImage(draws.NewXFlippedImage(rollEndImage))
			sprite.SetScaleToH(noteHeight)
			sprite.SetPoint(HitPosition, FieldPosition, draws.RightMiddle)
			skin.HeadSprites[i] = sprite
		}
		{
			paths := gosu.Paths(fmt.Sprintf("skin/drum/note/overlay/%s", sname))
			skin.OverlaySprites[i] = make(draws.Animation, len(paths))
			for j, path := range paths {
				sprite := draws.NewSprite(path)
				sprite.SetScaleToH(noteHeight)
				sprite.SetPoint(HitPosition, FieldPosition, draws.CenterMiddle)
				skin.OverlaySprites[i][j] = sprite
			}
		}
		{
			sprite := draws.NewSpriteFromImage(rollMidImage)
			sprite.SetScaleToH(noteHeight)
			sprite.SetPoint(HitPosition, FieldPosition, draws.LeftMiddle)
			sprite.Filter = ebiten.FilterNearest
			skin.BodySprites[i] = sprite
		}
	}
	{
		sprite := draws.NewSprite("skin/drum/note/roll/dot.png")
		sprite.ApplyScale(DotScale)
		sprite.SetPoint(HitPosition, FieldPosition, draws.CenterMiddle)
		skin.DotSprite = sprite
	}
	{
		const (
			scale     = 4.0
			thickness = 0.1
		)
		sw, sh := noteImage.Size()
		inner := draws.NewScaledImage(noteImage, scale)
		shake := ebiten.NewImage(inner.Size())
		{
			op := &ebiten.DrawImageOptions{}
			color := ColorPurple
			color.A = 128
			op.ColorM.ScaleWithColor(color)
			shake.DrawImage(inner, op)
		}
		{
			sprite := draws.NewSpriteFromImage(shake)
			sprite.SetScaleToH(scale * regularNoteHeight)
			sprite.SetPoint(HitPosition, FieldPosition, draws.CenterMiddle)
			skin.ShakeSprite = sprite
		}

		outer := draws.NewScaledImage(noteImage, scale+thickness)
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
			op.GeoM.Translate(thickness/2*float64(sw), thickness/2*float64(sh))
			border.DrawImage(inner, op)
		}
		{
			sprite := draws.NewSpriteFromImage(border)
			sprite.SetScaleToH((scale + thickness) * regularNoteHeight)
			sprite.SetPoint(HitPosition, FieldPosition, draws.CenterMiddle)
			skin.ShakeBorderSprite = sprite
		}
	}
	// Position of combo is dependent on widths of key sprite.
	// Key sprites are overlapped at each side.
	{
		sprite := draws.NewSprite("skin/drum/key/in.png")
		sprite.SetScaleToH(FieldInnerHeight)
		sprite.SetPoint(0, FieldPosition, draws.LeftMiddle)
		keyCenter = sprite.W()
		skin.KeySprites[LeftRed] = sprite
	}
	{
		sprite := draws.NewSprite("skin/drum/key/out.png")
		sprite.SetScaleToH(FieldInnerHeight)
		sprite.SetPoint(keyCenter, FieldPosition, draws.LeftMiddle)
		skin.KeySprites[RightBlue] = sprite
	}
	{
		src := draws.NewImage("skin/drum/key/out.png")
		sprite := draws.NewSpriteFromImage(draws.NewXFlippedImage(src))
		sprite.SetScaleToH(FieldInnerHeight)
		sprite.SetPoint(0, FieldPosition, draws.LeftMiddle)
		skin.KeySprites[LeftBlue] = sprite
	}
	{
		src := draws.NewImage("skin/drum/key/in.png")
		sprite := draws.NewSpriteFromImage(draws.NewXFlippedImage(src))
		sprite.SetScaleToH(FieldInnerHeight)
		sprite.SetPoint(keyCenter, FieldPosition, draws.LeftMiddle)
		skin.KeySprites[RightRed] = sprite
	}
	{
		w := keyCenter + skin.KeySprites[RightBlue].W()
		h := skin.KeySprites[RightBlue].H()
		src := ebiten.NewImage(int(w), int(h))
		src.Fill(color.NRGBA{0, 0, 0, uint8(255 * FieldDarkness)})
		sprite := draws.NewSpriteFromImage(src)
		sprite.SetPoint(0, FieldPosition, draws.LeftMiddle)
		skin.KeyFieldSprite = sprite
	}
	for i, name := range []string{"idle", "yes", "no", "high"} {
		fs, err := os.ReadDir(fmt.Sprintf("skin/drum/dancer/%s", name))
		if err != nil {
			continue
		}
		skin.DancerSprites[i] = make(draws.Animation, len(fs))
		for j := range fs {
			path := fmt.Sprintf("skin/drum/dancer/%s/%d.png", name, j)
			sprite := draws.NewSprite(path)
			sprite.ApplyScale(DancerScale)
			sprite.SetPoint(DancerPositionX, DancerPositionY, draws.CenterMiddle)
			skin.DancerSprites[i][j] = sprite
		}
	}
	skin.ScoreSprites = gosu.ScoreSprites
	var comboImages [10]*ebiten.Image
	for i := 0; i < 10; i++ {
		comboImages[i] = draws.NewImage(fmt.Sprintf("skin/combo/%d.png", i))
	}
	for i := 0; i < 10; i++ {
		sprite := draws.NewSpriteFromImage(comboImages[i])
		sprite.ApplyScale(ComboScale)
		sprite.SetPoint(keyCenter, FieldPosition, draws.CenterMiddle)
		skin.ComboSprites[i] = sprite
	}
}
