package drum

import (
	"fmt"
	"image/color"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/mode"
)

// Order of fields of Skin is roughly consistent with drawing order.
type Skin struct {
	Type int

	DrumSound         [2][2][]byte // First 2 is color, next 2 is size.
	DefaultBackground draws.Sprite
	Field             [2]draws.Sprite
	Hint              [2]draws.Sprite
	Bar               draws.Sprite
	Judgment          [2][3]draws.Animation

	Note    [2][4]draws.Sprite
	Overlay [2]draws.Animation
	Head    [2]draws.Sprite
	Tail    [2]draws.Sprite
	Body    [2]draws.Sprite
	Dot     draws.Sprite
	Shake   [2]draws.Sprite

	Key      [4]draws.Sprite
	KeyField draws.Sprite
	Dancer   [4]draws.Animation
	Score    [13]draws.Sprite
	Combo    [10]draws.Sprite
}

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

var (
	DefaultSkin = Skin{Type: mode.Default}
	UserSkin    = Skin{Type: mode.User}
	// PlaySkin    = Skin{Type: mode.Play}
)

// Todo: embed default skins to code for preventing panic when files are missing
func (skin *Skin) Load(fsys fs.FS) {
	for i, cname := range []string{"red", "blue"} {
		for j, sname := range []string{"", "-big"} {
			name := fmt.Sprintf("drum/sound/%s%s.wav", cname, sname)
			skin.DrumSound[i][j] = audios.NewSound(fsys, name)
		}
	}
	skin.DefaultBackground = mode.UserSkin.DefaultBackground
	for i, name := range []string{"idle", "high"} {
		s := draws.NewSprite(fsys, fmt.Sprintf("drum/stage/field-%s.png", name))
		s.SetSize(ScreenSizeX, S.FieldHeight)
		s.Locate(0, S.FieldPosition, draws.LeftMiddle)
		skin.Field[i] = s
	}
	var hintScale float64
	for i, name := range []string{"idle", "high"} {
		s := draws.NewSprite(fsys, fmt.Sprintf("drum/stage/hint-%s.png", name))
		if name == "idle" {
			hintScale = 1.2 * S.regularNoteHeight / s.H()
		}
		s.ApplyScale(hintScale)
		s.Locate(S.HitPosition, S.FieldPosition, draws.CenterMiddle)
		skin.Hint[i] = s
	}
	{
		src := draws.NewImage(1, S.FieldInnerHeight)
		src.Fill(color.White)
		s := draws.NewSpriteFromSource(src)
		s.Locate(S.HitPosition, S.FieldPosition, draws.CenterMiddle)
		skin.Bar = s
	}
	var (
		note = draws.LoadImage(fsys, "drum/note/note.png")
		end  = draws.LoadImage(fsys, "drum/note/end.png")
		head = draws.NewImageXFlipped(end)
		tail = end
		body = draws.LoadImage(fsys, "drum/note/mid.png")
	)
	for size, sname := range []string{"", "-big"} {
		noteHeight := S.regularNoteHeight
		if size == Big {
			noteHeight = S.bigNoteHeight
		}
		for kind, kname := range []string{"cool", "good", "miss"} {
			var name string
			if kname == "miss" {
				name = "drum/judgment/miss"
			} else {
				name = fmt.Sprintf("drum/judgment/%s%s", kname, sname)
			}
			animation := draws.NewAnimation(fsys, name)
			for i := range animation {
				animation[i].ApplyScale(S.JudgmentScale)
				animation[i].Locate(S.HitPosition, S.FieldPosition, draws.CenterMiddle)
			}
			skin.Judgment[size][kind] = animation
		}
		for kind, color := range []color.NRGBA{ColorRed, ColorBlue, ColorYellow, ColorPurple} {
			image := draws.NewImageColored(note, color)
			s := draws.NewSpriteFromSource(image)
			s.SetScaleToH(noteHeight)
			s.Locate(S.HitPosition, S.FieldPosition, draws.CenterMiddle)
			skin.Note[size][kind] = s
		}
		animation := draws.NewAnimation(fsys, fmt.Sprintf("drum/note/overlay%s", sname))
		for i := range animation {
			animation[i].SetScaleToH(noteHeight)
			animation[i].Locate(S.HitPosition, S.FieldPosition, draws.CenterMiddle)
		}
		skin.Overlay[size] = animation
		{
			s := draws.NewSpriteFromSource(head)
			s.SetScaleToH(noteHeight)
			s.Locate(S.HitPosition, S.FieldPosition, draws.RightMiddle)
			skin.Head[size] = s
		}
		{
			s := draws.NewSpriteFromSource(tail)
			s.SetScaleToH(noteHeight)
			s.Locate(S.HitPosition, S.FieldPosition, draws.LeftMiddle)
			skin.Tail[size] = s
		}
		{
			s := draws.NewSpriteFromSource(body)
			s.SetScaleToH(noteHeight)
			s.Locate(S.HitPosition, S.FieldPosition, draws.LeftMiddle)
			s.Filter = ebiten.FilterNearest
			skin.Body[size] = s
		}
	}
	{
		s := draws.NewSprite(fsys, "drum/note/dot.png")
		s.ApplyScale(S.DotScale)
		s.Locate(S.HitPosition, S.FieldPosition, draws.CenterMiddle)
		skin.Dot = s
	}
	skin.Shake = NewShake(note)
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
		s := draws.NewSpriteFromSource(image)
		s.SetScaleToH(S.FieldInnerHeight)
		if k < 2 { // Includes determining key field size.
			s.Locate(0, S.FieldPosition, draws.LeftMiddle)
			if w := s.W(); keyFieldSize.X < w*2 {
				keyFieldSize.X = w * 2
			}
			if h := s.H(); keyFieldSize.Y < h {
				keyFieldSize.Y = h
			}
		} else {
			s.Locate(keyFieldSize.X/2, S.FieldPosition, draws.LeftMiddle)
		}
		skin.Key[k] = s
	}
	{

		src := draws.NewImage(keyFieldSize.XY())
		src.Fill(color.NRGBA{0, 0, 0, uint8(255 * S.FieldOpaque)})
		s := draws.NewSpriteFromSource(src)
		s.Locate(0, S.FieldPosition, draws.LeftMiddle)
		skin.KeyField = s
	}
	for i, name := range []string{"idle", "yes", "no", "high"} {
		fs, err := fs.ReadDir(fsys, fmt.Sprintf("drum/dancer/%s", name))
		if err != nil {
			continue
		}
		skin.Dancer[i] = make(draws.Animation, len(fs))
		for j := range fs {
			name := fmt.Sprintf("drum/dancer/%s/%d.png", name, j)
			s := draws.NewSprite(fsys, name)
			s.ApplyScale(S.DancerScale)
			s.Locate(S.DancerPositionX, S.DancerPositionY, draws.CenterMiddle)
			skin.Dancer[i][j] = s
		}
	}
	skin.Score = mode.UserSkin.Score
	// Position of combo is dependent on widths of key sprite.
	{
		var comboImages [10]draws.Image
		for i := 0; i < 10; i++ {
			comboImages[i] = draws.LoadImage(fsys, fmt.Sprintf("combo/%d.png", i))
		}
		for i := 0; i < 10; i++ {
			s := draws.NewSpriteFromSource(comboImages[i])
			s.ApplyScale(S.ComboScale)
			s.Locate(keyFieldSize.X/2, S.FieldPosition, draws.CenterMiddle)
			skin.Combo[i] = s
		}
	}
}
func NewShake(note draws.Image) (sprites [2]draws.Sprite) {
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
		s := draws.NewSpriteFromSource(outerImage)
		s.SetScaleToH(scale * S.regularNoteHeight)
		s.Locate(S.HitPosition, S.FieldPosition, draws.CenterMiddle)
		sprites[outer] = s
	}
	{
		s := draws.NewSpriteFromSource(innerImage)
		s.SetScaleToH((scale + thickness) * S.regularNoteHeight)
		s.Locate(S.HitPosition, S.FieldPosition, draws.CenterMiddle)
		sprites[inner] = s
	}
	return
}

func (skin *Skin) Reset() {
	kind := skin.Type
	switch kind {
	case mode.User:
		*skin = DefaultSkin
	case mode.Play:
		*skin = UserSkin
	}
	skin.Type = kind
}
