package piano

import (
	"fmt"
	"image/color"

	"github.com/BurntSushi/toml"
	"github.com/hndada/gosu/defaultskin"
	"github.com/hndada/gosu/mode"
)

const (
	ScreenSizeX = mode.ScreenSizeX
	ScreenSizeY = mode.ScreenSizeY
)

// positionMargin should be large enough:
// more than MaxSize/2 of all note sprites' width or height.
const positionMargin = 100

type Settings struct {
	// Logic settings
	KeySettings   map[int][]string
	SpeedScale    float64
	HitPosition   float64 // The bottom y-value of Hint,  not a middle or top.
	maxPosition   float64
	minPosition   float64
	TailExtraTime float64
	ReverseBody   bool

	// Skin-independent settings
	NoteWidths        map[int][4]float64 // Fourth is a Scratch note.
	NoteHeigth        float64            // Applies to all types of notes.
	BodyStyle         int
	FieldPosition     float64
	ComboPosition     float64
	JudgmentPosition  float64
	ScratchColor      [4]uint8
	scratchColor      color.NRGBA
	FieldOpaque       float64
	KeyLightingOpaque float64
	HitLightingOpaque float64

	// Skin-dependent settings
	ComboScale    float64
	ComboDigitGap float64
	JudgmentScale float64
	HintHeight    float64
	LightingScale float64
}

var (
	DefaultSettings = Settings{
		KeySettings: map[int][]string{
			4:               {"D", "F", "J", "K"},
			5:               {"D", "F", "Space", "J", "K"},
			6:               {"S", "D", "F", "J", "K", "L"},
			7:               {"S", "D", "F", "Space", "J", "K", "L"},
			8 + LeftScratch: {"A", "S", "D", "F", "Space", "J", "K", "L"},
			8:               {"A", "S", "D", "F", "J", "K", "L", "Semicolon"},
			9:               {"A", "S", "D", "F", "Space", "J", "K", "L", "Semicolon"},
			10:              {"A", "S", "D", "F", "V", "N", "J", "K", "L", "Semicolon"},
		},
		SpeedScale:    1.0,
		HitPosition:   0.90,
		TailExtraTime: 0,
		ReverseBody:   false,

		NoteWidths: map[int][4]float64{
			4:  {0.065, 0.065, 0.065, 0.065},
			5:  {0.065, 0.065, 0.065, 0.065},
			6:  {0.065, 0.065, 0.065, 0.065},
			7:  {0.06, 0.06, 0.06, 0.06},
			8:  {0.06, 0.06, 0.06, 0.06},
			9:  {0.06, 0.06, 0.06, 0.06},
			10: {0.06, 0.06, 0.06, 0.06},
		},
		NoteHeigth:        0.05,
		BodyStyle:         BodyStyleStretch,
		FieldPosition:     0.50,
		ComboPosition:     0.40,
		JudgmentPosition:  0.66,
		ScratchColor:      [4]uint8{224, 0, 0, 255},
		FieldOpaque:       0.8,
		KeyLightingOpaque: 0.5,
		HitLightingOpaque: 1,

		ComboScale:    0.75,
		ComboDigitGap: -0.0008,
		JudgmentScale: 0.33,
		HintHeight:    0.04,
		LightingScale: 1.0,
	}
	UserSettings = DefaultSettings
)

// Generic function seems not allow to pass named type.
const (
	BodyStyleStretch int = iota
	BodyStyleAttach
)

func init() {
	DefaultSettings.process()
	DefaultSkins.Load(defaultskin.FS)
}

func (settings *Settings) Load(data string) {
	_, err := toml.Decode(data, settings)
	if err != nil {
		fmt.Println(err)
	}
	defer settings.process()

	for k := range settings.KeySettings {
		mode.NormalizeKeys(settings.KeySettings[k])
	}
	mode.Normalize(&settings.SpeedScale, 0.1, 2.0)
	mode.Normalize(&settings.HitPosition, 0, 1)
	mode.Normalize(&settings.TailExtraTime, -150, 150)
	// ReverseBody: bool

	for k, widths := range settings.NoteWidths {
		for kind := range widths {
			mode.Normalize(&widths[kind], 0, 0.15)
		}
		settings.NoteWidths[k] = widths
	}
	mode.Normalize(&settings.NoteHeigth, 0, 0.15)
	mode.Normalize(&settings.BodyStyle, 0, BodyStyleAttach)
	mode.Normalize(&settings.FieldPosition, 0, 1)
	mode.Normalize(&settings.ComboPosition, 0, 1)
	mode.Normalize(&settings.JudgmentPosition, 0, 1)
	// ScratchColor: [4]uint8
	mode.Normalize(&settings.FieldOpaque, 0, 1)
	mode.Normalize(&settings.KeyLightingOpaque, 0, 1)
	mode.Normalize(&settings.HitLightingOpaque, 0, 1)

	mode.Normalize(&settings.ComboScale, 0, 1.5)
	mode.Normalize(&settings.ComboDigitGap, -0.005, 0.005)
	mode.Normalize(&settings.JudgmentScale, 0, 1.5)
	mode.Normalize(&settings.HintHeight, 0, 0.1)
	mode.Normalize(&settings.LightingScale, 0, 1.5)
}

func (settings *Settings) process() {
	s := settings
	s.HitPosition *= ScreenSizeX
	max := ScreenSizeY * s.HitPosition
	s.maxPosition = max + positionMargin
	s.minPosition = max - ScreenSizeY - positionMargin
	if s.ReverseBody {
		max, min := s.maxPosition, s.minPosition
		s.maxPosition = -min
		s.minPosition = -max
	}

	for k, widths := range s.NoteWidths {
		for kind := range widths {
			widths[kind] *= ScreenSizeX
		}
		s.NoteWidths[k] = widths
	}
	s.NoteHeigth *= ScreenSizeY
	s.FieldPosition *= ScreenSizeX
	s.ComboPosition *= ScreenSizeY
	s.JudgmentPosition *= ScreenSizeY
	s.scratchColor = color.NRGBA{
		s.ScratchColor[0], s.ScratchColor[1],
		s.ScratchColor[2], s.ScratchColor[3]}
	s.ComboDigitGap *= ScreenSizeX
	s.HintHeight *= ScreenSizeY
}
