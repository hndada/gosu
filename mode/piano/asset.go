package piano

import (
	"fmt"
	"image/color"
	"io/fs"

	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/mode"
)

type Asset struct {
	KeyCount int // Each key count has different asset.

	Score    [13]draws.Sprite // numbers with sign (. , %)
	Combo    [10]draws.Sprite
	Judgment [4]draws.Animation
	Sound    []byte

	Field draws.Sprite
	Hint  draws.Sprite // bottom: hit position
	Bar   draws.Sprite // bottom: hit position

	Note         [][4]draws.Animation // bottom: hit position
	Key          [][2]draws.Sprite    // top: hit position
	KeyLighting  []draws.Sprite
	HitLighting  []draws.Animation
	HoldLighting []draws.Animation
}

// Todo: should Scratch be excluded from fieldWidth?
func NewAsset(fsys fs.FS, cfg *Config, keyCount int, sm ScratchMode) *Asset {
	fieldWidth := cfg.FieldWidth(keyCount, sm)
	keyWidths := cfg.KeyWidths(keyCount, sm)
	keyXs := cfg.KeyXs(keyCount, sm)
	keyTypes := KeyTypes(keyCount, sm)

	return &Asset{
		KeyCount: keyCount,
		Score:    newScoreSprites(fsys, cfg.ScreenSize, cfg.ScoreScale),
		Combo:    newComboSprites(fsys, cfg),
		Judgment: newJudgmentSprites(fsys, cfg),
		Sound:    newSound(fsys, cfg),

		Field: newFieldSprite(fsys, cfg, fieldWidth),
		Hint:  newHintSprite(fsys, cfg, fieldWidth),
		Bar:   newBarSprite(fsys, cfg, fieldWidth),

		Note:         newNoteSprites(fsys, cfg, keyWidths, keyXs, keyTypes),
		Key:          newKeySprites(fsys, cfg, keyWidths, keyXs),
		KeyLighting:  newKeyLightingSprites(fsys, cfg, keyWidths, keyXs),
		HitLighting:  newHitLightingSprites(fsys, cfg, keyWidths, keyXs),
		HoldLighting: newHoldLightingSprites(fsys, cfg, keyWidths, keyXs),
	}
}
func newScoreSprites(fsys fs.FS, cfg *Config) [13]draws.Sprite {
	return mode.NewScoreSprites(fsys, cfg)
}
func newComboSprites(fsys fs.FS, cfg *Config) [10]draws.Sprite {
	var sprites [10]draws.Sprite
	for i := 0; i < 10; i++ {
		s := draws.LoadSprite(fsys, fmt.Sprintf("combo/%d.png", i))
		s.MultiplyScale(cfg.ComboScale)
		s.Locate(cfg.FieldPosition, cfg.ComboPosition, draws.CenterMiddle)
		sprites[i] = s
	}
	return sprites
}
func newJudgmentSprites(fsys fs.FS, cfg *Config) [4]draws.Animation {
	var anims [4]draws.Animation
	for i, name := range []string{"kool", "cool", "good", "miss"} {
		a := draws.NewAnimation(fsys, fmt.Sprintf("piano/judgment/%s", name))
		for frame := range a {
			a[frame].MultiplyScale(cfg.JudgmentScale)
			a[frame].Locate(cfg.FieldPosition, cfg.JudgmentPosition, draws.CenterMiddle)
		}
		anims[i] = a
	}
	return anims
}
func newSound(fsys fs.FS, cfg *Config) []byte {
	return audios.NewSound(fsys, "piano/sound.wav")
}

func newBarSprite(fsys fs.FS, cfg *Config, w float64) draws.Sprite {
	img := draws.NewImage(w, 1)
	img.Fill(color.White)
	s := draws.NewSprite(img)
	s.Locate(cfg.FieldPosition, cfg.HitPosition, draws.CenterBottom)
	return s
}
func newHintSprite(fsys fs.FS, cfg *Config, w float64) draws.Sprite {
	img := draws.LoadImage(fsys, "piano/stage/hint.png")
	s := draws.NewSprite(img)
	s.SetSize(w, cfg.HintHeight)
	s.Locate(cfg.FieldPosition, cfg.HitPosition, draws.CenterBottom)
	return s
}
func newFieldSprite(fsys fs.FS, cfg *Config, w float64) draws.Sprite {
	img := draws.NewImage(w, ScreenSizeY)
	img.Fill(color.NRGBA{0, 0, 0, uint8(255 * cfg.FieldOpaque)})
	s := draws.NewSprite(img)
	s.Locate(cfg.FieldPosition, 0, draws.CenterTop)
	return s
}

// When note/normal image is not found, use default's note/normal.
// When note/head image is not found, use user's note/normal.
// When note/tail image is not found, let it be blank.
// When note/body image is not found, use user's note/normal.
func newNoteSprites(fsys fs.FS, cfg *Config, ws []float64, xs []float64, keyTypes []KeyType) [][4]draws.Animation {
	var keyNoteImgs [4][4][]draws.Image                            // key and note images
	for i, keyType := range []string{"one", "two", "mid", "mid"} { // Todo: 2nd mid -> tip
		for j, noteType := range []string{"normal", "head", "tail", "body"} {
			name := fmt.Sprintf("piano/note/%s/%s", keyType, noteType)
			keyNoteImgs[i][j] = draws.LoadImages(fsys, name)
		}
	}

	anims := make([][4]draws.Animation, len(ws))
	for k := range anims {
		keyType := keyTypes[k]
		noteImgs := keyNoteImgs[keyType]
		for j, imgs := range noteImgs {
			a := draws.NewAnimationFromImages(imgs[:])
			for frame := range a {
				a[frame].SetSize(ws[k], cfg.NoteHeigth)
				a[frame].Locate(xs[k], cfg.HitPosition, draws.CenterBottom)
			}
			anims[k][j] = a
		}
	}
	return anims
}

func newKeySprites(fsys fs.FS, cfg *Config, ws []float64, xs []float64) [][2]draws.Sprite {
	imgs := [2]draws.Image{
		draws.LoadImage(fsys, "piano/key/up.png"),
		draws.LoadImage(fsys, "piano/key/down.png"),
	}
	sprites := make([][2]draws.Sprite, len(ws))
	for k := range sprites {
		for i, img := range imgs {
			s := draws.NewSprite(img)
			s.SetSize(ws[k], ScreenSizeY-cfg.HitPosition)
			s.Locate(xs[k], cfg.HitPosition, draws.CenterTop)
			sprites[k][i] = s
		}
	}
	return sprites
}
func newKeyLightingSprites(fsys fs.FS, cfg *Config, ws []float64, xs []float64) []draws.Sprite {
	img := draws.LoadImage(fsys, "piano/key/lighting.png")
	sprites := make([]draws.Sprite, len(ws))
	for k := range sprites {
		s := draws.NewSprite(img)
		s.SetScaleToW(ws[k])
		s.Locate(xs[k], cfg.HitPosition, draws.CenterBottom) // -HintHeight
		sprites[k] = s
	}
	return sprites
}
func newHitLightingSprites(fsys fs.FS, cfg *Config, ws []float64, xs []float64) []draws.Animation {
	imgs := draws.LoadImages(fsys, "piano/lighting/hit")
	anims := make([]draws.Animation, len(ws))
	for k := range anims {
		a := draws.NewAnimationFromImages(imgs)
		for frame := range a {
			a[frame].MultiplyScale(cfg.LightingScale)
			a[frame].Locate(xs[k], cfg.HitPosition, draws.CenterMiddle) // -HintHeight
		}
		anims[k] = a
	}
	return anims
}
func newHoldLightingSprites(fsys fs.FS, cfg *Config, ws []float64, xs []float64) []draws.Animation {
	imgs := draws.LoadImages(fsys, "piano/lighting/hold")
	anims := make([]draws.Animation, len(ws))
	for k := range anims {
		a := draws.NewAnimationFromImages(imgs)
		for frame := range a {
			a[frame].MultiplyScale(cfg.LightingScale)
			a[frame].Locate(xs[k], cfg.HitPosition-cfg.HintHeight/2, draws.CenterMiddle)
		}
		anims[k] = a
	}
	return anims
}
