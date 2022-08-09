package main

var (
	MaxTPS      = 1000
	ScreenSizeX = 800
	ScreenSizeY = 600

	Volume      = 0.05
	Speed       = 0.16
	KeySettings = map[int][]Code{
		4:               {CodeD, CodeF, CodeJ, CodeK},
		5:               {CodeD, CodeF, CodeSpacebar, CodeJ, CodeK},
		6:               {CodeS, CodeD, CodeF, CodeJ, CodeK, CodeL},
		7:               {CodeS, CodeD, CodeF, CodeSpacebar, CodeJ, CodeK, CodeL},
		8 + LeftScratch: {CodeA, CodeS, CodeD, CodeF, CodeSpacebar, CodeJ, CodeK, CodeL},
		8:               {CodeA, CodeS, CodeD, CodeF, CodeJ, CodeK, CodeL, CodeSemicolon},
		9:               {CodeA, CodeS, CodeD, CodeF, CodeSpacebar, CodeJ, CodeK, CodeL, CodeSemicolon},
		10:              {CodeA, CodeS, CodeD, CodeF, CodeV, CodeN, CodeJ, CodeK, CodeL, CodeSemicolon},
	}
	// Scaled to 800 x 600.
	NoteWidths = map[int][4]float64{
		4:  {50, 50, 50, 50},
		5:  {50, 50, 50, 50},
		6:  {50, 50, 50, 50},
		7:  {50, 50, 50, 50},
		8:  {45, 45, 45, 45},
		9:  {45, 45, 45, 45},
		10: {45, 45, 45, 45},
	}
	NoteHeigth    float64 = 30 // Applies all notes
	ComboPosition float64 = 180
	JudgePosition float64 = 250
	HintPosition  float64 = 550 // The middle position of Judge line, not a topmost.
	FieldDark     float64 = 0.95
	BgDimness     float64 = 0.3
	// ScratchMode map[int]int
	ComboWidth    float64 = 40
	ScoreWidth    float64 = 33
	ComboGap      float64 = -2
	ScoreGap      float64 = -2
	HintHeight    float64 = 5
	JudgmentWidth float64 = 65
)

// Scale returns scaled value based on screen size
// func Scale(v float64) int     { return int(v * DisplayScale()) }
// func DisplayScale() float64   { return float64(ScreenSizeY) / 100 }
// func ScreenSize() image.Point { return image.Pt(ScreenSizeX, ScreenSiz eY) }
func Scale() float64 { return float64(ScreenSizeY) / 800 } // Value of Scale() is 1 in 800 x 600