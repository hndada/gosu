package drum

import (
	"github.com/hndada/gosu/defaultskin"
	"github.com/hndada/gosu/mode"
)

const (
	TPS = mode.TPS

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
	KeySettings map[int][]string
	SpeedScale  float64
	HitPosition float64
	maxPosition float64
	minPosition float64
	Reverse     bool

	// Skin-independent settings
	FieldOpaque       float64
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

		FieldOpaque:     0.7,
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
	DefaultSettings = NewSettings()
	UserSettings    = NewSettings()
	S               = &UserSettings
)

func init() {
	DefaultSettings.process()
	UserSettings.process()
	DefaultSkin.Load(defaultskin.FS)
	UserSkin.Load(defaultskin.FS)
}
func (settings *Settings) Load(src Settings) {
	*settings = src
	defer settings.process()

	for k := range settings.KeySettings {
		mode.NormalizeKeys(settings.KeySettings[k])
	}
	mode.Normalize(&settings.SpeedScale, 0.1, 2.0)
	mode.Normalize(&settings.HitPosition, 0, 1)
	// Reverse: bool

	mode.Normalize(&settings.FieldOpaque, 0, 1)
	mode.Normalize(&settings.FieldPosition, 0, 1)
	mode.Normalize(&settings.FieldHeight, 0, 1)
	mode.Normalize(&settings.DancerPositionX, 0, 1)
	mode.Normalize(&settings.DancerPositionY, 0, 1)

	mode.Normalize(&settings.FieldInnerHeight, 0, 1)
	mode.Normalize(&settings.JudgmentScale, 0, 2)
	mode.Normalize(&settings.DotScale, 0, 2)
	mode.Normalize(&settings.ShakeScale, 0, 2)
	mode.Normalize(&settings.DancerScale, 0, 2)
	mode.Normalize(&settings.ComboScale, 0, 2)
	mode.Normalize(&settings.ComboDigitGap, -0.005, 0.005)
}

func (settings *Settings) process() {
	s := settings

	s.volumeMusic = &mode.S.VolumeMusic
	s.volumeSound = &mode.S.VolumeSound
	s.offset = &mode.S.Offset
	s.backgroundBrightness = &mode.S.BackgroundBrightness

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
