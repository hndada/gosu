package drum

import (
	"github.com/hndada/gosu/assets"
	"github.com/hndada/gosu/mode"
)

const (
	TPS         = mode.TPS
	ScreenSizeX = mode.ScreenSizeX
	ScreenSizeY = mode.ScreenSizeY
)

const positionMargin = 100

type Settings struct {
	MusicVolume          float64
	musicVolume          *float64
	volumeSound          *float64
	offset               *int64
	backgroundBrightness *float64
	delayedJudge         *int64 // for HCI experiment
	debugPrint           *bool

	// Logic settings
	KeySettings map[int][]string
	SpeedScale  float64
	HitPosition float64
	maxPosition float64
	minPosition float64
	Reverse     bool

	// Skin-independent settings
	FieldOpacity      float64
	FieldPosition     float64
	FieldHeight       float64
	bigNoteHeight     float64
	regularNoteHeight float64
	DancerPositionX   float64
	DancerPositionY   float64

	// Skin-dependent settings
	FieldInnerHeight float64
	JudgmentScale    float64
	DotScale         float64
	ShakeScale       float64
	DancerScale      float64
	ComboScale       float64
	ComboDigitGap    float64
}

func NewSettings() Settings {
	return Settings{
		KeySettings: map[int][]string{4: {"D", "F", "J", "K"}},
		SpeedScale:  1.0,
		HitPosition: 0.1875,
		Reverse:     false,

		FieldOpacity:    0.7,
		FieldPosition:   0.4115,
		FieldHeight:     0.26,
		DancerPositionX: 0.1,
		DancerPositionY: 0.175,

		FieldInnerHeight: 0.95,
		JudgmentScale:    0.75,
		DotScale:         0.5,
		ShakeScale:       1,
		DancerScale:      0.75,
		ComboScale:       0.75,
		ComboDigitGap:    -0.001,
	}
}

var (
	UserSettings = NewSettings()
	S            = &UserSettings
)

func init() {
	S.process()
	DefaultSkin.Load(assets.FS)
	UserSkin.Load(assets.FS)
}
func (s *Settings) Load(src Settings) {
	*s = src
	defer s.process()

	for k := range s.KeySettings {
		S.KeySettings[k] = mode.NormalizeKeys(s.KeySettings[k])
	}
	mode.Normalize(&s.SpeedScale, 0.1, 2.0)
	mode.Normalize(&s.HitPosition, 0, 1)
	// Reverse: bool

	mode.Normalize(&s.FieldOpacity, 0, 1)
	mode.Normalize(&s.FieldPosition, 0, 1)
	mode.Normalize(&s.FieldHeight, 0, 1)
	mode.Normalize(&s.DancerPositionX, 0, 1)
	mode.Normalize(&s.DancerPositionY, 0, 1)

	mode.Normalize(&s.FieldInnerHeight, 0, 1)
	mode.Normalize(&s.JudgmentScale, 0, 2)
	mode.Normalize(&s.DotScale, 0, 2)
	mode.Normalize(&s.ShakeScale, 0, 2)
	mode.Normalize(&s.DancerScale, 0, 2)
	mode.Normalize(&s.ComboScale, 0, 2)
	mode.Normalize(&s.ComboDigitGap, -0.005, 0.005)
}

func (s *Settings) process() {
	s.musicVolume = &mode.S.MusicVolume
	s.volumeSound = &mode.S.SoundVolume
	s.offset = &mode.S.MusicOffset
	s.backgroundBrightness = &mode.S.BackgroundBrightness
	s.delayedJudge = &mode.S.DelayedJudge
	s.debugPrint = &mode.S.DebugPrint

	s.HitPosition *= ScreenSizeX
	s.minPosition = -s.HitPosition - positionMargin
	s.maxPosition = -s.HitPosition + ScreenSizeX + positionMargin
	if s.Reverse {
		max, min := s.maxPosition, s.minPosition
		s.maxPosition = -min
		s.minPosition = -max
	}

	s.FieldPosition *= ScreenSizeY
	s.FieldHeight *= ScreenSizeY
	s.bigNoteHeight = s.FieldHeight * 0.725
	s.regularNoteHeight = s.bigNoteHeight * 0.65
	s.DancerPositionX *= ScreenSizeX
	s.DancerPositionY *= ScreenSizeY

	s.FieldInnerHeight *= s.FieldHeight
	s.ComboDigitGap *= ScreenSizeX
}

func ExposureTime(speed float64) float64 {
	return (ScreenSizeX - S.HitPosition) / speed
}
