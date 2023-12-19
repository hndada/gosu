package piano

import mode "github.com/hndada/gosu/mode2"

// Objective: manage UI components with each own struct.
// NewXxxComponent() requires multiple arguments.
// XxxArgs is for wrapping required arguments.
// Config.NewXxxArgs() returns XxxArgs based on configuration values.
// Separating Config and Args is also a good idea for post-processing.

// Does a user have to know about XxxArgs?

// ComboPosition: 0.40 * ScreenSize.Y

// type Config struct { ... }
// func (cfg Config) Width() float64 { ... }

// Both Config and Args are exported.
// Args may be useful for showing actual values in-game.
// type XxxConfig struct { ... }
// type XxxArgs struct { ... }
// func (cfg XxxConfig) Args() XxxArgs { ... }
// func NewXxxComponent(args XxxArgs) XxxComponent { ... }
// }

// Config is for wrapping required arguments.
type Config struct {
	CommonConfig
	StageConfig
	BarConfig
	NoteConfig
	LightConfig
	JudgmentConfig
	mode.ComboConfig
	mode.ScoreConfig
}

func NewConfig() Config {
}

// virtual configuration; no component
// KeyCount, ScratchMode and some methods
type CommonConfig struct {
	mode.ScreenConfig
	mode.MusicConfig

	KeyCount       int
	ScratchModes   map[int]int
	StageWidths    map[int]float64
	KeyKindsList   map[int][]KeyKind
	KeyKindWidths  [4]float64 // * screen.Size.X
	StagePositionX float64    // * screen.Size.X
	HitPositionY   float64    // * screen.Size.Y
}

const (
	ScratchModeNone = iota
	ScratchModeLeft
	ScratchModeRight
)

type KeyKind int

const (
	One KeyKind = iota
	Two
	Mid
	Tip
)

func NewCommonConfig(keyCount int) CommonConfig {
	return CommonConfig{
		KeyCount: keyCount,
		ScratchModes: map[int]int{
			8: ScratchModeLeft,
		},
		StageWidths: map[int]float64{ // * screen.Size.X
			4: 0.50,
			5: 0.55,
			6: 0.60,
			7: 0.65,
			8: 0.70,
			9: 0.75,
		},
		KeyKindsList: map[int][]KeyKind{
			0:  {},
			1:  {Mid},
			2:  {One, One},
			3:  {One, Mid, One},
			4:  {One, Two, Two, One},
			5:  {One, Two, Mid, Two, One},
			6:  {One, Two, One, One, Two, One},
			7:  {One, Two, One, Mid, One, Two, One},
			8:  {Tip, One, Two, One, One, Two, One, Tip},
			9:  {Tip, One, Two, One, Mid, One, Two, One, Tip},
			10: {Tip, One, Two, One, Mid, Mid, One, Two, One, Tip},
		},
		KeyKindWidths: [4]float64{
			0.055, // One
			0.055, // Two
			0.055, // Mid
			0.055, // Tip
		},
		StagePositionX: 0.50,
		HitPositionY:   0.90,
	}
}

func (cfg CommonConfig) StageWidth() float64 {
	const defaultKeyCount = 4
	if w, ok := cfg.StageWidths[cfg.KeyCount]; ok {
		return w
	}
	return cfg.StageWidths[defaultKeyCount]
}

// I'm personally proud of this code.
func (cfg CommonConfig) keyKinds() []KeyKind {
	switch cfg.ScratchModes[cfg.KeyCount] {
	case ScratchModeNone:
		return cfg.KeyKindsList[cfg.KeyCount]
	case ScratchModeLeft:
		return append([]KeyKind{Tip}, cfg.KeyKindsList[cfg.KeyCount-1]...)
	case ScratchModeRight:
		return append(cfg.KeyKindsList[cfg.KeyCount-1], Tip)
	}
	return nil
}

func (cfg CommonConfig) KeyWidths() (ws []float64) {
	ws = make([]float64, cfg.KeyCount)
	for k, kk := range cfg.keyKinds() {
		w := cfg.KeyKindWidths[kk]
		ws[k] = w
	}

	// Adjust key width to fit the stage width.
	var rawSum float64
	for _, w := range ws {
		rawSum += w
	}
	scale := cfg.StageWidth() / rawSum
	for k := range ws {
		ws[k] *= scale
	}
	return ws
}

// KeyPositionXs returns centered x positions.
func (cfg CommonConfig) KeyPositionXs() (xs []float64) {
	xs = make([]float64, cfg.KeyCount)
	ws := cfg.KeyWidths()
	x := cfg.StagePositionX - cfg.StageWidth()/2
	for k, w := range ws {
		x += w / 2
		xs[k] = x
		x += w / 2
	}
	return
}

// NoteExposureDuration returns time in milliseconds
// that cursor takes to move 1 logical pixel.
func (cfg CommonConfig) NoteExposureDuration(speed float64) int32 {
	return int32(cfg.HitPositionY / speed)
}
