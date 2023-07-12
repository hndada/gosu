package piano

import (
	"fmt"
	"image/color"
	"io/fs"

	"github.com/faiface/beep"
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/mode"
)

// For Asset.KeysUpDowns
const (
	keyUp = iota
	keyDown
)

// Each key count has different asset.
var Assets = make(map[int]*Asset)

type Asset struct {
	// passed arguments
	cfg         Config
	fsys        fs.FS
	keyCount    int
	scratchMode ScratchMode

	// derived values from arguments
	fieldWidth float64
	keyWidths  []float64
	keyXs      []float64
	keyTypes   []KeyType

	// asset that are not affected by key count
	ScoreNumbers  [13]draws.Sprite // numbers with sign (. , %)
	ComboNumbers  [10]draws.Sprite
	JudgmentKinds [4]draws.Animation
	Sound         beep.Streamer // []byte

	// asset for a field
	Field draws.Sprite
	Hint  draws.Sprite // bottom: hit position
	Bar   draws.Sprite // bottom: hit position

	// asset for each key
	NoteTypes     [][4]draws.Animation // bottom: hit position
	KeysUpDowns   [][2]draws.Sprite    // top: hit position
	KeyLightings  []draws.Sprite
	HitLightings  []draws.Animation
	HoldLightings []draws.Animation
}

// Common argument is placed first
func NewAsset(cfg *Config, fsys fs.FS, keyCount int, scratchMode ScratchMode) *Asset {
	asset := &Asset{
		cfg:         *cfg,
		fsys:        fsys,
		keyCount:    keyCount,
		scratchMode: scratchMode,

		// Todo: should Scratch be excluded from fieldWidth?
		fieldWidth: cfg.FieldWidth(keyCount, scratchMode),
		keyWidths:  cfg.KeyWidths(keyCount, scratchMode),
		keyXs:      cfg.KeyXs(keyCount, scratchMode),
		keyTypes:   KeyTypes(keyCount, scratchMode),
	}

	asset.setScoreNumbers()
	asset.setComboNumbers()
	asset.setJudgmentKinds()
	asset.setSound()

	asset.setFieldSprite()
	asset.setHintSprite()
	asset.setBarSprite()

	asset.setNoteTypes()
	asset.setKeysUpDowns()
	asset.setKeyLightings()
	asset.setHitLightings()
	asset.setHoldLightings()
	return asset
}

func (asset *Asset) setScoreNumbers() {
	asset.ScoreNumbers = mode.NewScoreNumbers(asset.fsys, asset.cfg.ScreenSize, asset.cfg.ScoreScale)
}
func (asset *Asset) setComboNumbers() {
	var sprites [10]draws.Sprite
	for i := 0; i < 10; i++ {
		s := draws.NewSpriteFromFile(asset.fsys, fmt.Sprintf("combo/%d.png", i))
		s.MultiplyScale(asset.cfg.ComboScale)
		s.Locate(asset.cfg.FieldPosition, asset.cfg.ComboPosition, draws.CenterMiddle)
		sprites[i] = s
	}
	asset.ComboNumbers = sprites
}
func (asset *Asset) setJudgmentKinds() {
	var anims [4]draws.Animation
	for i, name := range []string{"kool", "cool", "good", "miss"} {
		a := draws.NewAnimationFromFile(asset.fsys, fmt.Sprintf("piano/judgment/%s", name))
		for frame := range a {
			a[frame].MultiplyScale(asset.cfg.JudgmentScale)
			a[frame].Locate(asset.cfg.FieldPosition, asset.cfg.JudgmentPosition, draws.CenterMiddle)
		}
		anims[i] = a
	}
	asset.JudgmentKinds = anims
}
func (asset *Asset) setSound() {
	streamer, _, _ := audios.DecodeFromFile(asset.fsys, "piano/sound.wav")
	asset.Sound = streamer
}

func (asset *Asset) setBarSprite() {
	img := draws.NewImage(asset.fieldWidth, 1)
	img.Fill(color.White)
	s := draws.NewSprite(img)
	s.Locate(asset.cfg.FieldPosition, asset.cfg.HitPosition, draws.CenterBottom)
	asset.Bar = s
}
func (asset *Asset) setHintSprite() {
	img := draws.NewImageFromFile(asset.fsys, "piano/stage/hint.png")
	s := draws.NewSprite(img)
	s.SetSize(asset.fieldWidth, asset.cfg.HintHeight)
	s.Locate(asset.cfg.FieldPosition, asset.cfg.HitPosition, draws.CenterBottom)
	asset.Hint = s
}
func (asset *Asset) setFieldSprite() {
	img := draws.NewImage(asset.fieldWidth, asset.cfg.ScreenSize.Y)
	img.Fill(color.NRGBA{0, 0, 0, uint8(255 * asset.cfg.FieldOpaque)})
	s := draws.NewSprite(img)
	s.Locate(asset.cfg.FieldPosition, 0, draws.CenterTop)
	asset.Field = s
}

// When note/normal image is not found, use default's note/normal.
// When note/head image is not found, use user's note/normal.
// When note/tail image is not found, let it be blank.
// When note/body image is not found, use user's note/normal.
func (asset *Asset) setNoteTypes() [][4]draws.Animation {
	var keyNoteImgs [4][4][]draws.Image                            // key and note images
	for i, keyType := range []string{"one", "two", "mid", "mid"} { // Todo: 2nd mid -> tip
		for j, noteType := range []string{"normal", "head", "tail", "body"} {
			name := fmt.Sprintf("piano/note/%s/%s", keyType, noteType)
			keyNoteImgs[i][j] = draws.NewImagesFromFile(asset.fsys, name)
		}
	}

	anims := make([][4]draws.Animation, len(asset.keyWidths))
	for k := range anims {
		keyType := asset.keyTypes[k]
		noteImgs := keyNoteImgs[keyType]
		for j, imgs := range noteImgs {
			a := draws.NewAnimation(imgs[:])
			for frame := range a {
				a[frame].SetSize(asset.keyWidths[k], asset.cfg.NoteHeigth)
				a[frame].Locate(asset.keyXs[k], asset.cfg.HitPosition, draws.CenterBottom)
			}
			anims[k][j] = a
		}
	}
	return anims
}
func (asset *Asset) setKeysUpDowns() [][2]draws.Sprite {
	imgs := [2]draws.Image{
		draws.NewImageFromFile(asset.fsys, "piano/key/up.png"),
		draws.NewImageFromFile(asset.fsys, "piano/key/down.png"),
	}
	sprites := make([][2]draws.Sprite, asset.keyCount)
	for k := range sprites {
		for i, img := range imgs {
			s := draws.NewSprite(img)
			s.SetSize(asset.keyWidths[k], asset.cfg.ScreenSize.Y-asset.cfg.HitPosition)
			s.Locate(asset.keyXs[k], asset.cfg.HitPosition, draws.CenterTop)
			sprites[k][i] = s
		}
	}
	return sprites
}
func (asset *Asset) setKeyLightings() []draws.Sprite {
	img := draws.NewImageFromFile(asset.fsys, "piano/key/lighting.png")
	sprites := make([]draws.Sprite, asset.keyCount)
	for k := range sprites {
		s := draws.NewSprite(img)
		s.SetScaleToW(asset.keyWidths[k])
		s.Locate(asset.keyXs[k], asset.cfg.HitPosition, draws.CenterBottom) // -HintHeight
		sprites[k] = s
	}
	return sprites
}
func (asset *Asset) setHitLightings() []draws.Animation {
	imgs := draws.NewImagesFromFile(asset.fsys, "piano/lighting/hit")
	anims := make([]draws.Animation, asset.keyCount)
	for k := range anims {
		a := draws.NewAnimation(imgs)
		for frame := range a {
			a[frame].MultiplyScale(asset.cfg.LightingScale)
			a[frame].Locate(asset.keyXs[k], asset.cfg.HitPosition, draws.CenterMiddle) // -HintHeight
		}
		anims[k] = a
	}
	return anims
}
func (asset *Asset) setHoldLightings() []draws.Animation {
	imgs := draws.NewImagesFromFile(asset.fsys, "piano/lighting/hold")
	anims := make([]draws.Animation, asset.keyCount)
	for k := range anims {
		a := draws.NewAnimation(imgs)
		for frame := range a {
			a[frame].MultiplyScale(asset.cfg.LightingScale)
			a[frame].Locate(asset.keyXs[k], asset.cfg.HitPosition-asset.cfg.HintHeight/2, draws.CenterMiddle)
		}
		anims[k] = a
	}
	return anims
}
