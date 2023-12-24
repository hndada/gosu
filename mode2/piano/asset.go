package piano

import (
	"fmt"
	"image/color"
	"io/fs"

	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/draws"
)

type Asset struct {
	// asset that are not affected by key count
	ScoreSprites            [13]draws.Sprite
	ComboSprites            [10]draws.Sprite
	JudgmentAnimations      [4]draws.Animation
	DefaultHitSoundStreamer audios.StreamSeekCloser
	DefaultHitSoundFormat   audios.Format

	// asset for each key
	KeyKindNoteTypeAnimations [][4]draws.Animation // bottom: hit position
	KeySprites                [][2]draws.Sprite    // top: hit position
	KeyLightingSprites        []draws.Sprite
	KeyLightingColors         []color.Color
	HitLightingAnimations     []draws.Animation
	HoldLightingAnimations    []draws.Animation
}

func NewAsset(cfg *Config, fsys fs.FS, keyCount int, scratchMode ScratchMode) *Asset {
	asset := &Asset{}

	KeyPositionXs := cfg.KeyPositionXs(keyCount, scratchMode)
	keyWidths := cfg.KeyWidths(keyCount, scratchMode)
	keyKinds := KeyKinds(keyCount, scratchMode)

	asset.setDefaultHitSound(cfg, fsys)

	asset.setKeyKindNoteTypeAnimations(cfg, fsys, KeyPositionXs, keyWidths, keyKinds)
	asset.setKeySprites(cfg, fsys, KeyPositionXs, keyWidths)
	asset.setKeyLightingSprites(cfg, fsys, KeyPositionXs, keyWidths)
	asset.setKeyLightingColors(cfg, fsys, keyKinds)
	asset.setHitLightingAnimations(cfg, fsys, KeyPositionXs)
	asset.setHoldLightingAnimations(cfg, fsys, KeyPositionXs)
	return asset
}

func (asset *Asset) setDefaultHitSound(cfg *Config, fsys fs.FS) {
	streamer, format, _ := audios.DecodeFromFile(fsys, "piano/sound/hit.wav")
	asset.DefaultHitSoundStreamer = streamer
	asset.DefaultHitSoundFormat = format
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
func (asset *Asset) setKeyKindNoteTypeAnimations(cfg *Config, fsys fs.FS, KeyPositionXs []float64, keyWidths []float64, keyKinds []KeyKind) {
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
				anim[frame].Locate(KeyPositionXs[k], cfg.HitPosition, draws.CenterBottom)
			}
			anims[k][noteType] = anim
		}
	}
	asset.KeyKindNoteTypeAnimations = anims
}

func (asset *Asset) setKeySprites(cfg *Config, fsys fs.FS, KeyPositionXs []float64, keyWidths []float64) {
	imgs := [2]draws.Image{
		draws.NewImageFromFile(fsys, "piano/key/up.png"),
		draws.NewImageFromFile(fsys, "piano/key/down.png"),
	}
	sprites := make([][2]draws.Sprite, len(KeyPositionXs))
	for k := range sprites {
		for i, img := range imgs {
			sprite := draws.NewSprite(img)
			sprite.SetSize(keyWidths[k], cfg.ScreenSize.Y-cfg.HitPosition)
			sprite.Locate(KeyPositionXs[k], cfg.HitPosition, draws.CenterTop)
			sprites[k][i] = sprite
		}
	}
	asset.KeySprites = sprites
}

func (asset *Asset) setKeyLightingSprites(cfg *Config, fsys fs.FS, KeyPositionXs []float64, keyWidths []float64) {
	img := draws.NewImageFromFile(fsys, "piano/key/lighting.png")
	sprites := make([]draws.Sprite, len(KeyPositionXs))
	for k := range sprites {
		s := draws.NewSprite(img)
		s.MultiplyScale(keyWidths[k] / s.Width())
		s.Locate(KeyPositionXs[k], cfg.HitPosition, draws.CenterBottom) // -HintHeight
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

func (asset *Asset) setHitLightingAnimations(cfg *Config, fsys fs.FS, KeyPositionXs []float64) {
	imgs := draws.NewFramesFromFilename(fsys, "piano/lighting/hit")
	anims := make([]draws.Animation, len(KeyPositionXs))
	for k := range anims {
		anim := draws.NewAnimation(imgs)
		for frame := range anim {
			anim[frame].MultiplyScale(cfg.LightingSpriteScale)
			anim[frame].Locate(KeyPositionXs[k], cfg.HitPosition, draws.CenterMiddle) // -HintHeight
		}
		anims[k] = anim
	}
	asset.HitLightingAnimations = anims
}

func (asset *Asset) setHoldLightingAnimations(cfg *Config, fsys fs.FS, KeyPositionXs []float64) {
	imgs := draws.NewFramesFromFilename(fsys, "piano/lighting/hold")
	anims := make([]draws.Animation, len(KeyPositionXs))
	for k := range anims {
		anim := draws.NewAnimation(imgs)
		for frame := range anim {
			anim[frame].MultiplyScale(cfg.LightingSpriteScale)
			anim[frame].Locate(KeyPositionXs[k], cfg.HitPosition-cfg.HintHeight/2, draws.CenterMiddle)
		}
		anims[k] = anim
	}
	asset.HoldLightingAnimations = anims
}
