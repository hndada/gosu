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
var Speed = 0.7
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
	BgDimness        float64 = 0.5
	ComboScale       float64 = 0.72
	ComboPosition    float64 = screenSizeY * 0.45
	ComboGap         float64 = screenSizeX * -0.001
	ScoreScale       float64 = 0.67
	JudgmentScale    float64 = 0.35
	JudgmentPosition float64 = screenSizeY * 0.67
	NoteHeigth       float64 = screenSizeY * 0.04 // Applies to all notes
	FieldDark        float64 = 0.8
	HintPosition     float64 = screenSizeY * 0.96 // The middle position of Judge line, not a topmost.
	HintHeight       float64 = screenSizeY * 0.04

	TimingMeterUnit = 5 // The number of pixels per 1ms
)
