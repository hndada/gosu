package piano

import (
	"image/color"

	"github.com/hndada/gosu/defaultskin"
	"github.com/hndada/gosu/mode"
)

const (
	TPS         = mode.TPS
	ScreenSizeX = mode.ScreenSizeX
	ScreenSizeY = mode.ScreenSizeY
)

const positionMargin = 100

type Settings struct {
	volumeMusic          *float64
	volumeSound          *float64
	offset               *int64
	backgroundBrightness *float64

	// Logic settings
	KeySettings   map[int][]string
	SpeedScale    float64
	HitPosition   float64 // The bottom y-value of Hint,  not a middle or top.
	maxPosition   float64
	minPosition   float64
	TailExtraTime float64
	ReverseBody   bool

	// Skin-independent settings
	NoteWidths         map[int][4]float64 // Fourth is a Scratch note.
	NoteHeigth         float64            // Applies to all types of notes.
	BodyStyle          int
	FieldPosition      float64
	ComboPosition      float64
	JudgmentPosition   float64
	ScratchColor       [4]uint8
	scratchColor       color.NRGBA
	FieldOpaque        float64
	KeyLightingColors  [4][4]uint8
	keyLightingColors  [4]color.NRGBA
	HitLightingOpaque  float64
	HoldLightingOpaque float64

	// Skin-dependent settings
	ComboScale    float64
	ComboDigitGap float64
	JudgmentScale float64
	HintHeight    float64
	LightingScale float64
}

// Fields which types are map should be explicitly make new map.
func NewSettings() Settings {
	return Settings{
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
		NoteHeigth:       0.05,
		BodyStyle:        BodyStyleStretch,
		FieldPosition:    0.50,
		ComboPosition:    0.40,
		JudgmentPosition: 0.66,
		ScratchColor:     [4]uint8{224, 0, 0, 255},
		FieldOpaque:      0.8,
		KeyLightingColors: [4][4]uint8{
			{224, 224, 224, 64}, // white
			{255, 170, 204, 64}, // pink
			{224, 224, 0, 64},   // yellow
			{255, 0, 0, 64},     // red
		},
		HitLightingOpaque:  0.5,
		HoldLightingOpaque: 0.8,

		ComboScale:    0.75,
		ComboDigitGap: -0.0008,
		JudgmentScale: 0.33,
		HintHeight:    0.055,
		LightingScale: 1.0,
	}
}

// Generic function seems not allow to pass named type.
const (
	BodyStyleStretch int = iota
	BodyStyleAttach
)

var (
	UserSettings = NewSettings()
	S            = &UserSettings
)

func init() {
	S.process()
	DefaultSkins.Load(defaultskin.FS)
	UserSkins.Load(defaultskin.FS)
}
func (s *Settings) Load(src Settings) {
	*S = src
	defer S.process()

	for k := range S.KeySettings {
		mode.NormalizeKeys(S.KeySettings[k])
	}
	mode.Normalize(&S.SpeedScale, 0.1, 2.0)
	mode.Normalize(&S.HitPosition, 0, 1)
	mode.Normalize(&S.TailExtraTime, -150, 150)
	// ReverseBody: bool

	for k, widths := range S.NoteWidths {
		for kind := range widths {
			mode.Normalize(&widths[kind], 0.01, 0.15)
		}
		S.NoteWidths[k] = widths
	}
	mode.Normalize(&S.NoteHeigth, 0, 0.15)
	mode.Normalize(&S.BodyStyle, 0, BodyStyleAttach)
	mode.Normalize(&S.FieldPosition, 0, 1)
	mode.Normalize(&S.ComboPosition, 0, 1)
	mode.Normalize(&S.JudgmentPosition, 0, 1)
	// ScratchColor: [4]uint8
	mode.Normalize(&S.FieldOpaque, 0, 1)
	// KeyLightingColors: [4][4]uint8
	mode.Normalize(&S.HitLightingOpaque, 0, 1)
	mode.Normalize(&S.HoldLightingOpaque, 0, 1)

	mode.Normalize(&S.ComboScale, 0, 2)
	mode.Normalize(&S.ComboDigitGap, -0.005, 0.005)
	mode.Normalize(&S.JudgmentScale, 0, 2)
	mode.Normalize(&S.HintHeight, 0, 0.1)
	mode.Normalize(&S.LightingScale, 0, 2)
}

// It is safe to use mode.UserSettings even for DefaultSettings:
// mode.UserSettings = mode.DefaultSettings when processing default.
func (s *Settings) process() {
	*s = NewSettings()
	s.volumeMusic = &mode.S.VolumeMusic
	s.volumeSound = &mode.S.VolumeSound
	s.offset = &mode.S.Offset
	s.backgroundBrightness = &mode.S.BackgroundBrightness

	s.HitPosition *= ScreenSizeY
	max := s.HitPosition
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
	{
		clr := s.ScratchColor
		s.scratchColor = color.NRGBA{clr[0], clr[1], clr[2], clr[3]}
	}
	for i, clr := range s.KeyLightingColors {
		s.keyLightingColors[i] = color.NRGBA{clr[0], clr[1], clr[2], clr[3]}
	}
	s.ComboDigitGap *= ScreenSizeX
	s.HintHeight *= ScreenSizeY
}

// 1 pixel is 1 millisecond.
func ExposureTime(speed float64) float64 { return S.HitPosition / speed }
