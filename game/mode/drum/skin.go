package drum

import (
	"fmt"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/framework/audios"
	"github.com/hndada/gosu/framework/draws"
	"github.com/hndada/gosu/game"
)

var (
	ColorRed    = color.NRGBA{235, 69, 44, 255}
	ColorBlue   = color.NRGBA{68, 141, 171, 255}
	ColorYellow = color.NRGBA{230, 170, 0, 255} // 252, 83, 6
	ColorPurple = color.NRGBA{150, 100, 200, 255}
)

const (
	DancerIdle = iota
	DancerYes
	DancerNo
	DancerHigh
)

var DefaultSkin Skin

// Order of fields of Skin is roughly consistent with drawing order.
type Skin struct {
	FieldSprites    [2]draws.Sprite
	HintSprites     [2]draws.Sprite
	BarSprite       draws.Sprite
	JudgmentSprites [2][3]draws.Animation

	NoteSprites    [2][4]draws.Sprite
	OverlaySprites [2]draws.Animation
	HeadSprites    [2]draws.Sprite
	TailSprites    [2]draws.Sprite
	BodySprites    [2]draws.Sprite
	DotSprite      draws.Sprite
	ShakeSprites   [2]draws.Sprite

	KeySprites     [4]draws.Sprite
	KeyFieldSprite draws.Sprite
	DancerSprites  [4]draws.Animation
	ScoreSprites   [10]draws.Sprite
	ComboSprites   [10]draws.Sprite

	SoundEffectBytes [2][2][]byte
}

// Todo: embed default skins to code for preventing panic when files are missing
func LoadSkin(fsys fs.FS) {
	var skin Skin
	defer func() { DefaultSkin = skin }()
	var note = draws.LoadImage(fsys, "drum/note/note.png")
	for i, name := range []string{"idle", "high"} {
		sprite := draws.NewSprite(fsys, fmt.Sprintf("drum/field/%s.png", name))
		sprite.SetSize(ScreenSizeX, FieldHeight)
		sprite.Locate(0, FieldPosition, draws.LeftMiddle)
		skin.FieldSprites[i] = sprite
	}
	var hintScale float64
	for i, name := range []string{"idle", "high"} {
		sprite := draws.NewSprite(fsys, fmt.Sprintf("drum/hint/%s.png", name))
		if name == "idle" {
			hintScale = 1.2 * regularNoteHeight / sprite.H()
		}
		sprite.ApplyScale(hintScale)
		sprite.Locate(HitPosition, FieldPosition, draws.CenterMiddle)
		skin.HintSprites[i] = sprite
	}
	{
		src := draws.NewImage(1, FieldInnerHeight)
		src.Fill(color.White)
		sprite := draws.NewSpriteFromSource(src)
		sprite.Locate(HitPosition, FieldPosition, draws.CenterMiddle)
		skin.BarSprite = sprite
	}
	var (
		end  = draws.LoadImage(fsys, "drum/note/roll/end.png")
		head = draws.NewImageXFlipped(end)
		tail = end
		body = draws.LoadImage(fsys, "drum/note/roll/mid.png")
	)
	for size, sizeName := range []string{"regular", "big"} {
		noteHeight := regularNoteHeight
		if size == Big {
			noteHeight = bigNoteHeight
		}
		for kind, kindName := range []string{"cool", "good", "miss"} {
			var name string
			if kindName == "miss" {
				name = "drum/judgment/miss"
			} else {
				name = fmt.Sprintf("drum/judgment/%s/%s", sizeName, kindName)
			}
			animation := draws.NewAnimation(fsys, name)
			for i := range animation {
				animation[i].ApplyScale(JudgmentScale)
				animation[i].Locate(HitPosition, FieldPosition, draws.CenterMiddle)
			}
			skin.JudgmentSprites[size][kind] = animation
		}
		for kind, color := range []color.NRGBA{ColorRed, ColorBlue, ColorYellow, ColorPurple} {
			image := draws.NewImageColored(note, color)
			sprite := draws.NewSpriteFromSource(image)
			sprite.SetScaleToH(noteHeight)
			sprite.Locate(HitPosition, FieldPosition, draws.CenterMiddle)
			skin.NoteSprites[size][kind] = sprite
		}
		animation := draws.NewAnimation(fsys, fmt.Sprintf("drum/note/overlay/%s", sizeName))
		for i := range animation {
			animation[i].SetScaleToH(noteHeight)
			animation[i].Locate(HitPosition, FieldPosition, draws.CenterMiddle)
		}
		skin.OverlaySprites[size] = animation
		{
			sprite := draws.NewSpriteFromSource(head)
			sprite.SetScaleToH(noteHeight)
			sprite.Locate(HitPosition, FieldPosition, draws.RightMiddle)
			skin.HeadSprites[size] = sprite
		}
		{
			sprite := draws.NewSpriteFromSource(tail)
			sprite.SetScaleToH(noteHeight)
			sprite.Locate(HitPosition, FieldPosition, draws.LeftMiddle)
			skin.TailSprites[size] = sprite
		}
		{
			sprite := draws.NewSpriteFromSource(body)
			sprite.SetScaleToH(noteHeight)
			sprite.Locate(HitPosition, FieldPosition, draws.LeftMiddle)
			sprite.Filter = ebiten.FilterNearest
			skin.BodySprites[size] = sprite
		}
	}
	{
		sprite := draws.NewSprite(fsys, "drum/note/roll/dot.png")
		sprite.ApplyScale(DotScale)
		sprite.Locate(HitPosition, FieldPosition, draws.CenterMiddle)
		skin.DotSprite = sprite
	}
	skin.ShakeSprites = NewShakeSprites(note)
	// Key sprites are overlapped at each side.
	var (
		in        = draws.LoadImage(fsys, "drum/key/in.png")
		out       = draws.LoadImage(fsys, "drum/key/out.png")
		keyImages = []draws.Image{
			draws.NewImageXFlipped(out),
			in,
			draws.NewImageXFlipped(in),
			out,
		}
		keyFieldSize draws.Vector2
	)
	for k, image := range keyImages {
		sprite := draws.NewSpriteFromSource(image)
		sprite.SetScaleToH(FieldInnerHeight)
		if k < 2 { // Includes determining key field size.
			sprite.Locate(0, FieldPosition, draws.LeftMiddle)
			if w := sprite.W(); keyFieldSize.X < w*2 {
				keyFieldSize.X = w * 2
			}
			if h := sprite.H(); keyFieldSize.Y < h {
				keyFieldSize.Y = h
			}
		} else {
			sprite.Locate(keyFieldSize.X/2, FieldPosition, draws.LeftMiddle)
		}
		skin.KeySprites[k] = sprite
	}
	{

		src := draws.NewImage(keyFieldSize.XY())
		src.Fill(color.NRGBA{0, 0, 0, uint8(255 * FieldDarkness)})
		sprite := draws.NewSpriteFromSource(src)
		sprite.Locate(0, FieldPosition, draws.LeftMiddle)
		skin.KeyFieldSprite = sprite
	}
	for i, name := range []string{"idle", "yes", "no", "high"} {
		fs, err := fs.ReadDir(fsys, fmt.Sprintf("drum/dancer/%s", name))
		if err != nil {
			continue
		}
		skin.DancerSprites[i] = make(draws.Animation, len(fs))
		for j := range fs {
			name := fmt.Sprintf("drum/dancer/%s/%d.png", name, j)
			sprite := draws.NewSprite(fsys, name)
			sprite.ApplyScale(DancerScale)
			sprite.Locate(DancerPositionX, DancerPositionY, draws.CenterMiddle)
			skin.DancerSprites[i][j] = sprite
		}
	}
	skin.ScoreSprites = game.ScoreSprites
	// Position of combo is dependent on widths of key sprite.
	{
		var comboImages [10]draws.Image
		for i := 0; i < 10; i++ {
			comboImages[i] = draws.LoadImage(fsys, fmt.Sprintf("combo/%d.png", i))
		}
		for i := 0; i < 10; i++ {
			sprite := draws.NewSpriteFromSource(comboImages[i])
			sprite.ApplyScale(ComboScale)
			sprite.Locate(keyFieldSize.X/2, FieldPosition, draws.CenterMiddle)
			skin.ComboSprites[i] = sprite
		}
	}
	for i, colorName := range []string{"regular", "big"} {
		for j, sizeName := range []string{"red", "blue"} {
			name := fmt.Sprintf("drum/sound/%s/%s.wav", colorName, sizeName)
			b, err := audios.NewBytes(fsys, name)
			if err != nil {
				panic(err)
			}
			skin.SoundEffectBytes[i][j] = b
		}
	}
}
func NewShakeSprites(note draws.Image) (sprites [2]draws.Sprite) {
	const (
		outer = iota
		inner
	)
	const (
		scale     = 4.0
		thickness = 0.1
	)
	var (
		outerImage = draws.NewImage(note.Size().Scale(scale + thickness).XY())
		innerImage = draws.NewImage(note.Size().Scale(scale).XY())
	)
	// Be careful that images goes sqaure when color the images by Fill().
	{
		op := draws.Op{}
		op.GeoM.Scale(scale+thickness, scale+thickness)
		op.ColorM.ScaleWithColor(color.NRGBA{255, 255, 255, 255})
		op.Filter = ebiten.FilterLinear
		note.Draw(outerImage, op)
	}
	{
		op := draws.Op{}
		purple := ColorPurple
		purple.A = 128 // 152
		op.GeoM.Scale(scale, scale)
		op.ColorM.ScaleWithColor(purple)
		note.Draw(innerImage, op)
	}
	{
		op := draws.Op{}
		op.ColorM.Scale(1, 1, 1, 1.5)
		op.CompositeMode = ebiten.CompositeModeDestinationOut
		op.GeoM.Translate(note.Size().Scale(thickness / 2).XY())
		innerImage.Draw(outerImage, op)
	}
	{
		sprite := draws.NewSpriteFromSource(outerImage)
		sprite.SetScaleToH(scale * regularNoteHeight)
		sprite.Locate(HitPosition, FieldPosition, draws.CenterMiddle)
		sprites[outer] = sprite
	}
	{
		sprite := draws.NewSpriteFromSource(innerImage)
		sprite.SetScaleToH((scale + thickness) * regularNoteHeight)
		sprite.Locate(HitPosition, FieldPosition, draws.CenterMiddle)
		sprites[inner] = sprite
	}
	return
}

// Deprecated.
// func NewHintGlowImage(skin Skin, noteImage draws.Image) {
// 	for i := range skin.HintSprites {
// 		const (
// 			padScale   = 1.1
// 			outerScale = 1.2
// 		)
// 		sw, sh := noteImage.Size()
// 		outer := draws.NewImageScaled(noteImage, outerScale)
// 		pad := draws.NewImageScaled(noteImage, padScale)
// 		inner := noteImage
// 		a := uint8(255 * FieldDarkness)
// 		img := ebiten.NewImage(outer.Size())
// 		{
// 			op := &ebiten.DrawImageOptions{}
// 			op.ColorM.ScaleWithColor(color.NRGBA{128, 128, 128, a})
// 			op.GeoM.Translate(0, 0)
// 			img.DrawImage(outer, op)
// 		}
// 		{
// 			op := &ebiten.DrawImageOptions{}
// 			op.ColorM.ScaleWithColor(color.NRGBA{255, 255, 0, a})
// 			if i == 0 { // Blank for idle, Yellow for highlight.
// 				op.CompositeMode = ebiten.CompositeModeDestinationOut
// 			}
// 			sd := outerScale - padScale // Size difference.
// 			op.GeoM.Translate(sd/2*float64(sw), sd/2*float64(sh))
// 			img.DrawImage(pad, op)
// 		}
// 		{
// 			op := &ebiten.DrawImageOptions{}
// 			op.ColorM.ScaleWithColor(color.NRGBA{60, 60, 60, a})
// 			sd := outerScale - 1 // Size difference.
// 			op.GeoM.Translate(sd/2*float64(sw), sd/2*float64(sh))
// 			img.DrawImage(inner, op)
// 		}
// 		sprite := draws.NewSpriteFromSource(img)
// 		sprite.SetScaleToH(1.2 * regularNoteHeight)
// 		sprite.Locate(HitPosition, FieldPosition, draws.CenterMiddle)
// 		skin.HintSprites[i] = sprite
// 	}
// }
