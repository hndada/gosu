package osu

import "image/color"

type FormatOsu struct {
	FormatVersion int
	General
	Editor
	Metadata
	Difficulty
	Events
	TimingPoints
	Colours
	HitObjects
}

type General struct { // delimiter:(space)
	AudioFilename            string
	AudioLeadIn              int
	AudioHash                string // deprecated
	PreviewTime              int
	Countdown                int // nofloat
	SampleSet                string
	StackLeniency            float64
	Mode                     int // nofloat
	LetterboxInBreaks        bool
	StoryFireInFront         bool // deprecated
	UseSkinSprites           bool
	AlwaysShowPlayfield      bool // deprecated
	OverlayPosition          string
	SkinPreference           string
	EpilepsyWarning          bool
	CountdownOffset          int
	SpecialStyle             bool
	WidescreenStoryboard     bool
	SamplesMatchPlaybackRate bool
}
type Editor struct { // delimiter:(space)
	Bookmarks       []int // delimiter,
	DistanceSpacing float64
	BeatDivisor     float64
	GridSize        int
	TimelineZoom    float64
}
type Metadata struct { // delimiter:
	Title         string
	TitleUnicode  string
	Artist        string
	ArtistUnicode string
	Creator       string
	Version       string
	Source        string
	Tags          []string // delimiter(space)
	BeatmapID     int
	BeatmapSetID  int
}
type Difficulty struct { // delimiter:
	HPDrainRate       float64
	CircleSize        float64
	OverallDifficulty float64
	ApproachRate      float64
	SliderMultiplier  float64
	SliderTickRate    float64
}
type Events []Event

// storyboard not implemented yet
type Event struct { // delimiter,
	Type      string
	StartTime int
	EndTime   int // optional
	Filename  string
	XOffset   int
	YOffset   int
}

type TimingPoints []TimingPoint
type TimingPoint struct { // delimiter,
	Time int
	// Bpm, SpeedScale float64 // todo: method
	BeatLength  float64
	Meter       int
	SampleSet   int // nofloat
	SampleIndex int // nofloat
	Volume      int
	Uninherited bool
	Effects     int // nofloat
	// Kiai        bool // todo: method
}
type Colours struct { // manual
	Combos              [4]color.RGBA
	SliderTrackOverride color.RGBA
	SliderBorder        color.RGBA
}
type HitObjects []HitObject
type HitObject struct { // delimiter,
	X            int
	Y            int
	Time    int
	NoteType     int          // nofloat
	HitSound     int          // nofloat
	EndTime      int          // optional
	SliderParams SliderParams // optional
	HitSample    HitSample    // optional
}
type SliderParams struct { // delimiter,
	CurveType   string   // one letter
	CurvePoints [][2]int // delimiter| // delimiter: // slice of paired integers
	Slides      int
	Length      float64
	EdgeSounds  []int    // delimiter|
	EdgeSets    [][2]int // delimiter| // delimiter:
}
type HitSample struct { // delimiter:
	NormalSet   int // nofloat
	AdditionSet int // nofloat
	Index       int // nofloat
	Volume      int
	Filename    string
}
