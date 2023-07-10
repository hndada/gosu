package mode

// ScreenSize is a logical size of in-game screen.
const (
	ScreenSizeX = 1600
	ScreenSizeY = 900
)

type SettingsType struct {
	MusicVolume          float64
	SoundVolume          float64
	BackgroundBrightness float64
	Offset               int32
	DebugPrint           bool

	ScoreScale    float64
	ScoreDigitGap float64

	MeterUnit   float64 // number of pixels per 1ms
	MeterHeight float64
}

var Settings = SettingsType{
	MusicVolume:          0.50,
	SoundVolume:          0.50,
	BackgroundBrightness: 0.6,
	Offset:               -20,
	DebugPrint:           true,

	ScoreScale:    0.65,
	ScoreDigitGap: 0,

	MeterUnit:   4,
	MeterHeight: 65,
}
