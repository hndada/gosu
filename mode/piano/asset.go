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
	ScoreSprites            [13]draws.Sprite // numbers with sign (. , %)
	ComboSprites            [10]draws.Sprite
	JudgmentAnimations      [4]draws.Animation
	DefaultHitSoundStreamer audios.StreamSeekCloser
	DefaultHitSoundFormat   audios.Format

	// asset for a field
	FieldSprite draws.Sprite
	HintSprite  draws.Sprite // bottom: hit position
	BarSprite   draws.Sprite // bottom: hit position

	// asset for each key
	KeyKindNoteTypeAnimations [][4]draws.Animation // bottom: hit position
	KeySprites                [][2]draws.Sprite    // top: hit position
	KeyLightingSprites        []draws.Sprite
	KeyLightingColors         []color.Color
	HitLightingAnimations     []draws.Animation
	HoldLightingAnimations    []draws.Animation
}

// Todo: should Scratch be excluded from fieldWidth?
func NewAsset(cfg *Config, fsys fs.FS, keyCount int, scratchMode ScratchMode) *Asset {
	asset := &Asset{}

	fieldWidth := cfg.FieldWidth(keyCount, scratchMode)
	keyXs := cfg.KeyXs(keyCount, scratchMode)
	keyWidths := cfg.KeyWidths(keyCount, scratchMode)
	keyKinds := KeyKinds(keyCount, scratchMode)

	asset.setScoreSprites(cfg, fsys)
	asset.setComboSprites(cfg, fsys)
	asset.setJudgmentAnimations(cfg, fsys)
	asset.setDefaultHitSound(cfg, fsys)

	asset.setFieldSprite(cfg, fsys, fieldWidth)
	asset.setHintSprite(cfg, fsys, fieldWidth)
	asset.setBarSprite(cfg, fsys, fieldWidth)

	asset.setKeyKindNoteTypeAnimations(cfg, fsys, keyXs, keyWidths, keyKinds)
	asset.setKeySprites(cfg, fsys, keyXs, keyWidths)
	asset.setKeyLightingSprites(cfg, fsys, keyXs, keyWidths)
	asset.setKeyLightingColors(cfg, fsys, keyKinds)
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
		sprite := draws.NewSpriteFromFile(fsys, fmt.Sprintf("combo/%d.png", i))
		sprite.MultiplyScale(cfg.ComboSpriteScale)
		sprite.Locate(cfg.FieldPosition, cfg.ComboPosition, draws.CenterMiddle)
		sprites[i] = sprite
	}
	asset.ComboSprites = sprites
}

func (asset *Asset) setJudgmentAnimations(cfg *Config, fsys fs.FS) {
	var anims [4]draws.Animation
	for i, name := range []string{"kool", "cool", "good", "miss"} {
		anim := draws.NewAnimationFromFile(fsys, fmt.Sprintf("piano/judgment/%s.png", name))
		for frame := range anim {
			anim[frame].MultiplyScale(cfg.JudgmentSpriteScale)
			anim[frame].Locate(cfg.FieldPosition, cfg.JudgmentPosition, draws.CenterMiddle)
		}
		anims[i] = anim
	}
	asset.JudgmentAnimations = anims
}

func (asset *Asset) setDefaultHitSound(cfg *Config, fsys fs.FS) {
	streamer, format, _ := audios.DecodeFromFile(fsys, "piano/sound/hit.wav")
	asset.DefaultHitSoundStreamer = streamer
	asset.DefaultHitSoundFormat = format
}

func (asset *Asset) setBarSprite(cfg *Config, fsys fs.FS, fieldWidth float64) {
	img := draws.NewImage(fieldWidth, 1)
	img.Fill(color.White)

	sprite := draws.NewSprite(img)
	sprite.Locate(cfg.FieldPosition, cfg.HitPosition, draws.CenterBottom)
	asset.BarSprite = sprite
}

func (asset *Asset) setHintSprite(cfg *Config, fsys fs.FS, fieldWidth float64) {
	img := draws.NewImageFromFile(fsys, "piano/stage/hint.png")

	sprite := draws.NewSprite(img)
	sprite.SetSize(fieldWidth, cfg.HintHeight)
	sprite.Locate(cfg.FieldPosition, cfg.HitPosition, draws.CenterBottom)
	asset.HintSprite = sprite
}

func (asset *Asset) setFieldSprite(cfg *Config, fsys fs.FS, fieldWidth float64) {
	img := draws.NewImage(fieldWidth, cfg.ScreenSize.Y)
	img.Fill(color.NRGBA{0, 0, 0, uint8(255 * cfg.FieldOpacity)})

	sprite := draws.NewSprite(img)
	sprite.Locate(cfg.FieldPosition, 0, draws.CenterTop)
	asset.FieldSprite = sprite
}

// When note/normal image is not found, use default's note/normal.
// When note/head image is not found, use user's note/normal.
// When note/tail image is not found, let it be blank.
// When note/body image is not found, use user's note/normal.
func (asset *Asset) setKeyKindNoteTypeAnimations(cfg *Config, fsys fs.FS, keyXs []float64, keyWidths []float64, keyKinds []KeyKind) {
	var keyKindNoteTypeFrames [4][4]draws.Frames
	// Todo: 2nd mid -> tip
	for keyKind, kkname := range []string{"one", "two", "mid", "mid"} {
		for noteType, ntname := range []string{"normal", "head", "tail", "body"} {
			name := fmt.Sprintf("piano/note/%s/%s.png", kkname, ntname)
			frames := draws.NewFramesFromFilename(fsys, name)
			keyKindNoteTypeFrames[keyKind][noteType] = frames
		}
	}

	anims := make([][4]draws.Animation, len(keyWidths))
	for k := range anims {
		noteTypeFrames := keyKindNoteTypeFrames[keyKinds[k]]
		for noteType, frames := range noteTypeFrames {
			anim := draws.NewAnimation(frames[:])
			for frame := range anim {
				anim[frame].SetSize(keyWidths[k], cfg.NoteHeigth)
				anim[frame].Locate(keyXs[k], cfg.HitPosition, draws.CenterBottom)
			}
			anims[k][noteType] = anim
		}
	}
	asset.KeyKindNoteTypeAnimations = anims
}

func (asset *Asset) setKeySprites(cfg *Config, fsys fs.FS, keyXs []float64, keyWidths []float64) {
	imgs := [2]draws.Image{
		draws.NewImageFromFile(fsys, "piano/key/up.png"),
		draws.NewImageFromFile(fsys, "piano/key/down.png"),
	}
	sprites := make([][2]draws.Sprite, len(keyXs))
	for k := range sprites {
		for i, img := range imgs {
			sprite := draws.NewSprite(img)
			sprite.SetSize(keyWidths[k], cfg.ScreenSize.Y-cfg.HitPosition)
			sprite.Locate(keyXs[k], cfg.HitPosition, draws.CenterTop)
			sprites[k][i] = sprite
		}
	}
	asset.KeySprites = sprites
}

func (asset *Asset) setKeyLightingSprites(cfg *Config, fsys fs.FS, keyXs []float64, keyWidths []float64) {
	img := draws.NewImageFromFile(fsys, "piano/key/lighting.png")
	sprites := make([]draws.Sprite, len(keyXs))
	for k := range sprites {
		s := draws.NewSprite(img)
		s.MultiplyScale(keyWidths[k] / s.Width())
		s.Locate(keyXs[k], cfg.HitPosition, draws.CenterBottom) // -HintHeight
		sprites[k] = s
	}
	asset.KeyLightingSprites = sprites
}

func (asset *Asset) setKeyLightingColors(cfg *Config, fsys fs.FS, keyKinds []KeyKind) {
	colors := make([]color.Color, len(keyKinds))
	for k := range colors {
		colors[k] = cfg.KeyKindLightingColors[keyKinds[k]]
	}
	asset.KeyLightingColors = colors
}

func (asset *Asset) setHitLightingAnimations(cfg *Config, fsys fs.FS, keyXs []float64) {
	imgs := draws.NewFramesFromFilename(fsys, "piano/lighting/hit")
	anims := make([]draws.Animation, len(keyXs))
	for k := range anims {
		anim := draws.NewAnimation(imgs)
		for frame := range anim {
			anim[frame].MultiplyScale(cfg.LightingSpriteScale)
			anim[frame].Locate(keyXs[k], cfg.HitPosition, draws.CenterMiddle) // -HintHeight
		}
		anims[k] = anim
	}
	asset.HitLightingAnimations = anims
}

func (asset *Asset) setHoldLightingAnimations(cfg *Config, fsys fs.FS, keyXs []float64) {
	imgs := draws.NewFramesFromFilename(fsys, "piano/lighting/hold")
	anims := make([]draws.Animation, len(keyXs))
	for k := range anims {
		anim := draws.NewAnimation(imgs)
		for frame := range anim {
			anim[frame].MultiplyScale(cfg.LightingSpriteScale)
			anim[frame].Locate(keyXs[k], cfg.HitPosition-cfg.HintHeight/2, draws.CenterMiddle)
		}
		anims[k] = anim
	}
	asset.HoldLightingAnimations = anims
}
