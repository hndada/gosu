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
	ColorGray   = color.NRGBA{67, 67, 67, 255}
)

const (
	ShakeNote = iota
	ShakeSpin
	ShakeLimit
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

// https://osu.ppy.sh/wiki/en/Skinning/osu%21taiko
type Skin struct {
	FieldSprite draws.Sprite
	HintSprite  draws.Sprite
	BarSprite   draws.Sprite

	// First [2] are for big notes.
	JudgmentSprites [2][3]draws.Sprite // 3 Judgments.
	// RedSprites      [2]draws.Sprite
	// BlueSprites     [2]draws.Sprite
	NoteSprites    [2][2]draws.Sprite // Second [2] are for Red / Blue.
	HeadSprites    [2]draws.Sprite    // Overlay will be drawn during game play.
	TailSprites    [2]draws.Sprite
	OverlaySprites [2][2]draws.Sprite // 2 Overlays.
	BodySprites    [2]draws.Sprite
	DotSprite      draws.Sprite
	ShakeSprites   [3]draws.Sprite

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
		s.SetScale(FieldHeight / s.H())
		s.SetPosition(0, FieldPosition, draws.OriginLeftCenter)
		skin.FieldSprite = s
	}
	{
		s := draws.NewSpriteFromImage(noteImage)
		s.SetScale(regularNoteHeight / s.H())
		s.SetPosition(HitPosition, FieldPosition, draws.OriginCenter)
		skin.HintSprite = s
	}
	{
		src := ebiten.NewImage(1, int(FieldInnerHeight))
		src.Fill(color.NRGBA{255, 255, 255, 255}) // White
		s := draws.NewSpriteFromImage(src)
		s.SetPosition(HitPosition, FieldPosition, draws.OriginCenter)
		skin.BarSprite = s
	}
	var (
		rollEndImage = draws.NewImage("skin/drum/note/roll/end.png")
		rollMidImage = draws.NewImage("skin/drum/note/roll/mid.png")
	)
	for i, sname := range []string{"regular", "big"} {
		noteHeight := regularNoteHeight
		if sname == "big" {
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
			s.SetPosition(HitPosition, FieldPosition, draws.OriginCenter)
			skin.JudgmentSprites[i][j] = s
		}
		for j, clr := range []color.NRGBA{ColorRed, ColorBlue} {
			s := draws.NewSpriteFromImage(noteImage)
			s.SetScale(noteHeight / s.H())
			s.SetColor(clr)
			s.SetPosition(HitPosition, FieldPosition, draws.OriginCenter)
			skin.NoteSprites[i][j] = s
		}
		// {
		// 	s := draws.NewSpriteFromImage(noteImage)
		// 	s.SetScale(noteHeight / s.H())
		// 	s.SetColor(ColorRed)
		// 	s.SetPosition(HitPosition, FieldPosition, draws.OriginCenter)
		// 	skin.RedSprites[i] = s
		// }
		// {
		// 	s := skin.RedSprites[i]
		// 	s.SetColor(ColorBlue)
		// 	skin.BlueSprites[i] = s
		// }
		{
			s := draws.NewSpriteFromImage(rollEndImage)
			s.SetScale(noteHeight / s.H())
			s.SetPosition(HitPosition, FieldPosition, draws.OriginLeftCenter)
			s.SetColor(ColorYellow)
			skin.TailSprites[i] = s
		}
		{
			s := draws.NewSpriteFromImage(rollEndImage)
			ratio := noteHeight / s.H()
			s.SetScaleXY(-ratio, ratio, ebiten.FilterLinear) // Goes flipped.
			s.SetPosition(HitPosition, FieldPosition, draws.OriginRightCenter)
			s.SetColor(ColorYellow)
			skin.HeadSprites[i] = s
		}
		{
			var path string
			for j := 0; j < 2; j++ {
				if _, err := os.Stat(path); !os.IsNotExist(err) { // Two overlays.
					path = fmt.Sprintf("skin/drum/overlay/%s/%d.png", sname, j)
				} else {
					path = fmt.Sprintf("/skin/drum/overlay/%s", sname)
				}
				s := draws.NewSprite(path)
				s.SetScale(noteHeight / s.H())
				s.SetPosition(HitPosition, FieldPosition, draws.OriginCenter)
				skin.OverlaySprites[i][j] = s
			}
		}
		{
			s := draws.NewSpriteFromImage(rollMidImage)
			s.SetScale(noteHeight / s.H())
			s.SetPosition(HitPosition, FieldPosition, draws.OriginLeftCenter)
			s.SetColor(ColorYellow)
			skin.BodySprites[i] = s
		}
	}
	{
		s := draws.NewSprite("skin/drum/note/roll/dot.png")
		s.SetScale(DotScale)
		s.SetPosition(HitPosition, FieldPosition, draws.OriginCenter)
		s.SetColor(ColorYellow)
		skin.DotSprite = s
	}
	for i, name := range []string{"note", "spin", "limit"} {
		path := fmt.Sprintf("skin/drum/note/shake/%s.png", name)
		s := draws.NewSprite(path)
		if name == "note" {
			s.SetScale(regularNoteHeight / s.H())
			s.SetPosition(HitPosition, FieldPosition, draws.OriginCenter)
		} else {
			s.SetScale(ShakeScale)
			s.SetPosition(ShakePosX, ShakePosY, draws.OriginCenter)
		}
		skin.ShakeSprites[i] = s
	}

	// Position of combo is dependent on widths of key sprite.
	// Key sprites are overlapped at each side.
	{
		s := draws.NewSprite("skin/drum/key/in.png")
		s.SetScale(KeyScale)
		s.SetPosition(0, FieldPosition, draws.OriginLeftCenter)
		skin.KeySprites[LeftRed] = s
		keyCenter = s.W()
	}
	{
		s := draws.NewSprite("skin/drum/key/out.png")
		s.SetScale(KeyScale)
		s.SetPosition(keyCenter, FieldPosition, draws.OriginLeftCenter)
		skin.KeySprites[RightBlue] = s
	}
	{
		s := skin.KeySprites[RightBlue]
		s.Flip(true, false)
		s.SetPosition(0, FieldPosition, draws.OriginLeftCenter)
		skin.KeySprites[LeftBlue] = s
	}
	{
		s := skin.KeySprites[LeftRed]
		s.Flip(true, false)
		s.SetPosition(keyCenter, FieldPosition, draws.OriginLeftCenter)
		skin.KeySprites[RightRed] = s
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
			s.SetPosition(DancerPosX, DancerPosY, draws.OriginCenter)
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
		s.SetPosition(keyCenter, FieldPosition, draws.OriginCenter)
		skin.ComboSprites[i] = s
	}
	for i := 0; i < 10; i++ {
		s := draws.NewSpriteFromImage(comboImages[i])
		s.SetScale(DotCountScale)
		s.SetPosition(HitPosition, FieldPosition, draws.OriginCenter)
		skin.DotCountSprites[i] = s
	}
	for i := 0; i < 10; i++ {
		s := draws.NewSpriteFromImage(comboImages[i])
		s.SetScale(ShakeCountScale)
		pos := ShakePosY + s.H()*ShakeCountPosition
		s.SetPosition(ShakePosX, pos, draws.OriginCenterTop)
		skin.ShakeCountSprites[i] = s
	}
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
