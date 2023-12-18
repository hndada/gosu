package piano

import mode "github.com/hndada/gosu/mode2"

type StageConfig struct {
	KeyCount       int
	ScratchMode    ScratchMode
	Widths         map[int]float64
	KeyKindWidths  [4]float64
	FieldPositionX float64
	FieldOpacity   float64
	HintHeight     float64
	HintPositionY  float64
}

func NewStageConfig(screen mode.ScreenConfig, keyCount int, scratchMode ScratchMode) StageConfig {
	return StageConfig{
		KeyCount:    keyCount,
		ScratchMode: scratchMode,
		Widths: map[int]float64{
			4: 0.50 * screen.Size.X,
			5: 0.55 * screen.Size.X,
			6: 0.60 * screen.Size.X,
			7: 0.65 * screen.Size.X,
			8: 0.70 * screen.Size.X,
			9: 0.75 * screen.Size.X,
		},
		KeyKindWidths: [4]float64{
			0.055 * screen.Size.X, // One
			0.055 * screen.Size.X, // Two
			0.055 * screen.Size.X, // Mid
			0.055 * screen.Size.X, // Tip
		},
		FieldPositionX: 0.50 * screen.Size.X,
		FieldOpacity:   0.8,
		HintHeight:     0.05 * screen.Size.Y,
		HintPositionY:  0.90 * screen.Size.Y,
	}
}

func (cfg StageConfig) Width() float64 {
	const defaultKeyCount = 4
	if w, ok := cfg.Widths[cfg.KeyCount]; ok {
		return w
	}
	return cfg.Widths[defaultKeyCount]
}

func (cfg StageConfig) KeyWidths() (ws []float64) {
	ws = make([]float64, cfg.KeyCount)
	for k, kk := range KeyKinds(cfg.KeyCount, cfg.ScratchMode) {
		w := cfg.KeyKindWidths[kk]
		ws[k] = w
	}

	// Adjust key width to fit the stage width.
	var rawSum float64
	for _, w := range ws {
		rawSum += w
	}
	scale := cfg.Width() / rawSum
	for k := range ws {
		ws[k] *= scale
	}
	return ws
}

// KeyPositionXs returns centered x positions.
func (cfg StageConfig) KeyPositionXs() (xs []float64) {
	xs = make([]float64, cfg.KeyCount)
	ws := cfg.KeyWidths()
	x := cfg.FieldPositionX - cfg.Width()/2
	for k, w := range ws {
		x += w / 2
		xs[k] = x
		x += w / 2
	}
	return
}

// NoteExposureDuration returns time in milliseconds
// that cursor takes to move 1 logical pixel.
func (cfg StageConfig) NoteExposureDuration(speed float64) int32 {
	return int32(cfg.HintPositionY / speed)
}
