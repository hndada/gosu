package gosu

// Logical size of in-game screen
const (
	screenSizeX = 800
	screenSizeY = 600
)

var (
	WindowSizeX = 800
	WindowSizeY = 600
	MaxTPS      = 1000 // MaxTPS should be 1000 or greater.
	Volume      = 0.05
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
var Speed = 0.75
var (
	BgDimness        float64 = 0.3
	ComboScale       float64 = 0.65
	ComboPosition    float64 = 240
	ComboGap         float64 = -2
	ScoreScale       float64 = 0.55
	JudgmentScale    float64 = 0.35
	JudgmentPosition float64 = 400
	ClearWidth       float64 = 225

	NoteWidths = map[int][4]float64{
		4:  {50, 50, 50, 50},
		5:  {50, 50, 50, 50},
		6:  {50, 50, 50, 50},
		7:  {50, 50, 50, 50},
		8:  {45, 45, 45, 45},
		9:  {45, 45, 45, 45},
		10: {45, 45, 45, 45},
	}
	NoteHeigth   float64 = 30 // Applies all notes
	FieldDark    float64 = 0.95
	HintPosition float64 = 590 // The middle position of Judge line, not a topmost.
	HintHeight   float64 = 5
)
