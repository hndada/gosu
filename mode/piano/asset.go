package piano

import (
	"fmt"
	"image/color"
	"io/fs"

	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/mode"
)

// For Asset.KeySprites
const (
	keyUp = iota
	keyDown
)

// All names of fields in Asset ends with their types.
type Asset struct {
	// asset that are not affected by key count
	ScoreSprites       [13]draws.Sprite // numbers with sign (. , %)
	ComboSprites       [10]draws.Sprite
	JudgmentAnimations [4]draws.Animation
	DefaultHitSound    audios.Sound

	// asset for a field
	FieldSprite draws.Sprite
	HintSprite  draws.Sprite // bottom: hit position
	BarSprite   draws.Sprite // bottom: hit position

	// asset for each key
	KeyNoteAnimations      [][4]draws.Animation // bottom: hit position
	KeySprites             [][2]draws.Sprite    // top: hit position
	KeyLightingSprites     []draws.Sprite
	HitLightingAnimations  []draws.Animation
	HoldLightingAnimations []draws.Animation
}

// Todo: should Scratch be excluded from fieldWidth?
func NewAsset(cfg *Config, fsys fs.FS, keyCount int, scratchMode ScratchMode) *Asset {
	asset := &Asset{}

	fieldWidth := cfg.FieldWidth(keyCount, scratchMode)
	keyXs := cfg.KeyXs(keyCount, scratchMode)
	keyWidths := cfg.KeyWidths(keyCount, scratchMode)
	keyTypes := KeyTypes(keyCount, scratchMode)

	asset.setScoreSprites(cfg, fsys)
	asset.setComboSprites(cfg, fsys)
	asset.setJudgmentAnimations(cfg, fsys)
	asset.setDefaultHitSound(cfg, fsys)

	asset.setFieldSprite(cfg, fsys, fieldWidth)
	asset.setHintSprite(cfg, fsys, fieldWidth)
	asset.setBarSprite(cfg, fsys, fieldWidth)

	asset.setKeyNoteAnimations(cfg, fsys, keyXs, keyWidths, keyTypes)
	asset.setKeySprites(cfg, fsys, keyXs, keyWidths)
	asset.setKeyLightingSprites(cfg, fsys, keyXs, keyWidths)
	asset.setHitLightingAnimations(cfg, fsys, keyXs)
	asset.setHoldLightingAnimations(cfg, fsys, keyXs)
	return asset
}

func (asset *Asset) setScoreSprites(cfg *Config, fsys fs.FS) {
	sprites := mode.NewScoreSprites(fsys, *cfg.ScreenSize, cfg.ScoreSpriteScale)
	asset.ScoreSprites = sprites
}

func (asset *Asset) setComboSprites(cfg *Config, fsys fs.FS) {
	var sprites [10]draws.Sprite
	for i := 0; i < 10; i++ {
		s := draws.NewSpriteFromFile(fsys, fmt.Sprintf("combo/%d.png", i))
		s.MultiplyScale(cfg.ComboSpriteScale)
		s.Locate(cfg.FieldPosition, cfg.ComboPosition, draws.CenterMiddle)
		sprites[i] = s
	}
	asset.ComboSprites = sprites
}

func (asset *Asset) setJudgmentAnimations(cfg *Config, fsys fs.FS) {
	var anims [4]draws.Animation
	for i, name := range []string{"kool", "cool", "good", "miss"} {
		a := draws.NewAnimationFromFile(fsys, fmt.Sprintf("piano/judgment/%s", name))
		for frame := range a {
			a[frame].MultiplyScale(cfg.JudgmentSpriteScale)
			a[frame].Locate(cfg.FieldPosition, cfg.JudgmentPosition, draws.CenterMiddle)
		}
		anims[i] = a
	}
	asset.JudgmentAnimations = anims
}

func (asset *Asset) setDefaultHitSound(cfg *Config, fsys fs.FS) {
	sound := audios.NewSound(fsys, "piano/sound.wav", cfg.SoundVolume)
	asset.DefaultHitSound = sound
}

func (asset *Asset) setBarSprite(cfg *Config, fsys fs.FS, fieldWidth float64) {
	img := draws.NewImage(fieldWidth, 1)
	img.Fill(color.White)
	s := draws.NewSprite(img)
	s.Locate(cfg.FieldPosition, cfg.HitPosition, draws.CenterBottom)
	asset.BarSprite = s
}

func (asset *Asset) setHintSprite(cfg *Config, fsys fs.FS, fieldWidth float64) {
	img := draws.NewImageFromFile(fsys, "piano/stage/hint.png")
	s := draws.NewSprite(img)
	s.SetSize(fieldWidth, cfg.HintHeight)
	s.Locate(cfg.FieldPosition, cfg.HitPosition, draws.CenterBottom)
	asset.HintSprite = s
}

func (asset *Asset) setFieldSprite(cfg *Config, fsys fs.FS, fieldWidth float64) {
	img := draws.NewImage(fieldWidth, cfg.ScreenSize.Y)
	img.Fill(color.NRGBA{0, 0, 0, uint8(255 * cfg.FieldOpaque)})
	s := draws.NewSprite(img)
	s.Locate(cfg.FieldPosition, 0, draws.CenterTop)
	asset.FieldSprite = s
}

// When note/normal image is not found, use default's note/normal.
// When note/head image is not found, use user's note/normal.
// When note/tail image is not found, let it be blank.
// When note/body image is not found, use user's note/normal.
func (asset *Asset) setKeyNoteAnimations(cfg *Config, fsys fs.FS,
	keyXs []float64, keyWidths []float64, keyTypes []KeyType) [][4]draws.Animation {

	var keyNoteImgs [4][4][]draws.Image                            // key and note images
	for i, keyType := range []string{"one", "two", "mid", "mid"} { // Todo: 2nd mid -> tip
		for j, noteType := range []string{"normal", "head", "tail", "body"} {
			name := fmt.Sprintf("piano/note/%s/%s", keyType, noteType)
			keyNoteImgs[i][j] = draws.NewImagesFromFile(fsys, name)
		}
	}

	anims := make([][4]draws.Animation, len(keyWidths))
	for k := range anims {
		keyType := keyTypes[k]
		noteImgs := keyNoteImgs[keyType]
		for j, imgs := range noteImgs {
			a := draws.NewAnimation(imgs[:])
			for frame := range a {
				a[frame].SetSize(keyWidths[k], cfg.NoteHeigth)
				a[frame].Locate(keyXs[k], cfg.HitPosition, draws.CenterBottom)
			}
			anims[k][j] = a
		}
	}
	return anims
}

func (asset *Asset) setKeySprites(cfg *Config, fsys fs.FS,
	keyXs []float64, keyWidths []float64) [][2]draws.Sprite {

	imgs := [2]draws.Image{
		draws.NewImageFromFile(fsys, "piano/key/up.png"),
		draws.NewImageFromFile(fsys, "piano/key/down.png"),
	}
	sprites := make([][2]draws.Sprite, len(keyXs))
	for k := range sprites {
		for i, img := range imgs {
			s := draws.NewSprite(img)
			s.SetSize(keyWidths[k], cfg.ScreenSize.Y-cfg.HitPosition)
			s.Locate(keyXs[k], cfg.HitPosition, draws.CenterTop)
			sprites[k][i] = s
		}
	}
	return sprites
}

func (asset *Asset) setKeyLightingSprites(cfg *Config, fsys fs.FS,
	keyXs []float64, keyWidths []float64) []draws.Sprite {

	img := draws.NewImageFromFile(fsys, "piano/key/lighting.png")
	sprites := make([]draws.Sprite, len(keyXs))
	for k := range sprites {
		s := draws.NewSprite(img)
		s.SetScaleToW(keyWidths[k])
		s.Locate(keyXs[k], cfg.HitPosition, draws.CenterBottom) // -HintHeight
		sprites[k] = s
	}
	return sprites
}

func (asset *Asset) setHitLightingAnimations(cfg *Config, fsys fs.FS,
	keyXs []float64) []draws.Animation {
	imgs := draws.NewImagesFromFile(fsys, "piano/lighting/hit")
	anims := make([]draws.Animation, len(keyXs))
	for k := range anims {
		a := draws.NewAnimation(imgs)
		for frame := range a {
			a[frame].MultiplyScale(cfg.LightingSpriteScale)
			a[frame].Locate(keyXs[k], cfg.HitPosition, draws.CenterMiddle) // -HintHeight
		}
		anims[k] = a
	}
	return anims
}

func (asset *Asset) setHoldLightingAnimations(cfg *Config, fsys fs.FS,
	keyXs []float64) []draws.Animation {
	imgs := draws.NewImagesFromFile(fsys, "piano/lighting/hold")
	anims := make([]draws.Animation, len(keyXs))
	for k := range anims {
		a := draws.NewAnimation(imgs)
		for frame := range a {
			a[frame].MultiplyScale(cfg.LightingSpriteScale)
			a[frame].Locate(keyXs[k], cfg.HitPosition-cfg.HintHeight/2, draws.CenterMiddle)
		}
		anims[k] = a
	}
	return anims
}
