package gosu

// Logical size of in-game screen
const (
	screenSizeX = 1600
	screenSizeY = 900
)

var (
	WindowSizeX = 1600
	WindowSizeY = 900
	MaxTPS      = 1000 // MaxTPS should be 1000 or greater.
	Volume      = 0.25
)
var MusicPath = "music"
var KeySettings = map[int][]Key{
	4:               {KeyD, KeyF, KeyJ, KeyK},
	5:               {KeyD, KeyF, KeySpace, KeyJ, KeyK},
	6:               {KeyS, KeyD, KeyF, KeyJ, KeyK, KeyL},
	7:               {KeyS, KeyD, KeyF, KeySpace, KeyJ, KeyK, KeyL},
	8 + LeftScratch: {KeyA, KeyS, KeyD, KeyF, KeySpace, KeyJ, KeyK, KeyL},
	8:               {KeyA, KeyS, KeyD, KeyF, KeyJ, KeyK, KeyL, KeySemicolon},
	9:               {KeyA, KeyS, KeyD, KeyF, KeySpace, KeyJ, KeyK, KeyL, KeySemicolon},
	10:              {KeyA, KeyS, KeyD, KeyF, KeyV, KeyN, KeyJ, KeyK, KeyL, KeySemicolon},
}
var BaseSpeed = 0.7
var NoteWidthsMap = map[int][4]float64{
	4:  {0.065, 0.065, 0.065, 0.065},
	5:  {0.065, 0.065, 0.065, 0.065},
	6:  {0.065, 0.065, 0.065, 0.065},
	7:  {0.06, 0.06, 0.06, 0.06},
	8:  {0.06, 0.06, 0.06, 0.06},
	9:  {0.06, 0.06, 0.06, 0.06},
	10: {0.06, 0.06, 0.06, 0.06},
}

func init() {
	ScaleNoteWidthsMap()
}

// Todo: generalize setting loading function
func ScaleNoteWidthsMap() {
	for k, ws := range NoteWidthsMap {
		ws2 := ws
		for i, w := range ws2 {
			ws2[i] = screenSizeX * w
		}
		NoteWidthsMap[k] = ws2
	}
}

var (
	BgDimness float64 = 0.5
	FieldDark float64 = 0.8
)

// Todo: Should NoteHeight be separated into NoteHeight, HeadHeight, TailHeight?
var (
	ComboPosition    float64 = screenSizeY * 0.45
	JudgmentPosition float64 = screenSizeY * 0.67
	NoteHeigth       float64 = screenSizeY * 0.04 // Applies to all notes
	HitPosition      float64 = screenSizeY * 0.96 // The bottom y-value of Hint,  not a middle or top.
	// HintPosition     float64 = screenSizeY * 0.96
)

const (
	BodySpriteStyleStretch = iota
	BodySpriteStyleAttach
)

// Skin scale settings
// Todo: make the struct SkinScaleSettings
var (
	ComboScale    float64 = 0.72
	ComboGap      float64 = screenSizeX * -0.0008
	ScoreScale    float64 = 0.67
	JudgmentScale float64 = 0.35
	HintHeight    float64 = screenSizeY * 0.04
	CursorScale   float64 = 0.1

	BodySpriteStyle = BodySpriteStyleStretch
)
var TimingMeterUnit = 5 // The number of pixels per 1ms

// 1 pixel is 1 millisecond.
func ExposureTime(speed float64) float64 { return (HitPosition + NoteHeigth/2) / speed }
