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
	ColorYellow = color.NRGBA{252, 83, 6, 255}
	ColorPurple = color.NRGBA{150, 100, 200, 255}
	ColorGray   = color.NRGBA{67, 67, 67, 255}
)

// const (
//
//	ShakeNote = iota
//	ShakeInner
//	ShakeOuter
//	// ShakeSpin
//	// ShakeLimit
//
// )
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

// https://osu.ppy.sh/wiki/en/Skinning/osu%21taiko
type Skin struct {
	FieldSprite draws.Sprite
	HintSprite  draws.Sprite
	BarSprite   draws.Sprite

	// First [2] are for big notes.
	JudgmentSprites [2][3]draws.Sprite // 3 Judgments.
	// RedSprites      [2]draws.Sprite
	// BlueSprites     [2]draws.Sprite
	// NoteSprites [2][3]draws.Sprite // [3] are for Red, Blue, Yellow each.
	NoteSprites [2][4]draws.Sprite // [4] are for Red, Blue, Yellow, Purple each.
	// HeadSprites    [2]draws.Sprite    // Overlay will be drawn during game play.
	TailSprites    [2]draws.Sprite
	OverlaySprites [2][2]draws.Sprite // 2 Overlays.
	BodySprites    [2]draws.Sprite
	DotSprite      draws.Sprite
	// ShakeSprites   [3]draws.Sprite
	ShakeBorderSprite draws.Sprite
	ShakeSprite       draws.Sprite

	KeySprites     [4]draws.Sprite // 4 Keys.
	KeyFieldSprite draws.Sprite
	DancerSprites  [4][]draws.Sprite // Dancer has 4 behaviors.

	ScoreSprites      [10]draws.Sprite
	ComboSprites      [10]draws.Sprite
	DotCountSprites   [10]draws.Sprite // For rolls.
	ShakeCountSprites [10]draws.Sprite // For shakes.
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
	{
		sw, sh := noteImage.Size()
		outer := draws.NewScaledImage(noteImage, 1.2)
		pad := draws.NewScaledImage(noteImage, 1.1)
		inner := noteImage
		img := ebiten.NewImage(outer.Size())
		a := uint8(255 * FieldDarkness)
		{
			op := &ebiten.DrawImageOptions{}
			op.ColorM.ScaleWithColor(color.NRGBA{128, 128, 128, a})
			op.GeoM.Translate(0, 0)
			img.DrawImage(outer, op)
		}
		{
			op := &ebiten.DrawImageOptions{}
			op.ColorM.ScaleWithColor(color.NRGBA{255, 255, 255, 255})
			op.CompositeMode = ebiten.CompositeModeDestinationOut
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
		skin.HintSprite = s
	}
	{
		src := ebiten.NewImage(1, int(FieldInnerHeight))
		src.Fill(color.NRGBA{255, 255, 255, 255}) // White
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
			// s.SetColor(ColorYellow)
			skin.TailSprites[i] = s
		}
		// {
		// 	s := draws.NewSpriteFromImage(rollEndImage)
		// 	ratio := noteHeight / s.H()
		// 	s.SetScaleXY(-ratio, ratio, ebiten.FilterLinear) // Goes flipped.
		// 	s.SetPosition(HitPosition, FieldPosition, draws.OriginRightMiddle)
		// 	s.SetColor(ColorYellow)
		// 	skin.HeadSprites[i] = s
		// }
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
			// s.SetColor(ColorYellow)
			skin.BodySprites[i] = s
		}
	}
	{
		s := draws.NewSprite("skin/drum/note/roll/dot.png")
		s.SetScale(DotScale)
		s.SetPosition(HitPosition, FieldPosition, draws.OriginCenterMiddle)
		// s.SetColor(ColorYellow)
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
	// for i, name := range []string{"note", "spin", "limit"} {
	// 	path := fmt.Sprintf("skin/drum/note/shake/%s.png", name)
	// 	s := draws.NewSprite(path)
	// 	if name == "note" {
	// 		s.SetScale(regularNoteHeight / s.H())
	// 		s.SetPosition(HitPosition, FieldPosition, draws.OriginCenterMiddle)
	// 	} else {
	// 		s.SetScale(ShakeScale)
	// 		s.SetPosition(HitPosition, FieldPosition, draws.OriginCenterMiddle)
	// 		// s.SetPosition(ShakePosX, ShakePosY, draws.OriginCenterMiddle)
	// 	}
	// 	skin.ShakeSprites[i] = s
	// }

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
		s := draws.NewSpriteFromImage(draws.FlipX(src))
		s.SetScale(FieldInnerHeight / s.H())
		s.SetPosition(0, FieldPosition, draws.OriginLeftMiddle)
		skin.KeySprites[LeftBlue] = s
	}
	{
		src := draws.NewImage("skin/drum/key/in.png")
		s := draws.NewSpriteFromImage(draws.FlipX(src))
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
			s.SetPosition(DancerPosX, DancerPosY, draws.OriginCenterMiddle)
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
	for i := 0; i < 10; i++ {
		s := draws.NewSpriteFromImage(comboImages[i])
		s.SetScale(DotCountScale)
		s.SetPosition(HitPosition, FieldPosition, draws.OriginCenterMiddle)
		skin.DotCountSprites[i] = s
	}
	// for i := 0; i < 10; i++ {
	// 	s := draws.NewSpriteFromImage(comboImages[i])
	// 	s.SetScale(ShakeCountScale)
	// 	pos := ShakePosY + s.H()*ShakeCountPosition
	// 	s.SetPosition(ShakePosX, pos, draws.OriginCenterTop)
	// 	skin.ShakeCountSprites[i] = s
	// }
}

// func IsKeyImageFlipped(keyType int) bool {
// 	return keyType == KeyLeftKat || keyType == KeyRightDon
// }
// func() {
// 	fs, err := os.ReadDir("skin/drum/overlay")
// 	if err != nil {
// 		return
// 	}
// 	for _, f := range fs {
// 		for i, name := range []string{"regular", "big"} {
// 			if name != f.Name() {
// 				continue
// 			}
// 			if f.IsDir() {
// 				for j := 0; j < 2; j++ {
// 					path := fmt.Sprintf("skin/drum/overlay/%s/%d.png", name, j)
// 					noteOverlays[i][j] = draws.NewImage(path)
// 				}
// 			} else {
// 				path := fmt.Sprintf("skin/drum/overlay/%s.png", name)
// 				noteOverlays[i][0] = draws.NewImage(path)
// 				noteOverlays[i][1] = noteOverlays[i][0]
// 			}
// 		}
// 	}
// }()

// func IsDir(path string) (bool, error) {
// 	f, err := os.Open(path)
// 	if err != nil {
// 		return false, err
// 	}
// 	defer f.Close()
// 	info, err := f.Stat()
// 	if err != nil {
// 		return false, err
// 	}
// 	return info.IsDir(), nil
// }
