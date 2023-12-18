package piano

import (
	"image/color"

	"github.com/hndada/gosu/draws"
)

type ScreenConfig struct {
	Size draws.Vector2
	TPS  int
	FPS  int
}
type MusicOffset struct {
	MusicVolume float64
	SoundVolume float64
	Offset      int32
}

type BarConfig struct {
	*StageConfig
	Height float64
}
type NotesConfig struct {
	*StageConfig
	SpeedScale            float64
	Heigth                float64 // Applies to all types of notes.
	Colors                [4]color.NRGBA
	TailNoteExtraDuration int32
	LongNoteBodyStyle     int // Stretch or Attach.
	UpsideDown            bool
}

func NewConfig(ScreenSize draws.Vector2) *Config {
	return &Config{
		SpeedScale:        1.0,
		HitPosition:       0.90 * ScreenSize.Y,
		TailExtraDuration: 0,

		NoteHeigth:    0.03 * ScreenSize.Y, // 0.03: 27px
		FieldPosition: 0.50 * ScreenSize.X,
		ComboPosition: 0.40 * ScreenSize.Y,
		ScratchColor:  color.NRGBA{224, 0, 0, 255},
		FieldOpacity:  0.8,

		ScoreSpriteScale: 0.65,
		ComboSpriteScale: 0.75,
		ComboDigitGap:    -1,                  // unit: pixel
		HintHeight:       0.05 * ScreenSize.Y, // 0.06: 45px
	}
}

func (cfg Config) KeyWidths(keyCount int, sm ScratchMode) []float64 {
	ws := make([]float64, keyCount)
	scale := cfg.FieldWidthScales[keyCount]
	for k, kk := range KeyKinds(keyCount, sm) {
		w := cfg.KeyKindWidths[kk] // Todo: math.Ceil()?
		ws[k] = w * scale
	}
	return ws
}

func (cfg Config) FieldWidth(keyCount int, sm ScratchMode) float64 {
	ws := cfg.KeyWidths(keyCount, sm)
	var fw float64
	for _, width := range ws {
		fw += width
	}
	return fw
}

// KeyXs returns centered x positions.
func (cfg Config) KeyXs(keyCount int, sm ScratchMode) []float64 {
	xs := make([]float64, keyCount)
	ws := cfg.KeyWidths(keyCount, sm)
	x := cfg.FieldPosition - cfg.FieldWidth(keyCount, sm)/2
	for k, w := range ws {
		x += w / 2
		xs[k] = x
		x += w / 2
	}
	return xs
}

// NoteExposureDuration returns time in milliseconds
// that cursor takes to move 1 logical pixel.
func (cfg Config) NoteExposureDuration(speed float64) int32 {
	return int32(cfg.HitPosition / speed)
}
