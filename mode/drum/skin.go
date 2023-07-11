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

// Todo: resolve ambiguity between image and derived sprites
// Order of fields of Skin is roughly consistent with drawing order.
type Skin struct {
	field   [2]draws.Image
	hint    [2]draws.Image
	note    draws.Image
	end     draws.Image
	mid     draws.Image
	dot     draws.Image
	overlay [2][]draws.Image
	key     [2]draws.Image
	dancer  [4][]draws.Image
	combo   [10]draws.Image

	defaultBackground draws.Sprite
	Field             [2]draws.Sprite
	Hint              [2]draws.Sprite
	Bar               draws.Sprite

	// [2] at last element stands for a size.
	DrumSound [2][2][]byte // color, size
	Judgment  [3][2]draws.Animation
	Note      [4][2]draws.Sprite
	Head      [2]draws.Sprite
	Tail      [2]draws.Sprite
	Body      [2]draws.Sprite
	Dot       draws.Sprite
	Shake     [2]draws.Sprite
	Overlay   [2]draws.Animation

	Key      [4]draws.Sprite
	KeyField draws.Sprite
	Dancer   [4]draws.Animation
	score    [13]draws.Sprite
	Combo    [10]draws.Sprite
}

var (
	ColorRed    = color.NRGBA{235, 69, 44, 255}
	ColorBlue   = color.NRGBA{68, 141, 171, 255}
	ColorYellow = color.NRGBA{230, 170, 0, 255} // 252, 83, 6
	ColorPurple = color.NRGBA{150, 100, 200, 255}
)

const (
	Idle = iota
	High
)
const (
	DancerIdle = iota
	DancerYes
	DancerNo
	DancerHigh
)

var (
	DefaultSkin = &Skin{}
	UserSkin    = &Skin{}
)

// Calling fillBlank() right after returning Load won't apply User settings.
func (skin *Skin) Load(fsys fs.FS) {
	// defer skin.fillBlank(DefaultSkin)
	skin.defaultBackground = mode.UserSkin.DefaultBackground
	for state, name := range []string{"idle", "high"} {
		img := draws.NewImageFromFile(fsys, fmt.Sprintf("drum/stage/field-%s.png", name))
		if !img.IsValid() {
			switch state {
			case Idle:
				img = DefaultSkin.field[state]
			case High: // Use user's idle one.
				img = skin.field[0]
			}
		}
		skin.field[state] = img
		s := draws.NewSprite(img)
		s.SetSize(ScreenSizeX, S.FieldHeight)
		s.Locate(0, S.FieldPosition, draws.LeftMiddle)
		skin.Field[state] = s
	}
	var hintScale float64 // For using idle image's size.
	for state, name := range []string{"idle", "high"} {
		img := draws.NewImageFromFile(fsys, fmt.Sprintf("drum/stage/hint-%s.png", name))
		if !img.IsValid() {
			switch state {
			case Idle:
				img = DefaultSkin.field[state]
			case High: // Use user's idle one.
				img = skin.field[0]
			}
		}
		skin.hint[state] = img
		s := draws.NewSprite(img)
		if name == "idle" {
			hintScale = 1.2 * S.regularNoteHeight / s.H()
		}
		s.MultiplyScale(hintScale)
		s.Locate(S.HitPosition, S.FieldPosition, draws.CenterMiddle)
		skin.Hint[state] = s
	}
	{
		src := draws.NewImage(1, S.FieldInnerHeight)
		src.Fill(color.White)
		s := draws.NewSprite(src)
		s.Locate(S.HitPosition, S.FieldPosition, draws.CenterMiddle)
		skin.Bar = s
	}

	skin.note = draws.NewImageFromFile(fsys, "drum/note/note.png")
	if !skin.note.IsValid() {
		skin.note = DefaultSkin.note
	}
	skin.end = draws.NewImageFromFile(fsys, "drum/note/end.png")
	if !skin.end.IsValid() {
		skin.end = DefaultSkin.end
	}
	skin.mid = draws.NewImageFromFile(fsys, "drum/note/mid.png")
	if !skin.mid.IsValid() {
		skin.mid = DefaultSkin.mid
	}
	skin.dot = draws.NewImageFromFile(fsys, "drum/note/dot.png")
	if !skin.dot.IsValid() {
		skin.dot = DefaultSkin.dot
	}
	for i, sname := range []string{"", "-big"} {
		name := fmt.Sprintf("drum/note/overlay%s", sname)
		imgs := draws.NewImagesFromFile(fsys, name)
		if len(imgs) == 1 && !imgs[0].IsValid() {
			skin.overlay = DefaultSkin.overlay
			break
		}
		skin.overlay[i] = imgs
	}
	var (
		head = draws.NewImageXFlipped(skin.end)
		tail = skin.end
		body = skin.mid
	)
	for size, sname := range []string{"", "-big"} {
		for i, cname := range []string{"red", "blue"} {
			name := fmt.Sprintf("drum/sound/%s%s.wav", cname, sname)
			s := audios.NewSound(fsys, name)
			if !s.IsValid() {
				s = DefaultSkin.DrumSound[i][size]
			}
			skin.DrumSound[i][size] = s
		}
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
			a := draws.NewAnimationFromFile(fsys, name)
			for frame := range a {
				a[frame].MultiplyScale(S.JudgmentScale)
				a[frame].Locate(S.HitPosition, S.FieldPosition, draws.CenterMiddle)
			}
			skin.Judgment[kind][size] = a
		}
		for kind, color := range []color.NRGBA{ColorRed, ColorBlue, ColorYellow, ColorPurple} {
			img := draws.NewImageColored(skin.note, color)
			s := draws.NewSprite(img)
			s.SetScaleToH(noteHeight)
			s.Locate(S.HitPosition, S.FieldPosition, draws.CenterMiddle)
			skin.Note[kind][size] = s
		}
		a := draws.NewAnimation(skin.overlay[size])
		for frame := range a {
			a[frame].SetScaleToH(noteHeight)
			a[frame].Locate(S.HitPosition, S.FieldPosition, draws.CenterMiddle)
		}
		skin.Overlay[size] = a
		{
			s := draws.NewSprite(head)
			s.SetScaleToH(noteHeight)
			s.Locate(S.HitPosition, S.FieldPosition, draws.RightMiddle)
			skin.Head[size] = s
		}
		{
			s := draws.NewSprite(tail)
			s.SetScaleToH(noteHeight)
			s.Locate(S.HitPosition, S.FieldPosition, draws.LeftMiddle)
			skin.Tail[size] = s
		}
		{
			s := draws.NewSprite(body)
			s.SetScaleToH(noteHeight)
			s.Locate(S.HitPosition, S.FieldPosition, draws.LeftMiddle)
			s.Filter = ebiten.FilterNearest
			skin.Body[size] = s
		}
	}
	{
		s := draws.NewSprite(skin.dot)
		s.MultiplyScale(S.DotScale)
		s.Locate(S.HitPosition, S.FieldPosition, draws.CenterMiddle)
		skin.Dot = s
	}
	skin.Shake = NewShakeSprites(skin.note)

	// Key sprites are overlapped at each side.
	for i, kname := range []string{"in", "out"} {
		name := fmt.Sprintf("drum/key/%s.png", kname)
		img := draws.NewImageFromFile(fsys, name)
		if !img.IsValid() {
			skin.key = DefaultSkin.key
			break
		}
		skin.key[i] = img
	}
	const (
		in = iota
		out
	)
	var (
		key = []draws.Image{
			draws.NewImageXFlipped(skin.key[out]),
			skin.key[in],
			draws.NewImageXFlipped(skin.key[in]),
			skin.key[out],
		}
		keyFieldSize draws.Vector2
	)
	for k, img := range key {
		s := draws.NewSprite(img)
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
		s := draws.NewSprite(src)
		s.Locate(0, S.FieldPosition, draws.LeftMiddle)
		skin.KeyField = s
	}
	for i, kname := range []string{"idle", "yes", "no", "high"} {
		name := fmt.Sprintf("drum/dancer/%s", kname)
		imgs := draws.NewImagesFromFile(fsys, name)
		if len(imgs) == 1 && !imgs[0].IsValid() {
			skin.dancer = DefaultSkin.dancer
			break
		}
		skin.dancer[i] = imgs
	}
	for i, imgs := range skin.dancer {
		a := draws.NewAnimation(imgs)
		for frame := range a {
			a[frame].MultiplyScale(S.DancerScale)
			a[frame].Locate(S.DancerPositionX, S.DancerPositionY, draws.CenterMiddle)
		}
		skin.Dancer[i] = a
	}
	skin.score = mode.UserSkin.Score
	{ // Position of combo is dependent on widths of key sprite.
		var imgs [10]draws.Image
		for i := 0; i < 10; i++ {
			img := draws.NewImageFromFile(fsys, fmt.Sprintf("combo/%d.png", i))
			if !img.IsValid() {
				skin.combo = DefaultSkin.combo
				break
			}
			imgs[i] = img
		}
		skin.combo = imgs
		for i := 0; i < 10; i++ {
			s := draws.NewSprite(imgs[i])
			s.MultiplyScale(S.ComboScale)
			s.Locate(keyFieldSize.X/2, S.FieldPosition, draws.CenterMiddle)
			skin.Combo[i] = s
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
		s := draws.NewSprite(outerImage)
		s.SetScaleToH(scale * S.regularNoteHeight)
		s.Locate(S.HitPosition, S.FieldPosition, draws.CenterMiddle)
		sprites[outer] = s
	}
	{
		s := draws.NewSprite(innerImage)
		s.SetScaleToH((scale + thickness) * S.regularNoteHeight)
		s.Locate(S.HitPosition, S.FieldPosition, draws.CenterMiddle)
		sprites[inner] = s
	}
	return
}
